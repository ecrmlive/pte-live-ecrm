SET NAMES utf8mb4;
-- 平台监管投影里已有、但店铺库尚未开通的商户：补齐 merchant/store/owner，
-- 供平台后台「店铺列表 → 登录」签发 store_console JWT。密码演示值均为 123456m。
-- 仅处理 merchant_id >= 1000 的 UI/演示夹具，不覆盖基线店铺 1–3。

INSERT INTO `qixi_crm_merchant`.`qixi_crm_m_merchant` (`id`,`name`,`status`,`region_id`)
SELECT v.`merchant_id`, v.`merchant_name`, IF(v.`status` = 1, 1, 0), NULLIF(v.`region_id`, 0)
FROM `qixi_crm_admin`.`qixi_crm_a_merchant_view` AS v
WHERE v.`merchant_id` >= 1000
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `status`=VALUES(`status`),
  `region_id`=VALUES(`region_id`);

INSERT INTO `qixi_crm_merchant`.`qixi_crm_m_store` (`id`,`merchant_id`,`app_id`,`name`,`status`)
SELECT
  v.`merchant_id`,
  v.`merchant_id`,
  CONCAT('qixi.store.demo.', v.`merchant_id`),
  v.`merchant_name`,
  IF(v.`status` = 1, 1, 0)
FROM `qixi_crm_admin`.`qixi_crm_a_merchant_view` AS v
WHERE v.`merchant_id` >= 1000
ON DUPLICATE KEY UPDATE
  `merchant_id`=VALUES(`merchant_id`),
  `app_id`=VALUES(`app_id`),
  `name`=VALUES(`name`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_merchant`.`qixi_crm_m_account`
  (`store_id`,`username`,`password_hash`,`role_code`,`display_name`,`phone`,`status`,`auth_version`)
SELECT
  v.`merchant_id`,
  IF(NULLIF(TRIM(v.`mer_account`), '') IS NULL,
     CONCAT('ui_mock_', LPAD(v.`merchant_id`, 4, '0')),
     TRIM(v.`mer_account`)),
  '$2a$10$7e1OmptO8l5P3lJ7ziIfOeC0GXY0MGUNNY/QS6LQKgXLNq2Z6TFJe',
  'owner',
  IF(NULLIF(TRIM(v.`contact_name`), '') IS NULL, '演示店长', TRIM(v.`contact_name`)),
  IFNULL(v.`contact_mobile`, ''),
  1,
  1
FROM `qixi_crm_admin`.`qixi_crm_a_merchant_view` AS v
WHERE v.`merchant_id` >= 1000
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_merchant`.`qixi_crm_m_account` AS a
    WHERE a.`store_id` = v.`merchant_id` AND a.`role_code` = 'owner'
  );
