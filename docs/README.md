# qixi-live-mergers 文档索引

本仓库目标：基于 CRMEB 多商户商城（MER v4.0）的功能基线，重建为可独立演进的多商户商城管理系统。

## 阅读顺序

0. [我们要做什么（产品全景）](./product-understanding.md) — **业务主链路与对象**
1. [功能解析完整度](./analysis-completeness.md) — 完整度与是否允许写代码
2. **[功能点清单（各端→按钮→CRUD）](./features/README.md)** — **验收主文档**
3. **[原项目接口文档](./api/README.md)** — Markdown + **OpenAPI 3.0 YAML**；正确性见 [`api/ACCURACY.md`](./api/ACCURACY.md)、真实功能见 [`api/FUNCTIONAL-TRUTH.md`](./api/FUNCTIONAL-TRUTH.md)
4. **[数据表 `qixi_`](./schema/README.md)** — 表前缀映射与字段
5. [项目总览](./overview.md) — 定位、三角色、建设范围
6. [功能地图（模块级摘要）](./feature-matrix.md) — 功能表 + 脑图摘要
7. [机器抽取统计](./generated/EXTRACT-STATS.md) — 菜单/路由/页面原始统计
8. [CRMEB 参考源码分析](./crmeb-reference.md) — 原系统结构、高风险链路
9. [目标技术架构](./architecture-target.md) — Go / Vben / uni-app x 目标栈
10. [领域模块拆分](./domain-modules.md) — 建议服务与边界
11. [角色与端入口](./roles-and-portals.md) — 平台 / 商户 / 用户端
12. **[服务命名 · Docker 网络/IP](./release/SERVICE-MATRIX.md)** — 各端命名、固定 IP、发布形态
13. [发布命令](./release/COMMANDS.md) · [Pack/Config](./release/PACK-AND-CONFIG.md)
14. **[全端开发计划](./dev-plan-full.md)** — 竖切阶段、各端交付与里程碑（阶段 1 起）

> 功能基线以 `docs/features/` 为准（**已锁定**，见 `analysis-completeness.md`）；表前缀一律 `qixi_`。缺口结案见 `docs/features/08-gaps.md`。

## 原始素材

| 素材 | 路径 |
| --- | --- |
| 前端脑图 | [assets/前端脑图.png](./assets/前端脑图.png) |
| 功能表 1（移动端/用户侧） | [assets/多商户商城系统功能表1.png](./assets/多商户商城系统功能表1.png) |
| 功能表 2（后台管理） | [assets/多商户商城系统功能表2.png](./assets/多商户商城系统功能表2.png) |
| 功能表 3（商城端标准版） | [assets/多商户商城系统功能表3.png](./assets/多商户商城系统功能表3.png) |
| 功能表 4（PC 端） | [assets/多商户商城系统功能表4.png](./assets/多商户商城系统功能表4.png) |
| CRMEB 源码（外部参考，不入库） | `~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0` |

## Agent 约定

- 仓库根目录 [`AGENTS.md`](../AGENTS.md)
- Cursor Skill：[`.cursor/skills/qixi-live-mergers/SKILL.md`](../.cursor/skills/qixi-live-mergers/SKILL.md)
- Codex Skill 镜像：[`codex-skills/qixi-live-mergers/SKILL.md`](../codex-skills/qixi-live-mergers/SKILL.md)

发布/打包相关文档待代码落地后补充到 `docs/release/`（遵循 unified-docker-release 约定）。
