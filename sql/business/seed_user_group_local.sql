-- 本地演示：用户运营分组（utf8mb4 中文；幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user_group` (`group_id`,`group_name`,`sort`,`is_del`,`create_time`) VALUES
  (9501,'精选会员',40,0,'2026-03-01 10:00:00'),
  (9502,'新品体验用户',30,0,'2026-03-02 11:00:00'),
  (9503,'高价值复购用户',20,0,'2026-03-03 12:00:00'),
  (9504,'沉睡召回用户',10,0,'2026-03-04 13:00:00'),
  (9505,'活动邀约用户',5,0,'2026-03-05 14:00:00')
ON DUPLICATE KEY UPDATE
  `group_name`=VALUES(`group_name`),
  `sort`=VALUES(`sort`),
  `is_del`=VALUES(`is_del`),
  `create_time`=VALUES(`create_time`);
