-- 平台「积分日志」本地演示数据（幂等，utf8mb4 中文可读）
SET NAMES utf8mb4;

-- 演示用户（若已存在则刷新昵称）
INSERT INTO `qixi_crm_b_user` (`id`,`nickname`,`mobile`,`status`,`group_id`,`auth_version`) VALUES
  (9101,'CRM Live体验用户','DEMO-USER-9101',1,0,1),
  (9102,'演示用户阿强','DEMO-USER-9102',1,0,1),
  (9103,'演示用户小林','DEMO-USER-9103',1,0,1)
ON DUPLICATE KEY UPDATE
  `nickname`=VALUES(`nickname`),
  `mobile`=VALUES(`mobile`),
  `status`=VALUES(`status`);

INSERT INTO `qixi_crm_b_member_account` (`user_id`,`level_id`,`points`,`balance`,`commission`) VALUES
  (9101,NULL,268,36.50,0.00),
  (9102,NULL,120,0.00,0.00),
  (9103,NULL,85,10.00,0.00)
ON DUPLICATE KEY UPDATE
  `points`=VALUES(`points`),
  `balance`=VALUES(`balance`);

-- 签到次数（统计卡「客户签到次数」）
INSERT INTO `qixi_crm_b_user_sign` (`id`,`user_id`,`sign_date`,`points`,`continuous_days`,`created_at`) VALUES
  (81201,9101,'2026-08-01',5,1,'2026-08-01 08:10:00'),
  (81202,9101,'2026-08-02',5,2,'2026-08-02 08:12:00'),
  (81203,9101,'2026-08-03',10,3,'2026-08-03 08:05:00'),
  (81204,9102,'2026-08-02',5,1,'2026-08-02 09:20:00'),
  (81205,9102,'2026-08-04',5,1,'2026-08-04 09:18:00'),
  (81206,9103,'2026-08-05',5,1,'2026-08-05 07:55:00')
ON DUPLICATE KEY UPDATE
  `points`=VALUES(`points`),
  `continuous_days`=VALUES(`continuous_days`),
  `created_at`=VALUES(`created_at`);

-- 积分流水（列表 + 汇总字段齐全）
INSERT INTO `qixi_crm_b_user_bill`
  (`bill_id`,`uid`,`link_id`,`pm`,`title`,`category`,`type`,`number`,`balance`,`mark`,`mer_id`,`status`,`create_time`)
VALUES
  (880001,9101,'SIGN-20260801',1,'签到赠送积分','integral','sign_integral',5.00,200.00,'连续签到第 1 天赠送 5 积分',0,1,'2026-08-01 08:10:00'),
  (880002,9101,'ORD-GIFT-1001',1,'下单赠送积分','integral','lock',30.00,230.00,'订单完成后赠送，冻结期内暂不可用',0,0,'2026-08-01 14:22:00'),
  (880003,9101,'SIGN-20260802',1,'签到赠送积分','integral','sign_integral',5.00,235.00,'连续签到第 2 天赠送 5 积分',0,1,'2026-08-02 08:12:00'),
  (880004,9101,'ORD-DED-2001',0,'购买商品','integral','deduction',20.00,215.00,'积分抵扣订单 ORD-DED-2001',0,1,'2026-08-02 16:40:00'),
  (880005,9102,'SIGN-20260802',1,'签到赠送积分','integral','sign_integral',5.00,80.00,'签到赠送 5 积分',0,1,'2026-08-02 09:20:00'),
  (880006,9101,'SIGN-20260803',1,'签到赠送积分','integral','sign_integral',10.00,225.00,'连续签到第 3 天额外奖励',0,1,'2026-08-03 08:05:00'),
  (880007,9101,'ORD-GIFT-1001',0,'扣除赠送积分','integral','refund_lock',10.00,215.00,'部分退款扣回冻结赠送积分',0,1,'2026-08-03 11:30:00'),
  (880008,9102,'ORD-DED-2101',0,'购买商品','integral','deduction',15.00,65.00,'积分抵扣订单 ORD-DED-2101',0,1,'2026-08-03 15:08:00'),
  (880009,9102,'MER-REF-2101',1,'订单退款','mer_integral','refund',8.00,73.00,'商户积分退回（抵扣返还）',1,1,'2026-08-03 18:20:00'),
  (880010,9103,'SIGN-20260805',1,'签到赠送积分','integral','sign_integral',5.00,60.00,'签到赠送 5 积分',0,1,'2026-08-05 07:55:00'),
  (880011,9103,'ORD-GIFT-3001',1,'下单赠送积分','integral','lock',25.00,85.00,'下单赠送积分（冻结中）',0,0,'2026-08-05 12:40:00'),
  (880012,9101,'SYS-INC-9001',1,'系统增加积分','integral','sys_inc',50.00,265.00,'平台活动补偿发放（演示）',0,1,'2026-08-06 10:00:00'),
  (880013,9101,'PTS-EX-4001',0,'兑换商品','integral','points',12.00,253.00,'积分商城兑换演示商品',0,1,'2026-08-07 13:15:00'),
  (880014,9102,'INVITE-5001',1,'邀请好友','integral','spread',20.00,93.00,'邀请好友注册成功赠送',0,1,'2026-08-08 09:45:00'),
  (880015,9103,'ORD-DED-3101',0,'购买商品','integral','deduction',8.00,77.00,'积分抵扣订单 ORD-DED-3101',0,1,'2026-08-09 16:05:00')
ON DUPLICATE KEY UPDATE
  `uid`=VALUES(`uid`),
  `link_id`=VALUES(`link_id`),
  `pm`=VALUES(`pm`),
  `title`=VALUES(`title`),
  `category`=VALUES(`category`),
  `type`=VALUES(`type`),
  `number`=VALUES(`number`),
  `balance`=VALUES(`balance`),
  `mark`=VALUES(`mark`),
  `mer_id`=VALUES(`mer_id`),
  `status`=VALUES(`status`),
  `create_time`=VALUES(`create_time`);
