-- 本地演示：转账记录 + 店铺余额（utf8mb4 中文）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 确保演示店铺有余额 / 冻结金额
UPDATE `qixi_crm_a_merchant_view`
SET `mer_money`=12880.50, `freeze_money`=320.00
WHERE `merchant_id`=1 AND EXISTS (
  SELECT 1 FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='mer_money'
);

UPDATE `qixi_crm_a_merchant_view`
SET `mer_money`=5620.00, `freeze_money`=0.00
WHERE `merchant_id`=2 AND EXISTS (
  SELECT 1 FROM information_schema.COLUMNS
  WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_merchant_view' AND COLUMN_NAME='mer_money'
);

INSERT INTO `qixi_crm_a_financial`
  (`financial_id`,`financial_sn`,`mer_money`,`extract_money`,`financial_type`,`financial_account`,`financial_status`,`status`,`refusal`,`mer_id`,`image`,`admin_id`,`create_time`,`status_time`,`update_time`,`is_del`,`mark`,`admin_mark`,`mer_admin_id`,`type`)
VALUES
  (9900101,'F202608100001',12880.50,2000.00,1,
   '{"name":"张演示","bank":"七禧演示银行","bank_code":"6222020000000001"}',
   0,0,'',1,'',NULL,'2026-08-08 10:20:00',NULL,NULL,0,'商户申请提现演示','',1,0),
  (9900102,'F202608100002',5620.00,1500.00,3,
   '{"name":"李验收","alipay":"demo_li@example.com","alipay_code":""}',
   0,1,'',2,'',1,'2026-08-07 14:05:00','2026-08-07 16:30:00',NULL,0,'支付宝收款演示','平台已审核待打款',1,0),
  (9900103,'F202608100003',12880.50,500.00,2,
   '{"name":"张演示","wechat":"wx_demo_zhang","wechat_code":""}',
   1,1,'',1,'https://example.com/fixture/transfer-voucher-9900103.png',1,'2026-08-05 09:12:00','2026-08-05 11:00:00','2026-08-05 15:40:00',0,'微信收款已打款','凭证已登记',1,0),
  (9900104,'F202608100004',5620.00,800.00,1,
   '{"name":"李验收","bank":"演示工商银行","bank_code":"6222020000000002"}',
   0,-1,'账户信息不符，请核对后重新申请',2,'',1,'2026-08-04 18:22:00','2026-08-04 19:10:00',NULL,0,'银行卡申请','审核拒绝演示',1,0),
  (9900105,'F202608100005',12880.50,3000.00,1,
   '{"name":"张演示","bank":"七禧演示银行","bank_code":"6222020000000001"}',
   0,1,'',1,'',1,'2026-08-09 11:45:00','2026-08-09 13:00:00',NULL,0,'大额提现待打款','',1,0)
ON DUPLICATE KEY UPDATE
  `financial_sn`=VALUES(`financial_sn`),
  `mer_money`=VALUES(`mer_money`),
  `extract_money`=VALUES(`extract_money`),
  `financial_type`=VALUES(`financial_type`),
  `financial_account`=VALUES(`financial_account`),
  `financial_status`=VALUES(`financial_status`),
  `status`=VALUES(`status`),
  `refusal`=VALUES(`refusal`),
  `mer_id`=VALUES(`mer_id`),
  `image`=VALUES(`image`),
  `admin_id`=VALUES(`admin_id`),
  `create_time`=VALUES(`create_time`),
  `status_time`=VALUES(`status_time`),
  `update_time`=VALUES(`update_time`),
  `is_del`=VALUES(`is_del`),
  `mark`=VALUES(`mark`),
  `admin_mark`=VALUES(`admin_mark`),
  `type`=VALUES(`type`);
