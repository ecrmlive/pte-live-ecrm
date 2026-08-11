USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 客服列表关联用户的本地演示数据：仅用于本地验收，不含真实手机号或登录凭据。
-- 客服账号在统一后台库维护；此处只提供可被选择和展示的 C 端用户资料。
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (991001,'微信用户U991001','CS-DEMO-USER-001',1,0,1),
  (991002,'微信用户U991002','CS-DEMO-USER-002',1,0,1),
  (991003,'微信用户U991003','CS-DEMO-USER-003',1,0,1),
  (991004,'微信用户U991004','CS-DEMO-USER-004',1,0,1),
  (991005,'微信用户U991005','CS-DEMO-USER-005',1,0,1),
  (991006,'微信用户U991006','CS-DEMO-USER-006',1,0,1),
  (991007,'微信用户U991007','CS-DEMO-USER-007',1,0,1),
  (991008,'微信用户U991008','CS-DEMO-USER-008',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`),
  `auth_version`=VALUES(`auth_version`);

-- 使用本地演示静态路径，头像为空时前端仍以默认头像回退展示。
INSERT INTO `qixi_crm_b_user_profile` (`user_id`,`avatar_url`,`real_name`,`gender`,`bio`,`source_channel`) VALUES
  (991001,'/demo/avatar-9101.png','微信用户U991001',0,'','wechat'),
  (991002,'/demo/avatar-9102.png','微信用户U991002',0,'','mini_program'),
  (991003,'/demo/avatar-9104.png','微信用户U991003',0,'','wechat'),
  (991004,'/demo/avatar-9105.png','微信用户U991004',0,'','mini_program'),
  (991005,'','微信用户U991005',0,'','wechat'),
  (991006,'','微信用户U991006',0,'','h5'),
  (991007,'','微信用户U991007',0,'','mini_program'),
  (991008,'','微信用户U991008',0,'','wechat')
ON DUPLICATE KEY UPDATE
  `avatar_url`=VALUES(`avatar_url`),
  `real_name`=VALUES(`real_name`),
  `gender`=VALUES(`gender`),
  `bio`=VALUES(`bio`),
  `source_channel`=VALUES(`source_channel`);
