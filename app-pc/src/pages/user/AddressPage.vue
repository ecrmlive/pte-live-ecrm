<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import AccountFrame from "@/components/AccountFrame.vue";
import { createAddress, fetchAddresses, removeAddress, updateAddress, type Address } from "@/api/trade";
import { ApiError } from "@/utils/request";

type AddressForm = Pick<Address, "real_name" | "phone" | "province" | "city" | "district" | "detail" | "region_code" | "post_code" | "is_default">;
const addresses = ref<Address[]>([]);
const editingID = ref<number | null>(null);
const editorOpen = ref(false);
const loading = ref(false);
const hint = ref("");
const form = reactive<AddressForm>({ real_name: "", phone: "", province: "", city: "", district: "", detail: "", region_code: "", post_code: 0, is_default: 0 });
const actionTitle = computed(() => editingID.value ? "编辑收货地址" : "新增收货地址");

function resetForm() {
  editingID.value = null;
  Object.assign(form, { real_name: "", phone: "", province: "", city: "", district: "", detail: "", region_code: "", post_code: 0, is_default: addresses.value.length ? 0 : 1 });
}

async function load() {
  loading.value = true;
  try { const data = await fetchAddresses(); addresses.value = data.list || []; hint.value = ""; }
  catch (error) { hint.value = error instanceof ApiError ? error.message : "收货地址加载失败"; }
  finally { loading.value = false; }
}

function create() { resetForm(); editorOpen.value = true; }
function edit(address: Address) { editingID.value = address.address_id; Object.assign(form, { ...address, region_code: address.region_code || "" }); editorOpen.value = true; }
function close() { editorOpen.value = false; resetForm(); }

async function save() {
  if (!form.real_name.trim() || !form.phone.trim() || !form.detail.trim()) { hint.value = "请填写收货人、手机号和详细地址"; return; }
  try {
    if (editingID.value) await updateAddress(editingID.value, { ...form });
    else await createAddress({ ...form });
    close(); await load();
  } catch (error) { hint.value = error instanceof ApiError ? error.message : "保存地址失败"; }
}

async function remove(address: Address) {
  if (!window.confirm(`确认删除“${address.real_name}”的收货地址吗？`)) return;
  try { await removeAddress(address.address_id); await load(); }
  catch (error) { hint.value = error instanceof ApiError ? error.message : "删除地址失败"; }
}

async function makeDefault(address: Address) {
  if (address.is_default) return;
  try { await updateAddress(address.address_id, { is_default: 1 }); await load(); }
  catch (error) { hint.value = error instanceof ApiError ? error.message : "设置默认地址失败"; }
}

onMounted(() => void load());
</script>

<template>
  <AccountFrame>
    <template #crumb><span>›</span> 地址管理</template>
    <header class="address-header"><h1>地址管理</h1><button class="pc-btn" type="button" @click="create">新增收货地址</button></header>
    <p v-if="hint" class="hint">{{ hint }}</p>
    <div v-if="loading" class="address-empty">正在加载地址…</div>
    <div v-else class="address-grid">
      <article v-for="address in addresses" :key="address.address_id" class="address-card" :class="{ default: address.is_default === 1 }">
        <b v-if="address.is_default" class="default-mark">默认</b>
        <h2>{{ address.real_name }} <span>{{ address.phone }}</span></h2>
        <p>{{ [address.province, address.city, address.district, address.detail].filter(Boolean).join(" ") }}</p>
        <footer><button v-if="!address.is_default" type="button" @click="makeDefault(address)">设为默认</button><button type="button" @click="edit(address)">编辑</button><button type="button" @click="remove(address)">删除</button></footer>
      </article>
      <button class="address-create" type="button" @click="create"><span>＋</span>添加新地址</button>
      <p v-if="!addresses.length" class="address-empty">还没有收货地址，点击右侧卡片添加。</p>
    </div>

    <div v-if="editorOpen" class="dialog-mask" @click.self="close">
      <form class="address-dialog" @submit.prevent="save">
        <header><h2>{{ actionTitle }}</h2><button type="button" @click="close">×</button></header>
        <label>收货人<input v-model="form.real_name" maxlength="32" placeholder="请输入收货人姓名" /></label>
        <label>手机号码<input v-model="form.phone" maxlength="32" placeholder="请输入手机号码" /></label>
        <div class="region-row"><label>省<input v-model="form.province" maxlength="32" placeholder="省" /></label><label>市<input v-model="form.city" maxlength="32" placeholder="市" /></label><label>区/县<input v-model="form.district" maxlength="32" placeholder="区/县" /></label></div>
        <label>详细地址<textarea v-model="form.detail" maxlength="255" placeholder="街道、门牌号等" /></label>
        <label class="checkbox"><input v-model="form.is_default" :true-value="1" :false-value="0" type="checkbox" /> 设为默认收货地址</label>
        <footer><button class="plain-btn" type="button" @click="close">取消</button><button class="pc-btn" type="submit">保存地址</button></footer>
      </form>
    </div>
  </AccountFrame>
</template>

<style scoped>
.address-header { display: flex; justify-content: space-between; align-items: center; padding-bottom: 19px; border-bottom: 1px solid #eee; }.address-header h1 { margin: 0; font-size: 20px; }.address-header .pc-btn { border-radius: 3px; padding: 9px 16px; }.hint { margin: 16px 0 0; color: #d9362b; }.address-grid { display: flex; flex-wrap: wrap; gap: 28px; padding-top: 32px; }.address-card, .address-create { position: relative; width: 290px; min-height: 205px; padding: 25px 28px; border: 1px solid #e7e7e7; background: #fff; text-align: left; }.address-card.default { border-color: #f8b2ac; }.default-mark { position: absolute; top: 0; right: 0; padding: 5px 13px; background: #f13728; color: #fff; font-size: 12px; }.address-card h2 { margin: 0 0 16px; font-size: 16px; font-weight: 500; }.address-card h2 span { padding-left: 9px; color: #777; font-size: 13px; font-weight: 400; }.address-card p { min-height: 44px; margin: 0; color: #777; font-size: 14px; line-height: 1.8; }.address-card footer { display: flex; gap: 15px; margin-top: 14px; padding-top: 12px; border-top: 1px dashed #e7e7e7; }.address-card button { border: 0; padding: 0; color: #555; background: transparent; font-size: 13px; }.address-card button:hover { color: #f13728; }.address-create { display: grid; place-content: center; justify-items: center; gap: 8px; color: #bbb; cursor: pointer; }.address-create span { font-size: 42px; line-height: 1; font-weight: 200; }.address-create:hover { border-color: #f13728; color: #f13728; }.address-empty { width: 100%; margin: 20px 0; color: #aaa; text-align: center; }.dialog-mask { position: fixed; z-index: 50; inset: 0; display: grid; place-items: center; padding: 20px; background: rgb(0 0 0 / 38%); }.address-dialog { width: min(100%, 520px); padding: 24px 28px; background: #fff; box-shadow: 0 16px 46px rgb(0 0 0 / 24%); }.address-dialog > header { display: flex; align-items: center; justify-content: space-between; padding-bottom: 18px; border-bottom: 1px solid #eee; }.address-dialog h2 { margin: 0; font-size: 18px; }.address-dialog header button { border: 0; background: transparent; color: #999; font-size: 28px; }.address-dialog > label { display: grid; gap: 7px; margin-top: 16px; color: #555; font-size: 14px; }.address-dialog input:not([type="checkbox"]), .address-dialog textarea { width: 100%; border: 1px solid #ddd; padding: 10px; outline: none; }.address-dialog textarea { min-height: 88px; resize: vertical; }.region-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }.region-row label { display: grid; gap: 7px; color: #555; font-size: 14px; }.address-dialog .checkbox { display: block; }.address-dialog footer { display: flex; justify-content: flex-end; gap: 10px; margin-top: 22px; }.plain-btn { border: 1px solid #ddd; padding: 8px 16px; background: #fff; color: #666; }
</style>
