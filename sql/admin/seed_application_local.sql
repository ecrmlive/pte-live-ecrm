-- 本地演示：系统表单 + 报名活动 + 报名记录（幂等，utf8mb4 中文）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 演示系统表单（参会报名）
INSERT INTO `qixi_crm_a_config_item`
  (`id`,`item_type`,`name`,`code`,`remark`,`payload`,`status`,`sort`,`is_del`)
VALUES
  (8105,'system_form','新品内测报名表','signup-early-access','报名活动演示表单',
   JSON_OBJECT(
     'fields', JSON_ARRAY(
       JSON_OBJECT('key','upload_image','label','上传图片','type','image'),
       JSON_OBJECT('key','attendee_name','label','参会姓名','type','text'),
       JSON_OBJECT('key','contact_phone','label','联系电话','type','text')
     )
   ),1,25,0),
  (8106,'system_form','线下沙龙报名表','signup-salon','报名活动演示表单',
   JSON_OBJECT(
     'fields', JSON_ARRAY(
       JSON_OBJECT('key','upload_image','label','上传图片','type','image'),
       JSON_OBJECT('key','attendee_name','label','参会姓名','type','text'),
       JSON_OBJECT('key','contact_phone','label','联系电话','type','text'),
       JSON_OBJECT('key','company','label','所在公司','type','text')
     )
   ),1,26,0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `code`=VALUES(`code`),
  `remark`=VALUES(`remark`),
  `payload`=VALUES(`payload`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `is_del`=0;

-- 报名活动
INSERT INTO `qixi_crm_a_signup_activity`
  (`id`,`name`,`info`,`cover_url`,`poster_url`,`color`,`form_id`,`quota`,`total`,`status`,`sort`,`starts_at`,`ends_at`,`is_del`)
VALUES
  (8201,'新品内测报名','欢迎参与七禧新品内测，填写资料即可报名。',
   'https://picsum.photos/seed/qixi-signup-cover-1/750/350',
   'https://picsum.photos/seed/qixi-signup-poster-1/750/1250',
   '#E8F5E9',8105,100,2,1,10,
   DATE_SUB(NOW(),INTERVAL 3 DAY),DATE_ADD(NOW(),INTERVAL 60 DAY),0),
  (8202,'线下香氛沙龙','线下体验活动，名额有限，请如实填写联系方式。',
   'https://picsum.photos/seed/qixi-signup-cover-2/750/350',
   'https://picsum.photos/seed/qixi-signup-poster-2/750/1250',
   '#FFF3E0',8106,50,1,1,20,
   DATE_ADD(NOW(),INTERVAL 7 DAY),DATE_ADD(NOW(),INTERVAL 90 DAY),0),
  (8203,'春季会员开放日','已结束的演示活动。',
   'https://picsum.photos/seed/qixi-signup-cover-3/750/350',
   'https://picsum.photos/seed/qixi-signup-poster-3/750/1250',
   '#E3F2FD',8105,0,1,0,30,
   DATE_SUB(NOW(),INTERVAL 90 DAY),DATE_SUB(NOW(),INTERVAL 10 DAY),0)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `info`=VALUES(`info`),
  `cover_url`=VALUES(`cover_url`),
  `poster_url`=VALUES(`poster_url`),
  `color`=VALUES(`color`),
  `form_id`=VALUES(`form_id`),
  `quota`=VALUES(`quota`),
  `total`=VALUES(`total`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `starts_at`=VALUES(`starts_at`),
  `ends_at`=VALUES(`ends_at`),
  `is_del`=0;

-- 报名记录
INSERT INTO `qixi_crm_a_signup_record`
  (`id`,`activity_id`,`user_id`,`nickname`,`mobile`,`avatar`,`form_value`,`is_del`,`created_at`)
VALUES
  (8301,8201,10001,'演示用户甲','13800001111','https://picsum.photos/seed/qixi-user-a/80/80',
   JSON_OBJECT(
     'upload_image','https://picsum.photos/seed/qixi-signup-form-img-1/200/200',
     'attendee_name','王小明',
     'contact_phone','13800001111'
   ),0,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (8302,8201,10002,'演示用户乙','13800002222','https://picsum.photos/seed/qixi-user-b/80/80',
   JSON_OBJECT(
     'upload_image','https://picsum.photos/seed/qixi-signup-form-img-2/200/200',
     'attendee_name','李小红',
     'contact_phone','13800002222'
   ),0,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (8303,8202,10003,'演示用户丙','13800003333','https://picsum.photos/seed/qixi-user-c/80/80',
   JSON_OBJECT(
     'upload_image','https://picsum.photos/seed/qixi-signup-form-img-3/200/200',
     'attendee_name','赵大海',
     'contact_phone','13800003333',
     'company','七禧科技'
   ),0,DATE_SUB(NOW(),INTERVAL 12 HOUR)),
  (8304,8203,10001,'演示用户甲','13800001111','https://picsum.photos/seed/qixi-user-a/80/80',
   JSON_OBJECT(
     'upload_image','https://picsum.photos/seed/qixi-signup-form-img-4/200/200',
     'attendee_name','王小明',
     'contact_phone','13800001111'
   ),0,DATE_SUB(NOW(),INTERVAL 40 DAY))
ON DUPLICATE KEY UPDATE
  `activity_id`=VALUES(`activity_id`),
  `user_id`=VALUES(`user_id`),
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `avatar`=VALUES(`avatar`),
  `form_value`=VALUES(`form_value`),
  `is_del`=0,
  `created_at`=VALUES(`created_at`);
