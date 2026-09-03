---
name: subagent-driven-verify
description: "Impetus verify 阶段只读并行验证调度器。并行派发 test/spec/review/artifact/diff verifier，汇总报告；禁止修代码、提交或归档。"
---

# Subagent-Driven Verify：只读并行验证

本 skill 只用于 `/impetus-verify` 阶段。它复用 subagent / Task / multi-agent 的后台派发能力，但**不复用** `subagent-driven-development`，因为 verify 只能审计，不能实现或修复。

## 硬约束

- 所有 verifier 只读；不得修改源码、测试、OpenSpec、Design Doc、tasks.md
- 不得执行 `git commit`、amend、merge、push、archive
- 不得启动 fix agent；失败项统一回到 `/impetus-verify` 的失败决策点
- 每个 verifier 只能写自己的报告文件
- 只有 verify coordinator 可以写 `summary.md` 和正式 verify report
- 平台没有真实后台 subagent / Task / multi-agent 能力时，降级为 standard verify，不得在主窗口伪装并行

## 输出目录

```text
openspec/changes/<change-name>/.impetus/verify/run-<YYYYMMDD-HHMMSS>/
├── tests.md
├── spec-coverage.md
├── implementation-review.md
├── artifact-consistency.md
├── diff-scope.md
└── summary.md
```

## Verifier

| Verifier | 职责 | 阻断规则 |
|----------|------|----------|
| `test-verifier` | 执行/核验构建、单测、集成测试证据 | 构建失败、测试失败阻断 |
| `spec-coverage-verifier` | 核验 tasks/spec/proposal/design 覆盖 | 核心验收场景缺失阻断 |
| `implementation-review-verifier` | 调用 `impetus-review` | `FAIL` 阻断，`WARN` 进入用户决策 |
| `artifact-consistency-verifier` | 核验 Design Doc、delta spec、tasks、报告一致性 | 矛盾或缺失阻断 |
| `diff-scope-verifier` | 核验 diff 范围、分支状态、无密钥/无越权改动 | 安全/范围扩张阻断 |

## 调度流程

1. coordinator 创建 run 目录。
2. 并行派发 verifier；每个 prompt 必须写明只读硬约束和唯一输出文件。
3. 等待所有 verifier 返回。
4. coordinator 汇总 `summary.md`。
5. 若存在 blocking fail，返回 `/impetus-verify` Step 1b；不得在 verify 内修复。
6. 全部通过后，`/impetus-verify` 继续分支收尾和正式报告。

## Summary 格式

```markdown
# Verify Summary

| Check | Status | Blocking | Evidence |
|---|---|---:|---|
| tests | PASS | no | tests.md |
| spec-coverage | PASS | no | spec-coverage.md |
| implementation-review | WARN | no | implementation-review.md |
| artifact-consistency | PASS | no | artifact-consistency.md |
| diff-scope | PASS | no | diff-scope.md |

Result: PASS
```

`Result` 只能是 `PASS`、`WARN`、`FAIL`。存在 blocking 项时必须为 `FAIL`。
