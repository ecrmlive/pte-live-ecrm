# 七禧 CRM 初始化 SQL

旧 `qixi_m_*` 与 `qixi_mergers` 单库初始化文件已删除，禁止执行。

| 数据库 | 顺序执行文件 |
| --- | --- |
| `qixi_crm_admin` | `admin/01_table.sql` → `02_data.sql` → `03_config.sql` → `04_key.sql` → `05_test_data.sql` |
| `qixi_crm_business` | `business/01_table.sql` → `02_data.sql` → `03_config.sql` → `04_key.sql` → `05_test_data.sql` |
| `qixi_crm_merchant` | `merchant/01_table.sql` → `02_data.sql` → `03_config.sql` → `04_key.sql` → `05_test_data.sql` |

所有脚本使用 utf8mb4、可重复执行的 `IF NOT EXISTS` / upsert 语义。密钥文件只保留结构和空值，真实 JWT、支付、云服务与 IM 凭证只能写入运行机 `app.yaml` 或后台加密配置。

## 完整 DDL 范围（未补齐前不得进入 Make/代码迁移）

当前 `01_table.sql` 仅为三库基础骨架，**不是完整功能 DDL**。必须补齐以下域后才能宣称 SQL 重构完成：

| 库 | 必须覆盖的完整域 |
| --- | --- |
| `qixi_crm_admin` | 后台账号/RBAC、角色菜单按钮、区域与商圈、商户审核、运营/内容/DIY/素材、全局商品治理、平台营销规则、后台通知、云服务配置、审计与导出 |
| `qixi_crm_business` | **PC、小程序&H5、iOS、Android、鸿蒙共同使用**：C 用户多端身份、地址、商品消费视图、收藏/足迹/搜索、购物车、订单/子单/明细/发票/物流、支付/回调/分账、退款/介入、券/秒杀/拼团/预售/助力、积分/会员/分销/佣金、充值/提现、社区/直播/预约、客服业务关联与业务事件 |
| `qixi_crm_merchant` | 商户/店铺/员工/RBAC、类目/品牌/参数/标签/保障、SPU/SKU/库存/回收站、店铺券与全部店铺活动、运费/物流/发货点/配送/打印、订单履约/核销/售后操作、店铺账本/提现/发票/分账视图、素材/装修/表单、操作日志与业务 outbox |

“平台、商户、区域、客服、运营”对业务数据的读写并不意味着数据进入 admin 库：订单、用户、支付、售后、营销、内容、客服业务关联等统一按业务归属进入 `qixi_crm_business`；admin 库只保存后台身份、权限、数据范围、后台配置和审计。

IM 必须直接使用 pte-live-im `sql/init_im_schema.sql` 的实际表规则、鉴权规则和 SDK；严禁七禧自建另一套 IM 表或规则。pte-live-im 的源码与初始化 SQL 由其仓库维护，七禧本仓不得修改。
