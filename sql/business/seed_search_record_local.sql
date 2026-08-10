-- 平台「搜索记录」本地演示数据（中文关键词 + 多渠道用户类型）
SET NAMES utf8mb4;
USE `qixi_crm_business`;

INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live体验用户','DEMO-USER-9101',1,9501,1),
  (9001,'晴空漫游者',NULL,1,0,1),
  (9002,'居家研究员',NULL,1,0,1),
  (9003,'通勤玩家',NULL,1,0,1),
  (9110,'微信逛逛酱','DEMO-USER-9110',1,0,1),
  (9111,'小程序尝鲜客','DEMO-USER-9111',1,0,1),
  (9112,'H5夜猫族','DEMO-USER-9112',1,0,1),
  (9113,'APP冲锋手','DEMO-USER-9113',1,0,1),
  (9114,'PC办公选购','DEMO-USER-9114',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`),
  `group_id`=VALUES(`group_id`),
  `auth_version`=VALUES(`auth_version`);

INSERT INTO `qixi_crm_b_user_profile` (`user_id`,`avatar_url`,`gender`,`bio`,`source_channel`) VALUES
  (9101,'https://via.placeholder.com/64x64.png?text=体验',0,'搜索记录演示','pc'),
  (9001,'https://via.placeholder.com/64x64.png?text=晴空',0,'搜索记录演示','wechat'),
  (9002,'https://via.placeholder.com/64x64.png?text=居家',0,'搜索记录演示','h5'),
  (9003,'https://via.placeholder.com/64x64.png?text=通勤',0,'搜索记录演示','mini_program'),
  (9110,'https://via.placeholder.com/64x64.png?text=微信',0,'搜索记录演示','wechat'),
  (9111,'https://via.placeholder.com/64x64.png?text=小程序',0,'搜索记录演示','mini_program'),
  (9112,'https://via.placeholder.com/64x64.png?text=H5',0,'搜索记录演示','h5'),
  (9113,'https://via.placeholder.com/64x64.png?text=APP',0,'搜索记录演示','ios'),
  (9114,'https://via.placeholder.com/64x64.png?text=PC',0,'搜索记录演示','pc')
ON DUPLICATE KEY UPDATE
  `avatar_url`=VALUES(`avatar_url`),
  `gender`=VALUES(`gender`),
  `bio`=VALUES(`bio`),
  `source_channel`=VALUES(`source_channel`);

INSERT INTO `qixi_crm_b_user_search_record` (`id`,`user_id`,`keyword`,`source`,`created_at`,`deleted_at`) VALUES
  (970001,9101,'夏季亚麻衬衫','pc','2026-08-03 10:20:00',NULL),
  (970002,9101,'居家香薰礼盒','h5','2026-08-03 11:35:00',NULL),
  (970003,9001,'儿童雨靴','mini','2026-08-04 09:10:00',NULL),
  (970004,9110,'防晒冰丝袖套','h5','2026-08-05 14:22:00',NULL),
  (970005,9111,'便携榨汁杯','mini','2026-08-05 15:08:00',NULL),
  (970006,9112,'深夜加班零食','h5','2026-08-06 22:41:00',NULL),
  (970007,9113,'运动速干T恤','h5','2026-08-07 08:15:00',NULL),
  (970008,9114,'人体工学腰靠','pc','2026-08-07 16:50:00',NULL),
  (970009,9002,'无火香薰套装','h5','2026-08-08 10:05:00',NULL),
  (970010,9003,'通勤双肩背包','mini','2026-08-08 18:30:00',NULL)
ON DUPLICATE KEY UPDATE
  `user_id`=VALUES(`user_id`),
  `keyword`=VALUES(`keyword`),
  `source`=VALUES(`source`),
  `created_at`=VALUES(`created_at`),
  `deleted_at`=NULL;
