-- 文章分类：补齐 CRMEB 字段（说明/图片/创建时间）+ 本地演示「生活」
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_article_category' AND COLUMN_NAME='info')=0,
    'ALTER TABLE `qixi_crm_a_article_category` ADD COLUMN `info` varchar(255) NOT NULL DEFAULT '''' COMMENT ''分类说明'' AFTER `title`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_article_category' AND COLUMN_NAME='image')=0,
    'ALTER TABLE `qixi_crm_a_article_category` ADD COLUMN `image` varchar(1024) NOT NULL DEFAULT '''' COMMENT ''分类图片'' AFTER `info`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_a_article_category' AND COLUMN_NAME='create_time')=0,
    'ALTER TABLE `qixi_crm_a_article_category` ADD COLUMN `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''创建时间'' AFTER `is_del`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

INSERT INTO `qixi_crm_a_article_category` (`cid`,`title`,`info`,`image`,`status`,`sort`,`is_del`,`create_time`) VALUES
  (501,'商城公告','平台公告与服务通知','https://picsum.photos/seed/qixi-article-cate-notice/120/120',1,20,0,NOW()),
  (502,'选购指南','商品选购与搭配建议','https://picsum.photos/seed/qixi-article-cate-guide/120/120',1,10,0,NOW()),
  (503,'售后须知','售后与退换货说明','https://picsum.photos/seed/qixi-article-cate-aftersale/120/120',1,5,0,NOW()),
  (504,'生活','生活灵感与日常好物','https://picsum.photos/seed/qixi-article-cate-life/120/120',1,30,0,NOW())
ON DUPLICATE KEY UPDATE
  `title`=VALUES(`title`),`info`=VALUES(`info`),`image`=VALUES(`image`),
  `status`=VALUES(`status`),`sort`=VALUES(`sort`),`is_del`=VALUES(`is_del`);
