# 目标技术架构

系统边界以 [SYSTEM-ARCHITECTURE.md](./SYSTEM-ARCHITECTURE.md) 为唯一口径。本文件只说明实现技术，不保留旧目录、旧库或旧服务拓扑。

## 固定技术

- Go 1.26.5、Gin、GORM、Swagger/OpenAPI、MySQL 8.4.10、Redis 8.8.0、etcd、NATS。
- 统一后台与店铺后台均使用 Vben 5.7+。
- H5 与小程序使用同一套 uni-app x（UTS / HBuilderX 5.0+）；PC 使用 Vue 3 + Vite + TypeScript。
- 所有运行配置使用 `.yaml`。
- IM 只使用 pte-live-im 与 pte-live-im-sdk。

## 应用边界

| 应用 | 前端/服务 | 数据所有权 |
| --- | --- | --- |
| 统一后台管理系统 | `admin-platform/` + `api-platform` | `qixi_crm_admin` |
| 店铺管理系统 | `admin-merchant/` + `api-merchant` | `qixi_crm_merchant` |
| PC / H5 / 小程序 / iOS / Android / 鸿蒙用户端 | `app-pc/`、`app-uni/`、`app-ios/`、`app-adnroid/`、`app-harmony/` + `api-business` | `qixi_crm_business` |
| 异步任务 | `job` | 事件消费、对账、通知；不越过各库所有权 |

统一后台按平台、商户、区域、客服、运营角色动态返回菜单；店铺后台独立部署、独立 JWT 与独立数据库边界。PC、小程序、H5、iOS、Android、鸿蒙六端共享用户 JWT 配置和用户主体，功能闭环必须对齐 H5。

## 数据与跨域

- 表前缀只能是 `qixi_crm_a_`、`qixi_crm_b_`、`qixi_crm_m_`、`pte_im_a_`、`pte_im_`。
- 服务只写自己所有的 qixi CRM 数据库；跨域调用使用 API / NATS 事件，跨库只保存稳定 ID，禁止跨库外键和事务。
- 支付、退款、库存、优惠、积分、佣金、结算在各自所有库中以幂等键和事务实现。

## 网络与发布目标

qixi 使用独立 Compose project `qixi_mergers`，其 API 加入 pte-live-im 已维护的 `pte_live_net` 并复用 MySQL、Redis、NATS、etcd。七禧不启动重复基础设施，且仅使用自己的 `qixi_crm_admin`、`qixi_crm_business`、`qixi_crm_merchant` 数据库，不得读写 pte-live-im 的数据库或表。前端本地开发直连 `127.0.0.1`，不安装本机 Nginx。构建、Make 和 Compose 具体命令在第 3 步目录迁移完成后重写。
