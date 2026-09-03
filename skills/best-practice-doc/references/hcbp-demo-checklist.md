# hcbp-demo 变更清单速查（中英双语）

来源约定见 [.cursor/rules/best-practice-style.mdc](https://github.com/Lance52259/hcbp-demo/blob/master/.cursor/rules/best-practice-style.mdc) 与模板：

- [templates/best_practice.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/best_practice.md)
- [templates/category_index.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/category_index.md)

## 强制顺序（1 → 8）

英文侧（2–5）先按**英文字母序**定目录；中文侧（6–8）**跟随**英文顺序。

| # | 文件 | 动作 |
|---|------|------|
| 1 | `docs/zh-cn/best-practices/{service}/{practice}.md` | create |
| 2 | `docs/en-us/best-practices/{service}/{practice}.md` | create |
| 3 | `docs/en-us/best-practices/{service}/index.md` | create（新服务）或 update 列表 |
| 4 | `docs/en-us/best-practices/README.md` | update 文档导航（**新服务必做**；仅新实践可跳过） |
| 5 | `docs/en-us/SUMMARY.md` | update：插入实践行；**新服务**再插入服务块 |
| 6 | `docs/zh-cn/best-practices/{service}/index.md` | create 或 update（顺序对齐英文 index） |
| 7 | `docs/zh-cn/best-practices/README.md` | update（顺序对齐英文 README；新服务必做） |
| 8 | `docs/zh-cn/SUMMARY.md` | update（对齐英文 SUMMARY；新服务含服务块） |

## 红线

- 中英 `SUMMARY.md` 一定存在，只定点插入，禁止整文件重写
- 禁止只生成中文、省略英文
- 禁止先改中文导航再改英文导航

## 示例对照

- 中文正文：`docs/zh-cn/best-practices/ecs/simple_instance.md`
- 英文正文：`docs/en-us/best-practices/ecs/simple_instance.md`
- 源脚本根：https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples
