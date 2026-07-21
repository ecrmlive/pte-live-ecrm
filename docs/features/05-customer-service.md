# 客服端 — 接口级功能点

> 来源：`route/service.php`。

> 每一条路由 = 一个可调用操作（等同按钮/动作）。CRUD 由 HTTP 方法 + 路径/处理器名推断，可人工校正。

## 统计

| 项 | 数量 |
| --- | ---: |
| 操作（路由）总数 | 18 |
| C | 4 |
| R | 12 |
| U | 3 |
| D | 0 |

## 模块 CRUD 覆盖

| 模块 | 操作数 | C | R | U | D |
| --- | ---: | --- | --- | --- | --- |
| `captcha` | 1 | — | ✓ | — | — |
| `config` | 1 | — | ✓ | — | — |
| `history` | 1 | — | ✓ | — | — |
| `info` | 1 | — | ✓ | — | — |
| `login` | 1 | ✓ | — | — | — |
| `login/scan` | 2 | ✓ | — | ✓ | — |
| `logout` | 1 | ✓ | — | — | — |
| `order` | 1 | — | ✓ | — | — |
| `order_express` | 1 | — | ✓ | — | — |
| `order_status` | 1 | — | ✓ | ✓ | — |
| `product` | 1 | — | ✓ | — | — |
| `refund` | 1 | — | ✓ | — | — |
| `refund_express` | 1 | — | ✓ | — | — |
| `upload` | 1 | ✓ | — | — | — |
| `user` | 1 | — | ✓ | — | — |
| `user/lst` | 1 | — | ✓ | — | — |
| `user/mark` | 1 | — | — | ✓ | — |

## 分模块操作明细

### `captcha`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 验证码 | `R` | GET | `captcha` | `service.Login/getCaptcha` | `` |

### `config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| config (config) | `R` | GET | `config` | `service.Common/config` | `` |

### `history`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 聊天记录 | `R` | GET | `history/:uid` | `service.Service/history` | `` |

### `info`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商户信息 | `R` | GET | `info` | `service.Common/info` | `` |

### `login`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 登录 | `C` | POST | `login` | `service.Login/login` | `` |

### `login/scan`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 扫码登录 | `C` | POST | `login/scan` | `service.Login/scanLogin` | `` |
| 登录 | `U` | POST | `login/scan/check` | `service.Login/checkScanLogin` | `` |

### `logout`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 退出登录 | `C` | POST | `logout` | `service.Login/logout` | `` |

### `order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 订单信息 | `R` | GET | `order/:id` | `service.Service/getOrderInfo` | `` |

### `order_express`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 快递 | `R` | GET | `order_express/:id` | `service.Service/orderExpress` | `` |

### `order_status`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order Status (order_status/:id) | `RU` | GET | `order_status/:id` | `service.Service/orderStatus` | `` |

### `product`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品 | `R` | GET | `product/:id` | `service.Service/product` | `` |

### `refund`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 退款单信息 | `R` | GET | `refund/:id` | `service.Service/getRefundOder` | `` |

### `refund_express`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 快递 | `R` | GET | `refund_express/:id` | `service.Service/refundOrderExpress` | `` |

### `upload`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 图片上传 | `C` | POST | `upload/:field` | `service.Service/upload` | `` |

### `user`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户信息 | `R` | GET | `user` | `service.Common/user` | `` |

### `user/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户聊天列表 | `R` | GET | `user/lst` | `service.Service/serviceUserList` | `` |

### `user/mark`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 用户备注 | `U` | POST | `user/mark/:uid` | `service.Service/mark` | `` |
