export type CaptchaAction = "login_password" | "login_sms" | "register";

interface CaptchaSDKOptions {
  endpoint: string;
  action: CaptchaAction;
  locale?: string;
  theme?: "light" | "dark";
  preferredMode?: "invisible" | "puzzle" | "image" | "click" | "rotate" | "restore";
  styles?: Record<string, string>;
  className?: string;
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

interface CaptchaClientSignals {
  fingerprint: string;
  event_count: number;
  duration_ms: number;
  language: "zh-CN";
}

interface CaptchaGatewayResponse {
  challenge_id: string;
  status: string;
  message?: string;
  verification_token?: string;
  image_png_base64?: string;
  image_svg?: string;
}

export interface InlineImageCaptcha {
  challengeId: string;
  imageUrl: string;
}

function clientSignals(): CaptchaClientSignals {
  return {
    fingerprint: `pc-${navigator.userAgent.length}-${screen.width}x${screen.height}`,
    event_count: 1,
    duration_ms: 0,
    language: "zh-CN",
  };
}

async function captchaRequest(path: string, body: Record<string, unknown>): Promise<CaptchaGatewayResponse> {
  const response = await fetch(path, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await response.json().catch(() => ({})) as CaptchaGatewayResponse & { error?: string };
  if (!response.ok) throw new Error(data.error || "验证码服务暂不可用，请稍后重试");
  return data;
}

function imageDataUrl(challenge: CaptchaGatewayResponse): string {
  if (challenge.image_png_base64) return `data:image/png;base64,${challenge.image_png_base64}`;
  if (challenge.image_svg) return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(challenge.image_svg)}`;
  throw new Error("图形验证码加载失败，请换一张");
}

/** 密码登录使用内嵌亮色图片验证码，不弹出二次窗口。 */
export async function createInlineImageCaptcha(action: "login_password"): Promise<InlineImageCaptcha> {
  const challenge = await captchaRequest("/api/v1/challenges", {
    preferred_mode: "image",
    action,
    locale: "zh-CN",
    theme: "light",
    client: clientSignals(),
  });
  if (!challenge.challenge_id) throw new Error("图形验证码加载失败，请换一张");
  return { challengeId: challenge.challenge_id, imageUrl: imageDataUrl(challenge) };
}

export async function verifyInlineImageCaptcha(challengeId: string, answer: string): Promise<string> {
  const result = await captchaRequest(`/api/v1/challenges/${encodeURIComponent(challengeId)}/verify`, {
    answer: answer.trim(),
    client: clientSignals(),
  });
  if (result.status !== "verified" || !result.verification_token) {
    throw new Error(result.message || "图形验证码错误，请换一张后重试");
  }
  return result.verification_token;
}

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
    theme: "light",
    // 登录和短信发送属于敏感动作，始终展示真实拼图挑战而不是前端放行。
    preferredMode: "puzzle",
    styles: { "--pte-captcha-primary": "#2878f0" },
    className: "pte-captcha--pc-puzzle",
  });
  const result = await captcha.verify();
  if (!result.verificationToken) throw new Error("安全验证未返回校验令牌");
  return result.verificationToken;
}
