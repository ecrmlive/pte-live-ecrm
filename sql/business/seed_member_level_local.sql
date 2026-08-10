-- 本地演示：青铜/白银/黄金/铂金/钻石会员等级（幂等）
USE `qixi_crm_business`;
SET NAMES utf8mb4;

-- 确保列存在（兼容未先跑 patch 的导入顺序）
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='icon_url')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `icon_url` varchar(1024) NOT NULL DEFAULT \'\' AFTER `rank`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;
SET @qixi_ddl := (
  SELECT IF(
    (SELECT COUNT(*) FROM information_schema.COLUMNS
      WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='qixi_crm_b_member_level' AND COLUMN_NAME='created_at')=0,
    'ALTER TABLE `qixi_crm_b_member_level` ADD COLUMN `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP AFTER `version`',
    'SELECT 1'
  )
);
PREPARE qixi_stmt FROM @qixi_ddl; EXECUTE qixi_stmt; DEALLOCATE PREPARE qixi_stmt;

INSERT INTO `qixi_crm_b_member_level`
  (`id`,`name`,`rank`,`icon_url`,`rules`,`benefits`,`status`,`version`,`created_at`,`deleted_at`)
VALUES
  (8101,'青铜会员',1,
   'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22g%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23CD7F32%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23A0522D%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2232%22%20r%3D%2228%22%20fill%3D%22url(%23g)%22/%3E%3Ctext%20x%3D%2232%22%20y%3D%2239%22%20text-anchor%3D%22middle%22%20font-size%3D%2218%22%20fill%3D%22%23fff%22%20font-weight%3D%22700%22%3E铜%3C/text%3E%3C/svg%3E',
   JSON_OBJECT('value',0,'image','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%22320%22%20height%3D%22120%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22b%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23CD7F32%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23F5DEB3%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Crect%20width%3D%22320%22%20height%3D%22120%22%20fill%3D%22url(%23b)%22/%3E%3C/svg%3E','description','注册即享基础会员服务'),
   JSON_ARRAY('积分兑换','优惠券领取'),1,1,'2026-03-01 10:00:00',NULL),
  (8102,'白银会员',2,
   'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22g%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23C0C0C0%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%2394A3B8%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2232%22%20r%3D%2228%22%20fill%3D%22url(%23g)%22/%3E%3Ctext%20x%3D%2232%22%20y%3D%2239%22%20text-anchor%3D%22middle%22%20font-size%3D%2218%22%20fill%3D%22%23fff%22%20font-weight%3D%22700%22%3E银%3C/text%3E%3C/svg%3E',
   JSON_OBJECT('value',100,'image','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%22320%22%20height%3D%22120%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22b%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23C0C0C0%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23E2E8F0%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Crect%20width%3D%22320%22%20height%3D%22120%22%20fill%3D%22url(%23b)%22/%3E%3C/svg%3E','description','累计成长值满 100 自动升级'),
   JSON_ARRAY('专属优惠提醒','会员活动优先参与'),1,1,'2026-03-01 10:05:00',NULL),
  (8103,'黄金会员',3,
   'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22g%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23F59E0B%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23D97706%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2232%22%20r%3D%2228%22%20fill%3D%22url(%23g)%22/%3E%3Ctext%20x%3D%2232%22%20y%3D%2239%22%20text-anchor%3D%22middle%22%20font-size%3D%2218%22%20fill%3D%22%23fff%22%20font-weight%3D%22700%22%3E金%3C/text%3E%3C/svg%3E',
   JSON_OBJECT('value',500,'image','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%22320%22%20height%3D%22120%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22b%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23F59E0B%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23FEF3C7%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Crect%20width%3D%22320%22%20height%3D%22120%22%20fill%3D%22url(%23b)%22/%3E%3C/svg%3E','description','累计成长值满 500 自动升级'),
   JSON_ARRAY('优先客服','会员专享活动'),1,1,'2026-03-01 10:10:00',NULL),
  (8104,'铂金会员',4,
   'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22g%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%2360A5FA%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%233B82F6%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2232%22%20r%3D%2228%22%20fill%3D%22url(%23g)%22/%3E%3Ctext%20x%3D%2232%22%20y%3D%2239%22%20text-anchor%3D%22middle%22%20font-size%3D%2218%22%20fill%3D%22%23fff%22%20font-weight%3D%22700%22%3E铂%3C/text%3E%3C/svg%3E',
   JSON_OBJECT('value',2000,'image','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%22320%22%20height%3D%22120%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22b%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%233B82F6%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23DBEAFE%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Crect%20width%3D%22320%22%20height%3D%22120%22%20fill%3D%22url(%23b)%22/%3E%3C/svg%3E','description','累计成长值满 2000 自动升级'),
   JSON_ARRAY('专属客服','生日礼遇'),1,1,'2026-03-01 10:15:00',NULL),
  (8105,'钻石会员',5,
   'data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%2264%22%20height%3D%2264%22%20viewBox%3D%220%200%2064%2064%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22g%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%23A78BFA%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%237C3AED%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle%20cx%3D%2232%22%20cy%3D%2232%22%20r%3D%2228%22%20fill%3D%22url(%23g)%22/%3E%3Ctext%20x%3D%2232%22%20y%3D%2239%22%20text-anchor%3D%22middle%22%20font-size%3D%2218%22%20fill%3D%22%23fff%22%20font-weight%3D%22700%22%3E钻%3C/text%3E%3C/svg%3E',
   JSON_OBJECT('value',10000,'image','data:image/svg+xml;charset=UTF-8,%3Csvg%20xmlns%3D%22http%3A//www.w3.org/2000/svg%22%20width%3D%22320%22%20height%3D%22120%22%3E%3Cdefs%3E%3ClinearGradient%20id%3D%22b%22%20x1%3D%220%22%20y1%3D%220%22%20x2%3D%221%22%20y2%3D%221%22%3E%3Cstop%20offset%3D%220%25%22%20stop-color%3D%22%237C3AED%22/%3E%3Cstop%20offset%3D%22100%25%22%20stop-color%3D%22%23EDE9FE%22/%3E%3C/linearGradient%3E%3C/defs%3E%3Crect%20width%3D%22320%22%20height%3D%22120%22%20fill%3D%22url(%23b)%22/%3E%3C/svg%3E','description','累计成长值满 10000 自动升级'),
   JSON_ARRAY('尊享礼遇','优先售后'),1,1,'2026-03-01 10:20:00',NULL)
ON DUPLICATE KEY UPDATE
  `name`=VALUES(`name`),
  `rank`=VALUES(`rank`),
  `icon_url`=VALUES(`icon_url`),
  `rules`=VALUES(`rules`),
  `benefits`=VALUES(`benefits`),
  `status`=VALUES(`status`),
  `version`=VALUES(`version`),
  `created_at`=VALUES(`created_at`),
  `deleted_at`=NULL;

-- 演示人数：多用户落到不同等级（虚构账号）
INSERT INTO `qixi_crm_b_member_account` (`user_id`,`level_id`,`points`,`balance`,`commission`) VALUES
  (9101,8102,268,36.50,0.00),
  (9001,8101,20,0.00,0.00),
  (9002,8103,620,12.00,0.00),
  (9003,8104,2100,88.00,0.00)
ON DUPLICATE KEY UPDATE
  `level_id`=VALUES(`level_id`),
  `points`=VALUES(`points`),
  `balance`=VALUES(`balance`),
  `commission`=VALUES(`commission`);

UPDATE `qixi_crm_b_member_level_log`
SET `note`='完成演示会员成长任务，已升级为白银会员。'
WHERE `id`=81001;
