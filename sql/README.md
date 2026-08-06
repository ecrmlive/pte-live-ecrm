# 七禧 CRM 初始化 SQL

旧 `qixi_m_*` 与编号碎片迁移文件（`01_table` / `ALTER` 补丁等）已废弃，禁止执行。

开发阶段以 **drop + recreate** 为准：改表只改 `init_table.sql` 里的 `CREATE TABLE`，不要再新增「加列」SQL 文件。

| 数据库 | 顺序执行文件（仓库） |
| --- | --- |
| `qixi_crm_admin` | `admin/init_table.sql` → `init_config.sql` → `init_data.sql` → `init_key.sql` → `init_test_data.sql` |
| `qixi_crm_business` | `business/init_table.sql` → `init_config.sql` → `init_data.sql` → `init_key.sql` → `init_test_data.sql` |
| `qixi_crm_merchant` | `merchant/init_table.sql` → `init_config.sql` → `init_data.sql` → `init_key.sql` → `init_test_data.sql` → **`init_menu_crmeb_full.sql`**（CRMEB `is_mer=1` 店铺菜单全量 693 条） |

## 密钥文件

| 文件 | Git |
| --- | --- |
| `*/init_key.sql.example` | ✅ 提交（仅结构/说明，无真实密钥） |
| `*/init_key.sql` | ❌ 本地/受控环境；`make local-db-init` 导入前从 example 复制并填入 |

```bash
cp sql/admin/init_key.sql.example sql/admin/init_key.sql
cp sql/business/init_key.sql.example sql/business/init_key.sql
cp sql/merchant/init_key.sql.example sql/merchant/init_key.sql
```

所有脚本使用 utf8mb4、可重复的 `CREATE TABLE IF NOT EXISTS` / upsert。真实 JWT、支付、云服务与 IM 凭证只能写入运行机 `app.yaml` 或后台加密配置 / 本地 `init_key.sql`。

## 完整 DDL 范围

| 库 | 必须覆盖的完整域 |
| --- | --- |
| `qixi_crm_admin` | 后台账号/RBAC、角色菜单按钮、区域与商圈、商户审核、运营/内容/DIY/素材、全局商品治理、平台营销规则、后台通知、云服务配置、审计与导出 |
| `qixi_crm_business` | **PC、小程序&H5、iOS、Android、鸿蒙共同使用**：C 用户多端身份、地址、商品消费视图、收藏/足迹/搜索、购物车、订单/子单/明细/发票/物流、支付/回调/分账、退款/介入、券/秒杀/拼团/预售/助力、积分/会员/分销/佣金、充值/提现、社区/直播/预约、客服业务关联与业务事件 |
| `qixi_crm_merchant` | 商户/店铺/员工/RBAC、类目/品牌/参数/标签/保障、SPU/SKU/库存/回收站、店铺券与全部店铺活动、运费/物流/发货点/配送/打印、订单履约/核销/售后操作、店铺账本/提现/发票/分账视图、素材/装修/表单、操作日志与业务 outbox |

“平台、商户、区域、客服、运营”对业务数据的读写并不意味着数据进入 admin 库：订单、用户、支付、售后、营销、内容、客服业务关联等统一按业务归属进入 `qixi_crm_business`；admin 库只保存后台身份、权限、数据范围、后台配置和审计。

IM 必须直接使用 pte-live-im `sql/init_im_schema.sql` 的实际表规则、鉴权规则和 SDK；严禁七禧自建另一套 IM 表或规则。pte-live-im 的源码与初始化 SQL 由其仓库维护，七禧本仓不得修改。
