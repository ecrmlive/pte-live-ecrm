package diy

import (
	_ "embed"
	"encoding/json"
)

//go:embed defaults.json
var defaultsJSON []byte

type defaultsFile struct {
	DefaultItems       map[string]any `json:"defaultItems"`
	DefaultPage        map[string]any `json:"defaultPage"`
	CenterDefaultItems map[string]any `json:"centerDefaultItems"`
}

func loadDefaults() defaultsFile {
	var d defaultsFile
	_ = json.Unmarshal(defaultsJSON, &d)
	if d.DefaultItems == nil {
		d.DefaultItems = map[string]any{}
	}
	if d.DefaultPage == nil {
		d.DefaultPage = map[string]any{}
	}
	d.DefaultItems["discountGroup"] = defaultDiscountGroup()
	d.DefaultItems["bargainProduct"] = defaultBargainProduct()
	d.DefaultItems["previewProduct"] = defaultPreviewProduct()
	d.DefaultItems["bottomNav"] = defaultBottomNav()
	d.DefaultItems["ranking"] = defaultRanking()
	d.DefaultItems["community"] = defaultCommunity()
	return d
}

func defaultCommunity() map[string]any {
	post := func(title, author string, images int) map[string]any {
		return map[string]any{"title": title, "author": author, "images": images, "image": ""}
	}
	return map[string]any{
		"name": "种草社区", "type": "community", "group": "marketing", "icon": "icon-zhongcaoshequ",
		"params": map[string]any{"title": "种草社区", "more": "好物分享", "titleType": "text", "titleImage": "", "layout": "scroll", "showNum": 3, "showTitle": true, "showAvatar": true, "showAuthor": true},
		"style":  map[string]any{"background": "#f5f5f5", "cardBackground": "#ffffff", "headStart": "#e93323", "headEnd": "#ff7931", "buttonColor": "#ffffff", "buttonSize": 12, "contentGap": 10, "contentRadius": 8, "contentShadow": "off", "paddingTop": 10, "paddingBottom": 10, "paddingLeft": 10, "marginTop": 10, "radius": 10, "shadow": "off"},
		"data":   []any{post("把春天装进生活里", "浅笑回眸", 8), post("像我这种乐形身材又怎样", "国宝小熊猫", 6), post("发现日常里的美好瞬间", "阿秋", 3)},
	}
}

func defaultRanking() map[string]any {
	product := func() map[string]any {
		return map[string]any{"name": "幸运美物", "image": "", "price": "350.00"}
	}
	board := func(title string) map[string]any {
		return map[string]any{"icon": "🔥", "title": title, "products": []any{product(), product(), product()}}
	}
	return map[string]any{
		"name": "排行榜", "type": "ranking", "group": "shop", "icon": "icon-paihangbang",
		"params": map[string]any{"title": "排行榜", "more": "更多", "titleType": "text", "titleImage": "", "boardCount": 2, "productCount": 3},
		"style": map[string]any{
			"background": "#f5f5f5", "cardBackground": "#ffffff", "boardBackground": "#fceae9", "boardTitleColor": "#ff4c8d", "moreColor": "#999999", "priceColor": "#ff4c8d",
			"paddingTop": 10, "paddingBottom": 10, "paddingLeft": 10, "marginTop": 10, "radius": 10, "boardRadius": 8, "shadow": "off",
		},
		"data": []any{board("销量榜"), board("好评榜")},
	}
}

func defaultBottomNav() map[string]any {
	return map[string]any{
		"name":   "底部导航",
		"type":   "bottomNav",
		"group":  "tools",
		"icon":   "icon-daohanglan",
		"params": map[string]any{"activeIndex": 0},
		"style": map[string]any{
			"navigationType": "icon-text", "positionType": "fixed", "themeMode": "system",
			"activeColor": "#f62c2c", "textColor": "#282828", "background": "rgba(255,255,255,0.96)",
			"paddingTop": 0, "paddingBottom": 0, "pagePadding": 0, "bottomSpacing": 0, "radius": 18,
		},
		"data": []any{
			map[string]any{"text": "首页", "linkUrl": "/pages/index/index", "selectedImgUrl": "", "unselectedImgUrl": "", "icon": "⌂", "hide": false},
			map[string]any{"text": "分类", "linkUrl": "/pages/category/index", "selectedImgUrl": "", "unselectedImgUrl": "", "icon": "□", "hide": false},
			map[string]any{"text": "购物车", "linkUrl": "/pages/cart/index", "selectedImgUrl": "", "unselectedImgUrl": "", "icon": "⌑", "hide": false},
			map[string]any{"text": "我的", "linkUrl": "/pages/user/index", "selectedImgUrl": "", "unselectedImgUrl": "", "icon": "◯", "hide": false},
		},
	}
}

func defaultBargainProduct() map[string]any {
	return map[string]any{
		"name":  "助力",
		"type":  "bargainProduct",
		"group": "marketing",
		"icon":  "icon-kanjia1",
		"params": map[string]any{
			"showNum":      3,
			"title":        "疯狂砍价",
			"desc":         "低至0元免费拿",
			"more":         "更多",
			"btntext":      "去砍价",
			"titleBgType":  2,
			"titleType":    1,
			"column":       1,
			"productName":  1,
			"productSales": 1,
			"productPrice": 1,
			"linePrice":    1,
			"product_btn":  1,
			"bgimage":      "",
			"titleimage":   "",
		},
		"style": map[string]any{
			"background":          "#F5F5F5",
			"bgcolor_color1":      "#FFFFFF",
			"bgcolor_color2":      "#FFFFFF",
			"paddingTop":          10,
			"paddingBottom":       0,
			"paddingLeft":         10,
			"paddingRight":        10,
			"marginTop":           10,
			"topRadio":            10,
			"bottomRadio":         10,
			"titleBg_color1":      "#FF3B30",
			"titleBg_color2":      "#FF6A2A",
			"titleColor":          "#FFFFFF",
			"descColor":           "rgba(255,255,255,.85)",
			"moreColor":           "#FFFFFF",
			"titleSize":           18,
			"titleWeight":         1,
			"moreSize":            12,
			"productTopRadio":     8,
			"productBottomRadio":  8,
			"productBgColor":      "#FFFFFF",
			"cardShadow":          "off",
			"productName_color":   "#333333",
			"nameWeight":          0,
			"productLine_color":   "#999999",
			"productPrice_color":  "#FF4C75",
			"product_sales_color": "#FF4C75",
			"productBtn_color1":   "#FF4C75",
			"productBtn_color2":   "#FF4C75",
			"btn_text_color":      "#FFFFFF",
		},
		"defaultData": []any{
			map[string]any{"product_name": "无线蓝牙耳机", "image": "", "bargain_price": "199.00", "ot_price": "299.00", "sales": 1288},
			map[string]any{"product_name": "家用护腰靠枕", "image": "", "bargain_price": "159.00", "ot_price": "229.00", "sales": 996},
			map[string]any{"product_name": "便携保温水杯", "image": "", "bargain_price": "89.00", "ot_price": "129.00", "sales": 768},
		},
		"data": []any{},
	}
}

func defaultPreviewProduct() map[string]any {
	return map[string]any{
		"name":  "预售",
		"type":  "previewProduct",
		"group": "marketing",
		"icon":  "icon-yushoucuifu",
		"params": map[string]any{
			"showNum":      3,
			"title":        "预售专区",
			"desc":         "火热预定中",
			"more":         "更多",
			"btntext":      "立即预定",
			"titleBgType":  2,
			"titleType":    1,
			"column":       1,
			"productName":  1,
			"productSales": 0,
			"productPrice": 1,
			"linePrice":    1,
			"product_btn":  1,
			"bgimage":      "",
			"titleimage":   "",
		},
		"style": map[string]any{
			"background":          "#F5F5F5",
			"bgcolor_color1":      "#FFFFFF",
			"bgcolor_color2":      "#FFFFFF",
			"paddingTop":          10,
			"paddingBottom":       0,
			"paddingLeft":         10,
			"paddingRight":        10,
			"marginTop":           10,
			"topRadio":            10,
			"bottomRadio":         10,
			"titleBg_color1":      "#FF3B7A",
			"titleBg_color2":      "#FF7A45",
			"titleColor":          "#FFFFFF",
			"descColor":           "rgba(255,255,255,.85)",
			"moreColor":           "#FFFFFF",
			"titleSize":           18,
			"titleWeight":         1,
			"moreSize":            12,
			"productTopRadio":     8,
			"productBottomRadio":  8,
			"productBgColor":      "#FFFFFF",
			"cardShadow":          "off",
			"productName_color":   "#333333",
			"nameWeight":          0,
			"productLine_color":   "#999999",
			"productPrice_color":  "#FF4C75",
			"product_sales_color": "#FF4C75",
			"productBtn_color1":   "#FF4C75",
			"productBtn_color2":   "#FF7A45",
			"btn_text_color":      "#FFFFFF",
		},
		"defaultData": []any{
			map[string]any{"product_name": "秋季新品轻薄羽绒服", "image": "", "presale_price": "3200.00", "ot_price": "3999.00", "sales": 286},
			map[string]any{"product_name": "轻奢羊毛大衣", "image": "", "presale_price": "2680.00", "ot_price": "3299.00", "sales": 169},
			map[string]any{"product_name": "简约通勤针织连衣裙", "image": "", "presale_price": "699.00", "ot_price": "899.00", "sales": 98},
		},
		"data": []any{},
	}
}

func defaultDiscountGroup() map[string]any {
	return map[string]any{
		"name": "折扣组",
		"params": map[string]any{
			"iconImage": "",
			"items": []any{
				map[string]any{"enabled": true, "image": "", "price": "", "productId": 0, "productName": "", "title": "全球美妆"},
				map[string]any{"enabled": true, "image": "", "price": "", "productId": 0, "productName": "", "title": "大牌鞋包"},
				map[string]any{"enabled": true, "image": "", "price": "", "productId": 0, "productName": "", "title": "数码产品"},
				map[string]any{"enabled": true, "image": "", "price": "", "productId": 0, "productName": "", "title": "精品腕表"},
			},
			"promotion": "券后低至7.3折",
			"slogan":    "真低价 放心买",
			"title":     "心动购物季",
		},
		"style": map[string]any{
			"background":    "rgba(255, 255, 255, 0)",
			"bgcolor":       "#f5f5f5",
			"cardRadius":    0,
			"cardShadow":    "off",
			"float":         0,
			"marginTop":     0,
			"paddingBottom": 0,
			"paddingLeft":   10,
			"paddingRight":  10,
			"paddingTop":    9,
			"radiusMode":    "all",
		},
		"type": "discountGroup",
	}
}
