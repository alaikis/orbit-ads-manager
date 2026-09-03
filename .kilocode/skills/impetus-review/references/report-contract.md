# 输入输出契约

## 输入参数

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `target`（位置参数） | string | 否 | 目标分支名。不传则使用当前分支；传入则使用指定分支。基准分支始终为 `test` |
| `--spec-dir` | string | 否 | OpenSpec 规格文件目录路径。指定后读取该目录下 `spec.md`（或 `*.md`）文件用于 9.2 节 OpenSpec 工作流审计；不指定则跳过 OpenSpec 审计 |
| `--project-rules` | string | 否 | 项目级规则 Markdown 文件路径 |
| `--output-dir` | string | 否 | 显式指定输出目录 |

始终使用本地分支对比模式（`test...<目标分支>`），不执行 `git fetch`，不更新远端分支。

默认输出目录：`<repo>/openspec/changes/<change-name>/.impetus/review/run-<YYYYMMDD-HHMMSS>/`

目录名安全处理：保留中文、英文、数字、点、横线、下划线；替换 Windows 非法字符 `< > : " / \ | ? *`。

## bundle.json

承载增量评审上下文：

```json
{
  "title": "代码评审报告",
  "generatedAt": "2026-03-13T01:23:45+08:00",
  "repoPath": "D:/workspace/java/ai",
  "mode": "local",
  "baseBranch": "test",
  "targetBranch": "feature-x",
  "compareRef": "test...feature-x",
  "specDir": "D:/workspace/java/ai/openspec",
  "specContent": "# API 规格定义\n...",
  "projectRulesPath": "D:/workspace/java/ai/docs/rules.md",
  "projectRulesContent": "# 团队规则...",
  "includedGlobs": ["*.java", "*.xml", "*.yml", "*.yaml", "*.properties", "*.sql"],
  "excludedGlobs": ["*.md", "*.txt", "docs/**", "README*"],
  "summary": {
    "totalChangedFiles": 3,
    "reviewableFiles": 2,
    "filteredOutFiles": 1,
    "totalAdditions": 40,
    "totalDeletions": 8
  },
  "files": [
    {
      "path": "backend/src/main/java/com/example/App.java",
      "status": "M",
      "additions": 12,
      "deletions": 3,
      "diff": "@@ ...",
      "content": "package ..."
    }
  ],
  "ignoredFiles": [
    {
      "path": "README.md",
      "reason": "excluded-by-pattern"
    }
  ]
}
```

- `mode` 固定为 `"local"`，始终使用本地分支对比
- `specDir` 和 `specContent`：仅在指定 `--spec-dir` 时存在；未指定时这两个字段**禁止出现**（而非设为 `null`），且**禁止读取任何 OpenSpec 相关文件**
- `files` 仅包含纳入评审范围的文件
- `content` 为目标分支版本内容，便于结合 diff 做判断
- `status` 使用 Git `name-status` 风格：`A`、`M`、`R100`

## findings.json

评审结果结构：

```json
{
  "summary": "本次改动整体风险可控，但存在 2 个需要处理的问题。",
  "conclusion": "建议修复问题后再合并。",
  "issues": [
    {
      "severity": "high",
      "file": "backend/src/main/java/com/example/App.java",
      "line": 42,
      "description": "error 日志缺少异常对象，失败原因会丢失。",
      "suggestion": "将异常对象作为最后一个参数传入日志方法，例如 log.error(\"处理失败\", e)。",
      "category": "日志与异常",
      "relatedRule": "arch:7.1.4",
      "codeSnippet": "log.error(\"处理失败\");",
      "fixExample": "log.error(\"处理失败\", e);"
    }
  ]
}
```

### 字段规范

| 字段 | 必填 | 说明 |
|------|------|------|
| `severity` | 是 | `high`（P0）/ `medium`（P1）/ `low`（P2），**必须与规则优先级固定绑定**，禁止按主观判断调整 |
| `file` | 是 | 相对仓库根目录的文件路径 |
| `line` | 是 | 行号 |
| `description` | 是 | **做了什么** + **为什么有问题** + **可能后果** |
| `suggestion` | 是 | **可操作的修复步骤**或代码示例，禁止"建议优化"等模糊描述 |
| `category` | 是 | 必须使用：`安全漏洞` / `消息可靠性` / `性能` / `数据正确性` / `代码规范` / `日志与异常` / `设计缺陷` / `OpenSpec合规` |
| `relatedRule` | 是 | 取 `references/` 目录下规则文件中的编号，格式 `{规则文件简称}:{编号}`（如 `arch:1.1.1`、`sec:1.1.1`、`perf:1.1.1`、`maint:1.1.1`、`test:1.1.1`、`style:1.1.1`）。简称映射：arch=架构与逻辑正确性，sec=安全与风险，perf=性能与资源，maint=可维护性可读性命名，test=测试覆盖，style=风格与格式 |
| `codeSnippet` | 是 | 问题代码片段（1-5行），**禁止省略** |
| `fixExample` | 是 | 修复后代码，**禁止省略** |

- `summary` 和 `conclusion` 使用中文
- 无问题时 `issues` 置空数组
- 禁止编造不存在的文件或行号
- **OpenSpec 审计条件（强制）**：仅在 `bundle.json` 包含 `specDir` 字段时，`category` 才允许使用 `OpenSpec合规`；未指定 `--spec-dir` 时，**禁止读取任何 OpenSpec 相关文件，禁止执行 arch-logic-rule 9.2 节审查，禁止产生 `OpenSpec合规` 类别的发现项**

### suggestion 编写指南

好的 suggestion 须满足：**具体**（改哪个方法/参数）、**可操作**（代码示例或明确步骤）、**有上下文**（给出类名或方法签名）。

| 差 | 好 |
|----|-----|
| 建议修复类型不匹配问题 | 将 `DTO.officeId` 类型从 `String` 改为 `Long`，或在 `copyProperties` 后手动 `record.setOfficeId(Long.valueOf(dto.getOfficeId()))` |
| 建议优化性能 | 在 `EmployeeService.findMapByIds` 中增加分片查询：`Lists.partition(ids, 500).stream().flatMap(partition -> mapper.findByIds(partition).stream()).collect(...)` |

## 报告输出

最终产物：`report.md`、`report.html`

中间产物：`bundle.json`、`findings.json`

> 所有产物共存于同一输出目录，均保留不自动清理。每次代码对比仅产生一个输出目录。

报告标题固定：`代码评审报告`

### 报告结构（仅4节，禁止添加其他章节）

1. **基础信息** — **基准分支**、**目标分支**、**生成时间**、**模式**、**OpenSpec文件目录**，每个指标一行
2. **改动统计** — 第一行展示问题总数及 HIGH/MEDIUM/LOW 各级别数量；第二行展示**变更文件总数**、**可审查文件数**、**新增行数**、**删除行数**，一行展示
3. **问题列表** — 按严重程度降序排列（high → medium → low），同级别内按文件路径排序；每个问题含：问题描述、文件、分类、规则、问题代码、修复建议
4. **结论** — 参考以下格式汇总：
   - **最严重的问题**：总结 HIGH 问题总数及明细，每个问题一行展示，建议修复上述问题后再合并到 test 分支。
   - **需要关注的事项**：MEDIUM、LOW问题进行总结。
   - 若某级别无问题则省略该条目；全量无问题时输出"无问题，可以合并。"

**禁止**输出"变更文件清单"章节。

### report.html 规范

`report.html` 必须是**独立可执行**的 HTML 文件（内联 CSS/JS，无外部依赖），浏览器直接打开即可渲染。

**配色方案**：天蓝色主题（Sky Blue），主色值参考：

| 用途 | 色值 | 说明 |
|------|------|------|
| 页面背景 | `#E0F2FE` | 天蓝色浅底 |
| 卡片背景 | `#FFFFFF` | 白色卡片 |
| 主色调 | `#0284C7` | 天蓝色主色（标题、徽章） |
| 高严重度 | `#DC2626` | 红色 |
| 中严重度 | `#D97706` | 橙色 |
| 低严重度 | `#059669` | 绿色 |
| 边框 | `#BAE6FD` | 天蓝色浅边框 |
| 代码背景 | `#1E293B` | 深色底（问题代码 + 修复建议） |
| 代码文字 | `#E2E8F0` | 浅灰白（深色底上的文字） |
| 正文 | `#1E293B` | 深灰蓝 |

**字号规范**（固定 `px`，禁止使用 `rem`/`em`/`vw` 等相对单位）：

| 元素 | 字号 | 字重 | 说明 |
|------|------|------|------|
| 页面标题 | `24px` | `700` | 报告主标题 |
| 卡片标题 | `18px` | `600` | 各节标题 |
| 正文 | `14px` | `400` | 描述、统计等 |
| 代码块 | `13px` | `400` | 问题代码、修复建议（等宽字体） |
| 徽章 | `12px` | `600` | 严重度标签 |
| 小字 | `12px` | `400` | 文件路径、规则编号 |

**HTML 结构要求**：

```html
<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>代码评审报告</title>
  <style>
    /* 内联完整样式，天蓝色主题 */
    /* 响应式布局，max-width: 960px 居中 */
    /* 卡片式分节，圆角 + 阴影 */
    /* 代码块使用等宽字体 + 黑色深色背景(#1E293B) + 浅色文字(#E2E8F0) */
    /* 严重度徽章：彩色圆角标签 */
    /* 表格：斑马纹 + 悬停高亮 */
  </style>
</head>
<body>
  <!-- 1. 基础信息卡片 -->
  <!-- 2. 改动统计卡片：问题总数（含各级别徽章）+ 变更统计 -->
  <!-- 3. 问题列表：卡片式，每个问题一个卡片，含严重度徽章 + 描述 + 文件路径 + 分类 + 规则 + 代码块 + 修复建议 -->
  <!-- 4. 结论卡片 -->
</body>
</html>
```

**交互功能**（纯 JS，无外部依赖）：
- 问题卡片可按严重度筛选（全部 / HIGH / MEDIUM / LOW）
- 代码块点击可复制
