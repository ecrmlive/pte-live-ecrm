-- 本地验收：社区话题演示（对齐截图：潮搭/彩妆分享/科技发展/说走就走的旅行）
-- 上级分类复用社区分类 seed：6401 穿搭 / 6404 彩妆 / 6403 数码 / 6402 旅游
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 隐藏旧演示话题，避免与截图样例混排
UPDATE `qixi_crm_b_social_topic`
SET `is_del`=1
WHERE `topic_id` IN (6501,6502,6503);

INSERT INTO `qixi_crm_b_social_topic`
  (`topic_id`,`topic_name`,`pic`,`status`,`is_hot`,`category_id`,`is_del`,`count_use`,`sort`) VALUES
  (6511,'潮搭','/demo/community-topic-fashion.png',1,1,6401,0,12,10),
  (6512,'彩妆分享','/demo/community-topic-makeup.png',1,1,6404,0,8,20),
  (6513,'科技发展','/demo/community-topic-digital.png',1,0,6403,0,5,30),
  (6514,'说走就走的旅行','/demo/community-topic-travel.png',1,1,6402,0,15,40)
ON DUPLICATE KEY UPDATE
  `topic_name`=VALUES(`topic_name`),
  `pic`=VALUES(`pic`),
  `status`=VALUES(`status`),
  `is_hot`=VALUES(`is_hot`),
  `category_id`=VALUES(`category_id`),
  `is_del`=0,
  `count_use`=VALUES(`count_use`),
  `sort`=VALUES(`sort`);
