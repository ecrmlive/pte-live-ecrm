<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import HomePage from "@/pages/home/HomePage.vue";
import { createInlineImageCaptcha, verifyInlineImageCaptcha, verifyWithPteCaptcha, type CaptchaAction, type InlineImageCaptcha } from "@/lib/pte-captcha";
import { useUserStore } from "@/stores/user";
import { ApiError } from "@/utils/request";

const user = useUserStore();
const router = useRouter();
const route = useRoute();
const page = ref<"login" | "register">("login");
const loginMode = ref<"password" | "sms">("password");
const phone = ref("");
const password = ref("");
const confirmPassword = ref("");
const smsCode = ref("");
const accepted = ref(true);
const loading = ref(false);
const captchaBusy = ref(false);
const message = ref("");
const captchaToken = ref("");
const captchaAction = ref<CaptchaAction | "">("");
const imageCaptcha = ref<InlineImageCaptcha>();
const imageCaptchaAnswer = ref("");
const imageCaptchaLoading = ref(false);

const phoneValid = computed(() => /^1\d{10}$/.test(phone.value.trim()));

function close() { void router.push({ name: "home" }); }
function resetFeedback() { message.value = ""; captchaToken.value = ""; captchaAction.value = ""; imageCaptchaAnswer.value = ""; }
function selectLoginMode(mode: "password" | "sms") {
  loginMode.value = mode;
  resetFeedback();
  if (mode === "password" && !imageCaptcha.value) void refreshImageCaptcha();
}
function openRegister() { page.value = "register"; resetFeedback(); password.value = ""; confirmPassword.value = ""; smsCode.value = ""; }
function openLogin() { page.value = "login"; resetFeedback(); password.value = ""; confirmPassword.value = ""; smsCode.value = ""; }
async function completeCaptcha(action: CaptchaAction) {
  message.value = "";
  captchaBusy.value = true;
  try {
    captchaToken.value = await verifyWithPteCaptcha(action);
    captchaAction.value = action;
    message.value = "安全验证通过";
  } catch (error) {
    const reason = (error as ApiError).message || "安全验证未完成";
    // 用户主动关闭验证码窗口不是错误，不在登录框里输出 SDK 的内部英文文案。
    if (reason !== "已取消安全验证" && reason !== "Captcha cancelled") message.value = reason;
  } finally {
    captchaBusy.value = false;
  }
}
async function refreshImageCaptcha() {
  imageCaptchaLoading.value = true;
  imageCaptchaAnswer.value = "";
  captchaToken.value = "";
  captchaAction.value = "";
  try {
    imageCaptcha.value = await createInlineImageCaptcha("login_password");
  } catch (error) {
    imageCaptcha.value = undefined;
    message.value = (error as ApiError).message || "图形验证码加载失败，请稍后重试";
  } finally {
    imageCaptchaLoading.value = false;
  }
}
async function completeInlineImageCaptcha(): Promise<boolean> {
  if (captchaAction.value === "login_password" && captchaToken.value) return true;
  if (!imageCaptcha.value) {
    await refreshImageCaptcha();
    message.value = "图形验证码正在加载，请输入后再登录";
    return false;
  }
  if (!imageCaptchaAnswer.value.trim()) { message.value = "请输入图形验证码"; return false; }
  captchaBusy.value = true;
  message.value = "";
  try {
    captchaToken.value = await verifyInlineImageCaptcha(imageCaptcha.value.challengeId, imageCaptchaAnswer.value);
    captchaAction.value = "login_password";
    message.value = "图形验证码已验证";
    return true;
  } catch (error) {
    message.value = (error as ApiError).message || "图形验证码错误，请换一张后重试";
    await refreshImageCaptcha();
    return false;
  } finally {
    captchaBusy.value = false;
  }
}
async function sendCode(action: "login_sms" | "register") {
  if (!phoneValid.value) { message.value = "请输入正确的手机号"; return; }
  await completeCaptcha(action);
  // 短信发送仅能在安全验证拿到服务端校验令牌后调用；目前不再由浏览器本地倒计时伪造发送成功。
  if (captchaAction.value === action && captchaToken.value) {
    message.value = "安全验证已完成，请由平台后台配置短信通道后发送验证码";
  }
}
function validAgreement() { if (accepted.value) return true; message.value = "请先阅读并同意用户协议和隐私政策"; return false; }

async function submit() {
  message.value = "";
  if (!validAgreement() || !phone.value.trim()) return;
  if (page.value === "login" && loginMode.value === "password") {
    if (!password.value) { message.value = "请输入登录密码"; return; }
    if (!(await completeInlineImageCaptcha())) return;
    loading.value = true;
    try {
      await user.login(phone.value.trim(), password.value, captchaToken.value);
      const redirect = typeof route.query.redirect === "string" ? route.query.redirect : "/user";
      await router.replace(redirect || "/user");
    } catch (error) { message.value = (error as ApiError).message || "登录失败"; }
    finally { loading.value = false; }
    return;
  }
  if (page.value === "register") {
    if (!password.value || password.value !== confirmPassword.value) { message.value = "请确认两次输入的密码一致"; return; }
    if (!smsCode.value) { message.value = "请输入短信验证码"; return; }
    message.value = "短信验证码服务接口正在接入，密码登录与账号注册接口已可用";
    return;
  }
  if (!smsCode.value) { message.value = "请输入短信验证码"; return; }
  message.value = "短信验证码登录服务接口正在接入";
}

onMounted(() => { void refreshImageCaptcha(); });

</script>

<template>
  <HomePage />
  <div class="login-mask" role="dialog" aria-modal="true" aria-label="用户登录">
    <form class="login-modal" @submit.prevent="submit">
      <button class="modal-close" type="button" aria-label="关闭登录弹窗" @click="close">×</button>
      <div v-if="page === 'login'" class="tabs" role="tablist">
        <button type="button" :class="{ active: loginMode === 'password' }" @click="selectLoginMode('password')">密码登录</button>
        <button type="button" :class="{ active: loginMode === 'sms' }" @click="selectLoginMode('sms')">验证码登录</button>
      </div>
      <h1 v-else>注册账号</h1>

      <div v-if="page === 'login' && loginMode === 'password'" class="fields">
        <label class="field"><input v-model="phone" autocomplete="username" placeholder="请输入手机号" /></label>
        <label class="field"><input v-model="password" type="password" autocomplete="current-password" placeholder="请输入登录密码" /></label>
        <label class="code-field image-captcha-field" :class="{ verified: captchaAction === 'login_password' }">
          <input v-model="imageCaptchaAnswer" maxlength="8" autocomplete="off" placeholder="请输入图形验证码" :disabled="captchaAction === 'login_password'" />
          <button type="button" class="captcha-image" :disabled="imageCaptchaLoading || captchaBusy" title="换一张图形验证码" @click="refreshImageCaptcha">
            <img v-if="imageCaptcha" :src="imageCaptcha.imageUrl" alt="图形验证码，点击换一张" />
            <span v-else>{{ imageCaptchaLoading ? '加载中…' : '换一张' }}</span>
          </button>
        </label>
      </div>
      <div v-else-if="page === 'login'" class="fields">
        <label class="phone-field"><span>+86</span><input v-model="phone" inputmode="numeric" maxlength="11" placeholder="请输入手机号" /></label>
        <label class="code-field"><input v-model="smsCode" maxlength="8" inputmode="numeric" placeholder="请输入短信验证码" /><button type="button" :disabled="captchaBusy" @click="sendCode('login_sms')">{{ captchaBusy ? '验证中…' : '获取验证码' }}</button></label>
      </div>
      <div v-else class="fields">
        <label class="phone-field"><span>+86</span><input v-model="phone" inputmode="numeric" maxlength="11" placeholder="请输入手机号" /></label>
        <label class="code-field"><input v-model="smsCode" maxlength="8" inputmode="numeric" placeholder="请输入短信验证码" /><button type="button" :disabled="captchaBusy" @click="sendCode('register')">{{ captchaBusy ? '验证中…' : '获取验证码' }}</button></label>
        <label class="field"><input v-model="password" type="password" autocomplete="new-password" placeholder="请输入您的登录密码" /></label>
        <label class="field"><input v-model="confirmPassword" type="password" autocomplete="new-password" placeholder="再次输入登录密码" /></label>
      </div>

      <p v-if="message" class="message">{{ message }}</p>
      <label class="agreement"><input v-model="accepted" type="checkbox" /> <span>我已阅读并同意</span><RouterLink to="/agreements/sys_user_agree">《用户协议》</RouterLink><span>与</span><RouterLink to="/agreements/sys_userr_privacy">《隐私政策》</RouterLink></label>
      <button class="submit" type="submit" :disabled="loading">{{ loading ? '处理中…' : page === 'register' ? '注册' : '登录' }}</button>
      <p class="bottom-switch"><template v-if="page === 'login'">没有账号？<button type="button" @click="openRegister">立即注册</button></template><template v-else>已有账号？<button type="button" @click="openLogin">立即登录</button></template></p>
    </form>
  </div>
</template>

<style scoped>
.login-mask { position: fixed; z-index: 50; inset: 0; display: grid; place-items: center; padding: 28px; background: rgb(0 0 0 / 48%); }.login-modal { position: relative; width: min(100%, 688px); min-height: 588px; padding: 72px 70px 48px; background: #fff; }.modal-close { position: absolute; top: 20px; right: 22px; border: 0; background: transparent; color: #777; font-size: 32px; line-height: 1; }.tabs { display: flex; justify-content: center; gap: 92px; }.tabs button { border: 0; padding: 0 0 11px; background: transparent; color: #333; font-size: 28px; font-weight: 600; }.tabs button.active { color: #f13728; }.login-modal h1 { margin: 0 0 42px; color: #333; font-size: 30px; text-align: center; }.fields { display: grid; gap: 20px; margin-top: 44px; }.field, .phone-field, .code-field { display: flex; height: 76px; border: 1px solid #d9d9d9; }.field input, .phone-field input, .code-field input { min-width: 0; width: 100%; border: 0; padding: 0 24px; color: #333; font-size: 18px; outline: none; }.field:focus-within,.phone-field:focus-within,.code-field:focus-within { border-color: #f13728; }.phone-field span { display: grid; width: 100px; flex: 0 0 auto; place-items: center; border-right: 1px solid #d9d9d9; color: #666; font-size: 18px; }.code-field button { flex: 0 0 auto; border: 0; border-left: 1px solid #d9d9d9; background: #fff; color: #f13728; font-size: 16px; font-weight: 600; }.code-field button:disabled { color: #aaa; }.image-captcha-field { background: #fff; }.image-captcha-field.verified { border-color: #39a866; }.image-captcha-field.verified input { color: #26844e; }.captcha-image { width: 220px; padding: 0; overflow: hidden; background: #f8fbff !important; }.captcha-image img { width: 100%; height: 100%; object-fit: contain; }.captcha-image span { display: grid; height: 100%; place-items: center; color: #999; font-size: 14px; }.agreement { display: flex; flex-wrap: wrap; align-items: center; gap: 7px; margin-top: 25px; color: #444; font-size: 15px; }.agreement input { width: 20px; height: 20px; accent-color: #f13728; }.agreement a { color: #f13728; }.submit { width: 100%; height: 72px; margin-top: 28px; border: 0; background: #f13728; color: #fff; font-size: 19px; font-weight: 700; }.submit:disabled { opacity: .65; }.message { margin: 15px 0 -4px; color: #d9362b; font-size: 14px; }.bottom-switch { margin: 28px 0 0; color: #bbb; font-size: 18px; text-align: center; }.bottom-switch button { border: 0; padding: 0 0 0 13px; background: transparent; color: #f13728; font-size: 18px; }
@media (max-width: 760px) { .login-modal { min-height: auto; padding: 54px 30px 38px; }.tabs { gap: 42px; }.tabs button { font-size: 22px; }.fields { margin-top: 30px; }.field,.phone-field,.code-field { height: 58px; }.modal-close { top: 15px; right: 16px; }.submit { height: 58px; }.captcha-image { width: 176px; }.code-field:not(.image-captcha-field) button { width: 130px; font-size: 14px; } }
</style>
