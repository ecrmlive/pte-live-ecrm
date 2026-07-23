<template>
	<div class="diy-option" :style="{
			paddingLeft: item.style.paddingLeft + 'px',
			paddingRight: item.style.paddingLeft + 'px',
			paddingTop: item.style.paddingTop + 'px',
			paddingBottom: item.style.paddingBottom + 'px',
			marginTop:item.style.marginTop + 'px',
			background: item.style.background
		}" :class="{ selected: index === selectedIndex }" @click.stop="diyEditer(index)">
		<div class="list" :class="`optionType${item.params.type || 1}`"
			:style="{borderTopLeftRadius: item.style.topRadio + 'px',
			borderTopRightRadius: item.style.topRadio + 'px',
			borderBottomLeftRadius: item.style.bottomRadio + 'px',
			borderBottomRightRadius: item.style.bottomRadio + 'px',
			backgroundImage: 'linear-gradient(to right, ' + (item.style.bgcolor_color1 || '#fff') + ', ' + (item.style.bgcolor_color2 || '#fff') + ')'}">
			<div class="diy-option item active" :style="{
					borderImageSource:item.params.type==1?`linear-gradient(to right,${item.style.active_color1 || '#fff'} ,${ item.style.active_color2 || '#fff'})`:'',
					borderImageSlice:1,
					backgroundImage: item.params.type==2?`linear-gradient(to right,${item.style.active_color1 || '#fff'} ,${ item.style.active_color2 || '#fff'})` : '',
					color:item.style.activeText || '#333'}">
				<div v-if="item.params.type==3" class="op3"
					:style="{borderColor:`${item.style.active_color1 || '#fff'}`}">
				</div>
				首页
			</div>
			<template v-for="(citem,cindex) in item.data" :key="cindex">
				<div class="diy-option item" v-if="cindex!=0">
					{{citem.text}}
				</div>
			</template>

		</div>
	</div>
</template>

<script>
	export default {
  inject: ['diyModel'],
		data() {
			return {};
		},
		props: ['item', 'index', 'selectedIndex'],
		methods: {
  diyEditer(index) {
    this.diyModel?.onEditer(index);
  },
  diyDeleteItem(index) {
    this.diyModel?.onDeleleItem(index);
  },
}
	};
</script>

<style lang="scss" scoped>
	.diy-option {
		.list {
			display: flex;
			flex-wrap: nowrap;
			overflow-x: auto;
			justify-content: flex-start;
			align-items: center;
			height: 43px;
			padding-left: 24px;

			.item {
				font-size: 16px;
				margin-right: 5px;
				color: #333;
				display: flex;
				justify-content: center;
				align-items: center;
				padding: 0 15px;
				line-height: 26px;
				height: 26px;
				position: relative;
				flex-shrink: 0;
			}


		}

		.optionType1 {
			.item {
				line-height: 32px;
				height: 32px;
				padding: 0;
				margin-right: 18px;
			}

			.item.active {
				font-weight: 600;
				font-size: 16px;
				border-bottom: 2px solid;
			}
		}

		.optionType2 {
			.item.active {
				color: #fff;
				border-radius: 100px;
			}
		}

		.optionType3 {
			.item.active {
				font-weight: 600;
				font-size: 16px;
			}
		}

		.op3 {
			position: absolute;
			bottom: 0;
			left: 0;
			right: 0;
			margin: auto;
			width: 30px;
			height: 30px;
			border: 3px solid;
			border-left: 3px solid transparent !important;
			border-top: 3px solid transparent !important;
			border-right: 3px solid transparent !important;
			border-radius: 50%;
			bottom: -4px;
		}
	}
</style>