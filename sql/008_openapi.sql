-- 阶段 6c：商户开放接口凭证（对照 CRMEB eb_open_auth → qixi_open_auth）
USE `qixi_mergers`;

CREATE TABLE IF NOT EXISTS `qixi_open_auth` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(50) DEFAULT NULL COMMENT '标题',
  `access_key` varchar(50) DEFAULT NULL,
  `secret_key` varchar(255) DEFAULT NULL,
  `status` tinyint(4) DEFAULT 1 COMMENT '状态 1启用',
  `mark` varchar(255) DEFAULT NULL COMMENT '备注',
  `mer_id` int(11) DEFAULT NULL COMMENT '商户ID',
  `auth` varchar(255) DEFAULT '1,2' COMMENT '权限：1商品 2订单',
  `sort` int(11) DEFAULT 0,
  `is_del` tinyint(4) DEFAULT 0 COMMENT '是否删除',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `delete_time` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  `last_ip` varchar(50) DEFAULT NULL COMMENT '最后登录的IP',
  `last_time` timestamp NULL DEFAULT NULL COMMENT '最后登录的时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_access_key` (`access_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商户开放接口凭证';

-- 演示凭证：商户1（明文 secret 仅本地演示，生产勿复用）
INSERT INTO `qixi_open_auth` (
  `id`, `title`, `access_key`, `secret_key`, `status`, `mark`, `mer_id`, `auth`, `sort`, `is_del`
)
SELECT 1, '演示开放凭证', 'demo_mer1_ak', 'demo_mer1_sk', 1, '本地演示', 1, '1,2', 1, 0
WHERE NOT EXISTS (SELECT 1 FROM `qixi_open_auth` WHERE `access_key` = 'demo_mer1_ak');

INSERT INTO `qixi_schema_meta` (`version`, `note`)
SELECT 'phase6c-openapi', '阶段6c：商户开放接口 qixi_open_auth'
WHERE NOT EXISTS (SELECT 1 FROM `qixi_schema_meta` WHERE `version` = 'phase6c-openapi');
