-- 平台「营销 - 优惠套餐」本地演示数据（幂等）
-- rules 对齐 CRMEB：type/package_type 0固定 1搭配；product.type 0主商品 1搭配；is_limit/is_time
USE `qixi_crm_business`;
SET NAMES utf8mb4;

INSERT INTO `qixi_crm_b_marketing_activity_view`
  (`activity_id`,`store_id`,`activity_type`,`name`,`rules`,`status`,`version`,`starts_at`,`ends_at`)
VALUES
  (
    5101,1,'discount','夏日香氛随行套餐',
    JSON_OBJECT(
      'package_price',199.00,
      'package_type',1,
      'type',1,
      'is_limit',0,
      'limit_num',0,
      'is_time',1,
      'free_shipping',TRUE,
      'remark','中文演示搭配套餐',
      'create_time','2026-07-01 09:30:00',
      'product_ids',JSON_ARRAY(1004,1006),
      'products',JSON_ARRAY(
        JSON_OBJECT('product_id',1004,'store_name','无火藤条香氛礼盒','image','https://picsum.photos/seed/qixi-discount-5101a/120/120','type',0,'spec','| 129.00'),
        JSON_OBJECT('product_id',1006,'store_name','香氛扩香石套装','image','https://picsum.photos/seed/qixi-discount-5101b/120/120','type',1,'spec','| 39.00')
      )
    ),
    1,1,
    '2026-07-01 00:00:00','2026-08-31 23:59:59'
  ),
  (
    5102,1,'discount','通勤数码固定套餐',
    JSON_OBJECT(
      'package_price',299.00,
      'package_type',0,
      'type',0,
      'is_limit',1,
      'limit_num',80,
      'is_time',0,
      'free_shipping',FALSE,
      'remark','中文演示固定套餐·不限时',
      'create_time','2026-06-18 14:20:00',
      'product_ids',JSON_ARRAY(1201,1202),
      'products',JSON_ARRAY(
        JSON_OBJECT('product_id',1201,'store_name','智能数显保温杯','image','https://picsum.photos/seed/qixi-discount-5102a/120/120','type',0,'spec','| 149.00'),
        JSON_OBJECT('product_id',1202,'store_name','便携数据线三件套','image','https://picsum.photos/seed/qixi-discount-5102b/120/120','type',0,'spec','| 59.00')
      )
    ),
    1,1,
    NULL,NULL
  ),
  (
    5103,2,'discount','居家香氛搭配套餐',
    JSON_OBJECT(
      'package_price',168.00,
      'package_type',1,
      'type',1,
      'is_limit',1,
      'limit_num',25,
      'is_time',1,
      'free_shipping',TRUE,
      'remark','中文演示搭配套餐·限时限量',
      'create_time','2026-07-10 11:05:00',
      'product_ids',JSON_ARRAY(1101,1102),
      'products',JSON_ARRAY(
        JSON_OBJECT('product_id',1101,'store_name','无火藤条香氛礼盒','image','https://picsum.photos/seed/qixi-discount-5103a/120/120','type',0,'spec','| 99.00'),
        JSON_OBJECT('product_id',1102,'store_name','棉柔香氛巾','image','https://picsum.photos/seed/qixi-discount-5103b/120/120','type',1,'spec','| 29.00')
      )
    ),
    1,1,
    '2026-07-15 00:00:00','2026-09-15 23:59:59'
  ),
  (
    5104,2,'discount','秋日居家固定套餐',
    JSON_OBJECT(
      'package_price',88.00,
      'package_type',0,
      'type',0,
      'is_limit',0,
      'limit_num',0,
      'is_time',0,
      'free_shipping',TRUE,
      'remark','中文演示固定套餐·不限时不限量',
      'create_time','2026-08-01 16:40:00',
      'product_ids',JSON_ARRAY(1103),
      'products',JSON_ARRAY(
        JSON_OBJECT('product_id',1103,'store_name','棉柔毛巾三件套','image','https://picsum.photos/seed/qixi-discount-5104a/120/120','type',0,'spec','| 88.00')
      )
    ),
    0,1,
    NULL,NULL
  )
ON DUPLICATE KEY UPDATE
  `store_id`=VALUES(`store_id`),
  `name`=VALUES(`name`),
  `rules`=VALUES(`rules`),
  `status`=VALUES(`status`),
  `version`=VALUES(`version`),
  `starts_at`=VALUES(`starts_at`),
  `ends_at`=VALUES(`ends_at`);
