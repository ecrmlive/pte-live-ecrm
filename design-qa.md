# 商城设置页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-4457a1f4-f109-460f-b8ec-d5fec6879872.png`
  （其余四个页签参考同一用户提供的截图）。
- Implementation URL：`http://127.0.0.1:15124/#/setting/shop`
- Target viewport：待以管理后台已登录状态捕获。
- State：商品设置页签，已登录的平台管理员。

## 浏览器渲染证据

已恢复本地预览并捕获到浏览器渲染状态，但当前浏览器没有平台管理员登录会话，只能到达登录页，无法进入 `/setting/shop`。因此没有与用户截图同一页面状态的实现截图，也不能把源码截图与实现截图放到同一输入中进行视觉判断。

已完成的非视觉验证：

- `pnpm build:platform` 通过；
- `go test ./internal/domain/content` 通过；
- 本机平台 API 已重启，`http://127.0.0.1:18081/healthz` 返回成功；
- 商城设置 SQL 已同步到本机 MySQL。

## 必查视觉维度

以下维度等待浏览器截图后再逐项判断：字体与层级、间距与布局节奏、颜色与令牌、图片/保证金标识选择器、产品文案与五个页签状态。

## 结论

final result: blocked

阻塞原因：浏览器缺少平台管理员登录会话，无法打开商城设置页完成截图级 UI 对照验收。

---

# 热门搜索页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-8828975e-9eb2-4a8f-be75-ea77f322d29d.png`（列表）与
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-0279c8cc-9199-4033-9480-c61a358860b1.png`（新增）。
- Implementation URL：`http://127.0.0.1:15126/#/group/config/67`。
- Target viewport：2048 × 864（列表）与 1248 × 558（新增）。两张参考均为 CSS 像素、密度未知；当前没有可用于归一化的已登录实现截图。
- State：已登录的平台管理员；热门搜索列表及打开「新增热门搜索」抽屉。

## 浏览器渲染证据

本机 Vite 预览已启动，平台 API 健康检查通过；但当前可用浏览器没有平台管理员会话，访问实现地址只会进入登录页。因此无法捕获与参考截图相同的页面、视口和交互状态，也不能进行并列视觉比较、控制台检查或主交互点击验证。

已完成的非视觉验证：

- `pnpm vite build --mode test` 通过；
- `go test ./internal/platform/nativeconfigitem` 通过；
- 热门搜索菜单权限与 8 条中文演示数据已同步到本机 MySQL；
- 平台 API 已重启，`/healthz` 返回成功。

## 必查视觉维度

待在已登录状态重新捕获后确认：列表卡片与列宽、分页贴底、是否显示开关无文字、抽屉字段行距与底部居中「取消 / 保存」按钮、字体层级、颜色令牌、排序步进器和窄屏适配。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，未能取得与用户提供截图同状态的实现渲染图进行强制视觉对照。

---

# 商品详情装修视觉验收

## 对照目标

- Source visual truth：用户提供的详情装修配置与预览截图：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-3d4fe468-4f9f-4d6c-8984-36b40c73f285.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-b809edff-ff99-43c2-bbb9-24bcefd35f88.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-f76ec7bf-c4ff-4856-b0cd-4952f559e55f.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-363abe3f-aa11-4502-9391-44c3f8ca26a7.png`。
- Implementation URL：`http://localhost:15124/#/setting/diy/product_detail`。
- Implementation screenshot：`/private/tmp/ecrm-product-detail-diy-20260816.png`。
- Viewport：当前已登录平台后台视口 1700 × 900 CSS px；源图为页面局部截图，按手机预览 375px 宽、右侧配置 400px 宽进行区域对齐。
- State：商品信息模块选中；导航菜单展开；全部默认配置启用。

## 比较与修正记录

- [P1 已修正] 预览模块与右侧配置没有模块级联动。
  修正：商品信息、排行榜、优惠券、商品参数、优惠套餐、商品评价、种草秀、店铺信息、底部菜单都成为独立可选模块；选中项切换对应配置。
- [P1 已修正] 顶部导航缺少参考页的展开菜单和可配置导航项。
  修正：复刻返回、菜单、刷新、更多、圆点控制区；菜单默认展开，并与“菜单内容”多选框同步。
- [P1 已修正] 右侧配置区字段结构、间距和字体层级不一致。
  修正：按参考顺序实现顶部导航、商品主图、收藏/分享、商品信息、SVIP 等分段，使用 400px 固定右栏、57px 标题行、15px 内容内距。
- [P1 已修正] 底部保存区及页面模块缺少固定/选择状态。
  修正：保存区固定于内容区底部；选中模块显示蓝色边框和左侧箭头提示。

## 交互验证

- 点击“优惠券”预览模块后，右栏标题切换为“优惠券”。
- 在“是否显示”切换为“隐藏”后，预览中的优惠券模块立即消失；页面刷新后默认状态恢复。
- 导航、轮播点、主图模式、收藏/分享、商品信息、各模块显示状态、数量滑杆和底部菜单均由同一响应式配置驱动。

## 必查保真面

- 字体与层级：后台既有 PingFang SC / Microsoft YaHei 回退栈；预览标题、模块标题、说明文字和右栏配置标题分别建立层级。
- 间距与布局：手机宽度 375px，右栏宽度 400px，预览模块与配置分段按参考页面固定节奏排列。
- 颜色与状态：选中边框 #1890ff、活动锚点 #4073fa、页面底色 #f0f2f5，与参考页面蓝灰体系一致。
- 图片与内容：使用 CRMEB 参考页面同源商品主图、菜单图、评价图、种草图及店铺图，不使用占位图。
- 文案：配置项与预览文案按参考页保留。

## 结论

final result: passed

---

# 协议设置 iPhone 17 DIY 预览视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-cc5e0ab5-59ea-4c0d-9776-3dc89c154ab6.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/agreements`。
- Target state：已登录的平台管理员，打开“用户协议”。预览需要与商品详情 DIY 共用同一 iPhone 17 Pro 设备壳，而不是自绘的通用手机外框。
- Source pixels：2886 × 1768；实现浏览器为本机会话默认视口。由于缺少同状态实现页面，无法做像素密度归一化或并列比较。

## 已实现的对照项

- 协议预览改为复用商品详情 DIY 的 `DevicePreviewFrame`；设备锁定为 iOS，关闭设备切换器、展开按钮和侧边工具条。
- 共用组件的 iOS 设备逻辑尺寸为 iPhone 17 Pro：402 × 874 CSS px，含动态岛、状态栏、44px 导航栏和底部安全区。
- 协议名称作为机内导航标题，正文在设备内容区内独立滚动，并随右侧富文本编辑实时更新。

## 浏览器渲染证据

- 本机 Vite 服务：`http://127.0.0.1:15126/`，HTTP 200；
- 已捕获浏览器状态：`/private/tmp/ecrm-agreement-preview-login-20260812.png`；
- 当前浏览器被重定向至登录页，没有平台管理员会话，无法进入受保护的协议设置路由。因此该截图不能与参考图做同状态比较，也不能判断字体、间距、颜色、图片资产和文案的最终屏幕呈现。

## 结论

final result: blocked

阻塞原因：当前自动化浏览器未登录本地平台后台，缺少“用户协议”页面的实现截图，无法完成截图级视觉对照；代码已通过 `pnpm typecheck` 与 `pnpm build:platform`。

---

# 消息管理页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-b1b0f459-ed19-43ad-b462-86833b51338b.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-fe7dce56-2533-45d8-ae1c-3b170e279cf4.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/notification/index`。
- Target state：已登录的平台管理员；分别打开“通知会员”与“通知店铺”。

## 已实现的对照项

- 顶部固定为“通知会员 / 通知店铺”两项标签；每项从本地 `qixi_crm_a_notification_config` 读取独立真实配置；
- 列表严格保留 ID、通知类型、通知场景说明、公众号模板、小程序订阅、发送短信、操作；发送渠道使用无文字开关；
- “设置”打开抽屉，维护三个渠道的启用状态与固定文本，底部按钮使用“取消 / 保存”；已启用渠道必须填写固定文本；
- 同步按钮会请求真实 API；外部微信凭据尚未配置时，API 明确返回未执行同步，不伪造同步成功；
- 本地数据库已同步 9 条会员通知和 9 条店铺通知配置，平台 API 已重新部署并通过健康检查。

## 浏览器渲染证据

当前会话没有可操作的已登录平台管理员浏览器会话，无法进入实现路由获取与参考图同一状态、同一视口的截图。因此不能进行截图并列对照或点击级视觉验收。

已完成的非视觉验证：

- `pnpm exec eslint`、`pnpm exec vue-tsc --noEmit`、`pnpm run build` 通过；
- `go test ./internal/domain/notification ./internal/infra/persist/notification ./internal/platform/notification` 通过；
- SQL 已同步，并确认本地 MySQL 中 `member=9`、`store=9`；
- 本地 `api-platform` 已重新部署，容器健康检查返回 `ok=true`。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录本地平台后台浏览器会话，未能完成同状态截图级视觉对照。

---

# 消息管理 - 公告管理页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-9d47780e-81df-497e-8c23-512b6d7a71eb.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-63a0a5e4-1a55-4f23-9138-5506f55b00fb.png`、
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-25381512-aae8-4130-8fd1-596c764bbfc6.png`。
- Implementation URL：`http://127.0.0.1:15126/#/content/notice`。
- Target state：已登录的平台管理员，公告列表及新增/编辑公告抽屉。

## 已实现的对照项

- 筛选区仅保留时间选择、启用状态、消息名称；按钮按全局强制规范为“重置 / 搜索”。
- 列表仅保留公告名称、店铺范围、启用状态、发送日期、操作；状态为无文字切换开关，操作为详情、编辑、删除。
- 工具栏为“发布公告”；表单抽屉标题统一为“新增公告”或“编辑公告”，底部使用居中的“取消 / 保存”。
- 表单字段仅包含消息名称、选择店铺、公告内容；公告内容采用富文本编辑器。店铺名称通过统一店铺选择弹窗关联，店铺类别/分类为真实关联下拉选择，未使用任何输入式 ID。
- 本地 MySQL 已同步 `qixi_crm_a_notice.scope_type` 与 `qixi_crm_a_notice_scope` 关联表；创建、编辑、删除、状态切换、列表筛选和详情均调用本地平台 API。

## 浏览器渲染证据

本次未获取可操作的已登录平台管理员浏览器会话。未认证请求本地列表接口返回 401，符合后台接口鉴权预期；因此无法取得与参考图同状态、同视口的实现截图，也不能做并列视觉比较。

已完成的非视觉验证：

- `pnpm exec eslint src/views/ecrm/content/notice.vue src/api/core/platform-content.ts` 通过；
- `pnpm exec vue-tsc --noEmit --pretty false` 通过；
- `go test ./internal/domain/content ./internal/infra/persist/content ./internal/platform/content` 通过；
- `make local-sync-sql` 已导入公告范围补丁，`make local-sync-api SVC=api-platform` 已完成本地 API 重建与重启；
- 容器 `pte_live_ecrm_api_platform` 为 running，数据库中公告范围列与关联表均存在。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，无法完成要求的同状态截图级视觉对照。

---

# 菜单管理页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-ee11d0b9-0e0d-4d24-8fa5-12e5d85e3b8c.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/menu`。
- State：已登录的平台管理员，打开“平台”菜单归属。

## 已实现的对照项

- 顶部固定“平台 / 商户 / 区域”三类菜单归属，切换后只读取相应归属的真实菜单；
- 工具栏为“新增菜单”和“展开/收起”，树表字段严格为菜单名称、菜单地址、菜单图标、创建时间、排序、操作；
- 操作区为新增子菜单、编辑、删除；新增与编辑使用“新增菜单 / 编辑菜单”抽屉及“取消 / 保存”；
- 上级菜单为树形关联选择，不允许输入 ID；后端校验菜单标识唯一、父子菜单归属一致、不能选择自身或子菜单为父级、有子菜单不可删除；
- 本地 MySQL 已迁移 `menu_scope`、`created_at` 与联合索引，平台 API 已重新构建启动。

## 浏览器渲染证据

当前会话没有可操作的已登录平台管理员浏览器会话，不能进入实现 URL 获取与参考图相同状态的实现截图，因而不能完成截图级视觉比较。

已完成的非视觉验证：

- `pnpm typecheck` 通过；
- `go test ./internal/admin/auth` 通过；
- `make local-sync-sql` 完成，迁移字段和索引已在本机 MySQL 复核；
- `make local-sync-api SVC=api-platform` 完成，`http://127.0.0.1:18081/healthz` 返回成功。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，未能取得与用户提供截图同状态的实现渲染图进行强制视觉对照。

---

# 管理员管理页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-4a91082c-25dc-484a-b417-d401e8d4634c.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/admin`。
- Target viewport：2048 × 684；已登录的平台管理员，管理员管理列表首屏。

## 已实现的对照项

- 筛选区提供选择时间、账号状态和账号/昵称关键字，按钮顺序固定为“重置 / 搜索”；
- 工具栏统一为“新增管理员”；表格字段严格为 ID、管理员姓名、身份、所属区域、账号、手机号、账号状态、创建时间、操作；
- 账号状态使用无文字开关；操作列固定在右侧，保留修改密码、编辑、删除；
- 新增和编辑使用抽屉，标题统一为“新增管理员 / 编辑管理员”，底部使用“取消 / 保存”；角色、区域、商户、店铺及代理均为关联选择，未提供 ID 输入；
- 列表、创建、编辑、状态切换、修改密码和删除均连接真实平台 API；新增区域选项接口直接从本地 MySQL 读取。

## 浏览器渲染证据

当前会话没有可操作的已登录平台管理员浏览器会话，无法进入实现 URL 获取与参考图相同状态、相同视口的渲染截图，因而不能完成截图级视觉比较。

已完成的非视觉验证：

- `pnpm exec eslint src/views/ecrm/setting/admin.vue` 通过；
- `pnpm --config.engine-strict=false run build` 通过；
- `go test ./internal/admin/auth` 通过；
- `make local-sync-sql` 完成，`make local-sync-api SVC=api-platform` 已重新构建并启动本地 API；
- `http://127.0.0.1:18081/healthz` 返回成功。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，未能取得与用户提供截图同状态的实现渲染图进行强制视觉对照。

---

# 角色权限页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-0db6c6df-eaf9-40c3-a242-3cc5ceff8047.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/role`。
- State：已登录的平台管理员，打开“平台”身份标签页。

## 已实现的对照项

- 顶部固定平台、商户、区域三类身份标签；工具栏统一为“新增身份管理”；
- 表格字段只保留 ID、身份名称、身份类型、是否开启、创建时间、更新时间、操作；
- 开关不显示启用/关闭文字；编辑与删除位于固定右侧操作列；新增与编辑抽屉使用“取消 / 保存”；
- 列表、状态切换、新增、编辑和删除均走本地平台 API 与 `qixi_crm_a_role` 真实数据；系统预置身份及已关联后台账号的身份不能删除。

## 浏览器渲染证据

当前会话没有可操作的已登录平台管理员浏览器会话，无法进入实现路由取得与参考图同状态的截图，也不能完成并列视觉比较和点击级核验。

已完成的非视觉验证：

- `pnpm --config.engine-strict=false run build` 通过；
- `go test ./internal/admin/auth` 通过；
- 本地 MySQL 已同步身份类型与时间字段，平台、商户、区域基础身份可读取；
- 本机平台 API 已重新部署，`http://127.0.0.1:18081/healthz` 返回成功。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，未能取得与用户提供截图同状态的实现渲染图进行强制视觉对照。

---

# 物流公司页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-12d2bcd3-57e6-47c0-9719-d9b9d9e1598e.png`。
- Implementation URL：`http://127.0.0.1:15126/#/freight/express`。
- Target viewport：2048 × 1244；已登录的平台管理员，物流公司列表首屏。

## 已实现的对照项

- 搜索字段为“物流公司名称或者编码”，筛选按钮按全局强制规范固定为“重置 / 搜索”；
- 工具栏仅保留“同步物流公司”；表格字段为 ID、物流公司名称、编码、排序、是否显示、操作；
- 是否显示使用无文字开关；编辑与删除保持右侧固定操作列；编辑抽屉使用“取消 / 保存”；
- 本地 MySQL 已同步 13 条 UTF-8 物流公司目录数据，搜索、显示状态、编辑、删除和目录同步均走真实平台 API。

## 浏览器渲染证据

参考图已打开并确认。当前 `127.0.0.1:15126` 对应的本地 Vite 监听端口无法从本会话发起连接，且没有可操作的已登录平台管理员浏览器会话；因此无法获取同视口、同状态的实现截图，也不能制作并列视觉比较。

已完成的非视觉验证：

- `pnpm --config.engine-strict=false run build` 通过；
- `go test ./internal/domain/logistics ./internal/infra/persist/logistics ./internal/platform/logistics` 通过；
- SQL 已同步，数据库中 1105–1117 物流公司名称的 UTF-8 十六进制编码正确；
- 平台 API 已重新部署，`http://127.0.0.1:18081/healthz` 返回成功。

## 必查视觉维度

待获得已登录实现截图后，按同一 2048 × 1244 视口对照：搜索栏宽度和间距、工具栏至表格距离、列宽、表头底色、行高、分页贴合、固定操作列、字体层级与颜色令牌。该页面无自定义图片资产，图像质量项不适用。

## 结论

final result: blocked

阻塞原因：本会话无法捕获已登录的本地平台后台渲染图，无法完成要求的同状态截图级视觉比较。

---

# 协议设置页面视觉验收

## 对照目标

- Source visual truth：
  `/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-57ff1f9a-b8ff-4a6c-ae59-aff54e160d0e.png`。
- Implementation URL：`http://127.0.0.1:15126/#/setting/agreements`。
- Target viewport：2048 × 1126；状态为已登录的平台管理员、打开“用户协议”。

## 已实现的对照项

- 左侧固定十项协议导航：用户协议、隐私政策、注销提示、平台规则、店铺入驻申请协议、代理入驻申请协议、商户入驻申请协议、注销声明、关于我们、资质证照；
- 中部移动端效果预览会随富文本实时更新；
- 右侧保留富文本工具栏和正文编辑区，底部按钮统一为居中的“重置 / 保存”；
- 后端列表接口仅返回上述十项，保存仍沿用既有协议缓存，避免影响其他独立协议维护页。

## 浏览器渲染证据

本机 Vite 预览和平台 API 均已可用，但当前浏览器没有平台管理员登录会话，无法进入 `/setting/agreements` 获取和参考图同一状态、同一视口的实现截图。因此不能做截图并列对照或完成点击级视觉验收。

已完成的非视觉验证：

- `pnpm vite build --mode test` 通过；
- `go test ./internal/domain/content` 通过；
- 协议设置菜单与 10 项本地演示协议已同步到本机 MySQL；
- `http://127.0.0.1:18081/healthz` 返回成功。

## 结论

final result: blocked

阻塞原因：缺少可操作的已登录平台管理员浏览器会话，未能取得与用户提供截图同状态的实现渲染图进行强制视觉对照。
