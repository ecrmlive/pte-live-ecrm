-- 秒杀配置（场次）演示数据，对齐 CRMEB 秒杀配置列表
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_seckill_time`
  (`seckill_time_id`,`title`,`start_time`,`end_time`,`status`,`pic`)
VALUES
  (1,'午夜',0,6,1,''),
  (2,'早上',7,12,1,''),
  (3,'下午',14,19,1,''),
  (4,'晚上',19,24,1,'')
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),
  `start_time`=VALUES(`start_time`),
  `end_time`=VALUES(`end_time`),
  `status`=VALUES(`status`);
