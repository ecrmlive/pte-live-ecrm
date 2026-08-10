-- 本地演示：反馈分类父子树（中文，幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_user_feedback_category` (`id`,`pid`,`name`,`sort`,`status`) VALUES
  (9711,0,'功能建议',10,1),
  (9712,0,'订单问题',20,1),
  (9713,0,'使用体验',30,1),
  (9714,0,'历史分类',90,0),
  (9720,0,'测试',1,1),
  (9721,9720,'测试下',1,1)
ON DUPLICATE KEY UPDATE
  `pid`=VALUES(`pid`),
  `name`=VALUES(`name`),
  `sort`=VALUES(`sort`),
  `status`=VALUES(`status`),
  `deleted_at`=NULL;
