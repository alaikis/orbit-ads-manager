---
name: impetus-verify
description: "Impetus 阶段 4：验证与收尾。用 /impetus-verify 调用。验证实现符合设计，处理开发分支。"
---

# Impetus 阶段 4：验证与收尾（Verify）

## 前置条件

- 代码已提交（阶段 3 完成）
- tasks.md 全部任务已完成

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
"$IMPETUS_BASH" "$IMPETUS_STATE" check <change-name> verify
```

验证通过后继续 Step 1。验证失败时脚本会输出具体失败原因。

**幂等性**：verify 阶段所有检查可安全重复执行。如 `verify_result` 已为 `pass` 且 `branch_status` 已为 `handled`，说明验证已完成，直接执行 guard 流转。如 `verify_result` 为 `pending`，从头开始验证。

### 1. 改动规模评估

执行规模评估：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" scale <change-name>
```

脚本自动统计任务数、增量规格数、变更文件数，判断使用 light 或 full 验证模式，并设置 verify_mode 字段。

`verify_mode` 只表示验证强度（`light` / `full`）。验证调度方式由 `verify_strategy` 控制：

- `standard`：串行执行验证检查，默认值
- `subagent-driven-verify`：使用只读并行 verifier；仅在平台具备真实后台 subagent / Task / multi-agent 派发能力，且用户已确认时使用

若 `.impetus.yaml` 中 `verify_strategy: subagent-driven-verify`，**立即使用 Skill 工具加载 `subagent-driven-verify`**，按其只读并行协议执行 Step 2/3 的验证检查并生成 `.impetus/verify/run-*/summary.md`。若平台能力不可用，必须降级为 `standard`，不得在主窗口伪装并行。

验证开始前，按 `impetus/reference/dirty-worktree.md` 协议检查并处理未提交改动。verify 阶段的特殊处理：

1. 若 dirty diff 属于当前 change 且涉及实现、测试、tasks、delta spec 或 design doc 变更，不在 verify 阶段直接修复或提交；报告失败项并进入 Step 1b 的验证失败决策阻塞点
2. 若 dirty diff 只是 verify 本阶段产物（例如验证报告草稿、分支处理记录），可继续在 verify 阶段完成并记录状态
3. 若 dirty diff 已实现但 tasks.md 未勾选，视为 build 状态滞后；报告失败项并进入 Step 1b，由用户决定回退修复或接受偏差

用户选择修复后，才允许回退到 build 阶段：

```bash
# 仅在用户确认修复后执行
"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail
```

注意：如果 build 阶段每个任务都已提交，脚本基于工作区 diff 的文件数可能低估改动规模。此时必须读取 plan 文件头的 `base-ref` 并用提交区间复核：

```bash
PLAN=$("$IMPETUS_BASH" "$IMPETUS_STATE" get <change-name> plan)
BASE_REF=$(grep '^base-ref:' "$PLAN" 2>/dev/null | head -1 | sed 's/^base-ref: *//')
git diff --stat "$BASE_REF"...HEAD
```

若提交区间显示改动超过轻量阈值（> 4 个文件、跨模块协调、或 delta spec 超过 1 个 capability），手动设置为完整验证：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <change-name> verify_mode full
```

### 1b. 验证失败决策（阻塞点）

验证不通过时**必须按 `impetus/reference/decision-point.md` 暂停并等待用户决定修复或接受偏差**。不得自动运行 `"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail`，也不得自动调用 `/impetus-build`。禁止仅输出文字提示后继续执行。

暂停时必须列出：
- 失败项
- 是否属于 CRITICAL（构建失败、测试失败、安全问题、核心验收场景失败）
- 推荐处理方式

**不确定性原则**：无法确定严重程度时，降级处理（SUGGESTION > WARNING > CRITICAL）。仅对构建失败、测试失败、安全问题使用 CRITICAL；模糊或不确定的问题标为 WARNING 或 SUGGESTION。

用户选择后按以下方式继续：
- **全部修复**：运行 `"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail`，然后调用 `/impetus-build` 修复
- **逐项处理**：CRITICAL 失败项必须修复；非 CRITICAL 失败项可选择接受偏差，但必须在验证报告中记录接受原因和影响范围。若存在任何 CRITICAL 失败项，不允许跳过修复直接全部接受

### 2a. 轻量验证（小改动）

当规模评估结果为"小"时，跳过 `openspec-verify-change`，直接执行以下检查：

1. tasks.md 全部任务已完成 `[x]`
2. 改动文件与 tasks.md 描述一致（`git diff --stat` / `git diff --cached --stat` / `git diff --stat <base-ref>...HEAD` 对照 tasks 内容）
3. 编译通过（执行项目对应的构建命令，如 `npm run build`、`mvn compile`、`cargo build` 等）
4. 相关测试通过
5. 无明显安全问题（无硬编码密钥、无新增 unsafe 操作）

**通过标准**：5 项全部 OK，无 CRITICAL 问题。

**不通过时**：报告失败项，进入 Step 1b 的验证失败决策阻塞点。用户选择修复后，才执行以下命令记录失败并回退到 build 阶段，然后调用 `/impetus-build` 修复：

```bash
# 仅在用户确认修复后执行
"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail
```

**报告格式**：简表列出 5 项检查结果 + PASS/FAIL。

**跳过项**（不在轻量验证中检查）：
- spec scenario 覆盖率
- design doc 一致性深度比对
- code pattern consistency 建议
- delta spec 与 design doc 漂移检测

### 2b. 完整验证（大改动）

当规模评估结果为"大"时：

**立即执行：** 使用 Skill 工具加载 `openspec-verify-change` 技能。禁止跳过此步骤。

技能加载后，按其指引验证。检查项：
1. tasks.md 全部任务已完成（`[x]`）
2. 实现符合 `openspec/changes/<name>/design.md` 高层设计决策
3. 实现符合 Design Doc（`docs/superpowers/specs/` 下的技术设计文档）
4. 能力规格场景全部通过
5. proposal.md 目标已满足
6. delta spec 与 design doc 无矛盾（若 Build 阶段有增量修改 spec，检查 design doc 是否有对应记录）
7. `docs/superpowers/specs/` 关联的设计文档可定位（文件存在且与当前 change 相关）

验证不通过时：报告缺失项，进入 Step 1b 的验证失败决策阻塞点。用户选择修复后，才执行以下命令记录失败并回退到 build 阶段，然后调用 `/impetus-build` 补充：

```bash
# 仅在用户确认修复后执行
"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail
```

**Spec 漂移处理**（用户决策点）：
- 若检查项 6 发现矛盾（delta spec 有内容但 design doc 未体现），**必须按 `impetus/reference/decision-point.md` 以单选决策形式暂停并等待用户选择处理方式**，不得自动选择。选项：
  - 选项 A：在 design doc 追加 "Implementation Divergence" 节记录偏差原因。选项 A 属于 verify 阶段允许产物；写入后不得因该 design doc 变更再次触发 Step 1b dirty-worktree 决策
  - 选项 B：用户选择 B 后，运行 `"$IMPETUS_BASH" "$IMPETUS_STATE" transition <change-name> verify-fail`，然后调用 `/impetus-build`；由 `/impetus-build` 的 Spec 增量更新规则加载 Superpowers `brainstorming` 更新 Design Doc + delta spec
  - 选项 C：确认偏差可接受，继续验证（归档时 design doc 将标记为 `superseded-by-main-spec`）

### 3. 实现审查门禁（impetus-review）

`impetus-review` 是 verify 阶段的正式实现审查门禁，也可由用户手动调用。它不绑定 Java；第一版内置 JVM/Java 规则集。

**立即执行：** 使用 Skill 工具加载 `impetus-review` 技能。禁止跳过实现审查门禁。

技能加载后，按其指引执行只读审查：

- 基准 ref 优先使用 plan 文件头 `base-ref`，其次回退到 `test`
- 目标 ref 为当前 `HEAD`
- 如存在 `openspec/changes/<change-name>/specs/`，传入以启用 OpenSpec 工作流合规审计
- 审查产物写入 `openspec/changes/<change-name>/.impetus/review/run-<timestamp>/`
- 产物必须包含 `bundle.json`、`findings.json`、`report.md`、`report.html`

审查完成后：

1. 将 `impetus-review` 产出报告路径记录到验证报告中
2. `findings.json` 中 `gate: FAIL` 视为 CRITICAL / blocking 失败项，进入 Step 1b
3. `gate: WARN` 记录到验证报告，由用户决定接受或回 build 修复
4. `gate: PASS` 继续后续 verify

**非 JVM/Java 项目**：仍调用 `impetus-review`。如果当前规则集没有可审查文件，skill 必须生成 `PASS` 空报告，而不是跳过审查门禁。

**失败处理**：如 `impetus-review` 不可用或产物不完整，记录为 verify 失败项；不得静默跳过实现审查门禁。

### 4. 收尾（Superpowers）

**分支按 issue 聚合，最后一个 change 才收尾**：同一个 issue 下的多个 change 共用一条分支（见 `/impetus-build` 的分支隔离）。在分支收尾前，先检测同 issue 是否还有其他活跃 change：

```bash
SIBLINGS=$("$IMPETUS_BASH" "$IMPETUS_STATE" issue-siblings <change-name>)
```

- `SIBLINGS > 0`：同 issue 还有未归档的 change 在用这条分支。**此时不得合并 / 推 PR / 删除 / 丢弃分支**，只能选择「保持分支（稍后处理）」，把真正的收尾留给该 issue 的最后一个 change。直接写入 `branch_status: handled` 并继续流转（分支保留是本 change 的合法收尾）。
- `SIBLINGS = 0`：本 change 是该 issue 的最后一个（或该 change 无 issue_id、为独立分支）。按下方完整流程让用户决策分支处理方式。

**立即执行：** 使用 Skill 工具加载 Superpowers `finishing-a-development-branch` 技能。禁止跳过此步骤。

如 Superpowers `finishing-a-development-branch` 技能不可用，停止流程并提示安装或启用 Superpowers 技能，不要用普通对话替代该步骤。

技能加载后，按其指引收尾。分支处理选项（仅当 `SIBLINGS = 0` 时开放完整选项）：
1. 本地合并到主分支
2. 推送并创建 PR
3. 保持分支（稍后处理）

不提供「删除分支」或「丢弃工作」选项——这类破坏性操作风险高且不可逆，如确需删除分支，由用户在收尾后自行执行。

这是用户决策点。**必须按 `impetus/reference/decision-point.md` 暂停并等待用户选择分支处理方式**，不得根据推荐、默认值或当前分支状态自行选择。禁止仅输出文字提示后继续执行。只有在用户完成选择且对应操作完成后，才允许写入 `branch_status: handled`。`SIBLINGS > 0` 时跳过本决策，直接按上方规则保留分支并写入 `branch_status: handled`。

**确认项**：
- 全部测试通过
- 无硬编码密钥或安全问题

### 5. 记录验证证据

验证报告必须落盘，并在 `.impetus.yaml` 中记录；分支处理完成后也必须写入状态字段。不要手动设置 `verify_result: pass`，通过 guard 自动流转。

```bash
mkdir -p docs/superpowers/reports
# 将本次验证结论写入报告文件，例如：
# docs/superpowers/reports/YYYY-MM-DD-<change-name>-verify.md

"$IMPETUS_BASH" "$IMPETUS_STATE" set <change-name> verification_report docs/superpowers/reports/YYYY-MM-DD-<change-name>-verify.md
"$IMPETUS_BASH" "$IMPETUS_STATE" set <change-name> branch_status handled
```

## 退出条件

- 验证报告通过
- 分支已处理
- `.impetus.yaml` 中 `verification_report` 指向已存在的验证报告文件
- `.impetus.yaml` 中 `branch_status: handled`
- **阶段守卫**：运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> verify --apply`，全部 PASS 后通过 `impetus-state transition verify-pass` 自动流转到 `phase: archive`

验证和分支处理均完成后，运行 guard 自动流转：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> verify --apply
```

状态文件自动更新为 `phase: archive`、`verify_result: pass`、`verified_at: YYYY-MM-DD`。

## 上下文压缩恢复

Verify 阶段可能触发上下文压缩。恢复时先运行：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" check <change-name> verify --recover
```

脚本输出结构化恢复上下文（phase、验证状态、分支状态、恢复动作），根据输出的 Recovery action 决定下一步。

## 自动流转

退出条件满足后（包括用户选择分支处理方式），自动流转到下一阶段：

> **REQUIRED NEXT SKILL:** 调用 `impetus-archive` skill 进入归档阶段。
