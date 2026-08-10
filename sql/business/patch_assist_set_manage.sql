-- 平台「助力 → 助力活动」实例列表演示数据（幂等）
-- 用户发起的助力单（qixi_crm_b_assist_set），不是活动商品配置。
SET NAMES utf8mb4;

-- 助力好友演示用户（仅昵称，无登录凭据）
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`auth_version`) VALUES
  (9201,'小林',NULL,1,1),
  (9202,'阿澈',NULL,1,1),
  (9203,'晚风',NULL,1,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`status`=VALUES(`status`);

-- 补充进行中 / 已满员 / 多店铺实例，供筛选与详情验收
INSERT INTO `qixi_crm_b_assist_set` (
  `product_assist_set_id`,`product_assist_id`,`product_id`,`uid`,`status`,
  `assist_count`,`assist_user_count`,`yet_assist_count`,`mer_id`,`is_del`,`create_time`
) VALUES
  (63101,6301,1005,9101,10,2,1,2,1,0,DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (63102,6301,1005,9203,1,2,1,1,1,0,DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (63103,6302,1107,9201,1,3,2,2,2,0,DATE_SUB(NOW(),INTERVAL 18 HOUR)),
  (63104,6303,1207,9202,10,4,3,4,3,0,DATE_SUB(NOW(),INTERVAL 6 HOUR))
ON DUPLICATE KEY UPDATE
  `status`=VALUES(`status`),
  `assist_count`=VALUES(`assist_count`),
  `assist_user_count`=VALUES(`assist_user_count`),
  `yet_assist_count`=VALUES(`yet_assist_count`),
  `mer_id`=VALUES(`mer_id`),
  `is_del`=VALUES(`is_del`),
  `create_time`=VALUES(`create_time`);

INSERT INTO `qixi_crm_b_assist_user` (`product_assist_set_id`,`product_assist_id`,`uid`,`nickname`,`avatar_img`) VALUES
  (63101,6301,9201,'小林',''),
  (63101,6301,9202,'阿澈',''),
  (63102,6301,9101,'CRM Live体验用户',''),
  (63103,6302,9101,'CRM Live体验用户',''),
  (63103,6302,9202,'阿澈',''),
  (63104,6303,9101,'CRM Live体验用户',''),
  (63104,6303,9201,'小林',''),
  (63104,6303,9203,'晚风',''),
  (63104,6303,9001,'晴空漫游者','')
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`avatar_img`=VALUES(`avatar_img`);
