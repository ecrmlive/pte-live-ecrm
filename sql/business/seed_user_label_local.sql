-- 本地演示：用户运营标签（utf8mb4 中文；幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user_label` (`label_id`,`label_name`,`sort`,`is_del`,`create_time`) VALUES
  (9401,'高频复购用户',30,0,'2026-03-01 10:00:00'),
  (9402,'香氛兴趣用户',20,0,'2026-03-02 11:00:00'),
  (9403,'售后关怀用户',10,0,'2026-03-03 12:00:00'),
  (9404,'提示',5,0,'2026-03-04 13:00:00'),
  (9405,'测试',1,0,'2026-03-05 14:00:00')
ON DUPLICATE KEY UPDATE
  `label_name`=VALUES(`label_name`),
  `sort`=VALUES(`sort`),
  `is_del`=VALUES(`is_del`),
  `create_time`=VALUES(`create_time`);
