---
name: best-practice-doc
description: 从 terraform-provider-huaweicloud/examples 生成中英双语文档并 PR 到 Lance52259/hcbp-demo
version: "0.3.0"
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
| **3** | **create** 或 **update** | `docs/en-us/best-practices/{service}/index.md` | 英文分类页：不存在则新建；已存在则只在 Best Practices List 中按 **practice 文件名** 字母序插入本条 |
| **4** | **update** | `docs/en-us/best-practices/README.md` | 英文文档导航：按 **服务目录名** 字母序维护入口（**新增服务时必须插入**；已有服务仅新增实践时可跳过，因 README 不列单条实践） |
| **5** | **update** | `docs/en-us/SUMMARY.md` | 英文 TOC：**一定存在**，只定点插入。已有服务 → 在该服务块内按 practice 文件名字母序插入实践链接；**新增服务** → 再按服务目录名字母序插入「服务节点 + Introduction + 本实践」 |
| **6** | **create** 或 **update** | `docs/zh-cn/best-practices/{service}/index.md` | 中文分类页：不存在则新建；已存在则只追加本实践。**列表顺序与英文 index 中 practice 文件顺序一致** |
| **7** | **update** | `docs/zh-cn/best-practices/README.md` | 中文文档导航：**按英文 README 已确定的服务目录顺序** 同步（新增服务时插入中文入口；已有服务仅新增实践时可跳过） |
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
2. **index.md**：服务目录已存在 → **update** 列表；服务目录首次创建 → **create**（不要 update 不存在的文件）。
3. **README.md**：维护服务级文档导航；新增服务时插入；已有服务只加实践时通常不改。
4. 所有导航 `update` 必须是「基线全文 + 最小插入」；非空行数不得明显少于基线。
5. 优先由 doc-craft 编排层对导航做确定性补丁；模型至少输出步骤 1–2 的正文，新服务时输出中英 `index.md` 的 create。

---

## 不要做的事

- 不要只生成中文、省略英文（步骤 2–5 为必做链路）。
- 不要先改中文导航再改英文（会破坏「英文定序、中文跟随」）。
- 不要编造源脚本中不存在的 resource / data source / 参数。
- 不要自由新增模板未定义的一级/二级标题。
- 不要把元注释写进最终 Markdown；敏感信息用占位符。
- 链接使用半角 `()`；每个 `.md` 文件末尾保留一个空行。
- 不要用精简版目录替换完整 `SUMMARY.md`。

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

HCL 注释：声明了 `region` 写「在指定 region 下…」；否则写「在指定 region（region 参数缺省时默认继承当前 provider 块中所指定的 region）下…」。

---

## 英文正文结构（步骤 2，固定标题顺序）

对齐 `docs/en-us/best-practices/ecs/simple_instance.md`（与中文镜像，标题用英文）：

1. `# {English practice title}`（如：Deploy Basic Instance）
2. `## Application Scenario`
3. `## Related Resources/Data Sources`
   - Fixed lead-in: `This best practice involves the following main resources and data sources:`
   - `### Data Sources` / `### Resources` / `### Resource/Data Source Dependencies`
4. `## Operation Steps`
5. `## Reference Information`（或与现有英文文档一致的 Reference 章节名）

内容须与中文版同一套资源/参数/步骤，禁止中英不一致或英文臆造。

---

## 分类 index.md 结构

### 中文（步骤 6，仅新建服务时 create 全文）

1. `# 简介`
2. `## 什么是{中文全称}（{简称}）`
3. `## 最佳实践简述`
4. `## 最佳实践列表`（顺序 = 英文 index 的 practice 文件序）
5. `## 参考资料`

### 英文（步骤 3，仅新建服务时 create 全文）

1. `# Introduction`
2. `## What is {Full Name} ({Abbr})`
3. `## Best Practices Overview`
4. `## Best Practices List`（按 `{practice}.md` 字母序）
5. `## Reference Materials`

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
