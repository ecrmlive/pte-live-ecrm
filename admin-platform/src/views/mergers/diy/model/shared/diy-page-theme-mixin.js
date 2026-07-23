import { getPageThemeApi } from '#/api/core/page';
import { getThemeColorByIndex } from '#/utils/shop-theme';

/** Editor preview: title + more text follow page theme (same as QixiLive). */
export default {
	data() {
		return {
			themeColor: getThemeColorByIndex('0'),
		};
	},
	created() {
		this.loadDiyPageThemeColor();
	},
	methods: {
		loadDiyPageThemeColor() {
			getPageThemeApi()
				.then((res) => {
					const theme = res.vars?.values?.theme ?? '0';
					this.themeColor = getThemeColorByIndex(theme);
				})
				.catch(() => {});
		},
	},
};
