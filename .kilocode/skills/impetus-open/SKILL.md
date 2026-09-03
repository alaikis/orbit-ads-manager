---
name: impetus-open
description: "Impetus 阶段 1：开启。用 /impetus-open 调用。通过 OpenSpec 探索想法、创建 change 结构（proposal + design + tasks）。"
---

# Impetus 阶段 1：开启（Open）

## 前置条件

- 无活跃 change，或用户希望创建新 change

## 步骤

### 1. 探索想法

**条件执行**：用户提供了变更描述时，使用 Skill 工具加载 `openspec-explore` 技能探索问题空间。**用户未提供描述时，跳过本步**——直接询问用户想做什么，然后进入 Step 2。

技能加载后，按其指引自由探索问题空间。

### 2. 创建 Change 结构 + 初始化状态

**立即执行：** 使用 Skill 工具加载 `openspec-new-change` 技能。若用户意图未明确、需要先形成建议，改为加载 `openspec-propose`。禁止跳过此步骤。

**命名与范围守卫**：change name 必须使用用户指定或通过 `impetus/reference/decision-point.md` 确认的名称，不得自动生成或推断。变更范围必须与用户描述一致，不得自行扩大或缩小。

确认以下产物已创建：

```
openspec/changes/<name>/
├── .openspec.yaml
├── .impetus.yaml
├── proposal.md       # Why + What：问题、目标、范围
├── design.md         # How（高层）：架构决策、方案选型
└── tasks.md          # 任务清单（勾选框）
```

**Issue ID 与类型确认（阻塞点，init 前必须完成）：**

issue_id 将同一个父任务下的多个 change 关联起来，支持合并归档。issue_type 决定提交信息的类型前缀（`story`/`task`/`bug`），默认 `story`。在执行下方 `impetus-state init` 之前：

- 如果用户在对话中已提供了需求/任务/Bug 跟踪 ID（例如 `REQ-123`、`BUG-456`），直接使用该值，并据此判定类型：`REQ-*` → `story`、`TASK-*` → `task`、`BUG-*` → `bug`，其他格式默认 `story`。
- 当前分支名（如 `ai/ling/123456`）、历史提交信息（如 `story#123456 ...`）、归档目录名（如 `archive/123456/`）或同 capability 既有 change 的 issue_id **只能作为候选推荐**，不得视为“用户已提供”。如果 issue_id 只来自这些上下文，仍必须通过 decision-point protocol 询问用户确认；推荐选项可以写成「使用 123456」。
- 如果用户没有提供，**仅使用一次 decision-point protocol 询问并等待回复**（禁止拆分多轮）：
  > 问题：「这个变更有对应的 issue/任务跟踪 ID 吗？」
  > 选项：
  > - 「没有，跳过」— issue_id 设为 null，issue_type 默认 `story`
  > - 「使用检测到的候选 ID：<ID>」— 仅当从分支/历史/归档检测到候选 ID 时提供，用户选择后才可使用
  > - 「有（其他）— 直接输入完整 ID，如 REQ-100」→ 用户在 Other 文本框直接输入 ID
  >
  > agent 收到 Other 文本后自动推断类型：`REQ-*` → `story`、`TASK-*` → `task`、`BUG-*` → `bug`，不匹配上述前缀的默认 `story`。

**不得在未确认 issue_id 与类型之前执行 `impetus-state init`。不得写“branch/history clearly belongs to issue X, no need to ask”之类推断。** 确认后继续执行下方初始化命令。

**初始化 `.impetus.yaml` 状态文件：**

```bash
IMPETUS_ENV="${IMPETUS_ENV:-$(find . "$HOME"/.*/skills "$HOME/.config" "$HOME/.gemini" -path '*/impetus/scripts/impetus-env.sh' -type f -print -quit 2>/dev/null)}"
if [ -z "$IMPETUS_ENV" ]; then
  echo "ERROR: impetus-env.sh not found. Ensure the impetus skill is installed." >&2
  return 1
fi
. "$IMPETUS_ENV"

if [ -z "$IMPETUS_STATE" ] || [ -z "$IMPETUS_GUARD" ]; then
  echo "ERROR: Impetus scripts not found. Ensure the impetus skill is installed." >&2
  return 1
fi

# issue_id 来自上方阻塞点确认结果，无则为空；issue_type 为 story|task|bug，默认 story
"$IMPETUS_BASH" "$IMPETUS_STATE" init <name> full [issue_id] [issue_type]
```

**创建分支隔离**：

full workflow 初始化时 `isolation` 已设为 `branch`。init 完成后立即创建工作分支。先解析目标仓库列表，再创建或切入分支：

```bash
REPOS_MODE=$("$IMPETUS_BASH" "$IMPETUS_REPOS" mode)   # single | multi | none
mapfile -t REPOS < <("$IMPETUS_BASH" "$IMPETUS_REPOS" list)

# none 模式时停止并提示
if [ "$REPOS_MODE" = "none" ]; then
  echo "ERROR: Not inside a git repository and no sub-repos found. Run in the correct directory or git init first." >&2
  exit 1
fi

# 从 .impetus.yaml 读取 issue_id（未设置时返回字符串 "null"）
ISSUE_ID=$("$IMPETUS_BASH" "$IMPETUS_STATE" get <change-name> issue_id)
# Windows git bash 下 whoami 可能返回 "DOMAIN\user"，只取最后一段
USER_SEG=$(whoami | sed 's#.*[\\/]##')
# 有 issue_id 时按 issue 聚合；无 issue_id 时回退到 change-name
if [ -n "$ISSUE_ID" ] && [ "$ISSUE_ID" != "null" ]; then
  BRANCH="ai/$USER_SEG/$ISSUE_ID"
else
  BRANCH="ai/$USER_SEG/<change-name>"
fi
for repo in "${REPOS[@]}"; do
  if git -C "$repo" rev-parse --verify "$BRANCH" >/dev/null 2>&1; then
    git -C "$repo" checkout "$BRANCH"
  else
    git -C "$repo" checkout -b "$BRANCH"
  fi
done
```

分支按 issue 聚合（`ai/<用户>/<issue_id>`），无 issue_id 时按 change-name 命名。同名 issue 的后续 change 自动复用已有分支。

### 3. 入口状态验证

验证状态机已正确初始化：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" check <name> open
```

验证通过后继续 Step 4。验证失败时脚本会输出具体失败原因。

**幂等性**：open 阶段所有操作可安全重复执行。如 `.impetus.yaml` 已处于 `phase: open` 且三个产物文件均已存在，跳过已完成步骤，从第一个缺失步骤继续。

### 4. 内容完整性检查

确认三个文档内容完整：
- **proposal.md**：问题背景、目标、范围、非目标
- **design.md**：高层架构决策、方案选型、数据流
- **tasks.md**：任务列表，每个任务有明确描述

**文件存在性验证**：逐个确认三个文件路径存在且非空。任一文件缺失或为空时，不得进入 Step 5 或执行阶段守卫，必须回到创建步骤补充。

### 5. 用户审视确认（阻塞点）

三个文档创建完成且内容完整性检查通过后，**必须按 `impetus/reference/decision-point.md` 暂停并等待用户确认**。不得在用户确认前执行阶段守卫或自动流转。

决策点必须以单选形式呈现，包含以下摘要和选项：

**摘要内容**：
- **proposal.md**：问题背景、目标、范围
- **design.md**：高层架构决策、方案选型
- **tasks.md**：任务数量和关键任务描述

**选项**：
- 「确认，继续下一阶段」— 产物符合预期，执行阶段守卫流转
- 「需要调整」— 附带调整说明，修改后重新请求确认

用户选择「确认」后继续执行退出条件。用户选择「需要调整」时，按其说明修改对应文件，然后重新使用 decision-point protocol 请求确认。

## 退出条件

- proposal.md、design.md、tasks.md 均已创建且内容完整
- **用户已确认** proposal、design、tasks 内容符合预期
- **阶段守卫**：运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> open --apply`，全部 PASS 后自动流转到下一阶段

退出前必须使用 `--apply`，否则 `.impetus.yaml` 仍停留在 `phase: open`，下一阶段入口检查会失败。

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> open --apply
```

完整流程会自动更新为 `phase: design`；hotfix/tweak preset 会自动更新为 `phase: build`。

## 自动流转

用户确认后，退出条件满足，自动流转到下一阶段：

> **REQUIRED NEXT SKILL（完整流程）:** 调用 `impetus-design` skill 进入深度设计阶段。
>
> hotfix/tweak preset 由对应 preset skill 控制后续流转（phase 直接进入 build），不经过本节。
