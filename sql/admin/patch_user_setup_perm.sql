-- 平台「用户设置」按钮权限 + 默认配置（幂等）
-- RequireAdminMenu 仅认 kind=button；page 码 user.setup 不能用于接口鉴权。
-- 用法（项目根目录）：
--   docker exec -i pte_live_mysql sh -ec 'MYSQL_PWD="$MYSQL_ROOT_PASSWORD" mysql -uroot --default-character-set=utf8mb4' < sql/admin/patch_user_setup_perm.sql

SET NAMES utf8mb4;

USE `qixi_crm_admin`;

-- 确保导航页存在且启用
INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (532,70,'user.setup','用户设置','lucide:settings-2','/user/setup_user','page',7,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=70,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `icon`=VALUES(`icon`),
  `route_path`=VALUES(`route_path`),
  `kind`='page',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT INTO `qixi_crm_a_menu` (`id`,`parent_id`,`code`,`title`,`icon`,`route_path`,`kind`,`sort`,`status`) VALUES
  (21072,532,'user.setup.read','查看用户设置','','user/setup_user','button',1,1),
  (20998,532,'user.setup.manage','保存用户设置','','user/setup_user','button',2,1)
ON DUPLICATE KEY UPDATE
  `parent_id`=532,
  `code`=VALUES(`code`),
  `title`=VALUES(`title`),
  `route_path`=VALUES(`route_path`),
  `kind`='button',
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`,`menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform','operations')
  AND m.code IN (
    'user.setup',
    'user.setup.read',
    'user.setup.manage'
  );

-- 旧版简版配置 / 缺 mark 字段时回写默认 JSON（API Get 也会注入缺失默认字段）
INSERT INTO `qixi_crm_a_setting_cache` (`key`,`expire_time`,`result`) VALUES
  ('user_setup_config',0,'{"user_default_avatar":"","fields":[{"id":1,"field":"real_name","title":"姓名","is_used":1,"is_require":0,"is_show":1,"type":"input","msg":"请填写真实姓名","is_default":1,"sort":0},{"id":2,"field":"sex","title":"性别","is_used":1,"is_require":0,"is_show":1,"type":"radio","msg":"请选择性别","content":["男","女","保密"],"is_default":1,"sort":1},{"id":3,"field":"birthday","title":"生日","is_used":1,"is_require":0,"is_show":1,"type":"date","msg":"请选择生日","is_default":1,"sort":2},{"id":4,"field":"address","title":"地址","is_used":1,"is_require":0,"is_show":1,"type":"address","msg":"请选择地址","is_default":1,"sort":3},{"id":5,"field":"mark","title":"备注","is_used":1,"is_require":0,"is_show":0,"type":"input","msg":"请填写备注","is_default":1,"sort":4},{"id":6,"field":"id_card","title":"身份证（实名认证）","is_used":0,"is_require":0,"is_show":0,"type":"id_card","msg":"请填写身份证号","is_default":1,"sort":5}],"is_phone_login":0,"first_avatar_switch":1,"open_update_info":1,"wechat_phone_switch":0,"newcomer_status":0,"register_popup_pic":"","register_money_status":0,"register_give_money":0,"register_integral_status":0,"register_give_integral":0,"register_coupon_status":0,"register_give_coupon":[]}')
ON DUPLICATE KEY UPDATE
  `expire_time`=VALUES(`expire_time`),
  `result`=IF(
    `result` IS NULL
      OR `result`=''
      OR `result` LIKE '%"register_enabled"%'
      OR `result` NOT LIKE '%"fields"%'
      OR `result` NOT LIKE '%"mark"%',
    VALUES(`result`),
    `result`
  );
