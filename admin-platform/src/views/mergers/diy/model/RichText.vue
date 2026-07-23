<template>


    <div @click.stop="diyEditer(index)">
      <div class="drag optional" :class="{selected: index === selectedIndex}">
        <div class="diy-richText" :style="{background: item.style.background, padding: item.style.paddingTop + 'px ' + item.style.paddingLeft + 'px'}"
          v-html="sanitizeHtml(item.params.content)">
        </div>
        <div class="btn-edit-del">
          <div class="btn-del" @click.stop="diyDeleteItem(index)">删除</div>
        </div>
      </div>
    </div>

</template>

<script>
import { sanitizeHtml } from '#/utils/xss-filter';

export default {
  inject: ['diyModel'],
	data() {
		return {

		};
	},
	props: ['item', 'index', 'selectedIndex'],
	methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },

		sanitizeHtml
	}
};
</script>

<style lang="scss" scoped>
.diy-richText video{ width: 100%;}
</style>
