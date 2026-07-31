export type CaptchaAction = "login_password" | "login_sms" | "register";

interface CaptchaSDKOptions {
  endpoint: string;
  action: CaptchaAction;
  locale?: string;
  theme?: "light" | "dark";
  preferredMode?: "invisible" | "puzzle" | "image" | "click" | "rotate" | "restore";
  styles?: Record<string, string>;
}

interface CaptchaSDKInstance {
  verify(): Promise<{ verificationToken: string }>;
}

interface CaptchaSDKConstructor {
  new(options: CaptchaSDKOptions): CaptchaSDKInstance;
}

declare global {
  interface Window {
    CaptchaSDK?: CaptchaSDKConstructor;
  }
}

let loading: Promise<CaptchaSDKConstructor> | undefined;

function loadStylesheet(href: string) {
  if (document.querySelector(`link[data-pte-captcha-sdk="${href}"]`)) return;
  const link = document.createElement("link");
  link.rel = "stylesheet";
  link.href = href;
  link.dataset.pteCaptchaSdk = href;
  document.head.appendChild(link);
}

async function loadSDK(): Promise<CaptchaSDKConstructor> {
  if (window.CaptchaSDK) return window.CaptchaSDK;
  if (!loading) {
    loading = new Promise((resolve, reject) => {
      loadStylesheet("/pte-tools-captcha/captcha-sdk.css");
      const script = document.createElement("script");
      script.src = "/pte-tools-captcha/captcha-sdk.js";
      script.async = true;
      script.onload = () => window.CaptchaSDK ? resolve(window.CaptchaSDK) : reject(new Error("验证码组件加载失败"));
      script.onerror = () => reject(new Error("验证码组件加载失败"));
      document.head.appendChild(script);
    });
  }
  return loading;
}

// SDK 请求同源 BFF；HMAC 密钥与 pte-tools-captcha 的内网地址始终留在服务端。
export async function verifyWithPteCaptcha(action: CaptchaAction): Promise<string> {
  const CaptchaSDK = await loadSDK();
  const captcha = new CaptchaSDK({
    // SDK 不接受空 endpoint；使用当前站点，由 Vite/Nginx 将 /api 转发到 api-business。
    endpoint: window.location.origin,
    action,
    locale: "zh-CN",
    theme: "dark",
    // 登录和短信发送属于敏感动作，始终展示真实拼图挑战而不是前端放行。
    preferredMode: "puzzle",
    styles: { "--pte-captcha-primary": "#2878f0" },
  });
  const result = await captcha.verify();
  if (!result.verificationToken) throw new Error("安全验证未返回校验令牌");
  return result.verificationToken;
}
