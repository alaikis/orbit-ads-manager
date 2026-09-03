---
name: impetus-solution
description: "Impetus 技术方案模式：通过 open → design → review → published 创建和持续打磨技术方案。包含 brainstorming，不写代码。"
---

# Impetus 技术方案模式

当用户想先做技术方案/架构方案，但暂时不想进入编码实现时，使用 `/impetus-solution`。

solution 工作流：

```
open → design → review → published
```

这个工作流没有 build 阶段、没有 verify 阶段、不修改源码、不提交代码。

## 入口检查

按当前阶段运行状态校验：

```bash
IMPETUS_ENV="${IMPETUS_ENV:-$(find . "$HOME"/.*/skills "$HOME/.config" "$HOME/.gemini" -path '*/impetus/scripts/impetus-env.sh' -type f -print -quit 2>/dev/null)}"
if [ -z "$IMPETUS_ENV" ]; then
  echo "ERROR: impetus-env.sh not found. Ensure the impetus skill is installed." >&2
  return 1
fi
. "$IMPETUS_ENV"
"$IMPETUS_BASH" "$IMPETUS_STATE" check <change-name> <phase>
```

如果还没有活跃 solution change，先按 open 方式创建 OpenSpec 产物，再初始化为 `workflow: solution`：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" init <change-name> solution <issue_id-or-null> task
```

方案型工作默认 `issue_type` 使用 `task`，除非用户明确提供其他类型。

## 阶段 1：Open

目标：澄清方案诉求，创建 OpenSpec 基线。

必需产物：

```
openspec/changes/<change-name>/proposal.md
openspec/changes/<change-name>/design.md
openspec/changes/<change-name>/tasks.md
openspec/changes/<change-name>/.impetus.yaml
```

可选 delta spec：

```
openspec/changes/<change-name>/specs/<capability>/spec.md
```

`tasks.md` 是未来开发可参考的实施任务草案，不是本工作流的 build checklist。

需要澄清：
- 问题与目标
- 范围与非范围
- 约束与假设
- 涉及的角色、系统或边界
- 这份技术方案要支撑什么决策

产物齐备且用户确认边界后：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> open --apply
```

然后遵循 `impetus/reference/auto-transition.md`。

## 阶段 2：Design With Brainstorming

目标：发散方案、比较选项、收敛推荐技术方案，并写入方案文档。

**必须能够头脑风暴，而且必须执行 brainstorming。** 立即使用 Skill 工具加载 Superpowers `brainstorming` 技能，并把 OpenSpec 产物作为上下文。不得用普通对话替代此步骤。

brainstorming 至少覆盖：
- 至少两个可行方案；除非该领域明显只有一个安全选项
- 关键取舍与风险
- 相关的数据、API、权限、安全影响
- 必要的迁移、发布和回滚考虑
- 测试与验证策略

brainstorming 后，必须按 `impetus/reference/decision-point.md` 暂停并等待用户确认采用方向。用户确认前不得写最终方案文档。

确认后，创建或更新：

```
docs/superpowers/solutions/YYYY-MM-DD-<change-name>-solution.md
openspec/changes/<change-name>/solution.md
```

`openspec/changes/<change-name>/solution.md` 可以是简短索引/指针；正式方案内容以 `docs/superpowers/solutions/...` 为准。

然后写入状态：

```bash
"$IMPETUS_BASH" "$IMPETUS_STATE" set <change-name> solution_doc docs/superpowers/solutions/YYYY-MM-DD-<change-name>-solution.md
"$IMPETUS_BASH" "$IMPETUS_STATE" set <change-name> solution_version 1
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> design --apply
```

### 技术方案模板

文档语言跟随用户语言。小需求可以简写内容，但保留标题，方便后续持续打磨同一份文档。

```markdown
---
impetus_change: <change-name>
workflow: solution
solution_version: 1
status: draft
---

# <标题> 技术方案

## 1. 背景与目标

## 2. 范围与非范围

## 3. 现状与约束

## 4. 需求澄清结果

## 5. 方案选项对比

## 6. 推荐方案

## 7. 架构与模块设计

## 8. API / 交互 / 数据流

## 9. 数据与存储

## 10. 安全、权限与合规

## 11. 异常、边界与降级

## 12. 测试与验证策略

## 13. 迁移、发布与回滚

## 14. 风险与取舍

## 15. 待确认问题

## 16. 后续实施任务草案

## 17. 修订记录
```

## 阶段 3：Review

目标：把方案打磨到可以发布/评审的状态。

自检清单：
- 除「待确认问题」章节外，无 `TBD`、`TODO`、`待定` 占位
- 范围与非范围明确
- 假设条件已列明
- 方案对比能支撑推荐结论
- 主要风险有缓解措施
- 后续实施任务草案足以支撑未来开发
- 没有修改源码

如需调整，继续更新同一份方案文档；有实质变化时递增 `solution_version`。

review 通过后：

```bash
"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> review --apply
```

然后遵循 `impetus/reference/auto-transition.md`。

## 阶段 4：Published

方案已可用于人工评审、共享或未来实施。

不要自动归档。不要把 delta spec 合并到主 capability spec。这是决策产物，不是已实现变更。

有效后续动作：
- 用户明确要求时，回到 design/review 继续打磨
- 之后创建独立开发 change，把已发布方案作为输入
- 用 `/impetus-knowledge` 导出到知识库

## 退出汇总

结束时汇报：
- change 名称和 `workflow: solution`
- 当前阶段
- 方案文档路径
- 方案版本
- 核心推荐结论
- 未决问题，如有
