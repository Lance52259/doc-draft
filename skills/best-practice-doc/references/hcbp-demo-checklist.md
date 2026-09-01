# hcbp-demo 变更清单速查

来源约定见仓库 [.cursor/rules/best-practice-style.mdc](https://github.com/Lance52259/hcbp-demo/blob/master/.cursor/rules/best-practice-style.mdc) 与模板：

- [templates/best_practice.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/best_practice.md)
- [templates/category_index.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/category_index.md)

## 已有服务 + 新实践

1. `docs/zh-cn/best-practices/{service}/{practice}.md` — create  
2. `docs/zh-cn/best-practices/{service}/index.md` — update「最佳实践列表」（文件名升序）  
3. `docs/zh-cn/SUMMARY.md` — update 增加实践链接  

## 新服务 + 新实践

1. 上表三项，其中 `index.md` 为 create  
2. `docs/zh-cn/best-practices/README.md` — update「文档导航」（服务链接英文升序）  
3. `SUMMARY.md` — 同时增加服务目录节点与简介节点  

## 示例对照

- 正文范例：`docs/zh-cn/best-practices/ecs/simple_instance.md`  
- 分类范例：`docs/zh-cn/best-practices/ecs/index.md`  
- 源脚本根：https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples  
