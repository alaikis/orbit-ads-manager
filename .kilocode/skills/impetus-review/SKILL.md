---
name: impetus-review
description: "Impetus 实现审查门禁。可手动调用，也会在 verify 阶段自动调用。只读审查当前 change/diff/commit range，生成 bundle.json、findings.json、report.md、report.html，不修代码、不提交、不归档。"
---

# Impetus Review：实现审查门禁

`impetus-review` 是独立审查能力，不绑定具体语言。第一版内置 JVM/Java 后端规则集，后续可扩展 Node、Python、Go 等规则。

## 硬约束

- 只读审查，不修改源码、测试、OpenSpec、Design Doc、tasks.md
- 不执行 `git commit`、不 amend、不创建归档
- 不启动 fix agent；发现问题只写报告
- 手动调用时由用户决定是否修复；verify 自动调用时由 `/impetus-verify` 统一处理失败
- 中文需求/中文 workflow 下，报告与建议使用中文

## 默认输出

输出目录由调用方提供；未提供时使用：

```text
openspec/changes/<change-name>/.impetus/review/run-<YYYYMMDD-HHMMSS>/
```

产物：

| 文件 | 说明 |
|------|------|
| `bundle.json` | 审查包：diff、文件内容、Spec 上下文、项目规则 |
| `findings.json` | 结构化发现项 |
| `report.md` | Markdown 审查报告 |
| `report.html` | HTML 审查报告 |

结果等级：

| 结果 | 含义 |
|------|------|
| `PASS` | 无阻断问题 |
| `WARN` | 仅有 low / 可接受问题 |
| `FAIL` | 存在 high / critical / blocking 问题 |

## Step 1：确定审查范围

优先使用调用方提供的范围：

| 参数 | 说明 |
|------|------|
| `--change <name>` | Impetus change 名称 |
| `--base <ref>` | 基准 ref；未提供时优先读取 plan 的 `base-ref`，再回退到 `test` |
| `--target <ref>` | 目标 ref；未提供时使用当前 `HEAD` |
| `--spec-dir <path>` | OpenSpec delta spec 目录；存在时启用需求/Spec 合规审查 |
| `--project-rules <path>` | 项目级 Markdown 规则 |
| `--output-dir <path>` | 输出目录 |
| `--profile <standard|strict|security|architecture>` | 审查 profile，默认 `standard` |

始终使用本地 Git 对比，不执行 `git fetch`，不更新远端分支。

## Step 2：构建审查包

从本地 Git 读取增量：

```bash
git diff --name-only --diff-filter=ACMR <base>...<target>
git diff --name-status --diff-filter=ACMR <base>...<target>
git diff --numstat --diff-filter=ACMR <base>...<target>
git diff --unified=3 <base>...<target> -- <file>
git show <target>:<file>
```

第一版审查 `*.java`、`*.xml`、`*.yml`、`*.yaml`、`*.properties`、`*.sql`；排除 `*.md`、`*.txt`、`docs/**`、`README*`。若未发现 JVM/Java 相关变更，仍生成空审查报告，结果为 `PASS`，并说明“当前规则集无可审查文件”。

如指定 `--spec-dir`，读取该目录下 `spec.md` 或 `*.md` 写入 `bundle.json`；未指定时不得读取 OpenSpec 内容，也不得产生 `OpenSpec合规` 发现项。

## Step 3：按规则集审查

阅读 `references/report-contract.md` + `bundle.json` + 项目规则（如有），仅审查目标相对基准的新增/修改内容。

第一版使用 JVM/Java 规则集：

| 轮次 | 规则文件 | 简称 | 审查目标 |
|------|----------|------|----------|
| 第一轮 | `references/arch-logic-rule.md` | `arch` | 架构、接口、数据库、业务逻辑、异常、事务、Spec 合规 |
| 第二轮 | `references/security-rule.md` | `sec` | SQL 注入、XSS、越权、敏感数据、接口安全 |
| 第三轮 | `references/performance-rule.md` | `perf` | 数据量、数据库、集合、并发、资源、IO |
| 第四轮 | `references/maintainability-rule.md` | `maint` | 命名、可读性、注释、职责、兼容演进 |
| 第五轮 | `references/test-rule.md` | `test` | 单测、异常场景、边界、自测证据 |
| 第六轮 | `references/style-rule.md` | `style` | 常量、格式、OOP、日志、版权 |

严重级别映射固定：

- P0 → `high`
- P1 → `medium`
- P2 → `low`

安全漏洞、核心验收失败、数据破坏风险、越权风险必须标为阻断。

## Step 4：写入结果

根据 `references/report-contract.md` 生成：

- `bundle.json`
- `findings.json`
- `report.md`
- `report.html`

`findings.json` 必须包含可机器读取的结果：

```json
{
  "gate": "PASS",
  "profile": "standard",
  "summary": "未发现阻断问题。",
  "issues": []
}
```

当存在 high / blocking 问题时，`gate` 必须为 `FAIL`；仅有 low/medium 非阻断问题时可为 `WARN`。

## Verify 阶段集成

`/impetus-verify` 自动调用本 skill 作为 Implementation Review Gate。verify 阶段使用时：

- `FAIL`：进入 verify 失败决策点，通常回 build 修复
- `WARN`：记录到验证报告，由用户决定接受或回 build 修复
- `PASS`：继续后续验证/归档

`impetus-review` 本身永远不修复。

## 参考资料

- 输入输出格式：`references/report-contract.md`
- 通用审查指引：`references/common-review-guide.md`
- JVM/Java 规则集：`references/*-rule.md`
