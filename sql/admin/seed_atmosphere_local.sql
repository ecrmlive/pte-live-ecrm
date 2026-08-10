-- 活动氛围图本地演示数据（幂等；图片用可访问占位 URL）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_a_marketing_decor`
  (`id`,`decor_type`,`name`,`code`,`cover_url`,`remark`,`payload`,`status`,`sort`,`starts_at`,`ends_at`,`is_del`)
VALUES
  (
    8001,
    'atmosphere',
    '夏日焕新氛围图',
    'summer-atmosphere',
    'https://picsum.photos/seed/qixi-atmosphere-summer/750/152',
    '中文演示氛围',
    JSON_OBJECT(
      'theme','summer',
      'scope_type',0,
      'spu_ids',JSON_ARRAY(),
      'cate_ids',JSON_ARRAY(),
      'mer_ids',JSON_ARRAY(),
      'label_ids',JSON_ARRAY()
    ),
    1,
    10,
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    DATE_ADD(NOW(), INTERVAL 180 DAY),
    0
  )
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `code`=VALUES(`code`),
  `cover_url`=VALUES(`cover_url`),
  `remark`=VALUES(`remark`),
  `payload`=VALUES(`payload`),
  `status`=VALUES(`status`),
  `sort`=VALUES(`sort`),
  `starts_at`=VALUES(`starts_at`),
  `ends_at`=VALUES(`ends_at`),
  `is_del`=0;
