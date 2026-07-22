# 用户端 — 接口级功能点（等同按钮操作）

> 来源：`route/api/notLogin.php` + `route/api/login.php`。页面清单见文末附录。

> 每一条路由 = 一个可调用操作（等同按钮/动作）。CRUD 由 HTTP 方法 + 路径/处理器名推断，可人工校正。

> **源码核对（2026-07-21）**：下列处理器在控制器中**不存在**，属死路由，**禁止按此开发**——`/createOrder`、`/checkOrder`、`/v3CreateOrder`、`/v3CheckOrder`、`PointsOrder/orderPay`。  
> 普通下单用 `v2CreateOrder` / `v2CheckOrder`；积分下单用 `PointsOrder/createOrder` / `beforCheck`。  
> 完整对照：[`docs/api/FUNCTIONAL-TRUTH.md`](../api/FUNCTIONAL-TRUTH.md)。

## 统计

| 项 | 数量 |
| --- | ---: |
| 操作（路由）总数 | 342 |
| C | 66 |
| R | 237 |
| U | 29 |
| D | 20 |

## 模块 CRUD 覆盖

| 模块 | 操作数 | C | R | U | D |
| --- | ---: | --- | --- | --- | --- |
| `account` | 1 | — | ✓ | — | — |
| `active/category` | 1 | — | ✓ | — | — |
| `address/create` | 1 | ✓ | — | — | — |
| `address/delete` | 1 | ✓ | — | — | ✓ |
| `address/detail` | 1 | ✓ | ✓ | — | — |
| `address/lst` | 1 | ✓ | ✓ | — | — |
| `address/update` | 1 | ✓ | — | ✓ | — |
| `again` | 1 | ✓ | — | — | — |
| `agreement` | 1 | — | ✓ | — | — |
| `agreement_lst` | 1 | — | ✓ | — | — |
| `ajcaptcha` | 1 | — | ✓ | — | — |
| `ajcheck` | 1 | — | — | ✓ | — |
| `appVersion` | 1 | — | ✓ | — | — |
| `apply` | 1 | — | — | — | ✓ |
| `assist` | 1 | — | ✓ | — | — |
| `assist/create` | 1 | ✓ | — | ✓ | — |
| `assist/detail` | 1 | — | ✓ | — | — |
| `assist/lst` | 1 | — | ✓ | — | — |
| `assist/set` | 3 | — | ✓ | ✓ | ✓ |
| `assist/share` | 1 | — | ✓ | — | — |
| `assist/user` | 1 | — | ✓ | — | — |
| `auth` | 1 | ✓ | — | — | — |
| `auth/app` | 1 | ✓ | — | — | — |
| `auth/apple` | 1 | ✓ | — | — | — |
| `auth/login` | 1 | ✓ | — | — | — |
| `auth/mp` | 1 | ✓ | — | — | — |
| `auth/mp_login_type` | 1 | ✓ | — | — | — |
| `auth/mp_phone` | 1 | ✓ | — | — | — |
| `auth/register` | 1 | ✓ | — | — | — |
| `auth/smslogin` | 1 | ✓ | — | — | — |
| `auth/verify` | 1 | — | — | ✓ | — |
| `auth/wechat` | 1 | — | ✓ | — | — |
| `back_goods` | 1 | ✓ | — | — | — |
| `bag` | 1 | — | ✓ | — | — |
| `bag/explain` | 1 | — | ✓ | — | — |
| `bag/recommend` | 1 | — | ✓ | — | — |
| `batchCreate` | 1 | ✓ | — | — | — |
| `batch_product` | 1 | — | ✓ | — | — |
| `bill` | 1 | — | ✓ | — | — |
| `binding` | 1 | — | — | ✓ | — |
| `brand/lst` | 1 | — | ✓ | — | — |
| `broadcast` | 1 | — | ✓ | — | — |
| `brokerage/all` | 1 | — | ✓ | — | — |
| `brokerage/info` | 1 | — | ✓ | — | — |
| `brokerage/notice` | 1 | — | ✓ | — | — |
| `brokerage_list` | 1 | — | ✓ | — | — |
| `brokerage_top` | 1 | — | ✓ | — | — |
| `cancel` | 2 | — | — | — | ✓ |
| `captcha` | 1 | — | ✓ | — | — |
| `cashier_order` | 1 | — | ✓ | — | — |
| `cate_hot` | 1 | — | ✓ | — | — |
| `category` | 2 | — | ✓ | — | — |
| `category/hotranking` | 1 | — | ✓ | — | — |
| `category/lst` | 4 | — | ✓ | — | — |
| `change` | 1 | — | — | ✓ | — |
| `change/info` | 1 | — | ✓ | ✓ | — |
| `change/password` | 1 | — | — | ✓ | — |
| `change/phone` | 1 | — | — | ✓ | — |
| `check` | 5 | — | — | ✓ | — |
| `clear` | 1 | — | — | — | ✓ |
| `command/copy` | 1 | — | ✓ | — | — |
| `common/base64` | 1 | — | ✓ | — | — |
| `common/commuunity` | 1 | — | ✓ | — | — |
| `common/express` | 1 | — | ✓ | — | — |
| `common/feedback_type` | 1 | — | ✓ | — | — |
| `common/home` | 1 | — | ✓ | — | — |
| `common/hot_banner` | 1 | — | ✓ | — | — |
| `common/hot_keyword` | 1 | — | ✓ | — | — |
| `common/menus` | 1 | — | ✓ | — | — |
| `common/pay_key` | 1 | ✓ | ✓ | — | — |
| `common/recharge_quota` | 1 | — | ✓ | — | — |
| `common/refund_message` | 1 | — | ✓ | — | — |
| `common/visit` | 1 | ✓ | — | — | — |
| `community` | 1 | — | ✓ | — | — |
| `compute` | 1 | ✓ | — | — | — |
| `config` | 1 | — | ✓ | — | — |
| `copy` | 1 | — | ✓ | — | — |
| `copyright` | 1 | — | ✓ | — | — |
| `count` | 1 | — | ✓ | — | — |
| `coupon` | 1 | — | ✓ | — | — |
| `coupon_lst` | 1 | — | ✓ | — | — |
| `coupon_product` | 1 | — | ✓ | — | — |
| `coupon_receive` | 1 | ✓ | — | — | — |
| `create` | 8 | ✓ | — | — | — |
| `del` | 2 | — | — | — | ✓ |
| `deleate` | 1 | — | — | — | ✓ |
| `delete` | 2 | — | — | — | ✓ |
| `delivery` | 1 | — | ✓ | — | — |
| `deliverySetings` | 1 | — | ✓ | — | — |
| `deliveryStation/list` | 1 | — | ✓ | — | — |
| `deliveryTrack` | 1 | — | ✓ | — | — |
| `detail` | 10 | — | ✓ | — | — |
| `detail/0` | 1 | — | ✓ | — | — |
| `diy` | 1 | — | ✓ | — | — |
| `excel/download` | 1 | — | ✓ | — | — |
| `express` | 2 | ✓ | ✓ | — | — |
| `extract/banklst` | 1 | — | ✓ | — | — |
| `extract/create` | 1 | ✓ | — | — | — |
| `extract/detail` | 1 | — | ✓ | — | — |
| `extract/history_bank` | 1 | — | ✓ | — | — |
| `extract/lst` | 1 | — | ✓ | — | — |
| `fab` | 1 | — | ✓ | — | — |
| `fans` | 1 | — | — | ✓ | — |
| `fans/lst` | 1 | — | ✓ | — | — |
| `feedback` | 1 | ✓ | — | — | — |
| `feedback/detail` | 1 | — | ✓ | — | — |
| `feedback/list` | 1 | — | ✓ | — | — |
| `fields/delete` | 1 | — | — | — | ✓ |
| `fields/info` | 1 | — | ✓ | — | — |
| `fields/save` | 1 | ✓ | — | — | — |
| `focus/lst` | 1 | — | ✓ | — | — |
| `focuslst` | 1 | — | ✓ | — | — |
| `get` | 1 | — | ✓ | — | — |
| `getVersion` | 1 | — | — | — | — |
| `get_attr_value` | 1 | — | ✓ | — | — |
| `get_hot_ranking` | 1 | — | ✓ | — | — |
| `get_spec` | 1 | — | ✓ | — | — |
| `getlst` | 1 | — | ✓ | — | — |
| `good_list` | 1 | — | ✓ | — | — |
| `group` | 1 | — | ✓ | — | — |
| `group/category` | 1 | — | ✓ | — | — |
| `group/count` | 1 | — | ✓ | — | — |
| `group/detail` | 1 | — | ✓ | — | — |
| `group/get` | 1 | ✓ | ✓ | — | — |
| `group/lst` | 1 | — | ✓ | — | — |
| `group_order_detail` | 1 | — | ✓ | — | — |
| `group_order_list` | 1 | — | ✓ | — | — |
| `guarantee` | 1 | — | ✓ | — | — |
| `has_service` | 1 | — | ✓ | — | — |
| `hist_product/lst` | 1 | — | ✓ | — | — |
| `history` | 1 | — | ✓ | — | — |
| `history/batch` | 1 | — | ✓ | — | ✓ |
| `history/delete` | 1 | — | ✓ | — | ✓ |
| `home` | 1 | — | ✓ | — | — |
| `hot` | 2 | — | ✓ | — | — |
| `hot_lst` | 1 | — | ✓ | — | — |
| `hot_top` | 2 | — | ✓ | — | — |
| `increase_take` | 1 | — | — | ✓ | — |
| `info` | 2 | — | ✓ | — | — |
| `integral/info` | 1 | — | ✓ | — | — |
| `integral/lst` | 1 | — | ✓ | — | — |
| `intention/business` | 1 | — | ✓ | — | — |
| `intention/cate` | 1 | — | ✓ | — | — |
| `intention/circles` | 1 | — | ✓ | — | — |
| `intention/create` | 1 | ✓ | — | — | — |
| `intention/detail` | 1 | — | ✓ | — | — |
| `intention/lst` | 1 | — | ✓ | — | — |
| `intention/type` | 1 | — | ✓ | — | — |
| `intention/update` | 1 | — | — | ✓ | — |
| `labels` | 1 | — | ✓ | — | — |
| `lbs/address` | 1 | ✓ | ✓ | — | — |
| `lbs/geocoder` | 1 | — | ✓ | — | — |
| `list` | 5 | — | ✓ | — | — |
| `local` | 2 | — | ✓ | — | — |
| `localDetail` | 1 | — | ✓ | — | — |
| `logout` | 1 | ✓ | — | — | — |
| `lst` | 12 | — | ✓ | — | — |
| `member/info` | 1 | — | ✓ | — | — |
| `member/log` | 1 | — | ✓ | — | — |
| `merchant` | 1 | — | ✓ | — | — |
| `micro` | 1 | — | ✓ | — | — |
| `mp/binding` | 1 | — | — | ✓ | — |
| `navigation` | 1 | — | ✓ | — | — |
| `new_people` | 1 | — | ✓ | — | — |
| `notice` | 1 | — | — | — | — |
| `notice/callback` | 1 | — | — | — | — |
| `notice/mchNotify` | 1 | — | — | — | — |
| `notice/pay` | 1 | ✓ | — | — | — |
| `number` | 1 | — | ✓ | — | — |
| `open_screen` | 1 | — | ✓ | — | — |
| `options` | 1 | — | ✓ | — | — |
| `order` | 1 | — | ✓ | — | — |
| `order_call_back` | 1 | — | — | — | — |
| `params` | 1 | — | ✓ | — | — |
| `params_value` | 1 | — | ✓ | — | — |
| `pay` | 3 | ✓ | — | — | — |
| `pay/config` | 1 | ✓ | ✓ | — | — |
| `pay_lst` | 1 | ✓ | ✓ | — | — |
| `pay_product/lst` | 1 | ✓ | ✓ | — | — |
| `platform` | 1 | ✓ | — | — | — |
| `points/pay` | 1 | ✓ | — | — | — |
| `presell` | 1 | — | ✓ | — | — |
| `presell/agree` | 1 | — | ✓ | — | — |
| `presell/detail` | 1 | — | ✓ | — | — |
| `presell/lst` | 1 | — | ✓ | — | — |
| `preview` | 1 | — | ✓ | — | — |
| `price_rule` | 1 | — | ✓ | — | — |
| `product` | 3 | — | ✓ | — | — |
| `product/lst` | 1 | — | ✓ | — | — |
| `productCategory` | 1 | — | ✓ | — | — |
| `product_detail` | 1 | — | ✓ | — | — |
| `product_lst` | 1 | — | ✓ | — | — |
| `qrcode` | 3 | — | ✓ | — | — |
| `receipt` | 1 | ✓ | — | — | — |
| `receipt/create` | 1 | ✓ | — | — | — |
| `receipt/delete` | 1 | — | — | — | ✓ |
| `receipt/detail` | 1 | — | ✓ | — | — |
| `receipt/is_default` | 1 | ✓ | — | — | — |
| `receipt/lst` | 1 | — | ✓ | — | — |
| `receipt/order` | 2 | — | ✓ | — | — |
| `receipt/update` | 1 | — | — | ✓ | — |
| `receive` | 1 | ✓ | — | — | — |
| `recharge` | 1 | ✓ | — | — | — |
| `recharge/brokerage` | 1 | ✓ | — | — | — |
| `recommend` | 2 | — | ✓ | — | — |
| `recommendProduct` | 1 | — | ✓ | — | — |
| `rela_product/lst` | 1 | — | ✓ | — | — |
| `relation/batch` | 2 | ✓ | — | — | ✓ |
| `relation/create` | 1 | ✓ | — | — | — |
| `relation/delete` | 1 | — | — | — | ✓ |
| `relation/merchant` | 1 | — | ✓ | — | — |
| `relation/product` | 1 | — | ✓ | — | — |
| `reply` | 1 | — | ✓ | — | — |
| `reply/create` | 1 | ✓ | — | ✓ | — |
| `reply/lst` | 1 | — | ✓ | — | — |
| `reply/start` | 1 | — | — | ✓ | — |
| `reservation/checkRange` | 1 | — | — | ✓ | — |
| `reservation/getDay` | 1 | — | ✓ | — | — |
| `reservation/getMonth` | 1 | — | ✓ | — | — |
| `revoke` | 1 | ✓ | — | — | — |
| `route/list` | 1 | — | ✓ | — | — |
| `scan_upload/image` | 1 | ✓ | — | — | — |
| `scope` | 1 | — | ✓ | — | — |
| `script` | 1 | — | ✓ | — | — |
| `seckill` | 1 | — | ✓ | — | — |
| `seckill/detail` | 1 | — | ✓ | — | — |
| `seckill/lst` | 1 | — | ✓ | — | — |
| `seckill/select` | 1 | — | ✓ | — | — |
| `self/cancel` | 1 | — | — | — | ✓ |
| `service/info` | 1 | — | ✓ | — | — |
| `services` | 1 | — | ✓ | — | — |
| `share_posters` | 1 | — | ✓ | — | — |
| `show` | 3 | — | ✓ | — | — |
| `sign/create` | 1 | ✓ | — | — | — |
| `sign/info` | 1 | — | ✓ | — | — |
| `sign/lst` | 1 | — | ✓ | — | — |
| `sign/month` | 1 | — | ✓ | — | — |
| `spread_image` | 1 | — | ✓ | — | — |
| `spread_info` | 1 | — | ✓ | — | — |
| `spread_level` | 1 | — | ✓ | — | — |
| `spread_list` | 1 | — | ✓ | — | — |
| `spread_order` | 1 | — | ✓ | — | — |
| `spread_top` | 1 | — | ✓ | — | — |
| `spu` | 1 | — | ✓ | — | — |
| `start` | 1 | ✓ | — | — | — |
| `start/lst` | 1 | — | ✓ | — | — |
| `status` | 1 | — | ✓ | ✓ | — |
| `store` | 2 | — | ✓ | — | — |
| `store/certificate` | 1 | — | ✓ | — | — |
| `store/expr` | 1 | — | ✓ | — | — |
| `store/product` | 2 | — | ✓ | — | ✓ |
| `store/test` | 1 | — | — | — | — |
| `subscribe` | 1 | — | ✓ | — | — |
| `svip/pay` | 1 | ✓ | — | — | — |
| `switch` | 1 | — | — | ✓ | — |
| `system/city` | 1 | — | ✓ | — | — |
| `take` | 2 | ✓ | — | — | — |
| `update` | 2 | — | — | ✓ | — |
| `upload/certificate` | 1 | ✓ | — | — | — |
| `upload/image` | 1 | ✓ | — | — | — |
| `upload/video` | 1 | ✓ | — | — | — |
| `user` | 1 | — | ✓ | — | — |
| `user/cancel` | 1 | — | — | — | ✓ |
| `user/change_pwd` | 1 | — | — | ✓ | — |
| `user/community` | 1 | — | ✓ | — | — |
| `user/community_video` | 1 | — | ✓ | — | — |
| `user/info` | 1 | — | ✓ | — | — |
| `user/spread` | 1 | — | ✓ | — | — |
| `user_info` | 1 | — | ✓ | — | — |
| `user_lst` | 1 | — | ✓ | — | — |
| `v2/spread_image` | 1 | — | ✓ | — | — |
| `v2/system` | 2 | — | ✓ | — | — |
| `verify_code` | 1 | — | ✓ | — | — |
| `version` | 1 | — | ✓ | — | — |
| `video_lst` | 1 | — | ✓ | — | — |
| `wechat/config` | 1 | — | ✓ | — | — |
| `登录用户` | 1 | — | — | ✓ | — |

## 分模块操作明细

### `account`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 切换账号 | `R` | GET | `account` | `User/account` | `` |

### `active/category`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 活动分类  product/spu/active/category/:type | `R` | GET | `/active/category/:type` | `StoreSpu/activeCategory` | `` |

### `address/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (/address/create) | `C` | POST | `/address/create` | `UserAddress/create` | `` |

### `address/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (/address/delete/:id) | `CD` | POST | `/address/delete/:id` | `UserAddress/delete` | `` |

### `address/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/address/detail/:id) | `CR` | GET | `/address/detail/:id` | `UserAddress/detail` | `` |

### `address/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 地址 | `CR` | GET | `/address/lst` | `UserAddress/lst` | `` |

### `address/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| edit Default (/address/update/:id) | `CU` | POST | `/address/update/:id` | `UserAddress/editDefault` | `` |

### `again`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| again (/again) | `C` | POST | `/again` | `StoreCart/again` | `` |

### `agreement`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 获取协议内容 | `R` | GET | `agreement/:key` | `admin.system.Cache/getAgree` | `` |

### `agreement_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 获取协议列表 | `R` | GET | `agreement_lst` | `admin.system.Cache/getKeyLst` | `` |

### `ajcaptcha`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 滑块验证码 | `R` | GET | `ajcaptcha` | `api.Auth/ajcaptcha` | `` |

### `ajcheck`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| ajcheck (ajcheck) | `U` | POST | `ajcheck` | `api.Auth/ajcheck` | `` |

### `appVersion`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| app Version (appVersion) | `R` | GET | `appVersion` | `api.Common/appVersion` | `` |

### `apply`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| refund (apply/:id) | `D` | POST | `apply/:id` | `/refund` | `` |

### `assist`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 助力 | `R` | GET | `/assist` | `/assist` | `` |

### `assist/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (/assist/create/:id) | `CU` | POST | `/assist/create/:id` | `StoreProductAssistSet/create` | `` |

### `assist/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/assist/detail/:id) | `R` | GET | `/assist/detail/:id` | `StoreProductAssistSet/detail` | `` |

### `assist/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 助力 | `R` | GET | `/assist/lst` | `StoreProductAssist/lst` | `` |

### `assist/set`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| set (/assist/set/:id) | `U` | POST | `/assist/set/:id` | `StoreProductAssistSet/set` | `` |
| delete (/assist/set/delete/:id) | `D` | POST | `/assist/set/delete/:id` | `StoreProductAssistSet/delete` | `` |
| lst (/assist/set/lst) | `R` | GET | `/assist/set/lst` | `StoreProductAssistSet/lst` | `` |

### `assist/share`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| share Num (/assist/share/:id) | `R` | GET | `/assist/share/:id` | `StoreProductAssistSet/shareNum` | `` |

### `assist/user`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| user List (/assist/user/:id) | `R` | GET | `/assist/user/:id` | `StoreProductAssistSet/userList` | `` |

### `auth`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 登录 | `C` | POST | `auth` | `api.Auth/authLogin` | `` |

### `auth/app`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| app授权 | `C` | POST | `auth/app` | `api.Auth/appAuth` | `` |

### `auth/apple`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| apple授权 | `C` | POST | `auth/apple` | `api.Auth/appleAuth` | `` |

### `auth/login`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 登录 | `C` | POST | `auth/login` | `api.Auth/login` | `` |

### `auth/mp`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 小程序授权 | `C` | POST | `auth/mp` | `api.Auth/mpAuth` | `` |

### `auth/mp_login_type`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 获取小程序登录是否需绑定手机号处理 | `C` | POST | `auth/mp_login_type` | `api.Auth/mpLoginType` | `` |

### `auth/mp_phone`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 小程序手机号注册 | `C` | POST | `auth/mp_phone` | `api.Auth/mpPhone` | `` |

### `auth/register`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 注册 | `C` | POST | `auth/register` | `api.Auth/register` | `` |

### `auth/smslogin`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 短信登录 | `C` | POST | `auth/smslogin` | `api.Auth/smsLogin` | `` |

### `auth/verify`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 验证码 | `U` | POST | `auth/verify` | `api.Auth/verify` | `` |

### `auth/wechat`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 微信授权 | `R` | GET | `auth/wechat` | `api.Auth/auth` | `` |

### `back_goods`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| back_goods (back_goods/:id) | `C` | POST | `back_goods/:id` | `/back_goods` | `` |

### `bag`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 礼包 product/spu/bag | `R` | GET | `/bag` | `StoreSpu/bag` | `` |

### `bag/explain`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Bag Explain (bag/explain) | `R` | GET | `bag/explain` | `StoreProduct/getBagExplain` | `` |

### `bag/recommend`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 礼包推荐  product/spu/bag/recommend | `R` | GET | `/bag/recommend` | `StoreSpu/bagRecommend` | `` |

### `batchCreate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Create (/batchCreate) | `C` | POST | `/batchCreate` | `StoreCart/batchCreate` | `` |

### `batch_product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Product (batch_product/:id) | `R` | GET | `batch_product/:id` | `/batchProduct` | `` |

### `bill`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 余额记录 | `R` | GET | `bill` | `User/bill` | `` |

### `binding`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 绑定手机号 | `U` | POST | `binding` | `User/binding` | `` |

### `brand/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (brand/lst) | `R` | GET | `brand/lst` | `StoreBrand/lst` | `` |

### `broadcast`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 直播 | `R` | GET | `/broadcast` | `/broadcast` | `` |

### `brokerage/all`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| brokerage_all (brokerage/all) | `R` | GET | `brokerage/all` | `User/brokerage_all` | `` |

### `brokerage/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| brokerage_info (brokerage/info) | `R` | GET | `brokerage/info` | `User/brokerage_info` | `` |

### `brokerage/notice`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| notice (brokerage/notice) | `R` | GET | `brokerage/notice` | `User/notice` | `` |

### `brokerage_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 佣金记录 | `R` | GET | `brokerage_list` | `User/brokerage_list` | `` |

### `brokerage_top`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 佣金排行榜 | `R` | GET | `brokerage_top` | `User/brokerage_top` | `` |

### `cancel`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cancel Group Order (cancel/:id) | `D` | POST | `cancel/:id` | `/cancelGroupOrder` | `` |
| cancel (cancel/:id) | `D` | POST | `cancel/:id` | `/cancel` | `` |

### `captcha`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 图片验证码 | `R` | GET | `captcha` | `api.Auth/getCaptcha` | `` |

### `cashier_order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Cashier Order (cashier_order/:id) | `R` | GET | `cashier_order/:id` | `/getCashierOrder` | `` |

### `cate_hot`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cate Hot List (/cate_hot) | `R` | GET | `/cate_hot` | `StoreProduct/cateHotList` | `` |

### `category`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 二级分类 | `R` | GET | `/category` | `/category` | `` |
| children (category) | `R` | GET | `category` | `StoreCategory/children` | `` |

### `category/hotranking`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cate Hot Ranking (category/hotranking) | `R` | GET | `category/hotranking` | `StoreCategory/cateHotRanking` | `` |

### `category/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (/category/lst) | `R` | GET | `/category/lst` | `ArticleCategory/lst` | `` |
| category List (/category/lst/:id) | `R` | GET | `/category/lst/:id` | `Merchant/categoryList` | `` |
| 分类&话题 | `R` | GET | `category/lst` | `CommunityCategory/lst` | `` |
| lst (category/lst) | `R` | GET | `category/lst` | `StoreCategory/lst` | `` |

### `change`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| change (/change/:id) | `U` | POST | `/change/:id` | `StoreCart/change` | `` |

### `change/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update Base Info (change/info) | `RU` | POST | `change/info` | `User/updateBaseInfo` | `` |

### `change/password`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| change Password (change/password) | `U` | POST | `change/password` | `User/changePassword` | `` |

### `change/phone`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 修改信息 | `U` | POST | `change/phone` | `User/changePhone` | `` |

### `check`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| check (/check/:id) ✅真实（路由名误写 checkCerate） | `U` | POST | `/check/:id` | `StoreCart/check` | `` |
| v3 Check Order (check) ⛔死路由 | `U` | POST | `check` | `/v3CheckOrder` | `` |
| v2 Check Order (check) ✅真实 | `U` | POST | `check` | `/v2CheckOrder` | `` |
| check Order (check) ⛔死路由 | `U` | POST | `check` | `/checkOrder` | `` |
| 积分 befor Check ✅真实 | `U` | POST | `check` | `PointsOrder/beforCheck` | `` |

### `clear`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| clear (/clear) | `D` | POST | `/clear` | `StoreCart/clear` | `` |

### `command/copy`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 复制口令 | `R` | GET | `command/copy` | `api.Common/getCommand` | `` |

### `common/base64`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 图片转 base64 | `R` | POST | `common/base64` | `api.Common/get_image_base64` | `` |

### `common/commuunity`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 社区热门搜索 | `R` | GET | `common/commuunity/hot_keyword` | `api.Common/hotKeyword` | `` |

### `common/express`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 快递公司 | `R` | GET | `common/express` | `api.Common/express` | `` |

### `common/feedback_type`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户反馈类型 | `R` | GET | `common/feedback_type` | `api.user.FeedBackCategory/lst` | `` |

### `common/home`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 首页数据 | `R` | GET | `common/home` | `api.Common/home` | `` |

### `common/hot_banner`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 推荐页 banner | `R` | GET | `common/hot_banner/:type` | `api.Common/hotBanner` | `` |

### `common/hot_keyword`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 热门搜索 | `R` | GET | `common/hot_keyword` | `api.Common/hotKeyword` | `` |

### `common/menus`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 个人中心菜单 | `R` | GET | `common/menus` | `api.Common/menus` | `` |

### `common/pay_key`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 获取支付宝支付链接 | `CR` | GET | `common/pay_key/:key` | `api.Common/pay_key` | `` |

### `common/recharge_quota`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 充值赠送 | `R` | GET | `common/recharge_quota` | `api.Common/userRechargeQuota` | `` |

### `common/refund_message`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 退款原因 | `R` | GET | `common/refund_message` | `api.Common/refundMessage` | `` |

### `common/visit`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 浏览 | `C` | POST | `common/visit` | `api.Common/visit` | `` |

### `community`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 社区 | `R` | GET | `/community` | `/community` | `` |

### `compute`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| compute (compute) | `C` | POST | `compute` | `/compute` | `` |

### `config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 公共配置 | `R` | GET | `config` | `api.Common/config` | `` |

### `copy`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 复制口令 | `R` | GET | `/copy` | `StoreSpu/copy` | `` |

### `copyright`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| copyright (copyright) | `R` | GET | `copyright` | `api.Common/copyright` | `` |

### `count`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cart Count (/count) | `R` | GET | `/count` | `StoreCart/cartCount` | `` |

### `coupon`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 优惠券 | `R` | GET | `/coupon` | `/coupon` | `` |

### `coupon_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| svip Coupon (coupon_lst) | `R` | GET | `coupon_lst` | `/svipCoupon` | `` |

### `coupon_product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 优惠券商品列表 | `R` | GET | `/coupon_product` | `StoreSpu/getProductByCoupon` | `` |

### `coupon_receive`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| receive Coupon (coupon_receive/:id) | `C` | POST | `coupon_receive/:id` | `/receiveCoupon` | `` |

### `create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (/create) | `C` | POST | `/create` | `StoreCart/create` | `` |
| create (/create) | `C` | POST | `/create` | `Community/create` | `` |
| create (/create/:id) | `C` | POST | `/create/:id` | `/create` | `` |
| v3 Create Order (create) ⛔死路由 | `C` | POST | `create` | `/v3CreateOrder` | `` |
| v2 Create Order (create) ✅真实 | `C` | POST | `create` | `/v2CreateOrder` | `` |
| create Order (create) ⛔死路由 | `C` | POST | `create` | `/createOrder` | `` |
| 积分 create Order ✅真实 | `C` | POST | `create` | `PointsOrder/createOrder` | `` |
| create (create) | `C` | POST | `create` | `/create` | `` |

### `del`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| del (del/:id) | `D` | POST | `del/:id` | `/del` | `` |
| del (del/:id) | `D` | POST | `del/:id` | `/del` | `` |

### `deleate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| del (deleate/:id) | `D` | POST | `deleate/:id` | `PointsOrder/del` | `` |

### `delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Delete (/delete) | `D` | POST | `/delete` | `StoreCart/batchDelete` | `` |
| delete (/delete/:id) | `D` | POST | `/delete/:id` | `Community/delete` | `` |

### `delivery`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Order Delivery (delivery/:id) | `R` | GET | `delivery/:id` | `/getOrderDelivery` | `` |

### `deliverySetings`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Config (deliverySetings) | `R` | GET | `deliverySetings` | `/deliveryConfig` | `` |

### `deliveryStation/list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Station List (deliveryStation/list) | `R` | GET | `deliveryStation/list` | `/deliveryStationList` | `` |

### `deliveryTrack`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Track (deliveryTrack/:id) | `R` | GET | `deliveryTrack/:id` | `/deliveryTrack` | `` |

### `detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/detail/:id) | `R` | GET | `/detail/:id` | `/detail` | `` |
| detail (/detail/:id) | `R` | GET | `/detail/:id` | `Merchant/detail` | `` |
| detail (/detail/:id) | `R` | GET | `/detail/:id` | `/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `PointsOrder/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `StoreProduct/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `Article/detail` | `` |
| detail (detail/:id) | `R` | GET | `detail/:id` | `/detail` | `` |

### `detail/0`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| system Detail (/detail/0) | `R` | GET | `/detail/0` | `Merchant/systemDetail` | `` |

### `diy`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| diy (diy) | `R` | GET | `diy` | `api.Common/diy` | `` |

### `excel/download`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| download (excel/download/:id) | `R` | GET | `excel/download/:id` | `merchant.store.order.Order/download` | `` |

### `express`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| express (express/:id) | `R` | GET | `express/:id` | `/express` | `` |
| express (express/:id) | `C` | POST | `express/:id` | `/express` | `` |

### `extract/banklst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| bank Lst (/extract/banklst) | `R` | GET | `/extract/banklst` | `UserExtract/bankLst` | `` |

### `extract/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (/extract/create) | `C` | POST | `/extract/create` | `UserExtract/create` | `` |

### `extract/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/extract/detail/:id) | `R` | GET | `/extract/detail/:id` | `UserExtract/detail` | `` |

### `extract/history_bank`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| history Bank (/extract/history_bank) | `R` | GET | `/extract/history_bank` | `UserExtract/historyBank` | `` |

### `extract/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 提现 | `R` | GET | `/extract/lst` | `UserExtract/lst` | `` |

### `fab`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 悬浮按钮 | `R` | GET | `/fab` | `/fab` | `` |

### `fans`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| set Focus (fans/:id) | `U` | POST | `fans/:id` | `Community/setFocus` | `` |

### `fans/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get User Fans (fans/lst) | `R` | GET | `fans/lst` | `Community/getUserFans` | `` |

### `feedback`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 反馈 | `C` | POST | `/feedback` | `Feedback/feedback` | `` |

### `feedback/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/feedback/detail/:id) | `R` | GET | `/feedback/detail/:id` | `Feedback/detail` | `` |

### `feedback/list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| feedback List (/feedback/list) | `R` | GET | `/feedback/list` | `Feedback/feedbackList` | `` |

### `fields/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (fields/delete) | `D` | DELETE | `fields/delete` | `UserFields/delete` | `` |

### `fields/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户表单数据操作 | `R` | GET | `fields/info` | `UserFields/info` | `` |

### `fields/save`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| save (fields/save) | `C` | POST | `fields/save` | `UserFields/save` | `` |

### `focus/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get User Focus (focus/lst) | `R` | GET | `focus/lst` | `Community/getUserFocus` | `` |

### `focuslst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| focuslst (/focuslst) | `R` | GET | `/focuslst` | `Community/focuslst` | `` |

### `get`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get (/get/:id) | `R` | GET | `/get/:id` | `StoreSpu/get` | `` |

### `getVersion`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Version (getVersion) | `O` | ANY | `getVersion` | `api.Common/getVersion` | `getVersion` |

### `get_attr_value`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Attr Value (/get_attr_value/:id) | `R` | GET | `/get_attr_value/:id` | `StoreProduct/getAttrValue` | `` |

### `get_hot_ranking`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 热卖排行 | `R` | GET | `/get_hot_ranking` | `StoreSpu/getHotRanking` | `` |

### `get_spec`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品列表获取商品规格 | `R` | GET | `/get_spec/:id` | `StoreProduct/getSpec` | `` |

### `getlst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get List (getlst) | `R` | GET | `getlst` | `api.store.product.StoreCoupon/getList` | `` |

### `good_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Good List (good_list/:id) | `R` | GET | `good_list/:id` | `StoreProduct/getGoodList` | `` |

### `group`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 拼团 | `R` | GET | `/group` | `/group` | `` |

### `group/category`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| category (group/category) | `R` | GET | `group/category` | `StoreProductGroup/category` | `` |

### `group/count`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| user Count (group/count) | `R` | GET | `group/count` | `StoreProductGroup/userCount` | `` |

### `group/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (group/detail/:id) | `R` | GET | `group/detail/:id` | `StoreProductGroup/detail` | `` |

### `group/get`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Buying (group/get/:id) | `CR` | GET | `group/get/:id` | `StoreProductGroup/groupBuying` | `` |

### `group/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 拼团 | `R` | GET | `group/lst` | `StoreProductGroup/lst` | `` |

### `group_order_detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Order Detail (group_order_detail/:id) | `R` | GET | `group_order_detail/:id` | `/groupOrderDetail` | `` |

### `group_order_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Order List (group_order_list) | `R` | GET | `group_order_list` | `/groupOrderList` | `` |

### `guarantee`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| guarantee Template (/guarantee/:id) | `R` | GET | `/guarantee/:id` | `StoreProduct/guaranteeTemplate` | `` |

### `has_service`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| has Service (has_service/:id) | `R` | GET | `has_service/:id` | `api.store.service.Service/hasService` | `` |

### `hist_product/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| history List (hist_product/lst) | `R` | GET | `hist_product/lst` | `Community/historyList` | `` |

### `history`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 浏览记录 | `R` | GET | `history` | `UserHistory/lst` | `` |

### `history/batch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete History Batch (history/batch/delete) | `RD` | POST | `history/batch/delete` | `UserHistory/deleteHistoryBatch` | `` |

### `history/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete History (history/delete/:id) | `RD` | POST | `history/delete/:id` | `UserHistory/deleteHistory` | `` |

### `home`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| home (home) | `R` | GET | `home` | `/home` | `` |

### `hot`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| hot (/hot) | `R` | GET | `/hot` | `BroadcastRoom/hot` | `` |
| 热门 product/spu/hot/:type | `R` | GET | `/hot/:type` | `StoreSpu/hot` | `` |

### `hot_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品热搜 | `R` | GET | `/hot_lst` | `StoreProduct/getHotList` | `` |

### `hot_top`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Hot Top (/hot_top) | `R` | GET | `/hot_top` | `StoreProduct/getHotTop` | `` |
| 热卖排行 | `R` | GET | `/hot_top` | `/hot_top` | `` |

### `increase_take`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| set Increase Take (/increase_take) | `U` | POST | `/increase_take` | `StoreProduct/setIncreaseTake` | `` |

### `info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Form Info (/info/:form_id) | `R` | GET | `/info/:form_id` | `/getFormInfo` | `` |
| activity Info (info/:id) | `R` | GET | `info/:id` | `api.Common/activityInfo` | `` |

### `integral/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 积分 | `R` | GET | `integral/info` | `User/integralInfo` | `` |

### `integral/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| integral List (integral/lst) | `R` | GET | `integral/lst` | `User/integralList` | `` |

### `intention/business`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| business (intention/business) | `R` | GET | `intention/business` | `api.store.merchant.MerchantIntention/business` | `` |

### `intention/cate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cate Lst (intention/cate) | `R` | GET | `intention/cate` | `api.store.merchant.MerchantIntention/cateLst` | `` |

### `intention/circles`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| circles (intention/circles) | `R` | GET | `intention/circles` | `api.store.merchant.MerchantIntention/circles` | `` |

### `intention/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 申请商户 | `C` | POST | `intention/create` | `api.store.merchant.MerchantIntention/create` | `` |

### `intention/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (intention/detail/:id) | `R` | GET | `intention/detail/:id` | `api.store.merchant.MerchantIntention/detail` | `` |

### `intention/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 申请商户 | `R` | GET | `intention/lst` | `api.store.merchant.MerchantIntention/lst` | `` |

### `intention/type`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| type Lst (intention/type) | `R` | GET | `intention/type` | `api.store.merchant.MerchantIntention/typeLst` | `` |

### `intention/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (intention/update/:id) | `U` | POST | `intention/update/:id` | `api.store.merchant.MerchantIntention/update` | `` |

### `labels`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 标签获取数据 | `R` | GET | `/labels` | `StoreSpu/labelsLst` | `` |

### `lbs/address`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 位置信息转经纬度 | `CR` | GET | `lbs/address` | `api.Common/lbs_address` | `` |

### `lbs/geocoder`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 经纬度转位置信息 | `R` | GET | `lbs/geocoder` | `api.Common/lbs_geocoder` | `` |

### `list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| list (/list) | `R` | GET | `/list` | `Article/list` | `` |
| lst (list) | `R` | GET | `list` | `/lst` | `` |
| lst (list) | `R` | GET | `list` | `/lst` | `` |
| list (list) | `R` | GET | `list` | `/list` | `` |
| lst (list) | `R` | GET | `list` | `api.store.product.StoreCoupon/lst` | `` |

### `local`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| local Lst (/local) | `R` | GET | `/local` | `Merchant/localLst` | `` |
| 本地生活商品 | `R` | GET | `/local/:id` | `StoreSpu/local` | `` |

### `localDetail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| local Detail (/localDetail/:id) | `R` | GET | `/localDetail/:id` | `Merchant/localDetail` | `` |

### `logout`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 退出登录 | `C` | POST | `logout` | `api.Auth/logout` | `` |

### `lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (/lst) | `R` | GET | `/lst` | `StoreCart/lst` | `` |
| lst (/lst) | `R` | GET | `/lst` | `/lst` | `` |
| 社区文章列表 | `R` | GET | `/lst` | `Community/lst` | `` |
| 商品 product/spu/lst | `R` | GET | `/lst` | `StoreSpu/lst` | `` |
| lst (/lst) | `R` | GET | `/lst` | `BroadcastRoom/lst` | `` |
| lst (/lst) | `R` | GET | `/lst` | `Merchant/lst` | `` |
| lst (/lst) | `R` | GET | `/lst` | `/lst` | `` |
| lst (/lst/:cid) | `R` | GET | `/lst/:cid` | `Article/lst` | `` |
| lst (lst) | `R` | GET | `lst` | `PointsOrder/lst` | `` |
| lst (lst) | `R` | GET | `lst` | `/lst` | `` |
| lst (lst) | `R` | GET | `lst` | `/lst` | `` |
| activity Lst (lst/:id) | `R` | GET | `lst/:id` | `api.Common/activityLst` | `` |

### `member/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| member Info (member/info) | `R` | GET | `member/info` | `User/memberInfo` | `` |

### `member/log`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Member Value (member/log) | `R` | GET | `member/log` | `Member/getMemberValue` | `` |

### `merchant`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商户商品  product/spu/merchant/:id | `R` | GET | `/merchant/:id` | `StoreSpu/merProductLst` | `` |

### `micro`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| micro (micro) | `R` | GET | `micro` | `api.Common/micro` | `getVersion` |

### `mp/binding`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 小程序获取手机号 | `U` | POST | `mp/binding` | `User/mpPhone` | `` |

### `navigation`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Navigation (navigation) | `R` | GET | `navigation` | `api.Common/getNavigation` | `` |

### `new_people`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| new People (new_people) | `R` | GET | `new_people` | `api.store.product.StoreCoupon/newPeople` | `` |

### `notice`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 微信支付回调:微信 服务商 收付通 | `O` | ANY | `notice/:type` | `api.Common/notify` | `wechatNotify` |

### `notice/callback`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Notify (notice/callback) | `O` | ANY | `notice/callback` | `api.Common/deliveryNotify` | `mchNotify` |

### `notice/mchNotify`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 企业付款到零钱回调 | `O` | ANY | `notice/mchNotify/:type` | `api.Common/mchNotify` | `mchNotify` |

### `notice/pay`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 支付宝支付回调 | `C` | ANY | `notice/pay/alipay` | `api.Common/alipayNotify` | `alipayNotify` |

### `number`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| number (number) | `R` | GET | `number` | `/number` | `` |

### `open_screen`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| open_screen (open_screen) | `R` | GET | `open_screen` | `api.Common/open_screen` | `` |

### `options`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| options (/options) | `R` | GET | `/options` | `/options` | `` |

### `order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Spu By Order (order/:id) | `R` | GET | `order/:id` | `Community/getSpuByOrder` | `` |

### `order_call_back`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| call Back Notify (order_call_back) | `O` | ANY | `order_call_back` | `api.Common/callBackNotify` | `mchNotify` |

### `params`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品参数 | `R` | GET | `/params` | `StoreParams/select` | `` |

### `params_value`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品参数值 | `R` | GET | `/params_value/:id` | `StoreParams/getValue` | `` |

### `pay`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Order Pay (pay/:id) ✅真实 | `C` | POST | `pay/:id` | `/groupOrderPay` | `` |
| pay (pay/:id) ✅真实（预售等） | `C` | POST | `pay/:id` | `/pay` | `` |
| order Pay (pay/:id) ⛔死路由 | `C` | POST | `pay/:id` | `PointsOrder/orderPay` | `` |

### `pay/config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| pay Config (pay/config) | `CR` | GET | `pay/config` | `api.store.order.StoreOrder/payConfig` | `` |

### `pay_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 价格列表 | `CR` | GET | `pay_lst` | `/getTypeLst` | `` |

### `pay_product/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| pay List (pay_product/lst) | `CR` | GET | `pay_product/lst` | `Community/payList` | `` |

### `platform`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| platform Intervene (platform/:id) | `C` | POST | `platform/:id` | `/platformIntervene` | `` |

### `points/pay`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Order Pay (points/pay/:id) | `C` | POST | `points/pay/:id` | `/groupOrderPay` | `` |

### `presell`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 预售 | `R` | GET | `/presell` | `/presell` | `` |

### `presell/agree`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 预售协议 | `R` | GET | `presell/agree` | `StoreProductPresell/getAgree` | `` |

### `presell/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/presell/detail/:id) | `R` | GET | `/presell/detail/:id` | `StoreProductPresell/detail` | `` |

### `presell/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 预售 | `R` | GET | `/presell/lst` | `StoreProductPresell/lst` | `` |

### `preview`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| preview (/preview) | `R` | GET | `/preview` | `StoreProduct/preview` | `` |

### `price_rule`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| price Rule (/price_rule/:id) | `R` | GET | `/price_rule/:id` | `StoreProduct/priceRule` | `` |

### `product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| coupon (product) | `R` | GET | `product` | `api.store.product.StoreCoupon/coupon` | `` |
| product (product/:id) | `R` | GET | `product/:id` | `/product` | `` |
| product (product/:id) | `R` | GET | `product/:id` | `/product` | `` |

### `product/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| product List (/product/lst/:id) | `R` | GET | `/product/lst/:id` | `Merchant/productList` | `` |

### `productCategory`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 平台商品分类 | `R` | GET | `/productCategory` | `/productCategory` | `` |

### `product_detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 普通商品详情 | `R` | GET | `/product_detail` | `/productDetail` | `` |

### `product_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| svip Product List (product_lst) | `R` | GET | `product_lst` | `/svipProductList` | `` |

### `qrcode`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| qrcode (/qrcode/:id) | `R` | GET | `/qrcode/:id` | `StoreProduct/qrcode` | `` |
| qrcode (/qrcode/:id) | `R` | GET | `/qrcode/:id` | `Merchant/qrcode` | `` |
| qrcode (qrcode/:id) | `R` | GET | `qrcode/:id` | `Community/qrcode` | `` |

### `receipt`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create Receipt (receipt/:id) | `C` | POST | `receipt/:id` | `/createReceipt` | `` |

### `receipt/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 发票 | `C` | POST | `receipt/create` | `UserReceipt/create` | `` |

### `receipt/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (receipt/delete/:id) | `D` | POST | `receipt/delete/:id` | `UserReceipt/delete` | `` |

### `receipt/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (receipt/detail/:id) | `R` | GET | `receipt/detail/:id` | `UserReceipt/detail` | `` |

### `receipt/is_default`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| is Default (receipt/is_default/:id) | `C` | POST | `receipt/is_default/:id` | `UserReceipt/isDefault` | `` |

### `receipt/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (receipt/lst) | `R` | GET | `receipt/lst` | `UserReceipt/lst` | `` |

### `receipt/order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order (receipt/order) | `R` | GET | `receipt/order` | `UserReceipt/order` | `` |
| order Detail (receipt/order/:id) | `R` | GET | `receipt/order/:id` | `UserReceipt/orderDetail` | `` |

### `receipt/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (receipt/update/:id) | `U` | POST | `receipt/update/:id` | `UserReceipt/update` | `` |

### `receive`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| receive Coupon (receive/:id) | `C` | POST | `receive/:id` | `api.store.product.StoreCoupon/receiveCoupon` | `` |

### `recharge`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 充值 | `C` | POST | `/recharge` | `UserRecharge/recharge` | `` |

### `recharge/brokerage`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| brokerage (/recharge/brokerage) | `C` | POST | `/recharge/brokerage` | `UserRecharge/brokerage` | `` |

### `recommend`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 推荐 product/spu/recommend | `R` | GET | `/recommend` | `StoreSpu/recommend` | `` |
| recommend List (/recommend) | `R` | GET | `/recommend` | `/recommendList` | `` |

### `recommendProduct`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品列表获取推荐商品 | `R` | GET | `/recommendProduct` | `StoreProduct/recommendProduct` | `` |

### `rela_product/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| relation List (rela_product/lst) | `R` | GET | `rela_product/lst` | `Community/relationList` | `` |

### `relation/batch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Create (/relation/batch/create) | `C` | POST | `/relation/batch/create` | `UserRelation/batchCreate` | `` |
| batch Delete (/relation/batch/delete) | `D` | POST | `/relation/batch/delete` | `UserRelation/batchDelete` | `` |

### `relation/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (/relation/create) | `C` | POST | `/relation/create` | `UserRelation/create` | `` |

### `relation/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (/relation/delete) | `D` | POST | `/relation/delete` | `UserRelation/delete` | `` |

### `relation/merchant`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| merchant List (/relation/merchant/lst) | `R` | GET | `/relation/merchant/lst` | `UserRelation/merchantList` | `` |

### `relation/product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 收藏 | `R` | GET | `/relation/product/lst` | `UserRelation/productList` | `` |

### `reply`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (/:id/reply) | `R` | GET | `/:id/reply` | `CommunityReply/lst` | `` |

### `reply/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (reply/create/:id) | `CU` | POST | `reply/create/:id` | `CommunityReply/create` | `` |

### `reply/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (/reply/lst/:id) | `R` | GET | `/reply/lst/:id` | `StoreReply/lst` | `` |

### `reply/start`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| start (reply/start/:id) | `U` | POST | `reply/start/:id` | `CommunityReply/start` | `` |

### `reservation/checkRange`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| check Range (/reservation/checkRange) | `U` | POST | `/reservation/checkRange` | `StorePrdouctReservation/checkRange` | `` |

### `reservation/getDay`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| show Day (/reservation/getDay/:id) | `R` | GET | `/reservation/getDay/:id` | `StorePrdouctReservation/showDay` | `` |

### `reservation/getMonth`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| show Month (/reservation/getMonth/:id) | `R` | GET | `/reservation/getMonth/:id` | `StorePrdouctReservation/showMonth` | `` |

### `revoke`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| revoke (revoke/:id) | `C` | POST | `revoke/:id` | `/revoke` | `` |

### `route/list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 可以查看所有路由的接口 - 开发者使用 | `R` | GET | `route/list` | `api.Route/list` | `` |

### `scan_upload/image`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| scan Upload Image (scan_upload/image/:field/:token) | `C` | POST | `scan_upload/image/:field/:token` | `api.Common/scanUploadImage` | `` |

### `scope`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| points_mall_scope (scope) | `R` | GET | `scope` | `/points_mall_scope` | `` |

### `script`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| script (script) | `R` | GET | `script` | `api.Common/script` | `` |

### `seckill`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 秒杀 | `R` | GET | `/seckill` | `/seckill` | `` |

### `seckill/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (seckill/detail/:id) | `R` | GET | `seckill/detail/:id` | `StoreProductSeckill/detail` | `` |

### `seckill/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (seckill/lst) | `R` | GET | `seckill/lst` | `StoreProductSeckill/lst` | `` |

### `seckill/select`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| select (seckill/select) | `R` | GET | `seckill/select` | `StoreProductSeckill/select` | `` |

### `self/cancel`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| cancel Order (self/cancel/:id) | `D` | POST | `self/cancel/:id` | `/cancelOrder` | `` |

### `service/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 获取商户基本信息 | `R` | GET | `service/info/:id` | `api.store.service.Service/merchantInfo` | `` |

### `services`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 客服列表 | `R` | GET | `services` | `User/services` | `` |

### `share_posters`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Share Posters (/share_posters/:id) | `R` | GET | `/share_posters/:id` | `/getSharePosters` | `` |

### `show`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| show (/show/:id) | `R` | GET | `/show/:id` | `/show` | `` |
| 详情 | `R` | GET | `/show/:id` | `Community/show` | `` |
| show (show/:id) | `R` | GET | `show/:id` | `StoreProduct/show` | `` |

### `sign/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (sign/create) | `C` | POST | `sign/create` | `UserSign/create` | `` |

### `sign/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 签到数据信息获取 | `R` | GET | `sign/info` | `UserSign/info` | `` |

### `sign/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 签到 | `R` | GET | `sign/lst` | `UserSign/lst` | `` |

### `sign/month`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| month (sign/month) | `R` | GET | `sign/month` | `UserSign/month` | `` |

### `spread_image`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 分销海报 | `R` | GET | `/spread_image` | `User/spread_image` | `` |

### `spread_info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| spread_info (spread_info) | `R` | GET | `spread_info` | `User/spread_info` | `` |

### `spread_level`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| spread_info (spread_level) | `R` | GET | `spread_level` | `User/spread_info` | `` |

### `spread_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 推广人列表 | `R` | GET | `/spread_list` | `User/spread_list` | `` |

### `spread_order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 推广人订单 | `R` | GET | `spread_order` | `User/spread_order` | `` |

### `spread_top`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 推广人排行榜 | `R` | GET | `spread_top` | `User/spread_top` | `` |

### `spu`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品列表 | `R` | GET | `/spu` | `/spu` | `` |

### `start`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| start Community (start/:id) | `C` | POST | `start/:id` | `Community/startCommunity` | `` |

### `start/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get User Start Community (start/lst) | `R` | GET | `start/lst` | `Community/getUserStartCommunity` | `` |

### `status`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| group Order Status (status/:id) | `RU` | GET | `status/:id` | `/groupOrderStatus` | `` |

### `store`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 品牌好店 | `R` | GET | `/store` | `/store` | `` |
| mer Coupon (store/:id) | `R` | GET | `store/:id` | `api.store.product.StoreCoupon/merCoupon` | `` |

### `store/certificate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Mer Certificate (store/certificate/:merId) | `R` | POST | `store/certificate/:merId` | `api.Auth/getMerCertificate` | `` |

### `store/expr`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Export Temp (store/expr/temps) | `R` | GET | `store/expr/temps` | `admin.system.serve.Export/getExportTemp` | `` |

### `store/product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| user Count (store/product/assist/count) | `R` | GET | `store/product/assist/count` | `api.store.product.StoreProductAssist/userCount` | `` |
| cancel (store/product/group/cancel) | `D` | POST | `store/product/group/cancel` | `api.store.product.StoreProductGroup/cancel` | `` |

### `store/test`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| test | `O` | ANY | `store/test` | `api.Test/test` | `` |

### `subscribe`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| subscribe (subscribe) | `R` | GET | `subscribe` | `api.Common/subscribe` | `` |

### `svip/pay`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 付费会员购买 | `C` | POST | `svip/pay/:id` | `api.user.Svip/createOrder` | `` |

### `switch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| switch User (switch) | `U` | POST | `switch` | `User/switchUser` | `` |

### `system/city`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 城市列表 | `R` | GET | `system/city/lst` | `merchant.store.shipping.City/getlist` | `` |

### `take`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| take (take/:id) | `C` | POST | `take/:id` | `/take` | `` |
| take (take/:id) | `C` | POST | `take/:id` | `PointsOrder/take` | `` |

### `update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (/update/:id) | `U` | POST | `/update/:id` | `Community/update` | `` |
| update (update/:id) | `U` | POST | `update/:id` | `/update` | `` |

### `upload/certificate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 入住商户上传证件接口 | `C` | POST | `upload/certificate/:field` | `api.Common/uploadCertificate` | `` |

### `upload/image`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 上传图片 | `C` | POST | `upload/image/:field` | `api.Common/uploadImage` | `` |

### `upload/video`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| upload Video (upload/video) | `C` | POST | `upload/video` | `admin.system.attachment.Attachment/uploadVideo` | `` |

### `user`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户信息 | `R` | GET | `user` | `api.Auth/userInfo` | `` |

### `user/cancel`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 注销用户 | `D` | POST | `user/cancel` | `api.Auth/cancel` | `` |

### `user/change_pwd`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 修改密码 | `U` | POST | `user/change_pwd` | `api.Auth/changePassword` | `` |

### `user/community`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户的文章 | `R` | GET | `/user/community/:id` | `Community/userCommunitylst` | `` |

### `user/community_video`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户的视频 | `R` | GET | `/user/community_video/:id` | `Community/userCommunityVideolst` | `` |

### `user/info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户页 | `R` | GET | `/user/info/:id` | `Community/userInfo` | `` |

### `user/spread`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 绑定推荐人 | `R` | POST | `user/spread` | `api.Auth/spread` | `` |

### `user_info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| svip User Info (user_info) | `R` | GET | `user_info` | `/svipUserInfo` | `` |

### `user_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| user List (/user_lst) | `R` | GET | `/user_lst` | `Community/userList` | `` |

### `v2/spread_image`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| spread_image_v2 (/v2/spread_image) | `R` | GET | `/v2/spread_image` | `User/spread_image_v2` | `` |

### `v2/system`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| city List (v2/system/city) | `R` | GET | `v2/system/city` | `merchant.store.shipping.City/cityList` | `` |
| lst V2 (v2/system/city/lst/:pid) | `R` | GET | `v2/system/city/lst/:pid` | `merchant.store.shipping.City/lstV2` | `` |

### `verify_code`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| verify Code (verify_code/:id) | `R` | GET | `verify_code/:id` | `/verifyCode` | `` |

### `version`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| version (version) | `R` | GET | `version` | `admin.Common/version` | `getVersion` |

### `video_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| video Show (/video_lst) | `R` | GET | `/video_lst` | `Community/videoShow` | `` |

### `wechat/config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 微信配置 | `R` | GET | `wechat/config` | `api.Wechat/jsConfig` | `` |

### `登录用户`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reply (:id) | `U` | POST | `:id` | `/reply` | `` |

## 附录：小程序页面路径（UI 入口）

共 183 个页面。每个页面上的按钮最终应对齐上文 API 操作。

- `pages/guide/index`
- `pages/index/index`
- `pages/order_addcart/order_addcart`
- `pages/user/index`
- `pages/goods_cate/goods_cate`
- `pages/auth/index`
- `pages/error/index`
- `pages/users/retrievePassword/index`
- `pages/users/user_setting/index`
- `pages/users/user_about/index`
- `pages/users/user_info/index`
- `pages/users/user_info_form/index`
- `pages/users/user_nickname/index`
- `pages/users/user_goods_collection/index`
- `pages/users/user_sgin/index`
- `pages/users/user_sgin_list/index`
- `pages/users/user_points_list/index`
- `pages/users/user_money/index`
- `pages/users/user_bill/index`
- `pages/users/user_integral/index`
- `pages/users/user_brokerage/index`
- `pages/users/user_grade/index`
- `pages/users/user_grade_list/index`
- `pages/users/user_coupon/index`
- `pages/users/user_spread_user/index`
- `pages/users/user_spread_code/index`
- `pages/users/user_spread_money/index`
- `pages/users/user_invoice_order/index`
- `pages/users/user_invoice_form/index`
- `pages/users/user_invoice_list/index`
- `pages/users/agreement_rules/index`
- `pages/users/user_spread_money/receiving`
- `pages/users/user_cash/index`
- `pages/users/user_address_list/index`
- `pages/users/user_address/index`
- `pages/users/user_phone/index`
- `pages/users/user_modify_phone/index`
- `pages/users/user_modify_pwd/index`
- `pages/users/user_payment/index`
- `pages/users/user_pwd_edit/index`
- `pages/users/order_payment/index`
- `pages/users/promoter-list/index`
- `pages/users/promoter-order/index`
- `pages/users/promoter_rank/index`
- `pages/users/commission_rank/index`
- `pages/users/order_list/index`
- `pages/users/order_list/search`
- `pages/users/presell_order_list/index`
- `pages/users/goods_logistics/index`
- `pages/users/user_return_list/index`
- `pages/users/goods_return/index`
- `pages/users/login/index`
- `pages/users/wechat_login/index`
- `pages/users/goods_comment_list/index`
- `pages/users/goods_comment_con/index`
- `pages/users/feedback/index`
- `pages/users/feedback/list`
- `pages/users/feedback/detail`
- `pages/users/refund/index`
- `pages/users/refund/confirm`
- `pages/users/refund/detail`
- `pages/users/refund/select`
- `pages/users/refund/goods/index`
- `pages/users/refund/list`
- `pages/users/refund/logistics`
- `pages/users/browsingHistory/index`
- `pages/users/distributor/index`
- `pages/users/privacy/index`
- `pages/users/order_confirm/index`
- `pages/store/home/index`
- `pages/store/home/goods`
- `pages/store/home/life`
- `pages/store/home/special`
- `pages/store/home/coupon`
- `pages/store/goods_cate/index`
- `pages/store/detail/index`
- `pages/store/list/index`
- `pages/store/settled/index`
- `pages/store/applicationRecord/index`
- `pages/store/merchantDetails/index`
- `pages/store/shopStreet/index`
- `pages/store/qualifications/index`
- `pages/admin/storeDiy/index`
- `pages/admin/order/index`
- `pages/admin/orderList/index`
- `pages/admin/orderRefund/index`
- `pages/admin/refundList/index`
- `pages/admin/business/index`
- `pages/admin/orderDetail/index`
- `pages/admin/refundDetail/index`
- `pages/admin/delivery/index`
- `pages/admin/statistics/index`
- `pages/admin/order_cancellation/index`
- `pages/admin/cancellate_result/index`
- `pages/admin/booking/detail`
- `pages/admin/booking/reschedule`
- `pages/goods_details/index`
- `pages/order_pay_status/index`
- `pages/order_details/index`
- `pages/order_details/stay`
- `pages/order_details/delivery`
- `pages/small_page/index`
- `pages/goods_recommend/index`
- `pages/news_list/index`
- `pages/news_details/index`
- `pages/order_pay_back/index`
- `pages/product/list/index`
- `pages/product/list/specs`
- `pages/product/goodsOnSale/index`
- `pages/product/soldOutGoods/index`
- `pages/product/recycleBin/index`
- `pages/product/storeClassification/index`
- `pages/product/storeClassification/addStoreClass`
- `pages/product/addGoods/index`
- `pages/plantGrass/plant_detail/index`
- `pages/plantGrass/plant_release/index`
- `pages/plantGrass/plant_show/index`
- `pages/plantGrass/plant_topic/index`
- `pages/plantGrass/plant_search/index`
- `pages/plantGrass/plant_search_list/index`
- `pages/plantGrass/plant_featured/index`
- `pages/plantGrass/plant_user/index`
- `pages/plantGrass/plant_user_attention/index`
- `pages/plantGrass/plant_user_fans/index`
- `pages/columnGoods/HotNewGoods/index`
- `pages/columnGoods/goods_list/index`
- `pages/columnGoods/goods_coupon_list/index`
- `pages/columnGoods/goods_search/index`
- `pages/columnGoods/goods_search_con/index`
- `pages/chat/customer_list/index`
- `pages/chat/customer_list/chat`
- `pages/chat/customer_login/index`
- `pages/chat/customer_info/index`
- `pages/activity/goods_seckill/index`
- `pages/activity/goods_seckill_details/index`
- `pages/activity/liveBroadcast/index`
- `pages/activity/presell/index`
- `pages/activity/presell_details/index`
- `pages/activity/combination/index`
- `pages/activity/combination_details/index`
- `pages/activity/combination_status/index`
- `pages/activity/assist/index`
- `pages/activity/assist_detail/index`
- `pages/activity/assist_record/index`
- `pages/activity/topic/index`
- `pages/activity/topic_detail/index`
- `pages/activity/lifeService/index`
- `pages/activity/collect_coupons/index`
- `pages/activity/rank/index`
- `pages/activity/registrate_activity/index`
- `pages/activity/registrate_list/index`
- `pages/activity/my_registrate/index`
- `pages/reservation/reservation/index`
- `pages/reservation/reservation_info/index`
- `pages/short_video/nvueSwiper/index`
- `pages/points_mall/index`
- `pages/points_mall/integral_goods_list`
- `pages/points_mall/goods_selection`
- `pages/points_mall/integral_goods_details`
- `pages/points_mall/integral_order`
- `pages/points_mall/exchange_record`
- `pages/points_mall/integral_order_details`
- `pages/annex/web_view/index`
- `pages/annex/vip_paid/index`
- `pages/annex/vip_center/index`
- `pages/annex/vip_clause/index`
- `pages/staff/index`
- `pages/staff/order_list`
- `pages/staff/order_detail`
- `pages/staff/checkin`
- `pages/staff/service_record`
- `pages/delivery/order_list`
- `pages/delivery/order_detail`
- `pages/plant_grass/index`
- `sub-packages/diy/index`
- `pages/agent/form`
- `pages/agent/records`
- `pages/agent/detail`
- `pages/circle/select`
- `pages/circle/search`
- `pages/merchant/apply/index`
- `pages/merchant/apply/records`
- `pages/merchant/apply/detail`
