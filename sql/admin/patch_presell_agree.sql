-- 平台「预售协议」协议键夹具（CRMEB CacheRepository::PRESELL_AGREE = sys_product_presell_agree）
-- 用法：make local-sync-sql
SET NAMES utf8mb4;

USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_product_presell_agree',0,'<p>1. 预售商品以页面公示的发货时间为准；</p><p>2. 定金支付后请在尾款支付期限内完成支付，逾期超时订单将按活动规则处理；</p><p>3. 最终解释权归平台所有。</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2109,'agreement','sys_product_presell_agree','','<p>1. 预售商品以页面公示的发货时间为准；</p><p>2. 定金支付后请在尾款支付期限内完成支付，逾期超时订单将按活动规则处理；</p><p>3. 最终解释权归平台所有。</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
