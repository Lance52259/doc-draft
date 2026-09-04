---
name: best-practice-doc
description: 从 terraform-provider-huaweicloud/examples 生成中英双语文档并 PR 到 Lance52259/hcbp-demo
version: "0.3.4"
source_repo: https://github.com/huaweicloud/terraform-provider-huaweicloud
source_repo_slug: huaweicloud/terraform-provider-huaweicloud
source_examples: https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples
source_default_branch: master
target_repo: https://github.com/Lance52259/hcbp-demo
target_repo_slug: Lance52259/hcbp-demo
target_default_branch: master
---

# Skill：华为云最佳实践文档生成（hcbp-demo）

## 何时使用

当 **B 仓库** [huaweicloud/terraform-provider-huaweicloud](https://github.com/huaweicloud/terraform-provider-huaweicloud) 的 `examples/` 出现尚未对接至 **C 仓库** [Lance52259/hcbp-demo](https://github.com/Lance52259/hcbp-demo) 的最佳实践时使用本 Skill。

## 目标

依据 B 仓实践目录中的 Terraform HCL，按 C 仓模板与导航约定，生成 **中文 + 英文** Markdown 变更，并由 doc-craft **向 [hcbp-demo](https://github.com/Lance52259/hcbp-demo) 推送分支并创建 PR**（默认基于 `master`）。

---

## 强制处理顺序（必须按 1 → 8 执行）

新增一条最佳实践时，**按下列顺序**生成/更新文件。  
**步骤 2–5 先定英文侧目录顺序（英文字母序）**；步骤 6–8 的中文导航必须 **跟随英文侧已确定的顺序**，不得按中文标题自行重排。

| 步骤 | 动作 | 路径 | 说明 |
|------|------|------|------|
| **1** | **create** | `docs/zh-cn/best-practices/{service}/{practice}.md` | 中文版最佳实践正文 |
| **2** | **create** | `docs/en-us/best-practices/{service}/{practice}.md` | 英文版最佳实践正文（与中文同一 `{service}/{practice}`） |
| **3** | **create** 或 **update** | `docs/en-us/best-practices/{service}/index.md` | 英文分类页：**新建**须按「分类 index.md 规范」写满（对齐 Anti-DDoS，禁止精简版）；已存在则只在 Best Practices List 按 practice 文件名字母序插入本条（含一句话说明） |
| **4** | **update** | `docs/en-us/best-practices/README.md` | 英文文档导航：按服务目录名字母序插入；**标题与简介必须从英文 `index.md` 截取**（规则见「README 导航条目规范」）。禁止占位句。新服务必做；已有服务仅新增实践时可跳过 |
| **5** | **update** | `docs/en-us/SUMMARY.md` | 英文 TOC：**一定存在**，只定点插入。已有服务 → 在该服务块内按 practice 文件名字母序插入实践链接；**新增服务** → 再按服务目录名字母序插入「服务节点 + Introduction + 本实践」 |
| **6** | **create** 或 **update** | `docs/zh-cn/best-practices/{service}/index.md` | 中文分类页：**新建**须按「分类 index.md 规范」写满（与英文同等篇幅）；已存在则只追加列表项。**列表顺序与英文 index 一致** |
| **7** | **update** | `docs/zh-cn/best-practices/README.md` | 中文文档导航：顺序跟随英文；**标题与简介必须从中文 `index.md` 截取**（规则见下节「README 导航条目规范」）。禁止占位句。新服务必做；已有服务仅新增实践时可跳过 |
| **8** | **update** | `docs/zh-cn/SUMMARY.md` | 中文 TOC：**一定存在**，只定点插入。结构/相对位置与英文 SUMMARY **对齐**（同一 `{service}/{practice}` 链接路径；**新增服务**时同步插入服务块） |

### 顺序与排序原则

1. **英文先定序**：index / README / SUMMARY 的插入位置一律按 **英文路径名（`{service}`、`{practice}.md`）字母序** 计算。  
2. **中文跟随**：中文 index 列表、README 服务块、SUMMARY 服务块与实践行的相对顺序，与英文侧保持一致；不要用中文标题拼音或汉字顺序重排。  
3. **双语路径对称**：中英文使用相同的 `{service}`、`{practice}` 文件名；仅标题与正文语言不同。

### 命名约定

- `{service}`：华为云服务简称，**全小写**（如 `ecs`、`anti-ddos`），与 C 仓目录一致；不等于 B 仓 `examples/` 照抄（见 `configs/practice_mapping.yaml` 的 `service_aliases`）。
- `{practice}`：实践文件名（无 `.md`），中英共用；优先下划线风格，可用 `practice_aliases` 覆盖。

---

## 导航文件硬性规则（中英 SUMMARY / index / README）

> 历史失败：把 `SUMMARY.md` 整文件重写成短目录。中英两侧均 **严禁** 再次发生。

1. `docs/en-us/SUMMARY.md` 与 `docs/zh-cn/SUMMARY.md` **一定存在**：只能基于当前全文 **定点插入**，禁止整文件重写；禁止改标题（英文保持 `# Summary`，中文保持 `# Summary`）。
2. **index.md**：服务目录已存在 → **update** 列表项；服务目录首次创建 → **create** 全文（不要 update 不存在的文件）。**新建全文必须达到 Anti-DDoS 同等结构与篇幅**（见「分类 index.md 规范」），禁止只写一句话简介 + 裸链接列表。
3. **README.md（中英均适用，新增服务时）**：必须按「README 导航条目规范」从对应语言的 `index.md` 截取标题与简介，**禁止占位描述**。
4. 所有导航 `update` 必须是「基线全文 + 最小插入」；非空行数不得明显少于基线。
5. 优先由 doc-craft 编排层对导航做确定性补丁；模型至少输出步骤 1–2 的正文，新服务时输出中英 `index.md` 的 create（须符合「分类 index.md 规范」，含完整「What is / 什么是」多段介绍，供 README 截取）。

### README 导航条目规范（`docs/{zh-cn|en-us}/best-practices/README.md`）

每个服务在「文档导航 / Documentation Navigation」下为一块，格式：

```markdown
### [{标题}]({service}/index.md)

{简介段落}
```

#### 英文

- **标题**：`{Name from ## What is …} Best Practices`  
  例：`### [Anti-DDoS Best Practices](anti-ddos/index.md)`、`### [Application Operations Management (AOM) Best Practices](aom/index.md)`
- **简介**：取英文 `index.md` 中 `## What is …` 下**第一段**（过长则截到首句）
- **禁止**：`AAD Terraform best practices.`、`XXX related best practices.` 等占位句

#### 中文（与英文同一问题，必须同等遵守）

- **标题**：`{名称 from ## 什么是…}最佳实践`（与现网 Anti-DDoS / AOM 条目一致）  
  例：  
  - `### [Anti-DDoS最佳实践](anti-ddos/index.md)`  
  - `### [应用运维管理（AOM）最佳实践](aom/index.md)`  
  - `### [DDoS高防（AAD）最佳实践](aad/index.md)`
- **简介**：取中文 `index.md` 中 `## 什么是…` 下**第一段**（过长则截到首句），须是产品介绍，不是实践列表说明  
  **正确示例（对齐 Anti-DDoS）：**  
  `Anti-DDoS（Anti-Distributed Denial of Service）是华为云提供的分布式拒绝服务攻击防护服务，能够有效防护针对公网IP的DDoS攻击，保障业务的稳定运行。`  
  **错误示例（禁止）：**  
  - `AAD 相关 Terraform 最佳实践。`  
  - `DDoS高防 Terraform 最佳实践。`  
  - `介绍如何使用 Terraform 完成本实践。`
- 中文简介内容应与英文 README 截取自同一产品介绍语义，仅语言不同；**不得**中英文一侧详述、另一侧占位。

---

## 不要做的事

- 不要只生成中文、省略英文（步骤 2–5 为必做链路）。
- 不要先改中文导航再改英文（会破坏「英文定序、中文跟随」）。
- 不要编造源脚本中不存在的 resource / data source / 参数。
- 不要自由新增模板未定义的一级/二级标题。
- 不要把元注释写进最终 Markdown；敏感信息用占位符。
- 链接使用半角 `()`；每个 `.md` 文件末尾保留一个空行。
- 不要用精简版目录替换完整 `SUMMARY.md`。
- 英文 `## Reference Information` 中源码链接锚文本必须为  
  `Best Practice Source Code Reference For {English practice title}`，  
  禁止写成 `{title} Best Practice Source Code Reference`。
- 中文 / 英文 `best-practices/README.md` 新增服务条目的简介必须从对应 `index.md` 的「什么是 / What is」首段截取，  
  禁止 `AAD 相关 Terraform 最佳实践。` / `AAD Terraform best practices.` 等占位句。
- 新建服务的 `index.md` 禁止精简版（单段 What is、一句 Overview、无列表导语/无条目说明、参考资料写成 Provider 文档）。

---

## 路径推断

典型源路径：`examples/{service}/{practice}/`。

| 输入 | 输出 |
|------|------|
| `practice_id` = `examples/ecs/basic` | 中文正文 → `docs/zh-cn/best-practices/ecs/{practice}.md`；英文正文 → `docs/en-us/best-practices/ecs/{practice}.md`（`{practice}` 以映射/C 仓约定为准） |
| 服务是否首次 | 分别扫描 `docs/zh-cn/best-practices/{service}/` 与 `docs/en-us/best-practices/{service}/`（通常应同步存在；以实际为准） |

---

## 中文正文结构（步骤 1，固定标题顺序）

对齐 `docs/zh-cn/best-practices/ecs/simple_instance.md` 与 `templates/best_practice_template.md`：

1. `# {实践中文标题}`
2. `## 应用场景`
3. `## 相关资源/数据源`
   - 固定句：`本最佳实践涉及以下主要资源和数据源：`
   - `### 数据源` / `### 资源` / `### 资源/数据源依赖关系`
4. `## 操作步骤`（含脚本准备、逐步 HCL、可选 tfvars、init/plan/apply/show）
5. `## 参考信息`
   - 华为云对应产品文档 index
   - 固定：`[华为云Provider文档](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs)`
   - 源码目录链接，锚文本固定格式：`[{实践中文标题}最佳实践源码参考](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/{b_service}/{b_practice})`  
     - **正确示例：** `[AAD黑白名单最佳实践源码参考](...)`、`[Anti-DDoS基础防护最佳实践源码参考](...)`  
     - **错误示例（禁止）：** `Best Practice Source Code Reference For ...`（中文文档勿用英文锚文本）、`源码参考`（过于简略）

HCL 注释：声明了 `region` 写「在指定 region 下…」；否则写「在指定 region（region 参数缺省时默认继承当前 provider 块中所指定的 region）下…」。

---

## 英文正文结构（步骤 2，固定标题顺序）

对齐 `docs/en-us/best-practices/ecs/simple_instance.md` / `anti-ddos/basic.md`（与中文镜像，标题用英文）：

1. `# {English practice title}`（如：Deploy Basic Instance）
2. `## Application Scenario`
3. `## Related Resources/Data Sources`
   - Fixed lead-in: `This best practice involves the following main resources and data sources:`
   - `### Data Sources` / `### Resources` / `### Resource/Data Source Dependencies`
4. `## Operation Steps`
5. `## Reference Information`
   - Huawei Cloud product documentation index for the service
   - Fixed: `[Huawei Cloud Provider Documentation](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs)`
   - **Source code link — mandatory English anchor text format（注意大小写与用词，勿自行改写）：**  
     `[Best Practice Source Code Reference For {English practice title}](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/{b_service}/{b_practice})`  
     - `{English practice title}` 与正文 `#` 标题语义一致（服务名 + 场景，可用名词短语；不必强行以 Deploy 开头）  
     - **正确示例：**  
       `[Best Practice Source Code Reference For AAD Black/White Lists](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/aad/black-white-lists)`  
       `[Best Practice Source Code Reference For Anti-DDoS Basic Protection](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples/antiddos/basic)`  
     - **错误示例（禁止）：**  
       `AAD Black/White Lists Best Practice Source Code Reference`（语序错误）  
       `Best practice source code reference for ...`（大小写错误）

内容须与中文版同一套资源/参数/步骤，禁止中英不一致或英文臆造。源码 URL 使用 B 仓真实 `examples/` 路径（可能与 C 仓 `{service}/{practice}` 映射名不同）。

---

## 分类 index.md 规范（新建服务必遵；对齐 `anti-ddos/index.md`）

> 权威对照：C 仓 `docs/{zh-cn|en-us}/best-practices/anti-ddos/index.md` 与 `templates/category_index.md`。  
> **禁止**生成 AAD 曾出现的精简版（见文末「错误示例」）。

仅 **服务目录首次出现** 时 `create` 全文；已有服务只 **update**「最佳实践列表 / Best Practices List」中的一条（含链接 + 一句话说明）。

### 英文（步骤 3）固定骨架与内容量

```markdown
# Introduction

## What is {ServiceName}

{两到三段产品介绍，勿压缩成一句}

## Best Practices Overview

{两段固定套话，见下}

## Best Practices List

This section contains the following best practices:

* [{Practice title}]({practice}.md) - {one-sentence description covering main resources/steps}.

## Reference Materials

- [Huawei Cloud {ServiceName} Product Documentation](https://support.huaweicloud.com/{product}/index.html)
- [Terraform Official Documentation](https://www.terraform.io/docs/index.html)
```

| 章节 | 硬性要求 |
|------|----------|
| `# Introduction` | 固定；不要改成服务名 |
| `## What is …` | 标题取产品英文名；全称与缩写相同时可写 `## What is Anti-DDoS`，否则 `## What is Advanced Anti-DDoS (AAD)`。正文 **至少两段**（宜 2–3 段）：职责/能力 → 典型能力或防护/部署模式 → 可选运维价值。须像产品文档介绍，**不是**「用 Terraform 部署本实践」的一句话 |
| `## Best Practices Overview` | **必须两段**，套用下列句式（仅替换服务名与资源表述，勿自行改写成一句短述）： |
| | ① `This section provides best practice examples for using Terraform to automatically deploy and manage Huawei Cloud {Name}, helping you understand how to efficiently manage cloud {Name} … resources using Infrastructure as Code (IaC).` |
| | ② `Through the best practices in this section, you can learn the main deployment processes for {Name} … resources. These best practices will help you quickly get started with automated {Name} deployment and lay a solid foundation for subsequent {Name} management and operation work.` |
| `## Best Practices List` | 固定导语一行：`This section contains the following best practices:`。列表用 `*`（不是仅 `-` 裸链接）。每条：`* [Title](file.md) - Introduces how to use Terraform to …`（说明须点出主要资源/步骤，对齐 Anti-DDoS 条目长度） |
| `## Reference Materials` | 两条：① 华为云该产品 Supports index；② **固定** `[Terraform Official Documentation](https://www.terraform.io/docs/index.html)`。**不要**在 index 放 Provider 文档链接（Provider 属于实践正文 Reference Information） |

### 中文（步骤 6）固定骨架与内容量

```markdown
# 简介

## 什么是{中文全称}（{简称}）

{两到三段产品介绍，勿压缩成一句}

## 最佳实践简述

{两段固定套话，见下}

## 最佳实践列表

本章节包含以下最佳实践：

* [{实践中文标题}]({practice}.md) - {一句话说明主要资源/步骤}。

## 参考资料

- [华为云{产品}产品文档](https://support.huaweicloud.com/{product}/index.html)
- [Terraform官方文档](https://www.terraform.io/docs/index.html)
```

| 章节 | 硬性要求 |
|------|----------|
| `# 简介` | 固定两字标题 |
| `## 什么是…` | **至少两段**产品介绍（对齐中文 Anti-DDoS），可引用华为云 Supports 语义，勿编造不存在的计费/规格细节 |
| `## 最佳实践简述` | **必须两段**：① `本章节提供了使用Terraform自动化部署和管理华为云{全称}（{简称}）的最佳实践示例，帮助您了解如何利用Infrastructure as Code（IaC）的方式高效地管理云上的{简称}…资源。` ② `通过本章节的最佳实践，您可以学习到主要的{简称}…资源的部署流程，这些最佳实践将帮助您快速上手{简称}的自动化部署，并为后续的…管理和运维工作奠定坚实基础。` |
| `## 最佳实践列表` | 固定导语：`本章节包含以下最佳实践：`。每条 `* [标题](file.md) - 介绍如何使用Terraform…`；顺序与英文 index 一致 |
| `## 参考资料` | 产品文档 + **固定** Terraform 官方文档；不要用 Provider 文档替代第二条 |

### 错误示例（AAD 精简版 — 禁止再现）

```markdown
## What is Advanced Anti-DDoS (AAD)
Advanced Anti-DDoS (AAD) is a professional DDoS protection service…（仅一段）

## Best Practices Overview
This section provides best practices for deploying and configuring AAD instances using Terraform, including…（仅一句，未用固定套话）

## Best Practices List
- [Deploy AAD Black and White Lists](black_white_lists.md)   ← 缺导语、缺「 - 说明」、列表符号不规范

## Reference Materials
- … Product Documentation
- [Huawei Cloud Provider Documentation](…)   ← 错误：index 应用 Terraform Official Documentation
```

中英 `index.md` 须语义对齐、篇幅相当；不得英文写满、中文精简（或反之）。

---

## 硬性约束（输出）

1. **只输出一个 JSON 对象**，不要 Markdown 围栏或前后解释。
2. `files[]` **至少**包含步骤 1、2 的正文 `create`；新服务时还须包含中英 `index.md` 的 `create`。
3. 导航类 `update` 若输出，必须是基线 + 最小插入；**禁止**缩成短目录。
4. `path` 相对 C 仓根，禁止 `..`，且落在 `docs/zh-cn/` 或 `docs/en-us/`。
5. 参数与依赖必须来自源 HCL。
6. `summary` 用 **英文** 一句话（供 PR 说明）。

## 输出 schema（示意）

```json
{
  "practice_id": "examples/ecs/basic",
  "summary": "Add bilingual ECS basic instance best-practice docs and navigation",
  "files": [
    {
      "path": "docs/zh-cn/best-practices/ecs/simple_instance.md",
      "action": "create",
      "content": "……中文正文……\n"
    },
    {
      "path": "docs/en-us/best-practices/ecs/simple_instance.md",
      "action": "create",
      "content": "……English body……\n"
    }
  ]
}
```

新服务时继续附带中英 `index.md`（create）。`SUMMARY.md` / `README.md` 优先由编排层按步骤 3–8 补丁；模型输出导航时须遵守「英文定序、中文跟随」。
