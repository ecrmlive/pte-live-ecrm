-- 在 pte-live-im 共享 MySQL（pte_live_mysql）上创建本仓库与账号。
-- IM 先起 db 后执行（示例）：
--   docker exec -i pte_live_mysql mysql -uroot -p"$MYSQL_ROOT_PASSWORD" < sql/000_shared_im_mysql_bootstrap.sql
-- 或本机：
--   mysql -h127.0.0.1 -P13306 -uroot -p < sql/000_shared_im_mysql_bootstrap.sql

CREATE DATABASE IF NOT EXISTS `qixi_mergers`
  DEFAULT CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS 'qixi'@'%' IDENTIFIED BY 'qixi_local';
CREATE USER IF NOT EXISTS 'qixi'@'localhost' IDENTIFIED BY 'qixi_local';
GRANT ALL PRIVILEGES ON `qixi_mergers`.* TO 'qixi'@'%';
GRANT ALL PRIVILEGES ON `qixi_mergers`.* TO 'qixi'@'localhost';
FLUSH PRIVILEGES;
