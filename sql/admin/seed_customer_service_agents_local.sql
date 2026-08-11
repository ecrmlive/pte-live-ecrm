USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 客服列表本地演示数据：与 C 端用户以 linked_user_id 真实关联。
-- password_hash 使用无效标记，账号不能据此登录；本夹具不写入任何密码或 IM 凭据。
INSERT INTO `qixi_crm_a_admin_user`
  (`id`,`username`,`password_hash`,`display_name`,`linked_user_id`,`avatar_url`,`phone`,`status`,`auth_version`,`data_scope_version`,`deleted_at`) VALUES
  (981001,'fixture_cs_agent_01','!local-fixture-login-disabled!','客服小林',991001,'/demo/avatar-9101.png','CS-DUTY-001',1,1,1,NULL),
  (981002,'fixture_cs_agent_02','!local-fixture-login-disabled!','客服阿澈',991002,'/demo/avatar-9102.png','CS-DUTY-002',1,1,1,NULL),
  (981003,'fixture_cs_agent_03','!local-fixture-login-disabled!','售后客服小夏',991003,'/demo/avatar-9104.png','CS-DUTY-003',1,1,1,NULL),
  (981004,'fixture_cs_agent_04','!local-fixture-login-disabled!','商品顾问小陆',991004,'/demo/avatar-9105.png','CS-DUTY-004',1,1,1,NULL),
  (981005,'fixture_cs_agent_05','!local-fixture-login-disabled!','订单客服小顾',991005,'','CS-DUTY-005',1,1,1,NULL),
  (981006,'fixture_cs_agent_06','!local-fixture-login-disabled!','VIP客服小宁',991006,'','CS-DUTY-006',1,1,1,NULL),
  (981007,'fixture_cs_agent_07','!local-fixture-login-disabled!','物流客服小周',991007,'','CS-DUTY-007',0,1,1,NULL),
  (981008,'fixture_cs_agent_08','!local-fixture-login-disabled!','夜间客服小麦',991008,'','CS-DUTY-008',0,1,1,NULL)
ON DUPLICATE KEY UPDATE
  `display_name`=VALUES(`display_name`),
  `linked_user_id`=VALUES(`linked_user_id`),
  `avatar_url`=VALUES(`avatar_url`),
  `phone`=VALUES(`phone`),
  `status`=VALUES(`status`),
  `deleted_at`=NULL;

-- 仅授予客服角色，平台账号仍可按现有 RBAC 全量管理这些坐席。
INSERT INTO `qixi_crm_a_admin_user_role` (`admin_user_id`,`role_id`)
SELECT fixture.admin_user_id, role.id
FROM (
  SELECT 981001 AS admin_user_id UNION ALL SELECT 981002 UNION ALL SELECT 981003 UNION ALL SELECT 981004
  UNION ALL SELECT 981005 UNION ALL SELECT 981006 UNION ALL SELECT 981007 UNION ALL SELECT 981008
) AS fixture
JOIN `qixi_crm_a_role` AS role ON role.`code`='customer_service'
ON DUPLICATE KEY UPDATE `role_id`=VALUES(`role_id`);

-- 为每名演示客服配置服务队列，确保客服身份登录时也能取得本人的可服务店铺范围。
INSERT INTO `qixi_crm_a_data_scope` (`admin_user_id`,`scope_type`,`scope_value`,`version`)
SELECT fixture.admin_user_id, 'service_queue', JSON_ARRAY(fixture.store_id), 1
FROM (
  SELECT 981001 AS admin_user_id, 1 AS store_id UNION ALL SELECT 981002, 1 UNION ALL SELECT 981003, 2 UNION ALL SELECT 981004, 2
  UNION ALL SELECT 981005, 3 UNION ALL SELECT 981006, 3 UNION ALL SELECT 981007, 1 UNION ALL SELECT 981008, 2
) AS fixture
WHERE NOT EXISTS (
  SELECT 1 FROM `qixi_crm_a_data_scope` AS scope
  WHERE scope.`admin_user_id`=fixture.admin_user_id AND scope.`scope_type`='service_queue'
);
