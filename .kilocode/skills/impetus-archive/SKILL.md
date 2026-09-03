---
name: impetus-archive
description: "Impetus 阶段 5：归档。用 /impetus-archive 调用。同步 delta spec 到主 spec，归档 change。"
---

# Impetus 阶段 5：归档（Archive）

## 前置条件

- 验证已通过（阶段 4 完成）
- 分支已处理
- `openspec/changes/<name>/.impetus.yaml` 中 `verify_result: pass`

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
"$IMPETUS_BASH" "$IMPETUS_STATE" check <name> archive
```

验证通过后继续 Step 1。验证失败时脚本会输出具体失败原因。

### 0b. Verify 边界

Archive 不执行验证或代码审查补跑。若 Step 0 入口验证失败，原因是 change 尚未处于 `phase: archive`、`verify_result` 不是 `pass`、`verification_report` 缺失，或分支尚未处理，则停止并回到 `/impetus-verify`。

实现审查门禁只属于 verify 阶段。不得在 archive 中加载 `impetus-review`，不得创建 catch-up 验证报告，也不得在 archive 中手动设置 `verify_result: pass`。

### 1. 执行归档

运行归档脚本，自动完成以下全部步骤：

```bash
"$IMPETUS_BASH" "$IMPETUS_ARCHIVE" "<change-name>"
```

脚本自动执行：
1. 入口状态验证（phase=archive, verify_result=pass, archived=false）
2. 同 issue 关联检查（如果 issue_id 已设置，扫描共享同一 issue_id 的其他活跃 change）
3. Delta spec 合并到主 spec（同名 `### Requirement:` 替换，新 requirement 追加，未涉及的既有 requirement 保留；如有同 issue 活跃 change，也合并其 delta spec）
4. Design doc 前置元数据标注（archived-with, status）
5. Plan 前置元数据标注（archived-with）
6. 移动 change 到归档目录（当 issue_id 存在时，归档到 `archive/<issue_id>/` 嵌套目录下）
7. 通过 `impetus-state transition <archive-name> archived` 更新 `archived: true`
8. Issue 分组 `_summary.md` 生成（issue_id 存在时，追踪该 issue 下所有已归档 change）

**合并归档**：当 `.impetus.yaml` 中有 `issue_id`（如 `REQ-123`）时，change 被归档到 `openspec/changes/archive/REQ-123/YYYY-MM-DD-<name>/` 而非扁平的 `archive/YYYY-MM-DD-<name>/`。如果多个 change 共享同一 `issue_id`，所有同 issue 活跃 change 的 delta spec 会被合并后再同步到主 spec。Issue 分组目录会生成 `_summary.md`，列出所有已归档 change。当最后一个同 issue 的 change 归档时，summary 标记为 **已完成**。

如脚本返回非零退出码，报告错误并停止。
如脚本返回零退出码，归档完成。
脚本摘要中的 `X/Y steps succeeded` 以真实执行步骤计数，不会因 delta spec 同步或文档标注重复累计。

当 change 有 `issue_id` 且同 issue 下还有其他待归档 change 时，需要对每个 change 分别执行归档。所有 change 归档完毕后，issue 分组即告完成。

当待同步的 delta spec 与已有主 spec 不一致时，脚本会在合并前打印 unified diff 预览，帮助确认归档同步内容。归档不得因为当前 delta 未包含某些旧 requirement，就从主 spec 删除这些既有内容。

如需预览而不实际执行，使用 `--dry-run` 参数。

### 2. 生命周期闭环

Spec 生命周期在此完成：
```
brainstorming → delta spec → 实施 → 验证 → 主 spec 合并 → design doc 标注 → 归档
```

## 退出条件

- 归档脚本执行成功（退出码 0）
- 归档目录 `openspec/changes/archive/[<issue_id>/]YYYY-MM-DD-<change-name>/` 存在
- 归档后的 `.impetus.yaml` 中 `archived: true`

归档脚本会把 `openspec/changes/<name>/` 移动到 `openspec/changes/archive/YYYY-MM-DD-<name>/`。归档成功后**不要再对原 change 名运行** `"$IMPETUS_BASH" "$IMPETUS_GUARD" <change-name> archive`，因为原活跃目录已经不存在。归档完整性以脚本退出码和归档目录状态为准。

## 完成

Impetus 流程全部完成。如需开始新工作，调用 `/impetus` 或 `/impetus-open`。
