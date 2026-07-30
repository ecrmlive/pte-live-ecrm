# 云服务配置中心

平台后台“系统设置 → 云服务配置”把第三方服务参数以数据库密文方式管理。表结构在 `sql/init_table.sql` 中统一初始化。

## 分组

- 微信支付、支付宝支付
- 腾讯云账号（共享 SecretId / SecretKey）与 COS
- 腾讯云 LVB 推流/拉流域名、URL 模板、鉴权 Key 与应用名
- 腾讯云 VOD
- 腾讯直播推流及播放器 License
- LVB 推流、断流、异常回调模板
- pte-live-im 接入地址及服务令牌

## 安全与生效规则

- 所有字段使用统一后台运行时密钥派生的 AES-GCM 密钥加密后写入 `qixi_crm_a_cloud_config`；数据库不保存明文。
- 读取后台页面时，密钥字段固定返回 `********`；提交空值或该掩码不会覆盖已有密钥。
- 仅拥有“保存云服务配置”按钮权限的平台管理员可以写入。
- COS 已接入运行时：后台把 COS 启用并补齐存储桶、地域及腾讯云账号后，下一次素材上传即时走数据库配置；未启用时保持 `app.yaml` 当前上传策略。
- IM 已接入运行时：C 端和客服端在开会话、发放 UserSig、创建 C2C 会话时通过服务端 `cloudconfig.Service.Values("im")` 读取；后台保存的非空字段优先于 YAML。`api_base` 只能填写服务端 S2S 地址，`api_public_url` 与 `ws_public_url` 必须填写客户端实际可访问地址。
- 支付、LVB/VOD、License 与回调配置的接入方必须通过服务端 `cloudconfig.Service.Values` 读取；该方法禁止在 HTTP 响应中使用。

## 运维要求

- 在每个同构环境按 `sql/README.md` 执行集中初始化后，再由后台录入对应环境的真实参数。
- `jwt.secret` 变化会使既有密文无法解密；生产环境必须在变更前导出并在变更后重新录入，或执行受控密钥轮换。
- 密钥禁止写入 Git、OpenAPI 示例、日志和截图。
