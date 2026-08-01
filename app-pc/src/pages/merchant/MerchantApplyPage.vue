<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { createMerchantApplication, fetchMyMerchantApplications, type MerchantApplication } from "@/api/merchant";
import { uploadMerchantLicense } from "@/api/upload";
import { ApiError } from "@/utils/request";
import { useUserStore } from "@/stores/user";

const user = useUserStore();
const route = useRoute();
const router = useRouter();
const hint = ref("");
const sending = ref(false);
const history = ref<MerchantApplication[]>([]);
const accepted = ref(false);
const form = reactive({ merchant_name: "", contact_name: "", contact_mobile: "", category_name: "", merchant_type: "", license_key: "" });
const uploading = ref(false);
const licenseName = ref("");

async function loadHistory() {
  if (!user.isLogin) return;
  try { history.value = (await fetchMyMerchantApplications()).list || []; } catch { /* 首次访问无需阻断申请页 */ }
}

async function submit() {
  if (!accepted.value) { hint.value = "请阅读并同意《入驻协议》"; return; }
  sending.value = true;
  try {
    const created = await createMerchantApplication({ ...form });
    history.value = [created, ...history.value];
    hint.value = "申请已提交，平台审核后会通过站内消息通知。";
  } catch (e) {
    hint.value = e instanceof ApiError ? e.message : "提交申请失败";
  } finally { sending.value = false; }
}

async function selectLicense(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;
  if (!/\.(jpe?g|png|webp)$/i.test(file.name) || file.size > 10 * 1024 * 1024) {
    hint.value = "营业执照仅支持 JPG、PNG、WebP，且不能超过 10MB。";
    return;
  }
  uploading.value = true;
  hint.value = "";
  try {
    form.license_key = await uploadMerchantLicense(file);
    licenseName.value = file.name;
  } catch (e) {
    form.license_key = "";
    licenseName.value = "";
    hint.value = e instanceof ApiError ? e.message : "营业执照上传失败";
  } finally {
    uploading.value = false;
    (event.target as HTMLInputElement).value = "";
  }
}

onMounted(() => {
  if (!user.isLogin) { router.replace({ name: "login", query: { redirect: route.fullPath } }); return; }
  void loadHistory();
});
</script>

<template>
  <div class="apply-page">
    <section class="apply-panel pc-container">
      <h1>入驻申请</h1>
      <p class="lead">请如实填写经营资料，平台将在审核后开通店铺管理权限。</p>
      <form class="apply-form" @submit.prevent="submit">
        <label><b>* 店铺名称</b><input v-model.trim="form.merchant_name" maxlength="128" placeholder="请输入店铺名称" /></label>
        <label><b>* 用户姓名</b><input v-model.trim="form.contact_name" maxlength="64" placeholder="请输入联系人姓名" /></label>
        <label><b>* 手机号码</b><input v-model.trim="form.contact_mobile" maxlength="32" inputmode="tel" placeholder="请输入手机号码" /></label>
        <label><b>* 店铺分类</b><select v-model="form.category_name"><option value="" disabled>请选择店铺分类</option><option>服饰鞋包</option><option>家居生活</option><option>数码家电</option><option>美妆个护</option><option>食品生鲜</option><option>运动户外</option></select></label>
        <label><b>* 店铺类型</b><select v-model="form.merchant_type"><option value="" disabled>请选择店铺类型</option><option>旗舰店</option><option>专营店</option><option>会员店</option><option>加盟店</option><option>工厂店</option><option>自营</option></select></label>
        <label class="license"><b>* 营业执照</b><span><input type="file" accept="image/jpeg,image/png,image/webp" :disabled="uploading" @change="selectLicense" /><small v-if="!licenseName">选择图片后将直传 COS，系统仅保存对象 Key。</small><small v-else class="uploaded">已上传：{{ licenseName }}</small></span></label>
        <label class="agreement"><input v-model="accepted" type="checkbox" /> 我已阅读并同意 <RouterLink to="/agreements/business_entry_agree">《入驻协议》</RouterLink></label>
        <p v-if="hint" class="hint">{{ hint }}</p>
        <button class="submit" type="submit" :disabled="sending">{{ sending ? "提交中…" : "提交申请" }}</button>
      </form>
    </section>
    <section v-if="history.length" class="history pc-container"><h2>我的入驻申请</h2><article v-for="item in history" :key="item.id"><b>{{ item.merchant_name }}</b><span>{{ item.category_name }} · {{ item.merchant_type }}</span><em :class="item.status">{{ item.status === "approved" ? "已通过" : item.status === "rejected" ? "已驳回" : "审核中" }}</em></article></section>
  </div>
</template>

<style scoped>
.apply-page { min-height: 680px; padding: 40px 0 76px; background: #f6f6f6; }.apply-panel { max-width: 1040px; padding: 34px 110px 44px; background: #fff; }.apply-panel h1 { margin: 0; color: #333; font-size: 24px; text-align: center; }.lead { margin: 17px 0 28px; padding-bottom: 23px; border-bottom: 1px solid #eee; color: #999; text-align: center; }.apply-form { width: min(100%, 620px); margin: 0 auto; }.apply-form label { display: grid; grid-template-columns: 116px minmax(0, 1fr); align-items: center; gap: 14px; margin: 17px 0; }.apply-form label > b { color: #555; font-weight: 400; text-align: right; }.apply-form label > b::first-letter { color: #f13728; }.apply-form input:not([type="checkbox"]), .apply-form select { box-sizing: border-box; width: 100%; height: 42px; border: 1px solid #dfe3ea; padding: 0 13px; color: #444; background: #fff; outline-color: #f13728; }.apply-form input[type="file"] { height: auto; padding: 9px 13px; }.apply-form select { appearance: auto; }.license > span { display: grid; gap: 7px; }.license small { color: #aaa; font-size: 12px; }.license .uploaded { color: #2c9b62; }.agreement { display: flex !important; grid-template-columns: none !important; justify-content: center; margin: 28px 0 19px !important; color: #777; font-size: 14px; }.agreement input { margin: 0 7px 0 0; }.agreement a { color: #f13728; }.hint { margin: 0 0 14px; color: #d9362b; text-align: center; }.submit { display: block; width: 160px; margin: 0 auto; padding: 12px; border: 0; border-radius: 3px; color: #fff; background: #f13728; font-weight: 600; cursor: pointer; }.submit:disabled { opacity: .65; cursor: wait; }.history { max-width: 1040px; margin-top: 18px; padding: 24px 40px; background: #fff; }.history h2 { margin: 0 0 16px; font-size: 18px; }.history article { display: grid; grid-template-columns: 1fr 1fr auto; gap: 14px; padding: 14px 0; border-top: 1px solid #eee; color: #666; }.history em { color: #e79a20; font-style: normal; }.history em.approved { color: #2c9b62; }.history em.rejected { color: #d9362b; }@media (max-width: 680px) { .apply-panel { padding: 28px 20px; }.apply-form label { grid-template-columns: 1fr; gap: 7px; }.apply-form label > b { text-align: left; }.history article { grid-template-columns: 1fr; gap: 6px; } }
</style>
