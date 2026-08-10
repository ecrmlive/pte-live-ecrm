-- 专场列表本地演示数据（幂等；图片用可访问占位 URL）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

-- 确保演示用商品标签存在（关联标签下拉）
INSERT INTO `qixi_crm_a_product_label` (`id`,`name`,`description`,`color`,`sort`,`status`)
VALUES
  (7501,'七天无理由','适用于演示商品的售后保障标签','#16a34a',100,1),
  (7502,'新品尝鲜','本地验收用新品标识','#2563eb',90,1)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `description`=VALUES(`description`),
  `color`=VALUES(`color`),
  `sort`=VALUES(`sort`),
  `status`=1;

INSERT INTO `qixi_crm_a_marketing_decor`
  (`id`,`decor_type`,`name`,`code`,`cover_url`,`remark`,`payload`,`status`,`sort`,`starts_at`,`ends_at`,`is_del`)
VALUES
  (
    8003,
    'topic',
    '居家香氛专场',
    'home-fragrance-topic',
    'https://picsum.photos/seed/qixi-topic-fragrance-list/710/340',
    '中文演示专场',
    JSON_OBJECT(
      'label_id', 7502,
      'banner', JSON_ARRAY(
        'https://picsum.photos/seed/qixi-topic-fragrance-b1/750/750',
        'https://picsum.photos/seed/qixi-topic-fragrance-b2/750/750'
      ),
      'image', 'https://picsum.photos/seed/qixi-topic-fragrance-theme/710/340',
      'color', '#F5E6D3',
      'type', 2
    ),
    1,
    30,
    NULL,
    NULL,
    0
  ),
  (
    8013,
    'topic',
    '夏日清凉专场',
    'summer-cool-topic',
    'https://picsum.photos/seed/qixi-topic-summer-list/710/340',
    '中文演示专场',
    JSON_OBJECT(
      'label_id', 7501,
      'banner', JSON_ARRAY(
        'https://picsum.photos/seed/qixi-topic-summer-b1/750/750'
      ),
      'image', 'https://picsum.photos/seed/qixi-topic-summer-theme/710/340',
      'color', '#E0F2FE',
      'type', 1
    ),
    1,
    20,
    NULL,
    NULL,
    0
  ),
  (
    8014,
    'topic',
    '品质生活专场',
    'quality-life-topic',
    'https://picsum.photos/seed/qixi-topic-life-list/710/340',
    '中文演示专场',
    JSON_OBJECT(
      'label_id', 7502,
      'banner', JSON_ARRAY(
        'https://picsum.photos/seed/qixi-topic-life-b1/750/750',
        'https://picsum.photos/seed/qixi-topic-life-b2/750/750',
        'https://picsum.photos/seed/qixi-topic-life-b3/750/750'
      ),
      'image', 'https://picsum.photos/seed/qixi-topic-life-theme/710/340',
      'color', '',
      'type', 3
    ),
    1,
    10,
    NULL,
    NULL,
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
