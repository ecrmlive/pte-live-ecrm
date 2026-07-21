# 商户手机端 / 店员端 — 接口级功能点

> 来源：`route/api/manager.php`（移动端商户管理/店员核销等）。

> 每一条路由 = 一个可调用操作（等同按钮/动作）。CRUD 由 HTTP 方法 + 路径/处理器名推断，可人工校正。

## 统计

| 项 | 数量 |
| --- | ---: |
| 操作（路由）总数 | 95 |
| C | 17 |
| R | 47 |
| U | 26 |
| D | 13 |

## 模块 CRUD 覆盖

| 模块 | 操作数 | C | R | U | D |
| --- | ---: | --- | --- | --- | --- |
| `attr/create` | 1 | ✓ | — | — | — |
| `attr/delete` | 1 | — | — | — | ✓ |
| `attr/detail` | 2 | — | ✓ | — | — |
| `attr/list` | 1 | — | ✓ | — | — |
| `attr/lst` | 1 | — | ✓ | — | — |
| `attr/update` | 1 | — | — | ✓ | — |
| `category/brandlist` | 1 | — | ✓ | — | — |
| `category/create` | 1 | ✓ | — | — | — |
| `category/delete` | 1 | — | — | — | ✓ |
| `category/detail` | 1 | — | ✓ | — | — |
| `category/list` | 1 | — | ✓ | — | — |
| `category/lst` | 1 | — | ✓ | — | — |
| `category/select` | 1 | — | ✓ | — | — |
| `category/status` | 1 | — | — | ✓ | — |
| `category/update` | 1 | — | — | ✓ | — |
| `delivery` | 1 | — | — | — | ✓ |
| `delivery/confirm` | 1 | — | — | ✓ | ✓ |
| `delivery/dispatch` | 1 | — | — | — | ✓ |
| `delivery/options` | 1 | — | ✓ | — | — |
| `delivery/person` | 1 | — | ✓ | — | — |
| `delivery/updateDispatch` | 1 | — | — | ✓ | ✓ |
| `delivery_config` | 1 | — | ✓ | — | — |
| `delivery_options` | 1 | — | ✓ | — | — |
| `dump_temp` | 1 | — | ✓ | — | — |
| `history` | 1 | — | ✓ | — | — |
| `list` | 1 | — | ✓ | — | — |
| `manager` | 1 | — | — | ✓ | — |
| `mark` | 2 | — | — | ✓ | — |
| `mer_form` | 1 | — | ✓ | — | — |
| `mer_history` | 1 | — | ✓ | — | — |
| `offline` | 1 | ✓ | — | — | — |
| `order` | 4 | — | ✓ | — | — |
| `order/check` | 1 | — | — | ✓ | — |
| `order/confirm` | 1 | — | — | ✓ | — |
| `order/dispatch` | 1 | ✓ | — | — | — |
| `order/mark` | 2 | — | — | ✓ | — |
| `order/receive` | 1 | ✓ | — | — | — |
| `order/trace` | 1 | ✓ | — | — | — |
| `order/verifier` | 1 | — | — | ✓ | — |
| `order_list` | 1 | — | ✓ | — | — |
| `order_lst` | 2 | — | ✓ | — | — |
| `order_price` | 1 | — | ✓ | — | — |
| `pay_number` | 1 | ✓ | ✓ | — | — |
| `pay_price` | 1 | ✓ | ✓ | — | — |
| `price` | 1 | ✓ | — | — | — |
| `product/config` | 1 | — | ✓ | — | — |
| `product/create` | 1 | ✓ | — | — | — |
| `product/delete` | 1 | — | — | — | ✓ |
| `product/destory` | 1 | ✓ | — | — | — |
| `product/detail` | 1 | — | ✓ | — | — |
| `product/edit_cate` | 1 | — | — | ✓ | — |
| `product/edit_mer_cate` | 1 | — | — | ✓ | — |
| `product/good` | 1 | — | — | ✓ | — |
| `product/lst` | 1 | — | ✓ | — | — |
| `product/restore` | 1 | ✓ | — | — | — |
| `product/status` | 1 | — | — | ✓ | — |
| `product/title` | 1 | — | ✓ | — | — |
| `product/update` | 1 | — | — | ✓ | — |
| `product/value` | 2 | — | ✓ | ✓ | — |
| `refund/check` | 1 | — | ✓ | — | — |
| `refund/compute` | 1 | — | — | — | ✓ |
| `refund/confirm` | 1 | — | — | ✓ | ✓ |
| `refund/create` | 1 | ✓ | — | — | ✓ |
| `refund/detail` | 1 | — | ✓ | — | — |
| `refund/express` | 1 | — | ✓ | — | — |
| `refund/get` | 1 | — | ✓ | — | — |
| `refund/lst` | 1 | — | ✓ | — | — |
| `refund/mark` | 1 | — | — | ✓ | ✓ |
| `refund/status` | 1 | — | — | ✓ | ✓ |
| `reservation/config` | 1 | — | ✓ | — | — |
| `reservation/staffs` | 1 | — | ✓ | — | — |
| `reservationconfig` | 1 | — | ✓ | — | — |
| `reservationdispatch` | 1 | ✓ | — | — | — |
| `reservationreschedule` | 1 | ✓ | — | — | — |
| `reservationupdateDispatch` | 1 | — | — | ✓ | — |
| `reservationverify` | 1 | — | — | ✓ | — |
| `scan_login` | 1 | ✓ | — | — | — |
| `statistics` | 1 | — | ✓ | — | — |
| `template/create` | 1 | ✓ | — | — | — |
| `template/delete` | 1 | — | — | — | ✓ |
| `template/detail` | 1 | — | ✓ | — | — |
| `template/lst` | 1 | — | ✓ | — | — |
| `template/select` | 1 | — | ✓ | — | — |
| `template/update` | 1 | — | — | ✓ | — |
| `user` | 1 | — | ✓ | — | — |
| `user_list` | 1 | — | ✓ | — | — |
| `verify` | 1 | — | — | ✓ | — |

## 分模块操作明细

### `attr/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (:merId/attr/create) | `C` | POST | `:merId/attr/create` | `StoreProductAttrTemplate/create` | `` |

### `attr/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Delete (:merId/attr/delete) | `D` | POST | `:merId/attr/delete` | `StoreProductAttrTemplate/batchDelete` | `` |

### `attr/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (:merId/attr/detail/:id) | `R` | GET | `:merId/attr/detail/:id` | `StoreProductAttrTemplate/detail` | `` |
| detail (:merId/attr/detail/:id) | `R` | GET | `:merId/attr/detail/:id` | `StoreProductAttrTemplate/detail` | `` |

### `attr/list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| getlist (:merId/attr/list) | `R` | GET | `:merId/attr/list` | `StoreProductAttrTemplate/getlist` | `` |

### `attr/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 品牌管理 | `R` | GET | `:merId/attr/lst` | `StoreProductAttrTemplate/lst` | `` |

### `attr/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (:merId/attr/update/:id) | `U` | POST | `:merId/attr/update/:id` | `StoreProductAttrTemplate/update` | `` |

### `category/brandlist`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| Brand List (:merId/category/brandlist) | `R` | GET | `:merId/category/brandlist` | `StoreCategory/BrandList` | `` |

### `category/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (:merId/category/create) | `C` | POST | `:merId/category/create` | `StoreCategory/create` | `` |

### `category/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (:merId/category/delete/:id) | `D` | POST | `:merId/category/delete/:id` | `StoreCategory/delete` | `` |

### `category/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (:merId/category/detail/:id) | `R` | GET | `:merId/category/detail/:id` | `StoreCategory/detail` | `` |

### `category/list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get List (:merId/category/list) | `R` | GET | `:merId/category/list` | `StoreCategory/getList` | `` |

### `category/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 商品分类 | `R` | GET | `:merId/category/lst` | `StoreCategory/lst` | `` |

### `category/select`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Tree List (:merId/category/select) | `R` | GET | `:merId/category/select` | `StoreCategory/getTreeList` | `` |

### `category/status`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| switch Status (:merId/category/status/:id) | `U` | POST | `:merId/category/status/:id` | `StoreCategory/switchStatus` | `` |

### `category/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (:merId/category/update/:id) | `U` | POST | `:merId/category/update/:id` | `StoreCategory/update` | `` |

### `delivery`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery (:merId/delivery/:id) | `D` | POST | `:merId/delivery/:id` | `/delivery` | `` |

### `delivery/confirm`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Confirm (:merId/delivery/confirm/:id) | `UD` | POST | `:merId/delivery/confirm/:id` | `/deliveryConfirm` | `` |

### `delivery/dispatch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Dispatch (:merId/delivery/dispatch/:id) | `D` | POST | `:merId/delivery/dispatch/:id` | `/deliveryDispatch` | `` |

### `delivery/options`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| options (:merId/delivery/options) | `R` | GET | `:merId/delivery/options` | `/options` | `` |

### `delivery/person`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 同城配送 | `R` | GET | `:merId/delivery/person` | `/deliveryPersonList` | `` |

### `delivery/updateDispatch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delivery Update Dispatch (:merId/delivery/updateDispatch/:id) | `UD` | POST | `:merId/delivery/updateDispatch/:id` | `/deliveryUpdateDispatch` | `` |

### `delivery_config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Delivery Config (:merId/delivery_config) | `R` | GET | `:merId/delivery_config` | `/getDeliveryConfig` | `` |

### `delivery_options`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Delivery Options (:merId/delivery_options) | `R` | GET | `:merId/delivery_options` | `/getDeliveryOptions` | `` |

### `dump_temp`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Form Data (:merId/dump_temp) | `R` | GET | `:merId/dump_temp` | `/getFormData` | `` |

### `history`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| chat History (history/:id) | `R` | GET | `history/:id` | `api.store.service.Service/chatHistory` | `` |

### `list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get List (list) | `R` | GET | `list` | `api.store.service.Service/getList` | `` |

### `manager`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| verify (/:merId/:id) | `U` | POST | `/:merId/:id` | `/verify` | `` |

### `mark`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| mark (:merId/mark/:id) | `U` | POST | `:merId/mark/:id` | `/mark` | `` |
| mark (mark/:merId/:uid) | `U` | POST | `mark/:merId/:uid` | `api.store.service.Service/mark` | `` |

### `mer_form`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Form Data (:merId/mer_form) | `R` | GET | `:merId/mer_form` | `/getFormData` | `` |

### `mer_history`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| service History (mer_history/:merId/:id) | `R` | GET | `mer_history/:merId/:id` | `api.store.service.Service/serviceHistory` | `` |

### `offline`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| offline (:merId/offline/:id) | `C` | POST | `:merId/offline/:id` | `/offline` | `` |

### `order`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (/:merId/order/:id) | `R` | GET | `/:merId/order/:id` | `/detail` | `` |
| order (:merId/order/:id) | `R` | GET | `:merId/order/:id` | `/order` | `` |
| order Detail (order/:id) | `R` | GET | `order/:id` | `/orderDetail` | `` |
| order Detail (order/:id) | `R` | GET | `order/:id` | `/orderDetail` | `deliveryOrderDetail` |

### `order/check`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| check In (order/:id/check) | `U` | POST | `order/:id/check` | `/checkIn` | `` |

### `order/confirm`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| confirm (order/:id/confirm) | `U` | POST | `order/:id/confirm` | `/confirm` | `` |

### `order/dispatch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Dispatch (order/:id/dispatch) | `C` | POST | `order/:id/dispatch` | `/reservationDispatch` | `` |

### `order/mark`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| mark (order/:id/mark) | `U` | POST | `order/:id/mark` | `/mark` | `` |
| mark (order/:id/mark) | `U` | POST | `order/:id/mark` | `/mark` | `` |

### `order/receive`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| receive (order/:id/receive) | `C` | POST | `order/:id/receive` | `/receive` | `` |

### `order/trace`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| add Trace (order/:id/trace) | `C` | POST | `order/:id/trace` | `/addTrace` | `` |

### `order/verifier`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| verify (order/:id/verifier) | `U` | POST | `order/:id/verifier` | `/verify` | `` |

### `order_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order List (:merId/order_list) | `R` | GET | `:merId/order_list` | `/orderList` | `` |

### `order_lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order_lst (order_lst) | `R` | GET | `order_lst` | `/order_lst` | `` |
| order_lst (order_lst) | `R` | GET | `order_lst` | `/order_lst` | `deliveryOrderLst` |

### `order_price`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order Detail (:merId/order_price) | `R` | GET | `:merId/order_price` | `/orderDetail` | `` |

### `pay_number`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| pay Number (:merId/pay_number) | `CR` | GET | `:merId/pay_number` | `/payNumber` | `` |

### `pay_price`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| pay Price (:merId/pay_price) | `CR` | GET | `:merId/pay_price` | `/payPrice` | `` |

### `price`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| price (:merId/price/:id) | `C` | POST | `:merId/price/:id` | `/price` | `` |

### `product/config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| config (:merId/product/config) | `R` | GET | `:merId/product/config` | `StoreProduct/config` | `` |

### `product/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 添加 | `C` | POST | `:merId/product/create` | `StoreProduct/create` | `` |

### `product/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| delete (:merId/product/delete/:id) | `D` | POST | `:merId/product/delete/:id` | `StoreProduct/delete` | `` |

### `product/destory`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| destory (:merId/product/destory/:id) | `C` | POST | `:merId/product/destory/:id` | `StoreProduct/destory` | `` |

### `product/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (:merId/product/detail/:id) | `R` | GET | `:merId/product/detail/:id` | `StoreProduct/detail` | `` |

### `product/edit_cate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 修改分类 | `U` | POST | `:merId/product/edit_cate/:id` | `StoreProduct/editCate` | `` |

### `product/edit_mer_cate`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| edit Mer Cate (:merId/product/edit_mer_cate/:id) | `U` | POST | `:merId/product/edit_mer_cate/:id` | `StoreProduct/editMerCate` | `` |

### `product/good`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update Good (:merId/product/good/:id) | `U` | POST | `:merId/product/good/:id` | `StoreProduct/updateGood` | `` |

### `product/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (:merId/product/lst) | `R` | GET | `:merId/product/lst` | `StoreProduct/lst` | `` |

### `product/restore`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| restore (:merId/product/restore/:id) | `C` | POST | `:merId/product/restore/:id` | `StoreProduct/restore` | `` |

### `product/status`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| switch Status (:merId/product/status/:id) | `U` | POST | `:merId/product/status/:id` | `StoreProduct/switchStatus` | `` |

### `product/title`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| title (:merId/product/title) | `R` | GET | `:merId/product/title` | `StoreProduct/title` | `` |

### `product/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 编辑 | `U` | POST | `:merId/product/update/:id` | `StoreProduct/update` | `` |

### `product/value`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Value (:merId/product/value/:id) | `R` | GET | `:merId/product/value/:id` | `StoreProduct/getValue` | `` |
| set Value (:merId/product/value/:id) | `U` | POST | `:merId/product/value/:id` | `StoreProduct/setValue` | `` |

### `refund/check`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 退款单 | `R` | GET | `:merId/refund/check/:id` | `/check` | `` |

### `refund/compute`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| compute (:merId/refund/compute) | `D` | POST | `:merId/refund/compute` | `/compute` | `` |

### `refund/confirm`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| refund Price (:merId/refund/confirm/:id) | `UD` | POST | `:merId/refund/confirm/:id` | `/refundPrice` | `` |

### `refund/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (:merId/refund/create) | `CD` | POST | `:merId/refund/create` | `/create` | `` |

### `refund/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (:merId/refund/detail/:id) | `R` | GET | `:merId/refund/detail/:id` | `/detail` | `` |

### `refund/express`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| express (:merId/refund/express/:id) | `R` | GET | `:merId/refund/express/:id` | `/express` | `` |

### `refund/get`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get Refund Price (:merId/refund/get/:id) | `R` | GET | `:merId/refund/get/:id` | `/getRefundPrice` | `` |

### `refund/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| lst (:merId/refund/lst) | `R` | GET | `:merId/refund/lst` | `/lst` | `` |

### `refund/mark`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| mark (:merId/refund/mark/:id) | `UD` | POST | `:merId/refund/mark/:id` | `/mark` | `` |

### `refund/status`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| switch Status (:merId/refund/status/:id) | `UD` | POST | `:merId/refund/status/:id` | `/switchStatus` | `` |

### `reservation/config`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Config (reservation/config) | `R` | GET | `reservation/config` | `/reservationConfig` | `` |

### `reservation/staffs`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 预约 | `R` | GET | `:merId/reservation/staffs` | `/staffList` | `` |

### `reservationconfig`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Config (:merId/reservationconfig) | `R` | GET | `:merId/reservationconfig` | `/reservationConfig` | `` |

### `reservationdispatch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Dispatch (:merId/reservationdispatch/:id) | `C` | POST | `:merId/reservationdispatch/:id` | `/reservationDispatch` | `` |

### `reservationreschedule`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Reschedule (:merId/reservationreschedule/:id) | `C` | POST | `:merId/reservationreschedule/:id` | `/reservationReschedule` | `` |

### `reservationupdateDispatch`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Update Dispatch (:merId/reservationupdateDispatch/:id) | `U` | POST | `:merId/reservationupdateDispatch/:id` | `/reservationUpdateDispatch` | `` |

### `reservationverify`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| reservation Verify (:merId/reservationverify/:id) | `U` | POST | `:merId/reservationverify/:id` | `/reservationVerify` | `` |

### `scan_login`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 客服扫码登录 | `C` | POST | `scan_login/:key` | `api.store.service.Service/scanLogin` | `` |

### `statistics`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| order Statistics (:merId/statistics) | `R` | GET | `:merId/statistics` | `/orderStatistics` | `` |

### `template/create`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| create (:merId/template/create) | `C` | POST | `:merId/template/create` | `ShippingTemplate/create` | `` |

### `template/delete`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| batch Delete (:merId/template/delete) | `D` | POST | `:merId/template/delete` | `ShippingTemplate/batchDelete` | `` |

### `template/detail`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| detail (:merId/template/detail/:id) | `R` | GET | `:merId/template/detail/:id` | `ShippingTemplate/detail` | `` |

### `template/lst`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| 运费模板 | `R` | GET | `:merId/template/lst` | `ShippingTemplate/lst` | `` |

### `template/select`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| get List (:merId/template/select) | `R` | GET | `:merId/template/select` | `ShippingTemplate/getList` | `` |

### `template/update`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| update (:merId/template/update/:id) | `U` | POST | `:merId/template/update/:id` | `ShippingTemplate/update` | `` |

### `user`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| user (user/:merId/:uid) | `R` | GET | `user/:merId/:uid` | `api.store.service.Service/user` | `` |

### `user_list`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| service User List (user_list/:merId) | `R` | GET | `user_list/:merId` | `api.store.service.Service/serviceUserList` | `` |

### `verify`

| 操作说明 | CRUD | HTTP | 路径 | 处理器 | 路由名 |
| --- | --- | --- | --- | --- | --- |
| verify (:merId/verify/:id) | `U` | POST | `:merId/verify/:id` | `/verify` | `` |
