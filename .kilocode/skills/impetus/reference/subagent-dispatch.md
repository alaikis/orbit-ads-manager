# Subagent 驱动开发的 Impetus 扩展

规范路径：`impetus/reference/subagent-dispatch.md`

本文档提供在 Superpowers `subagent-driven-development` 技能**之上**应用的 Impetus 专属扩展。Impetus 将它用作依赖图批量调度器：先构建 task 依赖图，将相互独立的 task 按 wave 并行派发，然后审查并勾选已完成 task，再解锁后续依赖 wave。本文档添加 Impetus 特有的真实后台调度、任务追踪、状态验证和上下文恢复。若 Superpowers 技能与本文档发生冲突时，以本文档中更具体的 Impetus 约束为准。

> **关键约束 — Wave 之间禁止暂停**
>
> 当一个 wave 中可运行 task 通过审查并被勾选后，**立即计算并派发下一个可运行 wave**，不得停止、总结或询问用户是否继续。用户期望所有符合条件的 task 自动执行。wave 之间暂停会中断工作流，导致用户每次都需要手动恢复。
>
> 仅在以下情况才停止并等待用户输入：
> - 任务处于 **BLOCKED** 状态（3 轮审查-修复仍未通过）
> - 存在无法从仓库、计划或既有上下文消除的真实歧义
> - 平台没有真实后台 agent 调度能力，需要用户改选 `executing-plans`
> - 剩余未勾选 task 全部被失败依赖阻塞
> - 用户**明确**要求暂停
>
> 此规则适用于整个 wave 调度器，而非单个任务。

## 开始前

1. 读取计划一次，按顺序提取所有未勾选 task 的完整文本。
2. 为每个 task 保存唯一标识：plan 中 checkbox 后的完整任务文本，以及它映射的 OpenSpec task 完整文本（若存在）。若文本不唯一，停止并先修正计划，禁止依赖"第一个匹配项"。
3. 构建依赖图：
   - 依赖边来自显式的 "depends on"、"after"、"requires"、编号顺序、架构约束，以及从 plan 推断出的文件/API 依赖。
   - 只有所有前置依赖都已勾选后，task 才能进入某个 wave。
   - 若依赖方向不明确，按 plan 顺序保持串行。
4. 为每个 task 分配写入范围，来源包括 plan 文件清单和推断会触碰的文件。写入范围重叠、共享生成文件、schema/API 契约耦合、迁移/测试 fixture 耦合的 task，不得进入同一个 wave，除非 plan 明确说明它们独立。
5. 确认平台对并行 implementer 的隔离模型。同一个 wave 内每个 implementer 必须运行在独立 worktree / branch / sandbox 中，或返回 patch 由协调者串行应用。若平台只能让 agent 在同一个实时工作树和分支上运行，则 wave size 必须降为 1。
6. 确认提交授权。默认情况下，`subagent-driven-development` 仍遵循 Impetus 的一次最终 build 提交策略：implementer 返回 patch 或隔离 branch/worktree 的结果，协调者在整个 build pass 与实现 Review 窗口通过后统一提交一次。只有用户在选择 `subagent-driven-development` 时明确预授权 task/wave checkpoint 提交，或协调者在该 checkpoint 提交前按 `impetus/reference/decision-point.md` 暂停并获得确认，才允许中间提交。后台 agent 不得绕过提交协议。
7. 计算第一个可运行 wave：依赖已完成、写入范围不冲突的未勾选 task。宁可 wave 小一点，也不要冒险并行；不得超过平台实际可承载的后台 agent 数量。

## Wave 调度

协调者按下列循环执行，直到所有 plan task 都勾选：

1. 从依赖图中选择当前可运行 wave。
2. 为 wave 中每个 task 派发一个全新的后台 implementer，每个 implementer 必须使用隔离工作区或 patch-return 模式。
3. 等待 wave 中所有 implementer 返回，或某个 task 进入 `BLOCKED`。
4. 将返回结果串行集成到协调者工作树。禁止多个后台 implementer 直接并发提交到同一个实时分支/工作树。
5. 确认每个返回的提交、patch 和变更文件在当前工作树可见。
6. 审查前先检测 wave 冲突：
   - 多个 implementer 修改同一文件
   - API/schema/config 变更不兼容
   - 合并 wave 输出后构建或测试失败
7. 对 wave 中每个已完成 task，派发 spec compliance review 与 code quality review。不同 task 的 reviewer 可以并行；同一 task 的两个 reviewer 也可在实现 diff 可见后并行。
8. 只修复失败的 task。失败 task 会阻塞其依赖 task，但依赖仍满足的独立 task 可继续进入下一 wave。
9. 审查通过的 task 勾选、持久化 checkpoint；只有在明确授权时才提交 checkpoint，然后计算下一可运行 wave。

## 每个 Task 的 Impetus 扩展

在 wave 内每个 task 上应用这些扩展，叠加在 Superpowers 技能的派发循环之上：

### 0. 派发强制约束（关键）

主会话**仅负责协调**，禁止直接执行 task。主会话禁止修改源代码。协调者唯一允许的文件修改是 plan、OpenSpec task 和 subagent 进度检查点的持久化更新。不得把多个 task 打包给同一个 agent。当前 wave 中每个 task 都派发一个全新的后台 implementer agent，spec reviewer、code quality reviewer、修复 agent 和 final reviewer 也必须分别使用全新的后台 agent：

- **Claude Code**：对每个 implementer、spec reviewer、code quality reviewer、修复 agent 和 final reviewer 使用 `Agent` 工具并设置 `run_in_background: true`。禁止内联执行 task，禁止错误进入需要预先创建 team 的团队模式。
- **其他平台**：使用平台等效的后台 agent / Task / 多 agent 派发机制。
- **禁止**跨 task 或角色复用 implementer、reviewer 或修复 agent。每个 agent 拥有全新的隔离上下文，并且只接收当前角色所需的单个 task 上下文。
- **禁止**为了增加并行度，把有依赖或写入冲突的 task 派到同一个 wave。
- **禁止**多个 implementer 运行在同一个实时工作树/分支上。没有隔离工作区或 patch-return 模式时，调度器必须降级为每个 wave 只有一个 task。
- 若平台无真实后台派发能力，不得继续；暂停并等待用户改选 `build_mode: executing-plans`。

### 1. 派发 Prompt 与回报契约

每个 implementer 或修复 agent prompt 必须包含：

- 当前单个 task 的完整文本、架构背景和依赖上下文
- `Language: 使用触发本次工作流的用户请求语言输出`
- `Code comment language: 新增/修改的代码注释、JavaDoc、TODO、测试说明跟随用户/Spec 语言；中文需求默认中文注释，除非项目既有规范要求英文`
- 允许修改的文件范围和禁止修改的范围
- 必须执行的测试命令和交付要求
- agent 应返回 patch 供协调者应用，还是仅在明确授权时提交到隔离 branch/worktree
- 修复 agent 还必须收到对应 reviewer 的完整反馈

agent 回报状态必须为 `DONE | DONE_WITH_CONCERNS | BLOCKED | NEEDS_CONTEXT`，并包含实现内容、测试结果、提交哈希或 patch 路径、变更文件和顾虑。进入审查前，主会话必须确认提交/patch 和文件已集成且在当前工作树可见；若平台使用隔离副本，先拉取、合并或应用变更。

每个 reviewer prompt 必须包含完整 task、实现提交或差异以及 RED/GREEN 证据（`tdd_mode: tdd` 时）。reviewer 不得只依据 implementer 的总结进行审查。

### 2. Implementer 范围与提交限制

implementer 只负责实现、测试，并返回可审计的变更产物。**implementer 不得勾选 plan 或 OpenSpec task**，也不得只更新内置 Todo 或对话 checklist。

Superpowers `subagent-driven-development/implementer-prompt.md` 模板里的 "Commit your work" 在 Impetus 中被覆盖：

- 默认：不要提交。返回 patch/diff、变更文件列表和测试证据。
- 若协调者明确选择隔离 branch/worktree 模式，且提交授权已经确认，implementer 只能在该隔离 branch/worktree 内提交，并回报 commit hash。
- 禁止直接提交到协调者当前实时分支/工作树。
- 禁止后台 implementer 自己向用户请求 per-task commit；所有提交决策都由协调者负责。
- 若确需在隔离 branch/worktree 内提交，提交前仍必须检查 Git author；不得使用 `Claude <claude@anthropic.local>`、`claude@anthropic.local` 或其他平台默认作者。

### 3. TDD 硬约束

若 `tdd_mode: tdd`，每个 implementer 和修复 agent 必须先使用 Skill 工具加载 Superpowers `test-driven-development` 技能，并在 prompt 中同时注入：

```text
You MUST follow TDD: write a failing test first, watch it fail, then write minimal code to pass. No production code without a failing test first.
```

implementer 或修复 agent 回报必须提供 **RED 失败命令与失败摘要**、**GREEN 通过命令与通过摘要**；缺少任一证据不得进入审查。spec compliance reviewer 和 code quality reviewer 都必须核验 RED/GREEN 证据与测试覆盖。

### 4. 持久进度检查点

主会话必须维护 `openspec/changes/<name>/.impetus/subagent-progress.md`，并在构建依赖图、选择 wave、每次派发、agent 回报、审查结果、修复轮次变化和 task 勾选后立即更新。检查点至少记录：

- 依赖图摘要和当前 wave 编号
- 当前 wave 的 task id、前置依赖和写入范围
- 当前 plan task 唯一文本及映射的 OpenSpec task 文本
- 每个 task 的当前阶段：`queued | implementing | merge-check | spec-review | quality-review | checkoff | done | blocked | final-review | final-fix`
- 实现提交哈希或 patch 路径、变更文件和 RED/GREEN 证据
- 已通过的审查阶段及尚未解决的 reviewer 反馈
- 当前 task 或 final review 的审查-修复轮次（最多 3 轮）

该文件只保存恢复所需的协调状态，不替代 plan 或 OpenSpec checkbox。已完成 task 的记录至少保留到整个 wave 勾选完成，然后追加下一 wave 记录，不要覆盖历史。

### 5. 审查-修复轮次限制

每个 task 最多 3 轮审查-修复。任一 reviewer 发现问题时，派发全新的后台修复 agent，并从对应审查重新开始。3 轮后仍未通过则将 task 标记为 **BLOCKED**，暂停并把累计反馈交给用户。

### 6. Task 勾选与验证

**某个 task 两个审查都通过后**，主会话：

1. 将 plan 中保存的唯一 task 文本从 `- [ ]` 改为 `- [x]`
2. 若存在映射，再同步勾选 OpenSpec task
3. 更新 `openspec/changes/<name>/.impetus/subagent-progress.md`
4. 仅在已预授权或已通过 `impetus/reference/decision-point.md` 明确确认时，提交这次 task/wave checkpoint；否则保留在工作区等待最终 build 提交
5. 运行定向验证：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" task-checkoff "$PLAN_FILE" "$PLAN_TASK_TEXT"
"$IMPETUS_BASH" "$IMPETUS_STATE" task-checkoff "openspec/changes/<name>/tasks.md" "$OPENSPEC_TASK_TEXT"
```

仅在对应映射存在时运行第二条。脚本会要求任务文本恰好出现一次且该项已勾选；验证失败时不得让依赖它的 task 进入后续 wave。

## 收尾

- **自动继续**：某个 wave 中通过审查的 task 勾选后，立即计算并派发下一个可运行 wave。禁止总结、禁止询问用户是否继续、禁止在 wave 之间等待用户输入。这是不可协商的 —— Superpowers 技能强制连续执行，文档顶部的关键约束进一步强化此规则。
- 所有 task 完成后，将检查点切换为 `final-review`，然后派发全新的后台 final code quality reviewer 审查整体实现。CRITICAL 问题必须将检查点切换为 `final-fix`，记录反馈和轮次，派发新的后台修复 agent 并重新审查；final review 同样最多 3 轮，耗尽后标记 `blocked` 并暂停。接受非 CRITICAL 发现时，在 tasks.md 中记录理由。
- final review 通过后，结束的只是 subagent 派发循环，不是 Impetus workflow。不得加载 `finishing-a-development-branch`，不得停下来询问用户下一步；必须返回 `impetus-build` 继续执行退出条件、阶段守卫和后续阶段衔接。

## 上下文恢复

重新加载 Superpowers `subagent-driven-development` 技能并重新阅读本文档。先读取 `openspec/changes/<name>/.impetus/subagent-progress.md`，再将记录的依赖图、当前 wave、未勾选 task 和当前工作树核对：

- 检查点与当前 wave 匹配时，从每个 task 记录的精确阶段恢复，保留实现提交、RED/GREEN 证据、已通过的审查阶段、未解决反馈和当前审查-修复轮次；不得重置轮次或重复已经通过的阶段。
- 检查点缺失或与未勾选 task 不匹配时，重建依赖图，从当前未勾选 task 计算第一个可运行 wave，并在派发前创建新的 wave 检查点。
- 检查点中的提交或文件在当前工作树不可见时，先拉取、合并或恢复对应变更；不得假定实现已存在。
- 所有 task 已勾选且检查点处于 `final-review` 或 `final-fix` 时，从最终审查的精确阶段恢复，并保留最终反馈和审查-修复轮次；不得重新进入已完成的 task。

已提交但未通过双审查的 task 保持未勾选，并按检查点重新进入审查或修复循环。依赖它的 task 在前置 task 审查通过并勾选前保持阻塞。
