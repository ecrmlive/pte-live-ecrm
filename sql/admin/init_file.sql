-- 系统默认素材引用。资源由运营上传至 COS 后再执行；本文件不包含 COS 凭据或真实地址。
SET NAMES utf8mb4;
USE `qixi_crm_admin`;

-- 仅接受本机 init_key.sql 的 bootstrap COS 地址，避免把运行时加密字段当成 URL 写入素材库。
SET @qixi_system_asset_base_url := (
  SELECT TRIM(TRAILING '/' FROM `ciphertext`)
  FROM `qixi_crm_a_cloud_config`
  WHERE `provider`='cos' AND `config_key`='base_url'
    AND `key_version`='bootstrap-local-v1' AND `ciphertext` <> ''
  LIMIT 1
);

-- 约定上传路径：system/store、system/service、system/product、system/member、system/common。
-- 未配置 COS 时不写入，避免产生不可用的默认素材；配置后重复导入可安全补齐。
INSERT INTO `qixi_crm_a_attachment_asset`
  (`attachment_category_id`,`attachment_name`,`attachment_src`,`upload_type`,`user_type`,`user_id`,`attachment_type`,`is_system`)
SELECT `seed`.`category_id`, `seed`.`name`, CONCAT(@qixi_system_asset_base_url, '/', `seed`.`object_key`), 1, 0, 0, 0, 1
FROM (
  SELECT 5101 AS `category_id`, '默认店铺封面' AS `name`, 'system/store/default-cover.png' AS `object_key`
  UNION ALL SELECT 5104, '默认客服头像', 'system/service/default-avatar.png'
  UNION ALL SELECT 5105, '默认商品图片', 'system/product/default-cover.png'
  UNION ALL SELECT 5107, '默认会员头像', 'system/member/default-avatar.png'
  UNION ALL SELECT 5107, '默认用户头像', 'system/member/default-user-avatar.png'
  UNION ALL SELECT 5106, '商城默认背景', 'system/common/default-background.png'
  UNION ALL SELECT 5102, '默认支付图标', 'system/common/default-payment.png'
  UNION ALL SELECT 5103, '默认物流图标', 'system/common/default-logistics.png'
  UNION ALL SELECT 5108, '平台默认标识', 'system/common/default-logo.png'
) AS `seed`
WHERE COALESCE(@qixi_system_asset_base_url, '') <> ''
  AND NOT EXISTS (
    SELECT 1 FROM `qixi_crm_a_attachment_asset` AS `asset`
    WHERE `asset`.`attachment_category_id`=`seed`.`category_id`
      AND `asset`.`attachment_name`=`seed`.`name`
      AND `asset`.`is_system`=1
  );
