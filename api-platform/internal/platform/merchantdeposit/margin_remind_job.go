package merchantdeposit

// TODO(margin-remind-job): 对齐 CRMEB SyncMerchantMarginStatusListen（约每 30 分钟）
// + MerchantRepository.changeMerchantStatus：
//
//  1. 读取 qixi_crm_a_setting_cache.key=margin_remind_config
//     （margin_remind_switch / margin_remind_day）
//  2. 开关关闭则直接返回
//  3. 扫描 qixi_crm_a_merchant_deposit_account.state=shortfall（及 pending 需补缴）商户
//  4. 首次不足时写入提醒截止时间（CRMEB: merchant.margin_remind_time），并发送补缴通知
//  5. 截止时间到期仍不足则关闭店铺经营状态（merchant_view.status=0）
//
// 配置 CRUD 已由 GET/PUT /setting/margin 完成；本文件仅作定时任务接入锚点，
// 完整关店与通知实现需补账户提醒截止字段、通知通道与 job 对 admin 库的访问。
