# 客服端 `/ser/`

> 对照文档。置信度：high=18 stale=0 unresolved=0。先读 [ACCURACY.md](./ACCURACY.md)。

合计 **18** 条。

## `ser/captcha`

### `GET /ser/captcha` — getCaptcha

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Login/getCaptcha`
- 源码：`app/controller/service/Login.php` :: `getCaptcha()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `ser/config`

### `GET /ser/config` — config

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Common/config`
- 源码：`app/controller/service/Common.php` :: `config()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `ser/history`

### `GET /ser/history/:uid` — history

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/history`
- 源码：`app/controller/service/Service.php` :: `history()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `last_id` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/info`

### `GET /ser/info` — info

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Common/info`
- 源码：`app/controller/service/Common.php` :: `info()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/login`

### `POST /ser/login` — login

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Login/login`
- 源码：`app/controller/service/Login.php` :: `login()`
- 请求参数：
- `account` (query/body, 可选)
- `password` (query/body, 可选)
- `key` (query/body, 可选)
- `code` (query/body, 可选)
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `POST /ser/login/scan` — scanLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Login/scanLogin`
- 源码：`app/controller/service/Login.php` :: `scanLogin()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data 对象字段: timeout, key, qrcode | 外层: {status,message,data}

### `POST /ser/login/scan/check` — checkScanLogin

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Login/checkScanLogin`
- 源码：`app/controller/service/Login.php` :: `checkScanLogin()`
- 请求参数：
- `key` (query/body, 可选)
- 返回：统一响应包 {status,message,data?}；详见控制器实现 | 外层: {status,message,data}


## `ser/logout`

### `POST /ser/logout` — logout

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Login/logout`
- 源码：`app/controller/service/Login.php` :: `logout()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `ser/order`

### `GET /ser/order/:id` — getOrderInfo

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/getOrderInfo`
- 源码：`app/controller/service/Service.php` :: `getOrderInfo()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `ser/order_express`

### `GET /ser/order_express/:id` — orderExpress

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/orderExpress`
- 源码：`app/controller/service/Service.php` :: `orderExpress()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/order_status`

### `GET /ser/order_status/:id` — orderStatus

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/orderStatus`
- 源码：`app/controller/service/Service.php` :: `orderStatus()`
- 请求参数：
- `id` (path, 必填) 路径参数
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `date` (query/body, 可选)
- `user_type` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}


## `ser/product`

### `GET /ser/product/:id` — product

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/product`
- 源码：`app/controller/service/Service.php` :: `product()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/refund`

### `GET /ser/refund/:id` — getRefundOder

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/getRefundOder`
- 源码：`app/controller/service/Service.php` :: `getRefundOder()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/refund_express`

### `GET /ser/refund_express/:id` — refundOrderExpress

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/refundOrderExpress`
- 源码：`app/controller/service/Service.php` :: `refundOrderExpress()`
- 请求参数：
- `id` (path, 必填) 路径参数
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}


## `ser/upload`

### `POST /ser/upload/:field` — upload

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/upload`
- 源码：`app/controller/service/Service.php` :: `upload()`
- 请求参数：
- `field` (path, 必填) 路径参数
- 返回：data 对象字段: src | 外层: {status,message,data}


## `ser/user`

### `GET /ser/user` — user

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Common/user`
- 源码：`app/controller/service/Common.php` :: `user()`
- 请求参数：
_（未见显式 param；或见 doc_note）_
- 返回：data: 业务数据对象/数组（见 Repository 返回） | 外层: {status,message,data}

### `GET /ser/user/lst` — serviceUserList

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/serviceUserList`
- 源码：`app/controller/service/Service.php` :: `serviceUserList()`
- 请求参数：
- `page` (query, 可选) 页码
- `limit` (query, 可选) 每页数量
- `keyword` (query/body, 可选)
- 返回：data: 列表结构（通常含 list/count 或分页数据，以实际返回为准） | 外层: {status,message,data}

### `POST /ser/user/mark/:uid` — mark

- 置信度：✅ high
- 说明：控制器方法已校验存在
- 处理器：`service.Service/mark`
- 源码：`app/controller/service/Service.php` :: `mark()`
- 请求参数：
- `uid` (path, 必填) 路径参数
- `mark` (query/body, 可选)
- 返回：data: 见控制器返回；成功时 status=200 | 外层: {status,message,data}

