# doc-craft

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)
[![CI](https://github.com/Lance52259/doc-draft/actions/workflows/test-workflow.yml/badge.svg)](https://github.com/Lance52259/doc-draft/actions/workflows/test-workflow.yml)

从 [Huawei Cloud Terraform Provider](https://github.com/huaweicloud/terraform-provider-huaweicloud) 的 `examples/` 自动发现新增最佳实践，按 Skill 约束调用 DeepSeek 生成 **中英双语文档**，并向文档仓提交 Pull Request。

---

## 目录

- [背景](#背景)
- [功能特性](#功能特性)
- [工作原理](#工作原理)
- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [配置](#配置)
- [CLI 用法](#cli-用法)
- [GitHub Actions](#github-actions)
- [项目结构](#项目结构)
- [开发](#开发)
- [许可证](#许可证)

---

## 背景

华为云最佳实践文档（目标仓 [hcbp-demo](https://github.com/Lance52259/hcbp-demo)）需要与 Provider 仓库中的 Terraform examples 保持同步。手工对照成本高、易漏改导航（尤其是 `SUMMARY.md`）。

**doc-craft** 把「探测 → 生成 → 开 PR」做成可重复流水线：本仓库负责编排与 Skill；源仓提供 examples；文档仓接收变更。

| 角色 | 仓库 | 职责 |
|------|------|------|
| **A** | 本仓库（doc-draft / doc-craft） | CLI、Skill、映射、Actions |
| **B** | [huaweicloud/terraform-provider-huaweicloud](https://github.com/huaweicloud/terraform-provider-huaweicloud) | `examples/` 最佳实践源 |
| **C** | [Lance52259/hcbp-demo](https://github.com/Lance52259/hcbp-demo) | 中英文档与 PR 目标 |

---

## 功能特性

- **增量探测**：对比 B 仓 `examples/` 与 C 仓已有文档（含服务别名、连字符/下划线模糊匹配）
- **Open PR 跳过**：C 仓若已存在对应 `doc-craft/...` 分支的 open PR，则跳过该实践
- **中英双语生成**：按 Skill 顺序产出 `docs/zh-cn/` 与 `docs/en-us/` 正文
- **安全导航补丁**：`SUMMARY.md` / `index.md` / `README.md` 仅定点插入；英文侧先按字母序定目录，中文侧跟随；禁止整文件重写
- **一实践一 PR**：提交信息与 PR 标题统一为 `docs({service}): support new best practice for {title}`
- **本地 Dry-run**：不 push、不开真实 PR，便于联调
- **定时 + 手动**：GitHub Actions 支持 schedule 与 `workflow_dispatch`

---

## 工作原理

```text
┌─────────────┐     clone/diff      ┌──────────────────────┐
│ B: examples │ ─────────────────►  │ detect new practices │
└─────────────┘                     └──────────┬───────────┘
                                               │
                                               ▼
                                    ┌──────────────────────┐
                                    │ Skill + DeepSeek     │
                                    │ → zh/en markdown     │
                                    └──────────┬───────────┘
                                               │
                                               ▼
┌─────────────┐     branch + PR     ┌──────────────────────┐
│ C: hcbp     │ ◄─────────────────  | nav patch + gitops   │
└─────────────┘                     └──────────────────────┘
```

1. 拉取 B / C 工作树  
2. 枚举 examples，过滤已对接文档与 open PR  
3. 按 `MAX_PRACTICES` 限流后调用 AI 生成  
4. 编排层补丁中英导航文件  
5. 推送分支并创建 PR（非 dry-run）

---

## 环境要求

- Go **1.22+**
- Git
- DeepSeek API Key（生成阶段）
- 对 C 仓有写权限的 GitHub Token（`Contents` + `Pull requests`；非 dry-run）

---

## 快速开始

```bash
git clone https://github.com/Lance52259/doc-draft.git
cd doc-draft

cp .env.example .env
# 编辑 .env：至少填写 AI_API_KEY；正式推 PR 还需 C_REPO_TOKEN

make build
make detect          # 仅探测，打印 JSON
make dry-run         # 全流程 dry-run（不 push / 不开真实 PR）
```

正式执行（会写 C 仓）：

```bash
make run
# 或
./bin/doc-craft run
./bin/doc-craft run --practice examples/ecs/basic
```

---

## 配置

环境变量优先于 `configs/default_config.yaml`。完整示例见 [`.env.example`](./.env.example)。

### B 仓

| 变量 | 默认 | 说明 |
|------|------|------|
| `B_REPO` | `huaweicloud/terraform-provider-huaweicloud` | 源仓库 |
| `B_REPO_TOKEN` | _(空)_ | 读权限；公开仓可省略 |
| `B_EXAMPLES_PATH` | `examples` | examples 根路径 |
| `B_DEFAULT_BRANCH` | `master` | 分支 |

### C 仓

| 变量 | 默认 | 说明 |
|------|------|------|
| `C_REPO` | `Lance52259/hcbp-demo` | 文档 / PR 目标仓 |
| `C_REPO_TOKEN` | _(空)_ | **写分支 + 开 PR（非 dry-run 必填）** |
| `C_DOCS_ROOT` | `docs/zh-cn/best-practices` | 中文文档根（探测用） |
| `C_DEFAULT_BRANCH` | `master` | PR base |
| `C_SYNCED_MANIFEST` | `synced-practices.json` | 可选已对接清单 |

### AI 与运行时

| 变量 | 默认 | 说明 |
|------|------|------|
| `AI_API_KEY` | — | DeepSeek API Key（生成必填） |
| `AI_BASE_URL` | `https://api.deepseek.com` | API 地址 |
| `AI_MODEL` | `deepseek-chat` | 模型 |
| `DRY_RUN` | `false` | `true` 时不 push / 不开真实 PR |
| `MAX_PRACTICES` | `0`（不限制） | 单次最多处理条数；联调建议 `1` |
| `SKILL_ID` | `best-practice-doc` | Skill 目录名 |

路径映射（B `examples` → C 服务目录 / 文件名）见 [`configs/practice_mapping.yaml`](./configs/practice_mapping.yaml)。  
文档生成规则见 [`skills/best-practice-doc/SKILL.md`](./skills/best-practice-doc/SKILL.md)。

---

## CLI 用法

```bash
doc-craft detect [--out new.json] [--no-refresh]
doc-craft generate [--practice ID] [--practices-file FILE] [--dry-run]
doc-craft run [--practice ID] [--dry-run=true|false] [--no-refresh]
```

| 命令 | 说明 |
|------|------|
| `detect` | 输出 `new_practices` / `synced_ids` / `skipped_open_pr` |
| `generate` | 仅生成（默认偏 dry-run 行为，见 flag） |
| `run` | 探测 → 生成 → 推送 → 开 PR |

---

## GitHub Actions

工作流：[`.github/workflows/monitor-and-generate.yml`](./.github/workflows/monitor-and-generate.yml)

| 触发 | 说明 |
|------|------|
| **schedule** | 每天 **东八区 02:00**（`cron: 0 18 * * *` UTC） |
| **workflow_dispatch** | 可选手动指定 `practice`、勾选 `dry_run` |

Secrets / Variables 建议放在 **Environment `Development`**（工作流已声明 `environment: Development`）：

- 必填：`AI_API_KEY`、`C_REPO_TOKEN`
- 可选：`B_REPO`、`C_REPO`、`MAX_PRACTICES`（Actions 默认常为 `1`）等

单元测试工作流：[`test-workflow.yml`](./.github/workflows/test-workflow.yml)（push / PR）。

---

## 项目结构

```text
doc-draft/
├── cmd/
│   └── doc-craft/              # CLI 入口
├── internal/
│   ├── monitor/                # clone、探测、状态
│   ├── mapping/                # B → C 服务/实践映射
│   ├── ai/                     # Skill、Prompt、生成
│   │   └── provider/           # DeepSeek OpenAI Compatible
│   ├── nav/                    # 中英 SUMMARY / index / README 手术式补丁
│   ├── gitops/                 # commit、push、PR
│   └── config/                 # .env + YAML
├── configs/                    # default_config.yaml、practice_mapping.yaml
├── skills/
│   └── best-practice-doc/      # 文档生成 Skill
├── templates/                  # 正文 / 分类 / PR body 模板
└── .github/
    └── workflows/              # 定时生成 + 单元测试
```

---

## 开发

```bash
make tidy
make test
make build
```

- 测试与源码同目录：`*_test.go`
- 本地联调建议：`MAX_PRACTICES=1` + `DRY_RUN=true`
- 生成规范变更请同步更新 `skills/best-practice-doc/SKILL.md`

---

## 许可证

本项目基于 [MIT License](./LICENSE) 发布。
