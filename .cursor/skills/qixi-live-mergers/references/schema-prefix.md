# 表前缀 `qixi_`

- 参考库：CRMEB `eb_*`（外部 SQL，不入库运行）
- 本仓库：一律 `qixi_*`
- 映射表：`docs/schema/table-prefix-map.tsv`
- 参考 DDL：`docs/schema/qixi_schema_reference.sql`
- 字段说明：`docs/schema/columns.tsv`、`docs/schema/domain-*.md`

示例：

```text
eb_store_order      → qixi_store_order
eb_merchant         → qixi_merchant
eb_store_product    → qixi_store_product
eb_user             → qixi_user
```

GORM：`TableName()` 返回 `qixi_...`。迁移与文档示例禁止再写 `eb_`。
