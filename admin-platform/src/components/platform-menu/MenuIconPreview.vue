<script lang="ts" setup>
import { computed } from 'vue';

import { Icon as IconifyIcon } from '@iconify/vue';

import { normalizePlatformMenuIcon } from '#/constants/platform-lucide-icons';

const props = defineProps<{
  icon?: string;
  size?: number;
}>();

const resolved = computed(() => normalizePlatformMenuIcon(props.icon));
</script>

<template>
  <IconifyIcon
    v-if="resolved"
    :icon="resolved"
    :style="{ fontSize: `${size ?? 18}px`, verticalAlign: 'middle' }"
    class="text-foreground/80"
  />
  <span v-else class="text-muted-foreground">—</span>
</template>
