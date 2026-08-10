-- 本地验收：社区分类中文演示（对齐 CRMEB 截图：旅游 / 数码 / 彩妆 / 穿搭）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 清理本地误插的重复演示分类（固定演示 ID 以外）
DELETE FROM `qixi_crm_b_social_category`
WHERE `category_id` IN (6411,6412,6413,6414);

INSERT INTO `qixi_crm_b_social_category` (`category_id`,`cate_name`,`pid`,`is_show`,`sort`) VALUES
  (6401,'穿搭',0,1,40),
  (6402,'旅游',0,1,10),
  (6403,'数码',0,1,20),
  (6404,'彩妆',0,1,30)
ON DUPLICATE KEY UPDATE
  `cate_name`=VALUES(`cate_name`),
  `pid`=VALUES(`pid`),
  `is_show`=VALUES(`is_show`),
  `sort`=VALUES(`sort`);
