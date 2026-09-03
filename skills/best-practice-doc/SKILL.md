---
name: best-practice-doc
description: 从 terraform-provider-huaweicloud/examples 生成文档并 PR 到 Lance52259/hcbp-demo
version: "0.2.5"
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

依据 B 仓实践目录中的 Terraform HCL，按 C 仓模板与导航约定生成 Markdown 变更，并由 doc-craft **向 [hcbp-demo](https://github.com/Lance52259/hcbp-demo) 推送分支并创建 PR**（默认基于 `master`）。

---

## 新增一条最佳实践时必须改动的内容

以中文文档为主（`docs/zh-cn/`）。英文目录 `docs/en-us/` **默认不生成**（除非调用方明确要求双语）。

### A. 始终需要（新建实践正文）

| 动作 | 路径 | 依据模板 |
|------|------|----------|
| **create** | `docs/zh-cn/best-practices/{service}/{practice}.md` | `templates/best_practice_template.md`（对齐 C 仓 `templates/best_practice.md`） |

约定：

- `{service}`：华为云服务简称，**全小写**（如 `ecs`、`vpc`、`anti-ddos`、`sfs-turbo`），与 C 仓现有目录名一致；**不等于** B 仓 `examples/` 一级目录照抄（如 B 的 `antiddos` → C 的 `anti-ddos`，见 `configs/practice_mapping.yaml` 的 `service_aliases`）。
- `{practice}`：实践文件名（无 `.md`），优先对齐 C 仓既有风格（下划线、产品语义名）；可用 `practice_aliases` 覆盖。detect 时会按别名 + 去分隔符模糊匹配判断是否已对接。
- 正文结构、固定标题 **不得自行增减章节标题**；严格按模板章节顺序。

### B. 导航规则（按服务目录是否已存在分支）

`docs/zh-cn/SUMMARY.md` **一定存在**，任何情况下都只做**定点插入**，禁止整文件重写。

#### 服务目录已存在（`docs/zh-cn/best-practices/{service}/` 已有）

| 动作 | 路径 | 改什么 |
|------|------|--------|
| **update** | `docs/zh-cn/best-practices/{service}/index.md` | **仅**在「最佳实践列表」描述中按文件名升序插入本实践一条链接说明 |
| **update** | `docs/zh-cn/SUMMARY.md` | **仅**在该服务小节的正确位置插入本实践一行导航链接 |

不改 `best-practices/README.md`（服务已在文档导航中）。

#### 服务目录首次创建（尚无 `{service}/`）

| 动作 | 路径 | 改什么 |
|------|------|--------|
| **create** | `docs/zh-cn/best-practices/{service}/index.md` | 目录不存在时 `index.md` 也不存在，按分类模板**新建**（含简介与本实践列表项） |
| **update** | `docs/zh-cn/best-practices/README.md` | **仅**在「文档导航」正确位置插入本服务入口（产品名 + 链接 + 一段简介） |
| **update** | `docs/zh-cn/SUMMARY.md` | **仅**在「最佳实践」下按服务目录名字母序插入：服务节点 + 简介 + 本实践链接 |

### C. 导航文件硬性规则 — 必须遵守

> 历史失败模式：把 `SUMMARY.md` 整文件重写成十几行（如把 `# Summary` 改成 `# 目录`）。**严禁再次发生。**

1. **`SUMMARY.md` 必存在、只追加**：永远基于 C 仓当前全文，在正确位置插入新行；不得删除、改写、重排已有条目；不得改标题。
2. **`index.md`**：目录已存在 → 只追加列表项；目录首次创建 → **create** 新文件（不要去“更新”一个不存在的文件）。
3. **`README.md` 文档导航**：仅服务首次接入时插入一块；已有服务不要动。
4. **优先编排层补丁**：doc-craft 对 SUMMARY / 已有 index 列表 / README 做确定性插入；AI 优先只生成实践正文，新服务时再 create `index.md`。

### D. 不要做的事

- 不要编造源脚本中不存在的 resource / data source / 参数。
- 不要自由新增模板未定义的一级/二级标题。
- 不要把生成说明性括号、`本括号中的内容…` 等元注释写进最终 Markdown。
- 不要在文档中写入真实 AK/SK/密码；示例使用占位符。
- 跳转链接的 Markdown 括号使用 **半角** `()`。
- 每个生成的 `.md` **文件末尾保留一个空行**。
- 不要用「精简版目录」替换完整 `SUMMARY.md`。

---

## 路径推断规则（从 practice_id / 源路径）

典型源路径：`examples/{service}/{practice}/`（来自 [terraform-provider-huaweicloud/examples](https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples)）。

| 输入 | 输出 |
|------|------|
| `practice_id` = `examples/ecs/basic` | 正文 → `docs/zh-cn/best-practices/ecs/basic.md`（若 C 仓已有约定别名如 `simple_instance.md`，以现有命名或调用方 `target_path` 为准） |
| 服务目录是否存在 | 扫描 C 仓是否已有 `docs/zh-cn/best-practices/{service}/` |

若调用方给出「建议写入路径」，优先采用；同时仍须按上表补齐 index / README / SUMMARY 的导航更新（优先由编排层补丁）。

---

## 正文文档结构（固定标题，顺序不可改）

对齐已落地文档（如 `docs/zh-cn/best-practices/ecs/simple_instance.md`）与模板：

1. `# {实践中文标题}`（如：部署基础实例）
2. `## 应用场景`
3. `## 相关资源/数据源`
   - 固定句：`本最佳实践涉及以下主要资源和数据源：`
   - `### 数据源`（列表项格式：`[释义（data.huaweicloud_xxx）](provider文档URL)`；顺序与脚本一致；无则写「本实践未使用数据源。」）
   - `### 资源`（`[释义（huaweicloud_xxx）](URL)`）
   - `### 资源/数据源依赖关系`（代码块树形依赖；须覆盖全部数据源与资源；无依赖则如实反映）
4. `## 操作步骤`
   - `### 1. 脚本准备`（固定说明 + 链到 `../../introductions/prepare_before_deploy.md`）
   - `### 2.` … 按脚本顺序逐步讲解数据源/资源（含 ` ```hcl ` 与 `**参数说明**：`）
   - 若有变量：`### n. 预设资源部署所需的入参（可选）`（含 `terraform.tfvars` 示例与使用方法；敏感值占位）
   - 最后一节：`### x. 初始化并应用Terraform配置`（init / plan / apply / show；apply/show 对象用本实践主体名称）
5. `## 参考信息`
   - 华为云对应产品文档 index
   - 固定：`[华为云Provider文档](https://registry.terraform.io/providers/huaweicloud/huaweicloud/latest/docs)`
   - 源码目录链接（B 仓 / examples 下该实践路径）

HCL 注释风格：

- 数据源/资源若声明了 `region`：写「在指定 region 下…」
- 未声明 `region`：写「在指定 region（region 参数缺省时默认继承当前 provider 块中所指定的 region）下…」

---

## 分类 index.md 结构（仅新建服务时）

固定标题顺序：

1. `# 简介`
2. `## 什么是{中文全称}（{简称}）`
3. `## 最佳实践简述`
4. `## 最佳实践列表`（含本实践链接，文件名升序）
5. `## 参考资料`（产品文档 + 固定 Terraform 官方文档链接）

---

## 硬性约束（输出）

1. **只输出一个 JSON 对象**，不要 Markdown 围栏或前后解释。
2. `files[]` **至少**包含实践正文的 `create`；新服务时包含服务 `index.md` 的 `create`。
3. 导航类 `update`（`SUMMARY.md` / 已有 `index.md` / `README.md`）可选；若输出，`content` 必须是 **基线全文 + 最小插入**后的完整正文，**禁止**缩成短目录。
4. `path` 相对 C 仓根，禁止 `..`，且落在 `docs/` 下。
5. 参数、资源类型、依赖关系必须来自源 HCL，禁止臆造。
6. 简体中文正文；专有名词可保留英文。
7. `summary` 用 **英文** 一句话（供 PR 说明）。

## 输出 schema

```json
{
  "practice_id": "examples/ecs/basic",
  "summary": "Add ECS basic instance best-practice doc and navigation entries",
  "files": [
    {
      "path": "docs/zh-cn/best-practices/ecs/basic.md",
      "action": "create",
      "content": "……完整 markdown……\n"
    }
  ]
}
```

可选（仅当基线已提供且能安全插入时）再附带导航 `update`。新建服务时，`files` 还应包含服务 `index.md`（create）。`SUMMARY.md` / `README.md` 导航优先由 doc-craft 编排层补丁，不强制模型输出。
