# PC 登录验证码视觉验收

- 参考问题截图：`/private/var/folders/yn/8j0yr4gs2pzg508td1q42cv00000gn/T/codex-clipboard-08334f61-1845-4a5e-b96d-0895149d3b22.png`
- 实现截图（密码登录）：`/private/tmp/qixi-pc-login-inline-captcha-final.png`
- 实现截图（验证码登录）：`/private/tmp/qixi-pc-login-puzzle-captcha.png`
- 视口：1280 × 720 CSS px，设备像素比 1。
- 状态：本地 `http://127.0.0.1:5183/login`；验证码服务走 `api-business` BFF，语言 `zh-CN`。

## 对比结论

原问题是密码登录将图形验证码做成二次弹窗，并把 SDK 的英文取消错误直接渲染到表单；验证码登录拼图窗口宽度为 840px，明显超过桌面登录流程所需的操作范围。

本次实现：

1. 密码登录在表单第三行直接显示亮色图片验证码和输入框；无独立刷新按钮，点击验证码图片即可刷新，图片完整等比显示。
2. 图形验证码请求明确传递 `preferred_mode: image`、`theme: light`、`locale: zh-CN`。
3. 验证码登录保留深色拼图交互，但收缩为最大 460px 宽、52px 滑块的紧凑弹窗。
4. 主动关闭拼图验证只结束当前验证，不再显示 `Captcha cancelled` 等英文内部错误。

## 验证记录

- [x] 浏览器实测密码登录页面已直接显示图片验证码。
- [x] 浏览器实测切换验证码登录后显示紧凑拼图窗口；未求解或绕过 CAPTCHA。
- [x] 关闭验证码窗口后，登录表单没有英文取消错误文案。
- [x] `pnpm -C app-pc type-check` 通过。
- [x] 浏览器控制台无 error/warn。
- [x] `git diff --check` 通过。

## 质量结论

- 字体与层级：与原 PC 登录弹窗的中文标题、表单层级保持一致。
- 间距与布局：验证码在原有表单行内完成，不增加二次窗口；拼图窗口不再覆盖主要页面区域。
- 色彩：图片验证码使用亮色背景；拼图验证保留深色安全验证层级。
- 图片：使用验证码服务返回的 PNG，等比完整显示，未裁剪字符。
- 文案：所有用户可见验证码文案为中文，取消行为不显示内部英文错误。

final result: passed
