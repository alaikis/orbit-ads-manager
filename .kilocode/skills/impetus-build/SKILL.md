---
name: impetus-build
description: "Impetus 阶段 3：计划与构建。用 /impetus-build 调用。制定计划并选择执行方式（subagent 或直接执行）实施。"
---

# Impetus 阶段 3：计划与构建（Build）

## 前置条件

- Design Doc 已创建（阶段 2 完成）
- 活跃 change 存在

## 步骤

### 0. 入口状态验证（Entry Check）

执行入口验证：

```bash
IMPETUS_ENV="${IMPETUS_ENV:-$(find . "$HOME"/.*/skills "$HOME/.config" "$HOME/.gemini" -path '*/impetus/scripts/impetus-env.sh' -type f -print -quit 2>/dev/null)}"
if [ -z "$IMPETUS_ENV" ]; then
  echo "ERROR: impetus-env.sh not found. Ensure the impetus skill is installed." >&2
  return 1
fi
. "$IMPETUS_ENV"
"$IMPETUS_BASH" "$IMPETUS_STATE" check <name> build
```

验证通过后继续 Step 1。验证失败时脚本会输出具体失败原因。

**幂等性**：build 阶段所有操作可安全重复执行。读取 `.impetus.yaml` 的 `phase` 字段确认仍在 build 阶段，读取 plan 文件头的 `base-ref`，再读取 tasks.md 找到第一个未勾选任务继续执行。已提交的任务不得重复提交。

### 1. 制定计划

**立即执行：** 使用 Skill 工具加载 Superpowers `writing-plans` 技能。禁止跳过此步骤。

技能加载后，按其指引制定计划。计划要求：
- 保存至 `docs/superpowers/plans/YYYY-MM-DD-<feature>.md`
- 引用设计文档，拆分为可执行任务
- **Impetus 提交策略覆盖 Superpowers writing-plans 的 frequent commits 默认偏好**：full workflow 的 `executing-plans` 计划中不得包含每个 task 的 `git commit` 步骤、`git add`/`git commit` 命令块，或英文提交信息示例。计划只能写明“完成所有 task、实现 Review 窗口通过后，按 `impetus/reference/commit-convention.md` 执行一次最终 build 提交”。
- **Plan 文件头必须包含关联元数据**：

```yaml
---
change: <openspec-change-name>
design-doc: docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md
base-ref: <git rev-parse HEAD before implementation>
---
```

`base-ref` 用于验证阶段跨提交统计改动规模。创建计划时先记录当前提交：

```bash
git rev-parse HEAD
```

计划写入后、记录到 `.impetus.yaml` 前，必须执行 plan 自检：

```bash
rg -n "git commit|git add|\\bCommit\\b|提交" docs/superpowers/plans/YYYY-MM-DD-<feature>.md
```

若命中 per-task 提交步骤，必须立即编辑 plan 删除这些步骤和命令，改为一个末尾的“最终 build 提交”说明；不得带着 per-task commit 计划进入 Step 2 或执行阶段。允许保留“最终 build 提交”说明和对 `impetus/reference/commit-convention.md` 的引用。

### 2. 更新计划状态并提供 plan-ready 暂停点

先记录 plan 路径：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> plan docs/superpowers/plans/YYYY-MM-DD-feature.md
```

无需手动更新 phase，guard 会在退出条件满足后自动流转。

计划写入后，立即提供一个新的用户决策点：

| 选项 | 行为 | 说明 |
|------|------|------|
| A | 继续执行 | 保持在当前模型中，进入 Step 3 选择工作区隔离和执行方式 |
| B | 暂停切换模型 | 记录 `build_pause: plan-ready`，本次 `/impetus-build` 停止，用户稍后可从 `/impetus` 或 `/impetus-build` 恢复 |

这是用户决策点。**必须按 `impetus/reference/decision-point.md` 暂停并等待用户明确选择**，不得自动继续，也不得把暂停写入 `build_mode`。

用户选择继续时：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_pause null
```

用户选择暂停时：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_pause plan-ready
```

设置 `build_pause: plan-ready` 后，当前调用停止。不要选择 `isolation` 或 `build_mode`，不要加载执行技能。

### 3. 选择执行方式

如果恢复时检测到 `build_pause: plan-ready` 且 `plan` 文件存在，不要重新运行 `writing-plans`。先告知用户当前停在 plan-ready 暂停点；用户确认继续后，设置：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_pause null
```

然后继续本步骤选择执行方式。

分支已在 open 阶段创建，`isolation` 已设为 `branch`。在开始执行前，**询问用户**选择执行方式：

**执行方式**：

| 选项 | 技能 | 适用场景 |
|------|------|---------|
| A | Superpowers `executing-plans` | 小/中型内聚需求的默认选择，轻量快速 |
| B | Superpowers `subagent-driven-development` | 大型独立任务批次，确实需要后台并行 agent 和双阶段审查 |

**执行方式推荐规则**：
- 默认推荐 → `executing-plans`
- 内聚的单功能、CRUD/API 增量、登录/鉴权切片、类似 bugfix 的 full workflow，或少于 10 个 plan task，推荐 `executing-plans`
- 只有同时满足以下条件才推荐 `subagent-driven-development`：当前平台具备真实后台 subagent/Task 派发能力；任务大多相互独立；计划有 10+ 个实质任务或多个模块可并行推进；用户接受 implementer → spec review → code quality review 的额外成本
- 来自 hotfix 路径 → 推荐 `executing-plans`

这是用户决策点。**必须按 `impetus/reference/decision-point.md` 暂停并等待用户明确选择执行方式**，不得根据推荐规则自行选择。推荐规则只能用于说明建议，不能替代用户确认。禁止仅输出文字提示后继续执行。

用户选择后，更新 `build_mode` 字段：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_mode <subagent-driven-development|executing-plans|direct>
```

选择 `subagent-driven-development` 时，必须先确认当前平台具备真实后台 subagent/Task 派发能力，且用户接受更重的 implementer → spec review → code quality review 流程成本。然后记录：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> subagent_dispatch confirmed
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_mode subagent-driven-development
```

如果能力不存在或不确定，使用 `executing-plans`。未记录 `subagent_dispatch: confirmed` 时，guard 和状态转换都会阻止 `subagent-driven-development`。

`build_mode` 默认仅 hotfix/tweak preset 使用 `direct`。full workflow 不得默认使用 `direct`。只有用户明确要求跳过计划执行技能，且你已记录显式 override 时，才允许：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> direct_override true
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_mode direct
```

没有 `direct_override: true` 时，full workflow 的 `build_mode=direct` 会被 guard 和状态转换同时拦截。

**加载执行技能**：使用 Skill 工具加载对应技能。禁止跳过此步骤。

如所选 Superpowers 技能不可用，停止流程并提示安装或启用对应技能，不要用普通对话替代该步骤。

技能加载后，按其指引执行：
- 按计划执行任务
- 每个任务通过本地验证后，完成 tasks.md 勾选（`- [ ]` → `- [x]`）
- 代码中新增或修改的注释、JavaDoc、TODO、错误说明和测试说明必须跟随用户/Spec 主要语言；中文需求/中文 workflow 下默认使用中文注释，除非项目既有规范或外部 API 文档要求英文。不要为了“看起来专业”把新增业务注释写成英文。
- 按 build mode 使用提交策略：
  - `executing-plans`：不要每完成一个 task 就请求提交。任务进度保留在工作区，只有全部任务通过、实现 Review 窗口确认符合预期、用户确认最终提交后，才做一次 change 级 build 提交。
  - `subagent-driven-development`：遵循 `impetus/reference/subagent-dispatch.md`；只有在 build-mode 选择时预授权，或协调者提交前再次明确确认时，才允许 task/wave 进度提交。
  - `direct`：默认结束时提交一次，除非用户明确要求拆分提交。

### 4. Spec 增量更新

实施过程中发现初版 spec 不完整时，按变更规模分级处理：

| 规模 | 触发条件 | 做法 |
|------|---------|------|
| 小 | 遗漏验收场景、边界条件 | 直接编辑 delta spec + design.md，追加 tasks.md 任务 |
| 中 | 接口变更、新增组件、数据流变化 | **按 `impetus/reference/decision-point.md` 暂停并等待用户确认后**，必须使用 Skill 工具加载 Superpowers `brainstorming` 更新 Design Doc + delta spec |
| 大 | 全新 capability 需求 | **必须按 `impetus/reference/decision-point.md` 暂停并等待用户确认拆分**；用户确认后，通过 `/impetus-open` 创建独立 change |

**50% 阈值判定**：以 tasks.md 初始任务总数为基准，若新增任务数超过该总数的一半，视为超出原计划范围，**必须按 `impetus/reference/decision-point.md` 暂停并等待用户决定是否拆分为新 change**。

创建独立 change 时必须调用 `/impetus-open`，不得直接调用 `/opsx:new`。`/impetus-open` 会同时创建 OpenSpec 产物和 `.impetus.yaml`，避免新 change 脱离 Impetus 状态机。

**原则**：
- delta spec 是活文档，本阶段期间随时可修改
- 每次更新必须在最终 build 提交前反映到产物中；`executing-plans` 默认并入一次 change 级 build 提交，除非用户明确要求单独提交 spec 更新。`subagent-driven-development` 遵循 subagent 派发提交策略。
- 不提前同步到 main spec，归档时统一同步
- 小规模增量直接改 delta spec 时，应在 commit message 中注明，便于归档时判断 design doc 漂移

### 5. 上下文管理

Build 是最长阶段，可能跨越大量任务。为支持上下文压缩后断点恢复：

- **每完成一个 task**：立即勾选 tasks.md，并保持测试/证据同步。`executing-plans` 不按 task 提交；恢复依赖 tasks.md、`.impetus.yaml` 与工作区。`subagent-driven-development` 按 `impetus/reference/subagent-dispatch.md` 持久化 checkpoint/wave 进度。
- **上下文压缩后恢复**：先运行 `"$IMPETUS_BASH" "$IMPETUS_STATE" check <change-name> build --recover`，脚本输出结构化恢复上下文（isolation/build_mode 状态、plan 路径、任务完成进度、恢复动作）。根据 Recovery action 决定下一步。
- **用户手动修改恢复**：按 `impetus/reference/dirty-worktree.md` 协议处理未提交改动。该协议定义了检查步骤、归因分类和禁令。build 阶段的特殊处理：
  1. 归因后，若 diff 暗示计划或 spec 已变化，按 Step 4「Spec 增量更新」分级处理
- **长任务拆分**：单任务超过 200 行代码变更时，考虑拆分为多个子任务。`executing-plans` 仍默认一个 change 一次提交，除非用户明确授权 checkpoint 提交。

### 6. 实现 Review 窗口（仅 full workflow）

所有 task 完成、构建/测试通过后，**不要立即提交，也不要立即运行 `guard build --apply` 跳转 verify**。full workflow 必须先停在 build 出口，让用户检查实现是否符合意图。复杂需求往往无法在 design 阶段一次想全，意图常常要在看到实现后才浮现——这个窗口就是捕获「需求不是这样的」的地方，让它在仍处于 build 阶段时被发现并走 Step 4，而不是拖到 verify 被当成 bug。

**仅 full workflow 执行此步骤。** hotfix/tweak 预设路径跳过本步骤，直接进入退出条件。

先标记暂停点：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_pause review-pending
```

然后**必须按 `impetus/reference/decision-point.md` 发起用户确认（阻塞点）**。先向用户展示本次实现的关键变更摘要（如 `git diff --stat <base-ref>...HEAD`），再调用 decision-point protocol：

> 问题：「以上实现是否符合你的预期？」
> 选项：
> - 「符合预期，继续」— 清空暂停点后进入退出条件
> - 「不符合，需要调整」— 路由到 Step 4「Spec 增量更新」

**禁止用会自动继续的文字提示替代 decision-point 协议。** 若平台不支持 decision-point protocol 工具，则以文本形式输出上述选项并要求用户回复对应文字。**必须收到用户明确包含「符合」或「不符合」的回复后才能继续，禁止解读、推断或假设用户意图。** 未收到明确回复前，停止一切操作。不得自行判定「应该符合」后继续。

**用户选择 A（符合）**：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> build_pause null
```

然后按 `impetus/reference/commit-convention.md` 执行最终 build 提交（仅当 `issue_id` 非 null）。这个提交是当前 build pass 的候选实现提交，不是最终归档信号。若 verify 后续失败，运行 `verify-fail` 回到 build，修复问题，并在安全时 amend 这个候选提交；否则经用户确认后追加一个小的修复提交。若 `issue_id` 为 null，改动留在工作区由用户手动处理。

提交策略满足后，再进入退出条件运行 guard。

**用户选择 B（不符合）**：这是「在某个阶段发现需求不是这样的」的发现，必须按 Step 4 分级回写产物，而不是进入 verify 当 bug 修：
- 取消受影响 task 的勾选（`- [x]` → `- [ ]`），或按 Step 4 追加新任务
- 按 Step 4 规模分级更新 delta spec 的验收场景（scenario）+ design.md；中/大规模需 decision-point protocol 再确认并加载 `brainstorming` 或拆新 change
- 回到本阶段从第一个未勾选 task 继续实施
- 改完后重新回到本 Review 窗口（循环），直到用户选择 A 才离开 build

整个过程 `phase` 始终停在 build，仅 `build_pause` 与 tasks.md 变化，不进入 verify、不触发 verify-fail。`build_pause: review-pending` 未清空为 null 前，`guard build --apply` 与 `transition build-complete` 都会被阻止离开 build。

## 退出条件

- tasks.md 全部勾选
- 本轮 build pass 已经用户确认后按 `impetus/reference/commit-convention.md` 提交一次；若 `issue_id` 为 null，改动可保留在工作区由用户手动处理，不阻塞流转
- 已显式运行项目对应的构建/测试命令并通过（不要只依赖 guard 自动猜测）
- `isolation` 已在 open 阶段设为 `branch`
- `build_mode` 已写为 `subagent-driven-development`、`executing-plans` 或带显式 override 的 `direct`
- `build_pause` 已清空为 `null`（full workflow 须先通过 Step 6 实现 Review 窗口）
- **阶段守卫**：运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> build --apply`，全部 PASS 后自动流转到 `phase: verify`

Guard 会优先读取项目配置中的命令：

```yaml
build_command: <build command>
verify_command: <verify command>
```

配置位置可为 change 的 `.impetus.yaml`，也可为仓库根目录的 `.impetus.yaml` / `impetus.yaml` / `.impetus.yml` / `impetus.yml`。
未配置时才回退到 `npm run build`、Maven 或 Cargo 的默认探测。构建失败时 guard 会打印失败命令输出，作为排查证据。

退出前运行 guard 自动流转：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> build --apply
```

状态文件自动更新为 `phase: verify`、`verify_result: pending`。

## 自动流转

退出条件满足后（包括用户选择工作方式），自动流转到下一阶段：

> **REQUIRED NEXT SKILL:** 调用 `impetus-verify` skill 进入验证与收尾阶段。
