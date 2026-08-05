# H5 商户入驻申请审核闭环验收记录

日期：2026-08-04

## 闭环范围

- 用户端提交申请后，业务库以事务 outbox 发布 `business.merchant_application.submitted`；平台端继续只保存监管投影。
- 平台审核通过或驳回时，审核状态和说明与平台审核记录在同一事务内写入 `qixi_crm_a_outbox`。
- 平台 dispatcher 发布 `platform.merchant_application.reviewed`；业务端 NATS 订阅只更新同一申请 ID 且仍为 `pending` 的业务记录，重复消息不会覆盖已终态结果。
- 用户端 `GET /merchant-applications/mine/:id` 通过当前用户范围读取申请详情，不返回营业执照对象键或存储 URL。
- H5 入驻页的申请记录可进入详情页，展示审核状态、申请资料、执照提交状态与驳回说明。

## 自动验证

| 检查 | 结果 |
| --- | --- |
| `api-platform go test ./...` | 通过 |
| `api-business go test ./...` | 通过 |
| 平台/业务端审核事件单元测试 | 通过；无效、非终态或未关联申请的消息在访问数据库前被忽略 |
| H5 构建 | 退出成功，新增页面无 UTS 警告，产物存在 |

## 待受控环境验收

- 需要在真实 MySQL + NATS 环境执行“提交 → 平台驳回/通过 → H5 刷新详情”的端到端验证，确认 outbox 失败重投和重复消息幂等。
- 营业执照对象上传与下载授权仍依赖受控对象存储配置，H5 详情不直接暴露对象键或 URL。
- 本机 DCloud CLI 的 `TS6059 tslib rootDir` 是既有工具链诊断，非本页业务诊断；关闭前不能标为发布级类型验收。
