-- 平台「拼团 → 拼团活动列表」开团记录演示数据（幂等）
-- 对应 qixi_crm_b_combination_buying / qixi_crm_b_combination_member，不是拼团商品配置。
SET NAMES utf8mb4;

-- 团长/参团演示用户（仅昵称，无登录凭据）
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`auth_version`) VALUES
  (9101,'CRM Live体验用户',NULL,1,1),
  (9201,'小林',NULL,1,1),
  (9202,'阿澈',NULL,1,1),
  (9203,'晚风',NULL,1,1)
ON DUPLICATE KEY UPDATE `nickname`=VALUES(`nickname`),`status`=VALUES(`status`);

-- 开团记录：已完成 / 未完成 / 多店铺
INSERT INTO `qixi_crm_b_combination_buying` (
  `group_buying_id`,`product_group_id`,`status`,`buying_count_num`,`buying_num`,
  `yet_buying_num`,`is_del`,`mer_id`,`end_time`,`create_time`
) VALUES
  (61101,6101,10,2,1,2,0,1,UNIX_TIMESTAMP(DATE_SUB(NOW(),INTERVAL 12 HOUR)),DATE_SUB(NOW(),INTERVAL 2 DAY)),
  (61102,6101,0,2,1,1,0,1,UNIX_TIMESTAMP(DATE_ADD(NOW(),INTERVAL 20 HOUR)),DATE_SUB(NOW(),INTERVAL 6 HOUR)),
  (61103,6102,0,3,1,2,0,1,UNIX_TIMESTAMP(DATE_ADD(NOW(),INTERVAL 18 HOUR)),DATE_SUB(NOW(),INTERVAL 1 DAY)),
  (61104,6103,10,2,1,2,0,2,UNIX_TIMESTAMP(DATE_SUB(NOW(),INTERVAL 3 HOUR)),DATE_SUB(NOW(),INTERVAL 30 HOUR)),
  (61105,6103,0,2,1,1,0,2,UNIX_TIMESTAMP(DATE_ADD(NOW(),INTERVAL 10 HOUR)),DATE_SUB(NOW(),INTERVAL 3 HOUR))
ON DUPLICATE KEY UPDATE
  `product_group_id`=VALUES(`product_group_id`),
  `status`=VALUES(`status`),
  `buying_count_num`=VALUES(`buying_count_num`),
  `buying_num`=VALUES(`buying_num`),
  `yet_buying_num`=VALUES(`yet_buying_num`),
  `is_del`=VALUES(`is_del`),
  `mer_id`=VALUES(`mer_id`),
  `end_time`=VALUES(`end_time`),
  `create_time`=VALUES(`create_time`);

-- 参团成员（团长 is_initiator/is_leader=1）
INSERT INTO `qixi_crm_b_combination_member` (
  `id`,`group_buying_id`,`product_group_id`,`status`,`is_initiator`,`order_id`,
  `uid`,`nickname`,`avatar`,`is_del`,`create_time`,`is_leader`
) VALUES
  (6110101,61101,6101,1,1,0,9101,'CRM Live体验用户','',0,DATE_SUB(NOW(),INTERVAL 2 DAY),1),
  (6110102,61101,6101,1,0,0,9201,'小林','',0,DATE_SUB(NOW(),INTERVAL 47 HOUR),0),
  (6110201,61102,6101,1,1,0,9203,'晚风','',0,DATE_SUB(NOW(),INTERVAL 6 HOUR),1),
  (6110301,61103,6102,1,1,0,9201,'小林','',0,DATE_SUB(NOW(),INTERVAL 1 DAY),1),
  (6110302,61103,6102,1,0,0,9202,'阿澈','',0,DATE_SUB(NOW(),INTERVAL 20 HOUR),0),
  (6110401,61104,6103,1,1,0,9202,'阿澈','',0,DATE_SUB(NOW(),INTERVAL 30 HOUR),1),
  (6110402,61104,6103,1,0,0,9101,'CRM Live体验用户','',0,DATE_SUB(NOW(),INTERVAL 28 HOUR),0),
  (6110501,61105,6103,1,1,0,9203,'晚风','',0,DATE_SUB(NOW(),INTERVAL 3 HOUR),1)
ON DUPLICATE KEY UPDATE
  `group_buying_id`=VALUES(`group_buying_id`),
  `product_group_id`=VALUES(`product_group_id`),
  `status`=VALUES(`status`),
  `is_initiator`=VALUES(`is_initiator`),
  `uid`=VALUES(`uid`),
  `nickname`=VALUES(`nickname`),
  `is_del`=VALUES(`is_del`),
  `create_time`=VALUES(`create_time`),
  `is_leader`=VALUES(`is_leader`);
