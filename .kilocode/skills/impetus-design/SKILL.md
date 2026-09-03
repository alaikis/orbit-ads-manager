---
name: impetus-design
description: "Impetus 阶段 2：深度设计。用 /impetus-design 调用。通过 brainstorming 产出 Design Doc 和 delta spec。"
---

# Impetus 阶段 2：深度设计（Design）

## 前置条件

- 活跃 change 已存在（proposal.md、design.md、tasks.md）
- 无 Design Doc（`docs/superpowers/specs/` 下无对应文件）

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
"$IMPETUS_BASH" "$IMPETUS_STATE" check <name> design
```

验证通过后继续 Step 1。验证失败时脚本会输出具体失败原因。

**幂等性**：所有 design 阶段操作可以安全重试。如果 `handoff_context` 和 `handoff_hash` 已存在，先确认它们与当前产物一致再决定是否重新生成。

### 1a. 生成 OpenSpec → Superpowers 交接包

**必须由脚本生成，不允许 agent 临场手写 summary 代替。**

```bash
"$IMPETUS_BASH" "$IMPETUS_HANDOFF" <change-name> design --write
```

脚本会读取 change 级 `.impetus.yaml` 的 `context_compression` 快照，并回退到项目级 `.impetus/config.yaml` 和 `IMPETUS_CONTEXT_COMPRESSION`。

默认 `context_compression: beta` 生成并记录压缩后的 spec 投影。显式设置 `context_compression: off` 时生成并记录：

```
openspec/changes/<name>/.impetus/handoff/design-context.json
openspec/changes/<name>/.impetus/handoff/design-context.md
```

并在 `.impetus.yaml` 写入：

```yaml
handoff_context: openspec/changes/<name>/.impetus/handoff/design-context.json
handoff_hash: <sha256>
```

默认交接包是 **compact 可追溯摘录**，不是 agent summary：
- `design-context.json`：机器索引，包含 change、phase、canonical spec、source paths、hash
- `design-context.md`：供 Superpowers 阅读的上下文，包含脚本标记、source path、line range、sha256、确定性摘录
- 超出摘录预算时标记 `[TRUNCATED]`，并保留 Full source 路径

Beta 模式（`context_compression: beta`）生成 **全文 spec 投影**：

```
openspec/changes/<name>/.impetus/handoff/spec-context.json
openspec/changes/<name>/.impetus/handoff/spec-context.md
```

并写入：

```yaml
context_compression: beta
handoff_context: openspec/changes/<name>/.impetus/handoff/spec-context.json
handoff_hash: <sha256>
```

`spec-context.md` 是确定性、可追溯的投影文件。OpenSpec delta spec 仍是权威来源；该文件只把 proposal、design notes、tasks 和 delta specs 投影到一个更适合 Build 阶段读取的低 token 交接文件中。

如确实需要全文上下文，可显式运行：

```bash
"$IMPETUS_BASH" "$IMPETUS_HANDOFF" <change-name> design --write --full
```

交接包来源来自 OpenSpec open 阶段产物：
- `proposal.md`：目标、动机、范围、非目标
- `design.md`：高层架构决策、方案约束
- `tasks.md`：初始任务边界
- `specs/*/spec.md`：delta 能力规格

### 1b. 执行 Brainstorming（带上下文）

**立即执行：** 使用 Skill 工具加载 Superpowers `brainstorming` 技能，ARGUMENTS 包含：

```
Change: <change-name>
OpenSpec Context Pack: openspec/changes/<name>/.impetus/handoff/design-context.md
Machine handoff: openspec/changes/<name>/.impetus/handoff/design-context.json

如果 context_compression 为 beta，使用：
OpenSpec Context Pack: openspec/changes/<name>/.impetus/handoff/spec-context.md
Machine handoff: openspec/changes/<name>/.impetus/handoff/spec-context.json

OpenSpec 产物是上游事实源，不要重新定义需求，不要重写 proposal/spec。
你的任务是基于交接包做深度技术设计：实现方案、技术风险、测试策略、边界条件。
如发现 OpenSpec delta spec 缺少验收场景，只能提出 Spec Patch，并回写 OpenSpec delta spec；不要在 Design Doc 中创建第二份需求 spec。

Design Doc frontmatter 必须最小化，只包含：
---
impetus_change: <change-name>
role: technical-design
canonical_spec: openspec
---

跳过重复上下文探索，直接进入设计提问。
```

禁止跳过此步骤，禁止在未加载该技能的情况下继续。

如 Superpowers `brainstorming` 技能不可用，停止流程并提示安装或启用 Superpowers 技能，不要用普通对话替代该步骤。

技能加载后，按其指引产出设计方案（以对话形式呈现）：
- 技术方案：架构、数据流、关键技术选型与风险
- 测试策略
- 如需补充验收场景，标明将回写的 delta spec 变更

brainstorming 阶段不写入 Design Doc 文件，仅产出设计方案供 Step 1c 用户确认。确认后才创建 `docs/superpowers/specs/YYYY-MM-DD-<topic>-design.md` 并回写 delta spec。

### 1c. 用户确认设计方案（阻塞点）

brainstorming 产出设计方案后，**必须按 `impetus/reference/decision-point.md` 暂停并等待用户明确确认设计方案**。不得在用户确认前创建最终 Design Doc、写入 `design_doc`、运行 design guard，或进入 `/impetus-build`。也不得仅输出文字提示后继续执行。

暂停时只展示必要摘要：
- 采用的技术方案
- 关键取舍与风险
- 测试策略
- 如有 Spec Patch，列出将回写的 delta spec 变更

用户明确确认后，才继续 Step 2。若用户要求调整，继续 brainstorming 迭代，直到用户确认。

### 1d. Brainstorming 恢复检查点

用户确认设计方案后，在创建 Design Doc 前，创建或更新：

```
openspec/changes/<name>/.impetus/handoff/brainstorm-summary.md
```

写入已确认方案、关键取舍、测试策略和 Spec Patches。它是恢复检查点，不是最终 Design Doc。

### 1e. 主动上下文压缩门

写入 `brainstorm-summary.md` 后、创建 Design Doc 前，主动释放早期 OpenSpec 和 brainstorming 上下文：

- 如果当前平台提供原生上下文压缩/清理能力，触发一次。
- 恢复提示必须包含 change name、当前步骤、`brainstorm-summary.md` 和交接文件（`design-context.*` 或 `spec-context.*`）。
- 如果平台无法由 agent 程序化触发压缩，暂停并提示用户手动执行平台压缩操作，或确认无需压缩继续。

### 2. 更新 Impetus 状态

记录 `design_doc` 并运行 design guard 前，先做最终 OpenSpec 一致性自检。此时 change 仍处于 design 阶段，因此允许编辑 `proposal.md`、`design.md`、`tasks.md` 和 delta specs，用于消除与已确认 Design Doc 不一致的二义性、占位或矛盾。不要在 build plan 创建前声称这些文件已只读。实现基线只在 `/impetus-build` 记录 plan 路径后才开始。

先记录 design_doc 路径。如果 Step 1c 回写了 delta spec（新增或修改了 `specs/*/spec.md`），必须重新生成 handoff 以更新 hash：

```bash
# 记录 design_doc 路径
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> design_doc docs/superpowers/specs/YYYY-MM-DD-topic-design.md

# 如有 delta spec 变更，重新生成 handoff（更新 hash）
"$IMPETUS_BASH" "$IMPETUS_HANDOFF" <change-name> design --write

# 自动流转到下一阶段
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> design --apply
```

如果没有 delta spec 变更，跳过 handoff 重新生成步骤。状态文件自动更新，无需手动编辑其他字段。

## 退出条件

- Design Doc 已创建并保存
- Design Doc frontmatter 包含 `impetus_change`、`role: technical-design`、`canonical_spec: openspec`
- `handoff_context` 和 `handoff_hash` 已写入 `.impetus.yaml`（由 guard 强制校验）
- `handoff_hash` 与当前 OpenSpec open 阶段产物一致（由 guard 强制校验）
- `design-context.md` 或 beta `spec-context.md` 必须是脚本生成，且包含 source path、mode、sha256 等可追溯标记（由 guard 强制校验）
- 如有新能力或补充验收场景，OpenSpec delta spec 已创建/更新
- `design_doc` 已写入 `.impetus.yaml`
- **阶段守卫**：运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> design --apply`，全部 PASS 后自动流转到 `phase: build`

退出前必须使用 `--apply`：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> design --apply
```

## 上下文压缩恢复

design 阶段在 brainstorming 过程中可能触发上下文压缩。恢复时先运行：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" check <change-name> design --recover
```

脚本输出结构化恢复上下文（阶段、已完成字段、待完成字段、恢复动作）。按 Recovery action 判断下一步。

## 自动流转

退出条件满足后，遵循共享自动流转协议：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" next <change-name>
```

- `NEXT: auto` → 调用 `SKILL` 指向的 skill
- `NEXT: manual` → 暂停并提示用户手动运行显示的 skill
