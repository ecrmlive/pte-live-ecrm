-- 本地演示：付费会员权益（中文 + 可访问 80x80 图标）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_svip_interest`
  (`id`,`name`,`display_name`,`description`,`icon_url`,`on_icon_url`,`link`,`status`,`sort`,`version`,`deleted_at`)
VALUES
  (981001,'会员专属价','会员尊享专属价','会员尊享专属价',
   'https://picsum.photos/seed/qixi-svip-price-off/80/80',
   'https://picsum.photos/seed/qixi-svip-price-on/80/80',
   'https://mer.crmeb.net',1,10,1,NULL),
  (981002,'签到返利','返多倍积分','返多倍积分',
   'https://picsum.photos/seed/qixi-svip-sign-off/80/80',
   'https://picsum.photos/seed/qixi-svip-sign-on/80/80',
   '/pages/users/user_sgin/index',1,20,1,NULL),
  (981003,'消费返利','返多倍积分','返多倍积分',
   'https://picsum.photos/seed/qixi-svip-pay-off/80/80',
   'https://picsum.photos/seed/qixi-svip-pay-on/80/80',
   'https://mer.crmeb.net',1,30,1,NULL),
  (981004,'专属客服','尊享客服','尊享客服',
   'https://picsum.photos/seed/qixi-svip-service-off/80/80',
   'https://picsum.photos/seed/qixi-svip-service-on/80/80',
   '/pages/chat/customer_list/chat?mer_id=0',1,40,1,NULL),
  (981005,'经验翻倍','返多倍经验','返多倍经验',
   'https://picsum.photos/seed/qixi-svip-exp-off/80/80',
   'https://picsum.photos/seed/qixi-svip-exp-on/80/80',
   '/pages/users/user_grade/index',1,50,1,NULL),
  (981006,'会员优惠券','专属优惠券','专属优惠券',
   'https://picsum.photos/seed/qixi-svip-coupon-off/80/80',
   'https://picsum.photos/seed/qixi-svip-coupon-on/80/80',
   '/pages/columnGoods/goods_coupon_list/index',1,60,1,NULL)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `display_name`=VALUES(`display_name`),
  `description`=VALUES(`description`),
  `icon_url`=VALUES(`icon_url`),
  `on_icon_url`=VALUES(`on_icon_url`),
  `link`=VALUES(`link`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `version`=VALUES(`version`),
  `deleted_at`=NULL,
  `updated_at`=NOW();

