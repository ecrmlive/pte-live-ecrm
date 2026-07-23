/**
 * P6b pass2: el-radio-group / el-checkbox / inline el-color-picker in DIY params.
 */
import fs from 'node:fs';
import path from 'node:path';

const ROOT = path.resolve(import.meta.dirname, '../src/views/native/page/diy/params');

const FILES = [
  'Title.vue',
  'BargainProduct.vue',
  'Product.vue',
  'TopMerge.vue',
  'Preview.vue',
  'Seckill.vue',
  'assembleProduct.vue',
  'Option.vue',
  'NewActivity.vue',
];

function ensureDiyColorImport(content) {
  if (content.includes("import DiyColorField from './shared/diy-color-field.vue'")) {
    return content;
  }
  return content.replace(
    /import DiyInputField from '\.\/shared\/diy-input-field\.vue';/,
    "import DiyColorField from './shared/diy-color-field.vue';\nimport DiyInputField from './shared/diy-input-field.vue';",
  );
}

function transformTemplate(t) {
  // radio-group
  t = t.replace(/<el-radio-group\b/g, '<component :is="RadioGroup"');
  t = t.replace(/<\/el-radio-group>/g, '</component>');

  // checkbox (self-closing and paired)
  t = t.replace(/<el-checkbox\b/g, '<component :is="Checkbox"');
  t = t.replace(/<\/el-checkbox>/g, '</component>');

  // standalone color row: picker + input + reset -> DiyColorField (common pattern)
  t = t.replace(
    /<div class="flex-1 d-s-c" style="height: 36px;">\s*<el-color-picker(?: show-alpha)? size="default" v-model="([^"]+)"><\/el-color-picker>\s*<DiyInputField class="ml10" v-model="\1" placeholder="([^"]*)" \/?>\s*<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"?\s*@click\.stop="editor\.onEditorResetColor\([^)]+\)"[^>]*>\s*重置\s*<\/component>\s*<\/div>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" placeholder="$2" />',
  );

  // picker + input without reset
  t = t.replace(
    /<div class="flex-1 d-s-c" style="height: 36px;">\s*<el-color-picker(?: show-alpha)? size="default" v-model="([^"]+)"><\/el-color-picker>\s*<DiyInputField class="ml10" v-model="\1" placeholder="([^"]*)" \/?>\s*<\/div>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" placeholder="$2" />',
  );

  // inline picker in view wrapper (gradient rows) - replace picker only, keep input+reset
  t = t.replace(
    /<view class="ml10"><el-color-picker show-alpha size="default"\s+v-model="([^"]+)"><\/el-color-picker><\/view>/g,
    '<view class="ml10"><DiyColorField v-model="$1" default-color="#FF4C01" class="diy-inline-color" /></view>',
  );

  // single standalone picker with adjacent input (Title.vue style)
  t = t.replace(
    /<el-color-picker size="default" v-model="([^"]+)"><\/el-color-picker>\s*<DiyInputField class="ml10" v-model="\1" placeholder="([^"]*)" \/?>\s*<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"?\s*@click\.stop="editor\.onEditorResetColor\([^)]+\)"[^>]*>\s*重置\s*<\/component>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" placeholder="$2" />',
  );

  // remaining standalone pickers
  t = t.replace(
    /<el-color-picker show-alpha size="default" v-model="([^"]+)"><\/el-color-picker>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" />',
  );
  t = t.replace(
    /<el-color-picker size="default" v-model="([^"]+)"><\/el-color-picker>/g,
    '<DiyColorField v-model="$1" default-color="#ffffff" />',
  );

  // multiline color picker + input + reset row
  t = t.replace(
    /<el-color-picker show-alpha size="default"\s+v-model="([^"]+)"><\/el-color-picker>\s*<DiyInputField class="ml10" v-model="\1" placeholder="([^"]*)" \/?>\s*<component :is="PrimaryButton" link type="primary" style="margin-left: 10px;"\s+@click\.stop="editor\.onEditorResetColor\([^,]+,\s*'[^']+',\s*'([^']+)'\)">重置<\/component>/gs,
    '<DiyColorField v-model="$1" default-color="$3" placeholder="$2" />',
  );

  // multiline standalone show-alpha picker
  t = t.replace(
    /<el-color-picker show-alpha size="default"\s+v-model="([^"]+)"><\/el-color-picker>/gs,
    '<DiyColorField v-model="$1" default-color="#ffffff" />',
  );

  // multiline default size picker in div.ml10
  t = t.replace(
    /<div class="ml10"><el-color-picker size="default"\s+v-model="([^"]+)"><\/el-color-picker><\/div>/gs,
    '<DiyColorField v-model="$1" default-color="#ffffff" class="ml10" />',
  );

  // view.ml10 picker wrappers (Option.vue gradient rows — inputs already present)
  t = t.replace(
    /<view class="ml10"><el-color-picker size="default"\s+v-model="([^"]+)"><\/el-color-picker><\/view>\s*/gs,
    '',
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
  content = ensureDiyColorImport(content);
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
