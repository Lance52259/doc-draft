# hcbp-demo 变更清单速查

来源约定见仓库 [.cursor/rules/best-practice-style.mdc](https://github.com/Lance52259/hcbp-demo/blob/master/.cursor/rules/best-practice-style.mdc) 与模板：

- [templates/best_practice.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/best_practice.md)
- [templates/category_index.md](https://github.com/Lance52259/hcbp-demo/blob/master/templates/category_index.md)

## 共性

1. `docs/zh-cn/best-practices/{service}/{practice}.md` — create  
2. `docs/zh-cn/SUMMARY.md` — **一定存在**；只在正确位置插入链接（禁止整文件重写）

## 服务目录已存在

1. `docs/zh-cn/best-practices/{service}/index.md` — **仅**在最佳实践列表中追加本实践  
2. 不改 `best-practices/README.md`

## 服务目录首次创建

1. `docs/zh-cn/best-practices/{service}/index.md` — **create**（此时文件尚不存在）  
2. `docs/zh-cn/best-practices/README.md` — 在文档导航正确位置插入本服务入口  
3. `SUMMARY.md` — 插入服务节点 + 简介 + 本实践

## 导航红线

- **禁止**重写整个 `SUMMARY.md`
- doc-craft 编排层对 SUMMARY / 已有 index / README 做确定性插入

## 示例对照

- 正文范例：`docs/zh-cn/best-practices/ecs/simple_instance.md`  
- 分类范例：`docs/zh-cn/best-practices/ecs/index.md`  
- 源脚本根：https://github.com/huaweicloud/terraform-provider-huaweicloud/tree/master/examples  
