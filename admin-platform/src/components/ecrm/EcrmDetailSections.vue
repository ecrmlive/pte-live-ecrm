<script setup lang="ts">
export interface EcrmDetailSection {
  title: string;
  items: Array<{ label: string; value: string | number | null | undefined }>;
}

defineProps<{
  column?: number;
  sections: EcrmDetailSection[];
}>();
</script>

<template>
  <div class="ecrm-detail-sections flex flex-col gap-4">
    <div v-for="(section, index) in sections" :key="index">
      <h3 class="mb-2 text-sm font-medium text-foreground">{{ section.title }}</h3>
      <el-descriptions :column="column || 2" border>
        <el-descriptions-item
          v-for="(item, itemIndex) in section.items"
          :key="itemIndex"
          :label="item.label"
        >
          {{ item.value ?? '—' }}
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>
