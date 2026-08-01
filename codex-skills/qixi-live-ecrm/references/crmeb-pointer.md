# CRMEB 对照指针

详细分析：`docs/crmeb-reference.md`。

## 路径

```text
~/Downloads/CRMEB多商户系统/CRMEB_MER_v4.0
```

另有素材目录：`~/Downloads/CRMEB多商户系统/`（功能表 PNG、前端脑图）。本仓库已复制素材到 `docs/assets/`。

## 优先打开

| 目的 | 路径 |
| --- | --- |
| 业务编排 | `app/common/repositories/` |
| 路由权限别名 | `route/admin` `route/merchant` `route/api` |
| 表结构 | `install/crmeb_merchant.sql`（约 165 张 `eb_*`） |
| 支付回调 | `crmeb/listens/pay/` `crmeb/services/PayService.php` |
| 下单退款 | `repositories/store/order/*` |
| 小程序 IA | `extend/mp-weixin/v4.0/pages/` |
| PHP 二开约定 | 外部 `AGENTS.md` 与 `codex-skills/crmeb-merchant-extension-guide/` |

## 禁止

- 将 CRMEB 源码、vendor、证书键文件提交进本仓库
- 在本仓库继续以 ThinkPHP/Swoole 作为默认实现栈
