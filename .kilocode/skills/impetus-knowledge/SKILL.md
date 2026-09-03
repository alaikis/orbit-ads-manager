---
name: impetus-knowledge
description: "Impetus 知识库整理。用 /impetus-knowledge 调用，将归档变更或某个主题整理到 Obsidian 等知识库。"
---

# Impetus 知识库整理

当用户想把 Impetus 已归档 change、spec、设计文档、验证报告，或某个主题整理成可复用知识库时，使用本 skill。

## 支持的存储目标

第一版支持：

- `obsidian` — 输出 Markdown 笔记到 Obsidian vault

## 模式

### 1. 归档导出

把已归档 change 导出为 Obsidian 笔记：

```bash
impetus knowledge . --target obsidian --source archive --vault <obsidian-vault-path>
```

只导出指定 change：

```bash
impetus knowledge . --target obsidian --source archive --change <archive-name-or-change-name> --vault <obsidian-vault-path>
```

命令会写入：

```text
<vault>/Impetus/Archive/<issue_id-or-no-issue>/<archive-name>.md
```

### 2. 主题整理

从归档中找出匹配来源，创建一个 AI 可继续整理的主题工作台笔记：

```bash
impetus knowledge . --target obsidian --source topic --topic "<主题>" --vault <obsidian-vault-path>
```

命令会写入：

```text
<vault>/Impetus/Topics/<主题>.md
```

创建后，读取笔记列出的来源，在 `## 整理结果` 中补充结构化知识：关键概念、设计决策、行为约束、测试经验、踩坑点、后续问题。

## 规则

- 不修改 archive 中的源产物。
- 保留指向归档来源的链接或路径。
- 整理时保留 requirement 名称、API 路径、错误码、明确的设计决策。
- 如果用户要求尚未支持的存储后端，说明当前只实现 `obsidian`，并询问是否先导出 Markdown 供手动导入。
