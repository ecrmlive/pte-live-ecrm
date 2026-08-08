-- 虚拟评论显示头像（对齐 CRMEB StoreProductReply.avatar）。
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_product_comment' AND COLUMN_NAME='virtual_author_avatar')=0,
    'ALTER TABLE `qixi_crm_b_product_comment` ADD COLUMN `virtual_author_avatar` varchar(1024) NOT NULL DEFAULT '''' AFTER `virtual_author_name`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
