---
name: impetus-hotfix
description: "Impetus 预设路径：Bug fix / 热修复。跳过 brainstorming，直接 open → build → verify → archive。适用于行为修复、不涉及新 capability 设计的场景。"
---

# Impetus 预设路径：Hotfix

快速 bug fix 工作流：open → build → verify → archive。跳过 brainstorming 和完整 plan，适用于行为修复、不涉及新 capability 设计的场景。

**适用条件**（必须全部满足）：
1. 修复已有功能的 bug，不新增 capability
2. 不涉及接口变更或架构调整
3. 改动范围可预估（通常 ≤ 2 个文件）

**不适用**：如修复过程发现需要架构调整，应升级为完整 `/impetus` 流程。

---

## 流程（preset workflow，5 阶段）

执行链路：open → build → verify → archive。Hotfix 为每个阶段提供默认决策：精简开启、直接构建、按规模验证、验证通过后归档。

开始前先定位 Impetus 脚本：

```bash
IMPETUS_ENV="${IMPETUS_ENV:-$(find . "$HOME"/.*/skills "$HOME/.config" "$HOME/.gemini" -path '*/impetus/scripts/impetus-env.sh' -type f -print -quit 2>/dev/null)}"
if [ -z "$IMPETUS_ENV" ]; then
  echo "ERROR: impetus-env.sh not found. Ensure the impetus skill is installed." >&2
  return 1
fi
. "$IMPETUS_ENV"
```

### 1. 快速开启（preset open）

复用 Impetus open 能力创建 change，但使用 hotfix 默认值：不执行 `openspec-explore` 长探索，直接进入精简 change 创建。

**立即执行：** 使用 Skill 工具加载 `openspec-new-change` 技能。禁止跳过此步骤。

技能加载后，按其指引创建精简版产物：
  - `proposal.md` — 问题描述 + 根因分析 + 修复目标（无需方案对比）
  - `design.md` — 修复方案（1 个即可，无需多方案对比）
  - `tasks.md` — 修复任务清单
- **无需 delta spec**（除非修复改变了已有 spec 的验收场景）

初始化 Impetus 状态文件：

```bash
# 有 issue_id 时：
"$IMPETUS_BASH" "$IMPETUS_STATE" init <name> hotfix <issue_id> bug
# 无 issue_id 时（用空串占位，确保 bug 落在第四个参数 issue_type 上）：
"$IMPETUS_BASH" "$IMPETUS_STATE" init <name> hotfix "" bug
```

如果用户提供了 Bug 跟踪 ID（例如 `BUG-456`），作为第三个参数传入。hotfix 的 issue_type 固定为 `bug`（第四个参数），用于提交信息的类型前缀。
```

初始化后验证状态：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" check <name> open
```

阶段守卫完成 open → build 过渡：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> open --apply
```

**创建分支隔离**：

hotfix 初始化时 `isolation` 已设为 `branch`，guard 通过后立即创建工作分支。先解析目标仓库列表，再创建或切入分支：

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

### 2. 直接构建（preset build）

使用 hotfix 默认值：`build_mode: direct`。跳过 Superpowers `brainstorming` 和 `writing-plans`（除非任务 > 3 个；若超过 3 个任务，转入 `/impetus-build` 的计划与执行方式选择）。

分支已在 open 阶段创建。继续或开始修改前，按 `impetus/reference/dirty-worktree.md` 协议处理未提交改动。若归因后发现修复范围超出 hotfix，按本文件"升级条件"处理。

**立即执行：** 按 tasks.md 逐个执行任务：

1. 读取 `openspec/changes/<name>/tasks.md`，获取未完成任务列表
2. 对每个未完成任务：
   - 根据任务描述修改代码
   - 运行项目格式化命令（如 `mvn spotless:apply`、`npm run format` 等）
   - 运行相关测试确认通过
   - 将 tasks.md 中对应 `- [ ]` 勾选为 `- [x]`
   - 按 `impetus/reference/commit-convention.md` 协议提交代码（提交前必须 `impetus/reference/decision-point.md` 确认；commit message 格式 `<issue_type>#<issue_id> 简述修复`，hotfix 默认 issue_type 为 `bug`；`issue_id` 为 null 时不自动提交，改动保留在工作区由用户处理）
3. 全部任务完成后，显式运行项目相关测试和构建命令

**如修复影响已有 spec 验收场景**：
- 在 `openspec/changes/<name>/specs/<capability>/spec.md` 创建 delta spec
- 仅包含 `## MODIFIED Requirements` 部分

### 3. 根因消除检查

**在运行 build guard 之前执行**，确保修复确实消除了问题根因：

1. 读取 proposal.md 中的 bug 描述和根因
2. 搜索验证问题代码不再存在
3. 如根因未消除，回到 Step 2 继续修复（此时仍在 build 阶段，无需状态回退）

**升级条件**：
- 根因消除检查发现深层架构问题 → 停止 hotfix，按升级条件阻塞确认处理
- 修复需要额外接口变更 → 停止 hotfix，按升级条件阻塞确认处理

根因确认消除后，运行阶段守卫完成 build → verify 过渡：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> build --apply
```

状态文件自动更新为 `phase: verify`、`verify_result: pending`，然后进入验证。

### 4. 验证（preset verify）

复用 `/impetus-verify`，由 impetus-verify 的规模评估决定轻量或完整验证。

**立即执行：** 使用 Skill 工具加载 `impetus-verify` 技能。禁止跳过此步骤。

无 delta spec 的小范围 hotfix 通常满足轻量验证条件（≤ 3 tasks、≤ 2 files），impetus-verify 的规模评估会选择轻量验证路径（5 项快速检查）。若 hotfix 创建了 delta spec，则根据 impetus-verify 的规模评估规则进入完整验证路径。

验证通过后，按 `/impetus-verify` 的规则将 `.impetus.yaml` 的 `verify_result` 记录为 `pass`，归档前不得跳过该状态。

### 5. 归档（preset archive）

复用 `/impetus-archive`。归档前必须满足 `.impetus.yaml` 中 `verify_result: pass`。

**立即执行：** 使用 Skill 工具加载 `impetus-archive` 技能进行归档。禁止跳过此步骤。
如有 delta spec，按 impetus-archive 规则同步到 main spec，并处理关联 Design Doc 与 Plan 的归档标注。

---

## 连续执行模式

<IMPORTANT>
Hotfix 流程为 **一次性连续执行**。调用 `/impetus-hotfix` 后，agent 在 hotfix 自有步骤间自动推进，不主动停顿。但以下情况必须暂停等待用户确认：

1. 遇到升级条件（见"升级条件"章节），**必须按 `impetus/reference/decision-point.md` 暂停并等待用户明确确认**升级为完整流程
2. 任务超过 3 个转入 `/impetus-build` 时的工作区隔离和执行方式选择
3. 验证阶段（impetus-verify）的验证失败决策和分支处理决策

执行顺序：快速开启 → 直接构建 → 根因消除检查 → 验证 → 归档 → 完成

每个阶段完成后立即进入下一阶段。阶段内部仍必须按上文要求调用对应 Impetus/OpenSpec/Superpowers skill，被调用的 skill 如有自己的用户决策点，按该 skill 规则执行。
</IMPORTANT>

---

## 升级条件

满足以下**任一**条件时，停止 hotfix 流程，升级为完整 `/impetus`：

| 条件 | 说明 |
|------|------|
| 改动涉及 **3+ 文件** | 超出单点修复范围 |
| 架构变更 | 新模块、新接口、新依赖 |
| 数据库 schema 变更 | 结构性调整 |
| 引入新的 public API | 修复产生了新的对外接口 |
| 修复范围超出单一函数/模块 | 需要多处协调修改 |

满足升级条件时**必须按 `impetus/reference/decision-point.md` 暂停并等待用户明确确认**升级为完整 `/impetus` 流程。不得直接进入 `/impetus-design`，不得自动补充 Design Doc。不得仅输出文字提示后继续执行。

用户确认升级后，**必须先更新 workflow 字段**再进入完整流程：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <name> workflow full
```

然后在当前 change 基础上补充 Design Doc：**立即使用 Skill 工具加载 `impetus-design` skill**，后续正常走完整流程。若用户不确认升级，停止 hotfix 并报告当前变更已超出 hotfix 适用范围。

---

## 退出条件

- Bug 已修复，测试通过
- change 已归档
- 如有 spec 变更，已同步到 main spec
- **阶段守卫**：build → verify 前运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> build --apply`，verify → archive 前按 `/impetus-verify` 规则运行 `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> verify --apply`
