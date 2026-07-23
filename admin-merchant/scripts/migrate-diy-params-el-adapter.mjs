/**
 * One-time: replace el-input / el-button / el-slider in DIY param panels with adapter/shared fields.
 */
import fs from 'node:fs';
import path from 'node:path';

const ROOT = path.resolve(import.meta.dirname, '../src/views/native/page/diy/params');

const FILES = [
  'Seckill.vue',
  'Preview.vue',
  'BargainProduct.vue',
  'Option.vue',
  'NewActivity.vue',
  'assembleProduct.vue',
  'Product.vue',
  'TopMerge.vue',
  'Title.vue',
  'Banner.vue',
  'Hotspot.vue',
  'ImageSingle.vue',
  'NavBar.vue',
  'Window.vue',
  'Order.vue',
  'QixiLive.vue',
  'Surface.vue',
];

const IMPORT_BLOCK = `import DiyColorField from './shared/diy-color-field.vue';
import DiyInputField from './shared/diy-input-field.vue';
import DiyLinkInputField from './shared/diy-link-input-field.vue';
import DiySliderField from './shared/diy-slider-field.vue';
import { useDiyAdapterComponents } from './shared/use-diy-adapter-components';
`;

function ensureImports(content) {
  if (content.includes('useDiyAdapterComponents')) {
    return content;
  }
  content = content.replace(
    /import \{ useDiyEditor \} from '\.\/shared\/use-diy-editor';/,
    `${IMPORT_BLOCK}import { useDiyEditor } from './shared/use-diy-editor';`,
  );
  if (!content.includes('useDiyAdapterComponents')) {
    content = content.replace(
      /import \{ computed,/,
      `${IMPORT_BLOCK}import { computed,`,
    );
  }
  return content;
}

function ensureAdapterHook(content) {
  if (content.includes('useDiyAdapterComponents()')) {
    return content;
  }
  const insertAfter = content.match(
    /const editor = useDiyEditor\(\);|defineOptions\([^)]+\);/,
  );
  if (insertAfter) {
    const idx = content.indexOf(insertAfter[0]) + insertAfter[0].length;
    return (
      content.slice(0, idx) +
      '\n\nconst { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();' +
      content.slice(idx)
    );
  }
  return content.replace(
    /defineOptions\(\{ name: '[^']+' \}\);/,
    (m) =>
      `${m}\n\nconst { PrimaryButton, RadioGroup, Checkbox } = useDiyAdapterComponents();`,
  );
}

function stripElImports(content) {
  content = content.replace(
    /\nimport \{([^}]*)\} from 'element-plus';/g,
    (match, inner) => {
      const kept = inner
        .split(',')
        .map((s) => s.trim())
        .filter(
          (s) =>
            s &&
            !['ElButton', 'ElInput', 'ElSlider', 'ElColorPicker'].includes(s),
        );
      if (kept.length === 0) {
        return '';
      }
      return `\nimport { ${kept.join(', ')} } from 'element-plus';`;
    },
  );
  return content;
}

function transformTemplate(content) {
  let t = content;

  // el-slider -> DiySliderField
  t = t.replace(/<el-slider\b/g, '<DiySliderField');
  t = t.replace(/<\/el-slider>/g, '</DiySliderField>');

  // el-input -> DiyInputField or DiyLinkInputField (link rows with suffix icon)
  t = t.replace(
    /<ElInput([^>]*?)v-model="([^"]+)"([^>]*?)@click="([^"]+)"([^>]*?)>\s*<template #suffix>[\s\S]*?<\/template>\s*<\/ElInput>/g,
    '<DiyLinkInputField$1v-model="$2"$3@click="$4"$5>\n<template #suffix>\n<ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>\n</template>\n</DiyLinkInputField>',
  );
  t = t.replace(
    /<el-input([^>]*?)v-model="([^"]+)"([^>]*?)@click="([^"]+)"([^>]*?)>\s*<template #suffix>[\s\S]*?<\/template>\s*<\/el-input>/g,
    '<DiyLinkInputField$1v-model="$2"$3@click="$4"$5>\n<template #suffix>\n<ElIcon color="#333" size="16px"><ArrowRight /></ElIcon>\n</template>\n</DiyLinkInputField>',
  );

  t = t.replace(/<ElInput\b/g, '<DiyInputField');
  t = t.replace(/<\/ElInput>/g, '</DiyInputField>');
  t = t.replace(/<el-input\b/g, '<DiyInputField');
  t = t.replace(/<\/el-input>/g, '</DiyInputField>');

  // el-button link reset -> PrimaryButton
  t = t.replace(
    /<el-button([^>]*?)type="primary"\s+link([^>]*?)@click\.stop="([^"]+)"([^>]*?)>\s*重置\s*<\/el-button>/g,
    '<component :is="PrimaryButton" link type="primary"$1$2 @click.stop="$3"$4>重置</component>',
  );
  t = t.replace(
    /<ElButton([^>]*?)type="primary"\s+link([^>]*?)@click\.stop="([^"]+)"([^>]*?)>\s*重置\s*<\/ElButton>/g,
    '<component :is="PrimaryButton" link type="primary"$1$2 @click.stop="$3"$4>重置</component>',
  );
  t = t.replace(
    /<el-button([^>]*?)@click\.stop="([^"]+)"\s+type="primary"\s+link([^>]*?)>\s*重置\s*<\/el-button>/g,
    '<component :is="PrimaryButton" link type="primary"$1 @click.stop="$2"$3>重置</component>',
  );

  // plain primary add buttons
  t = t.replace(
    /<ElButton plain type="primary" @click="([^"]+)">([^<]+)<\/ElButton>/g,
    '<component :is="PrimaryButton" plain @click="$1">$2</component>',
  );
  t = t.replace(
    /<el-button plain type="primary" @click="([^"]+)">([^<]+)<\/el-button>/g,
    '<component :is="PrimaryButton" plain @click="$1">$2</component>',
  );
  t = t.replace(
    /<ElButton plain size="small" type="primary" @click="([^"]+)">([^<]+)<\/ElButton>/g,
    '<component :is="PrimaryButton" plain size="small" @click="$1">$2</component>',
  );

  // style picker buttons (Order.vue pattern)
  t = t.replace(
    /<ElButton\s+v-for="([^"]+)"\s+:key="([^"]+)"\s+size="small"\s+:type="([^"]+)"\s+@click="([^"]+)"\s*>\s*\{\{([^}]+)\}\}\s*<\/ElButton>/g,
    '<component :is="PrimaryButton" v-for="$1" :key="$2" size="small" :type="$3" @click="$4">{{ $5 }}</component>',
  );

  // danger link delete
  t = t.replace(
    /<ElButton link size="small" type="danger" @click\.stop="([^"]+)">\s*删除\s*<\/ElButton>/g,
    '<component :is="PrimaryButton" link size="small" class="!text-destructive" @click.stop="$1">删除</component>',
  );

  // single-row color: picker + input + reset -> DiyColorField
  t = t.replace(
    /<div class="flex-1 d-s-c" style="height: 36px;">\s*<el-color-picker show-alpha size="default" v-model="([^"]+)"><\/el-color-picker>\s*<el-input class="ml10" v-model="([^"]+)"(?: placeholder="[^"]*")? \/?>\s*<el-button style="margin-left: 10px;"[\s\S]*?@click\.stop="([^"]+)"[\s\S]*?>\s*重置\s*<\/el-button>\s*<\/div>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" @click.stop />',
  );

  return t;
}

for (const file of FILES) {
  const fp = path.join(ROOT, file);
  if (!fs.existsSync(fp)) {
    console.warn('skip missing', file);
    continue;
  }
  let content = fs.readFileSync(fp, 'utf8');
  const before = content;
  content = ensureImports(content);
  content = ensureAdapterHook(content);
  content = stripElImports(content);
  const parts = content.split('<template>');
  if (parts.length === 2) {
    content = parts[0] + '<template>' + transformTemplate(parts[1]);
  }
  if (content !== before) {
    fs.writeFileSync(fp, content);
    console.log('updated', file);
  } else {
    console.log('unchanged', file);
  }
}
