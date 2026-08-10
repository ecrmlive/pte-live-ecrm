-- 秒杀活动页：补齐平台/运营角色菜单授权（修复 local 运营进页无写权限、前端曾被拦成空表）
USE `qixi_crm_admin`;
SET NAMES utf8mb4;

INSERT IGNORE INTO `qixi_crm_a_role_menu` (`role_id`, `menu_id`)
SELECT r.id, m.id
FROM `qixi_crm_a_role` AS r
CROSS JOIN `qixi_crm_a_menu` AS m
WHERE r.code IN ('platform', 'operations')
  AND r.status = 1
  AND m.status = 1
  AND m.code IN (
    'marketing.seckill.dir',
    'marketing.seckill.config',
    'marketing.seckill.manage.page',
    'marketing.seckill.activity',
    'marketing.seckill.manage'
  );
