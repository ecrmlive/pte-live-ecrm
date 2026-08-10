-- 平台优惠券「使用说明」（CRMEB CacheRepository::COUPON_AGREE = sys_coupon_agree）
-- 用法：make local-sync-sql
SET NAMES utf8mb4;

USE `qixi_crm_admin`;
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('sys_coupon_agree',0,'<p>1. 优惠券领取后请在有效期内使用；</p><p>2. 每张优惠券限使用一次，不可叠加；</p><p>3. 最终解释权归平台所有。</p>')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(`result` IS NULL OR `result`='',VALUES(`result`),`result`);

USE `qixi_crm_business`;
INSERT INTO `qixi_crm_b_content_view` (`content_id`,`content_type`,`title`,`cover_url`,`body`,`status`,`version`,`published_at`,`updated_at`) VALUES
  (2108,'agreement','sys_coupon_agree','','<p>1. 优惠券领取后请在有效期内使用；</p><p>2. 每张优惠券限使用一次，不可叠加；</p><p>3. 最终解释权归平台所有。</p>',1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE
  `content_type`=VALUES(`content_type`),
  `title`=VALUES(`title`),
  `body`=IF(`body` IS NULL OR `body`='',VALUES(`body`),`body`),
  `status`=1,
  `updated_at`=NOW();
