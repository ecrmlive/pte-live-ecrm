<script setup lang="ts">
import { computed } from 'vue';

import formImg from '#/assets/diy/img/apply-form/formImg.png';
import uploadImg from '#/assets/diy/img/apply-form/uploadImg.png';

const props = withDefaults(
  defineProps<{
    /** 示例缩略图数量（对齐 CRMEB numConfig → picList） */
    count?: number;
  }>(),
  { count: 1 },
);

const slots = computed(() =>
  Math.max(0, Math.min(Math.floor(props.count || 0), 8)),
);
</script>

<template>
  <div class="upload-slots">
    <div v-for="n in slots" :key="n" class="upload-slots__pic">
      <img :src="formImg" alt="" />
    </div>
    <div class="upload-slots__pic">
      <img :src="uploadImg" alt="上传照片" />
    </div>
  </div>
</template>

<style scoped lang="scss">
.upload-slots {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  width: 100%;
}

.upload-slots__pic {
  width: 74px;
  height: 74px;
  margin: 10px 10px 0 0;

  img {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}
</style>
