-- 平台报名活动表 + 报名记录（幂等）
-- 对齐 CRMEB eb_store_activity(activity_type=form) / eb_store_activity_related
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_activity` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(128) NOT NULL COMMENT '活动名称',
  `info` varchar(500) NOT NULL DEFAULT '' COMMENT '活动简介',
  `cover_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '封面图 750*350',
  `poster_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '分享海报 750*1250',
  `color` varchar(32) NOT NULL DEFAULT '' COMMENT '活动背景色',
  `form_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '关联系统表单 qixi_crm_a_config_item.id',
  `quota` int unsigned NOT NULL DEFAULT 0 COMMENT '人数上限，0=不限制',
  `total` int unsigned NOT NULL DEFAULT 0 COMMENT '已报名人数',
  `status` tinyint NOT NULL DEFAULT 1 COMMENT '是否显示/开启 1是 0否',
  `sort` int NOT NULL DEFAULT 0 COMMENT '排序',
  `starts_at` datetime DEFAULT NULL COMMENT '开始时间',
  `ends_at` datetime DEFAULT NULL COMMENT '结束时间',
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_form` (`form_id`),
  KEY `idx_visible` (`is_del`,`status`,`sort`),
  KEY `idx_time` (`starts_at`,`ends_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台报名活动';

CREATE TABLE IF NOT EXISTS `qixi_crm_a_signup_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` bigint unsigned NOT NULL COMMENT '报名活动ID',
  `user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '用户ID',
  `nickname` varchar(128) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `mobile` varchar(32) NOT NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(1024) NOT NULL DEFAULT '' COMMENT '头像',
  `form_value` json NOT NULL COMMENT '动态表单提交内容',
  `is_del` tinyint NOT NULL DEFAULT 0,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_activity` (`activity_id`,`is_del`,`created_at`),
  KEY `idx_user` (`user_id`),
  KEY `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='平台报名活动用户记录';
