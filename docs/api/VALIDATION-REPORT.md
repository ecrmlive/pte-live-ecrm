# API 文档自动校验报告

置信度：high=1999 stale=26 unresolved=2 total=2027。

## 本轮核对结论

1. **方法映射**：1999 条 high 均在控制器中存在对应方法；0 条 false-high。
2. **下单三条线**已写入 FUNCTIONAL-TRUTH.md（普通 v2 / 积分 order/v3 / 死路由）。
3. **购物车/地址**参数已按 checkParams 补全（此前文档几乎为空，会误导开发）。
4. **openapi/product/create** 字段对齐 `ProductRepository::CREATE_PARAMS`（63+鉴权）。

## 开发规则

1. 禁止按 stale/unresolved 实现。
2. 普通下单只用 `/api/v2/order/*`；积分单用 `/api/order/v3/*`；勿混用 `/api/v3/order/*`。

## stale 全表

| method | path | declared | file |
| --- | --- | --- | --- |
| GET | `api/copyright` | `copyright` | `app/controller/api/Common.php` |
| GET | `api/excel/download/:id` | `download` | `app/controller/merchant/store/order/Order.php` |
| POST | `api/order/check` | `checkOrder` | `app/controller/api/store/order/StoreOrder.php` |
| POST | `api/order/create` | `createOrder` | `app/controller/api/store/order/StoreOrder.php` |
| POST | `api/order/v3/pay/:id` | `orderPay` | `app/controller/api/store/order/PointsOrder.php` |
| GET | `api/user/form/detail/:id` | `detail` | `app/controller/api/store/form/FormRelated.php` |
| POST | `api/v3/order/check` | `v3CheckOrder` | `app/controller/api/store/order/StoreOrder.php` |
| POST | `api/v3/order/create` | `v3CreateOrder` | `app/controller/api/store/order/StoreOrder.php` |
| POST | `mer/broadcast/assistant/mark/:id` | `mark` | `app/controller/merchant/store/broadcast/BroadcastAssistant.php` |
| GET | `mer/product/cdkey/library/detail/:id` | `detail` | `app/controller/merchant/store/product/CdkeyLibrary.php` |
| GET | `mer/product/cdkey/library/excel` | `excel` | `app/controller/merchant/store/product/CdkeyLibrary.php` |
| GET | `mer/serve/detail/:id` | `detail` | `app/controller/merchant/system/serve/Serve.php` |
| POST | `mer/store/productcopy/save` | `save` | `app/controller/merchant/store/product/ProductCopy.php` |
| GET | `mer/test` | `test` | `app/controller/merchant/system/admin/Login.php` |
| GET | `sys/auth` | `auth` | `app/controller/admin/Common.php` |
| GET | `sys/check_auth` | `check_auth` | `app/controller/admin/Common.php` |
| GET | `sys/community/category/detail/:id` | `detail` | `app/controller/admin/community/CommunityCategory.php` |
| GET | `sys/community/topic/detail/:id` | `detail` | `app/controller/admin/community/CommunityTopic.php` |
| GET | `sys/copyright/auth` | `authCopyright` | `app/controller/admin/Common.php` |
| GET | `sys/member/interests/options` | `options` | `app/controller/admin/user/MemberInterests.php` |
| GET | `sys/notice/config/detail/:id` | `detail` | `app/controller/admin/system/notice/SystemNoticeConfig.php` |
| GET | `sys/points/cate/detail/:id` | `detail` | `app/controller/admin/points/Category.php` |
| POST | `sys/seckill/product/delete/:id` | `delete` | `app/controller/admin/store/StoreProductSeckill.php` |
| POST | `sys/seckill/product/destory/:id` | `destory` | `app/controller/admin/store/StoreProductSeckill.php` |
| POST | `sys/store/product/batch_cate_hot` | `batchCateHot` | `app/controller/admin/store/StoreProduct.php` |
| POST | `sys/system/form/status/:id` | `statusSwitch` | `app/controller/admin/system/form/Form.php` |

## unresolved

- `GET /sys/micro/recovery/:id` — 未能可靠映射到控制器，开发时勿直接照抄，需对照 route 源码
- `ANY /api/store/test` — 未能可靠映射到控制器，开发时勿直接照抄，需对照 route 源码
