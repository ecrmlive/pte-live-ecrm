-- 系统设置：移除一号通/呼叫系统入口，仅保留服务配置和平台短信验证码配置。
-- 短信 App Key 只通过后台加密保存，禁止在 SQL 中写入明文。

SET NAMES utf8mb4;
USE `qixi_crm_admin`;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (1511,1500,'setting.system.service','服务配置','lucide:settings-2','/service/settings','page',1,1),
  (1512,1500,'setting.system.sms','短信配置','lucide:mail','/setting/sms','page',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=VALUES(`parent_id`),
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

UPDATE `qixi_crm_a_menu`
SET `parent_id`=1512,
    `route_path`='setting/sms',
    `title`='维护平台短信验证码配置',
    `status`=1
WHERE `id`=20936 OR `code`='setting.sms.manage';

DELETE FROM `qixi_crm_a_role_menu` WHERE `menu_id` IN (1501,1510);
DELETE FROM `qixi_crm_a_menu` WHERE `id` IN (1501,1510) OR `code` IN ('setting.serve','setting.serve.login');

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code='platform'
  AND m.id IN (1500,1511,1512,20936);

-- 可公开的腾讯云短信标识用于本地页面验收；App Key 只能由平台管理员在页面中加密保存。
INSERT INTO `qixi_crm_a_cloud_config`
  (`provider`,`config_key`,`ciphertext`,`key_version`,`updated_by`)
VALUES
  ('tencent_sms','enabled','true','bootstrap-public-v1',0),
  ('tencent_sms','sdk_app_id','1401165606','bootstrap-public-v1',0),
  ('tencent_sms','sign_id','711884','bootstrap-public-v1',0),
  ('tencent_sms','sign_content','杭州乐成体育','bootstrap-public-v1',0),
  ('tencent_sms','template_id','2701987','bootstrap-public-v1',0)
ON DUPLICATE KEY UPDATE
  `ciphertext`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`ciphertext`),`ciphertext`),
  `key_version`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`key_version`),`key_version`),
  `updated_by`=IF(`config_key`='enabled' AND `key_version`='bootstrap-public-v1',VALUES(`updated_by`),`updated_by`);
