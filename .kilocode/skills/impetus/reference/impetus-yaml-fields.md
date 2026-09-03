# .impetus.yaml 字段说明

规范路径：`impetus/reference/impetus-yaml-fields.md`

本文件是 `.impetus.yaml` 状态文件的字段参考。按需查阅，不随 skill 一次性加载。

## 示例

```yaml
workflow: full
phase: build
design_doc: docs/superpowers/specs/YYYY-MM-DD-topic-design.md
solution_doc: null
solution_version: 1
solution_status: active
plan: docs/superpowers/plans/YYYY-MM-DD-feature.md
base_ref: a1b2c3d4e5f6...
build_mode: subagent-driven-development
build_pause: null
subagent_dispatch: confirmed
tdd_mode: tdd
isolation: branch
verify_mode: light
verify_strategy: standard
review_profile: standard
verify_result: pending
verification_report: null
branch_status: pending
created_at: 2026-05-26
verified_at: null
archived: false
```

## 必需字段

| 字段 | 含义 |
|------|------|
| `workflow` | `full`、`hotfix`、`tweak` 或 `solution` |
| `phase` | 当前阶段：`open`、`design`、`build`、`verify`、`archive`、`review` 或 `published`（init 统一设为 `open`，guard 负责过渡） |
| `design_doc` | 关联的 Superpowers Design Doc 路径，可为空 |
| `solution_doc` | `workflow: solution` 使用的技术方案文档路径，方案设计完成前可为空 |
| `solution_version` | 技术方案文档版本号，从 `1` 开始 |
| `solution_status` | `active` 或 `archived`；第一版中 published solution change 保持 active，不自动归档 |
| `plan` | 关联的 Superpowers Plan 路径，可为空 |
| `base_ref` | init 时记录的 git commit SHA，用于 scale 评估。无 plan 时作为改动文件数统计基准 |
| `build_mode` | 已选择的执行方式，可为空 |
| `build_pause` | build 阶段内部暂停点。`null` 表示无暂停，`plan-ready` 表示 plan 已生成，用户选择切换模型后暂停 |
| `subagent_dispatch` | `null` 或 `confirmed`。仅当已确认当前平台存在真实后台 subagent / Task / multi-agent 调度能力时，`build_mode: subagent-driven-development` 才能写入并用于离开 build 阶段 |
| `tdd_mode` | `tdd` 或 `direct`。full workflow 离开 build 阶段前必须已选择。`tdd` 强制每个任务先写失败测试再实现；`direct` 不强制 TDD。hotfix/tweak 默认 `direct` |
| `isolation` | `branch` 或 `worktree`，工作区隔离方式。full 初始化可为 `null`，但只允许持续到 `/impetus-build` Step 3 前；hotfix/tweak 默认 `branch` |
| `verify_mode` | 验证强度：`light` 或 `full`，scale 评估前可为空 |
| `verify_strategy` | 验证调度策略：`standard` 或 `subagent-driven-verify`。并行策略需要真实后台 subagent / Task / multi-agent 派发能力 |
| `review_profile` | `impetus-review` 使用的实现审查 profile：`standard`、`strict`、`security` 或 `architecture` |
| `auto_transition` | `true` 或 `false`。只控制阶段守卫推进 phase 后是否自动调用下一个 skill；`false` 时由 `impetus-state next` 输出 `manual`，暂停下一 skill 调用，但不阻止 phase 字段更新 |
| `verify_result` | `pending`、`pass` 或 `fail` |
| `verification_report` | 验证报告文件路径，verify 通过前必须指向已存在文件 |
| `branch_status` | `pending` 或 `handled`，分支处理完成后设为 `handled` |
| `created_at` | change 创建日期（init 时自动写入），格式 `YYYY-MM-DD` |
| `verified_at` | 验证通过时间，可为空 |
| `archived` | change 是否已归档 |

## 可选字段

| 字段 | 含义 |
|------|------|
| `direct_override` | `true`/`false`。full workflow 如需使用 `build_mode: direct`，必须显式设为 `true` |
| `build_command` | 项目构建命令。guard 优先运行该命令，失败时打印命令输出 |
| `verify_command` | 项目验证命令。verify guard 优先运行该命令，未配置时回退到构建命令 |

## 状态机硬约束

- `build → verify` 前，`isolation` 必须是 `branch` 或 `worktree`
- `build → verify` 前，`build_mode` 必须已选择
- `build_mode: subagent-driven-development` 必须同时有 `subagent_dispatch: confirmed`
- full workflow 离开 build 阶段前 `tdd_mode` 必须已选择为 `tdd` 或 `direct`
- `build_mode: direct` 默认只允许 `hotfix` / `tweak`；full workflow 需要 `direct_override: true`
- `build_pause` 不是执行方式，不得写入 `build_mode`
- solution 工作流使用 `open → design → review → published`，不得进入 `build`、`verify` 或 `archive`
- solution `design → review` 前，`solution_doc` 必须指向已存在文件
- 这些约束同时存在于 `impetus-guard.sh build --apply` 和 `impetus-state.sh transition <name> build-complete`
