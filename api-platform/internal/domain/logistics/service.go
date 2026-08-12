package logistics

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Store interface {
	ListExpress(ctx context.Context, page, limit int, showOnly bool, keyword, sortOrder string) ([]Express, int64, error)
	GetExpress(ctx context.Context, id uint) (*Express, error)
	ExistsExpress(ctx context.Context, name, code string, excludeID uint) (bool, error)
	CreateExpress(ctx context.Context, row *Express) error
	UpdateExpress(ctx context.Context, row *Express) error
	SoftDeleteExpress(ctx context.Context, id uint) error
	SyncExpressCatalog(ctx context.Context, rows []Express) (created, updated int, err error)

	ListCity(ctx context.Context, parentID *uint) ([]City, error)

	ListTemplate(ctx context.Context, merID uint, page, limit int) ([]ShippingTemplate, int64, error)
	GetTemplate(ctx context.Context, merID, id uint) (*ShippingTemplate, error)
	ListRegions(ctx context.Context, templateID uint) ([]Region, error)
	CreateTemplate(ctx context.Context, row *ShippingTemplate, regions []Region) error
	UpdateTemplate(ctx context.Context, row *ShippingTemplate, regions []Region) error
	SetDefaultTemplate(ctx context.Context, merID, id uint) error
	SoftDeleteTemplate(ctx context.Context, merID, id uint) error
}

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ListExpress(ctx context.Context, page, limit int, showOnly bool, keyword, sortOrder string) (*PageResult[Express], error) {
	page, limit = normalize(page, limit)
	sortOrder = strings.ToLower(strings.TrimSpace(sortOrder))
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = ""
	}
	list, total, err := s.store.ListExpress(ctx, page, limit, showOnly, strings.TrimSpace(keyword), sortOrder)
	if err != nil {
		return nil, err
	}
	return &PageResult[Express]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetExpressName(ctx context.Context, id uint) (string, error) {
	row, err := s.store.GetExpress(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrNotFound
		}
		return "", err
	}
	return row.Name, nil
}

func (s *Service) CreateExpress(ctx context.Context, in ExpressInput) (*Express, error) {
	name := strings.TrimSpace(in.Name)
	code := strings.TrimSpace(in.Code)
	if name == "" || code == "" || in.Sort < 0 {
		return nil, ErrBadParam
	}
	exists, err := s.store.ExistsExpress(ctx, name, code, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrExpressExists
	}
	row := &Express{
		Name: name, Code: code, Sort: in.Sort,
		IsShow: 1, CreateTime: time.Now(),
	}
	if in.IsShow != nil {
		row.IsShow = *in.IsShow
	}
	if err := s.store.CreateExpress(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *Service) UpdateExpress(ctx context.Context, id uint, in ExpressInput) (*Express, error) {
	if id == 0 || strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Code) == "" || in.Sort < 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetExpress(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	name := strings.TrimSpace(in.Name)
	code := strings.TrimSpace(in.Code)
	exists, err := s.store.ExistsExpress(ctx, name, code, id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrExpressExists
	}
	row.Name = name
	row.Code = code
	row.Sort = in.Sort
	if in.IsShow != nil {
		row.IsShow = *in.IsShow
	}
	if err := s.store.UpdateExpress(ctx, row); err != nil {
		return nil, err
	}
	return row, nil
}

// SyncExpressCatalog 将内置物流公司目录同步至本地数据库。该系统不接入已下线的
// 一号通服务；同步目标是平台本地可维护的标准承运商目录。
func (s *Service) SyncExpressCatalog(ctx context.Context) (*ExpressSyncResult, error) {
	rows := standardExpressCatalog()
	created, updated, err := s.store.SyncExpressCatalog(ctx, rows)
	if err != nil {
		return nil, err
	}
	return &ExpressSyncResult{Created: created, Updated: updated, Total: len(rows)}, nil
}

// standardExpressCatalog 使用 CRMEB_MER_v4.0 的 express.xlsx 标准目录，共 419 条。
// 本地初始化 SQL 与该目录保持一致，避免后台“同步物流公司”回退为最小示例集。
func standardExpressCatalog() []Express {
	return []Express{
		{Name: "A2U速递", Code: "a2u", Sort: 0, IsShow: 1},
		{Name: "AAE快递", Code: "aae", Sort: 0, IsShow: 1},
		{Name: "爱彼西快递", Code: "abc", Sort: 0, IsShow: 1},
		{Name: "德方物流", Code: "ahdf", Sort: 0, IsShow: 1},
		{Name: "航空快递", Code: "airgtc", Sort: 0, IsShow: 1},
		{Name: "阿里物流", Code: "ALP", Sort: 0, IsShow: 1},
		{Name: "安得物流", Code: "ande", Sort: 0, IsShow: 1},
		{Name: "安捷快递", Code: "anjie", Sort: 0, IsShow: 1},
		{Name: "安能物流", Code: "anneng", Sort: 0, IsShow: 1},
		{Name: "安信达快递", Code: "anxinda", Sort: 0, IsShow: 1},
		{Name: "安迅物流", Code: "anxun", Sort: 0, IsShow: 1},
		{Name: "AOL快递", Code: "aol", Sort: 0, IsShow: 1},
		{Name: "AOL澳通速递", Code: "aolau", Sort: 0, IsShow: 1},
		{Name: "Aramex", Code: "aramex", Sort: 0, IsShow: 1},
		{Name: "方舟速递", Code: "arke", Sort: 0, IsShow: 1},
		{Name: "澳邮中国快运", Code: "auexpress", Sort: 0, IsShow: 1},
		{Name: "卡行天下", Code: "B2B", Sort: 0, IsShow: 1},
		{Name: "百千诚国际物流", Code: "baiqian", Sort: 0, IsShow: 1},
		{Name: "百腾物流", Code: "baitengwuliu", Sort: 0, IsShow: 1},
		{Name: "八梁物流", Code: "BALIANGWL", Sort: 0, IsShow: 1},
		{Name: "巴伦支快递", Code: "balunzhi", Sort: 0, IsShow: 1},
		{Name: "邦送物流", Code: "bangsongwuliu", Sort: 0, IsShow: 1},
		{Name: "宝通达物流", Code: "baotongda", Sort: 0, IsShow: 1},
		{Name: "BCWELT", Code: "bcwelt", Sort: 0, IsShow: 1},
		{Name: "奔腾物流", Code: "benteng", Sort: 0, IsShow: 1},
		{Name: "滨发物流", Code: "BFWL", Sort: 0, IsShow: 1},
		{Name: "布谷鸟快递", Code: "bgn", Sort: 0, IsShow: 1},
		{Name: "挂号信", Code: "bgpyghx", Sort: 0, IsShow: 1},
		{Name: "BHT", Code: "bht", Sort: 0, IsShow: 1},
		{Name: "华慧快递", Code: "BHTEXP", Sort: 0, IsShow: 1},
		{Name: "笨鸟海淘", Code: "birdex", Sort: 0, IsShow: 1},
		{Name: "速方国际物流", Code: "bphchina", Sort: 0, IsShow: 1},
		{Name: "百事亨通", Code: "bsht", Sort: 0, IsShow: 1},
		{Name: "百世快运", Code: "bsky", Sort: 0, IsShow: 1},
		{Name: "博源恒通", Code: "byht", Sort: 0, IsShow: 1},
		{Name: "河南次晨达", Code: "ccd", Sort: 0, IsShow: 1},
		{Name: "CCES快递", Code: "cces", Sort: 0, IsShow: 1},
		{Name: "长通物流", Code: "changtong", Sort: 0, IsShow: 1},
		{Name: "程光快递", Code: "chengguang", Sort: 0, IsShow: 1},
		{Name: "城际速递", Code: "chengji", Sort: 0, IsShow: 1},
		{Name: "城市100快递", Code: "chengshi100", Sort: 0, IsShow: 1},
		{Name: "同舟行物流", Code: "chinatzx", Sort: 0, IsShow: 1},
		{Name: "传志快递", Code: "chuanzhi", Sort: 0, IsShow: 1},
		{Name: "出口易", Code: "chukouyi", Sort: 0, IsShow: 1},
		{Name: "CityLink快递", Code: "citylink", Sort: 0, IsShow: 1},
		{Name: "CE易欧通国际速递", Code: "cloudexpress", Sort: 0, IsShow: 1},
		{Name: "GLS快递", Code: "CNGLS", Sort: 0, IsShow: 1},
		{Name: "中环快递", Code: "cnpex", Sort: 0, IsShow: 1},
		{Name: "东方快递", Code: "coe", Sort: 0, IsShow: 1},
		{Name: "城市之星", Code: "cszx", Sort: 0, IsShow: 1},
		{Name: "云南中诚", Code: "czwlyn", Sort: 0, IsShow: 1},
		{Name: "大达物流", Code: "dada", Sort: 0, IsShow: 1},
		{Name: "大金物流", Code: "dajin", Sort: 0, IsShow: 1},
		{Name: "大顺物流", Code: "dashun", Sort: 0, IsShow: 1},
		{Name: "达速物流", Code: "dasu", Sort: 0, IsShow: 1},
		{Name: "大田物流", Code: "datian", Sort: 0, IsShow: 1},
		{Name: "大洋物流快递", Code: "dayang", Sort: 0, IsShow: 1},
		{Name: "大众佐川急便", Code: "dazhong", Sort: 0, IsShow: 1},
		{Name: "德邦物流", Code: "debang", Sort: 0, IsShow: 1},
		{Name: "德创物流", Code: "dechuangwuliu", Sort: 0, IsShow: 1},
		{Name: "德中快递", Code: "decnlh", Sort: 0, IsShow: 1},
		{Name: "德坤供应链", Code: "dekuncn", Sort: 0, IsShow: 1},
		{Name: "达方物流", Code: "dfpost", Sort: 0, IsShow: 1},
		{Name: "DHL快递", Code: "dhl", Sort: 0, IsShow: 1},
		{Name: "店通快递", Code: "diantong", Sort: 0, IsShow: 1},
		{Name: "递达快递", Code: "dida", Sort: 0, IsShow: 1},
		{Name: "叮咚澳洲转运", Code: "dindon", Sort: 0, IsShow: 1},
		{Name: "递四方速递", Code: "disifang", Sort: 0, IsShow: 1},
		{Name: "东瀚物流", Code: "donghanwl", Sort: 0, IsShow: 1},
		{Name: "东红物流", Code: "donghong", Sort: 0, IsShow: 1},
		{Name: "东骏快捷物流", Code: "dongjun", Sort: 0, IsShow: 1},
		{Name: "DPEX快递", Code: "dpex", Sort: 0, IsShow: 1},
		{Name: "D速快递", Code: "dsu", Sort: 0, IsShow: 1},
		{Name: "易满客", Code: "ecmscn", Sort: 0, IsShow: 1},
		{Name: "益递物流", Code: "edlogistics", Sort: 0, IsShow: 1},
		{Name: "百福东方快递", Code: "ees", Sort: 0, IsShow: 1},
		{Name: "易联通达物流", Code: "el56", Sort: 0, IsShow: 1},
		{Name: "EMS", Code: "ems", Sort: 0, IsShow: 1},
		{Name: "俄顺达", Code: "eshunda", Sort: 0, IsShow: 1},
		{Name: "欧亚专线", Code: "euasia", Sort: 0, IsShow: 1},
		{Name: "EWE全球快递", Code: "ewe", Sort: 0, IsShow: 1},
		{Name: "安鲜达", Code: "exfresh", Sort: 0, IsShow: 1},
		{Name: "E邮宝", Code: "eyoubao", Sort: 0, IsShow: 1},
		{Name: "伍圆速递", Code: "F5XM", Sort: 0, IsShow: 1},
		{Name: "颿达国际快递", Code: "fandaguoji", Sort: 0, IsShow: 1},
		{Name: "方方达物流", Code: "fangfangda", Sort: 0, IsShow: 1},
		{Name: "凡宇速递", Code: "fanyu", Sort: 0, IsShow: 1},
		{Name: "泛远国际物流", Code: "farlogistis", Sort: 0, IsShow: 1},
		{Name: "FedEx英国", Code: "fedexuk", Sort: 0, IsShow: 1},
		{Name: "飞邦物流", Code: "feibang", Sort: 0, IsShow: 1},
		{Name: "飞豹快递", Code: "feibao", Sort: 0, IsShow: 1},
		{Name: "原飞航快递", Code: "feihang", Sort: 0, IsShow: 1},
		{Name: "飞狐快递", Code: "feihu", Sort: 0, IsShow: 1},
		{Name: "飞快达物流", Code: "feikuaida", Sort: 0, IsShow: 1},
		{Name: "飞特物流", Code: "feite", Sort: 0, IsShow: 1},
		{Name: "飞洋快递", Code: "feiyang", Sort: 0, IsShow: 1},
		{Name: "飞远物流", Code: "feiyuan", Sort: 0, IsShow: 1},
		{Name: "丰达快递", Code: "fengda", Sort: 0, IsShow: 1},
		{Name: "风行天下", Code: "fengxingtianxia", Sort: 0, IsShow: 1},
		{Name: "飞康达物流", Code: "fkd", Sort: 0, IsShow: 1},
		{Name: "飞力士物流", Code: "flysman", Sort: 0, IsShow: 1},
		{Name: "FOX国际速递", Code: "fox", Sort: 0, IsShow: 1},
		{Name: "港快速递", Code: "gangkuai", Sort: 0, IsShow: 1},
		{Name: "GATI快递", Code: "gaticn", Sort: 0, IsShow: 1},
		{Name: "广东ems快递", Code: "gdems", Sort: 0, IsShow: 1},
		{Name: "国际包裹", Code: "gjbg", Sort: 0, IsShow: 1},
		{Name: "英脉物流", Code: "gml", Sort: 0, IsShow: 1},
		{Name: "国内小包", Code: "gnxb", Sort: 0, IsShow: 1},
		{Name: "共速达物流", Code: "gongsuda", Sort: 0, IsShow: 1},
		{Name: "GSM", Code: "gsm", Sort: 0, IsShow: 1},
		{Name: "万通快递", Code: "gswtkd", Sort: 0, IsShow: 1},
		{Name: "GTS快递", Code: "gts", Sort: 0, IsShow: 1},
		{Name: "高铁速递", Code: "gtsd", Sort: 0, IsShow: 1},
		{Name: "冠达快递", Code: "guada", Sort: 0, IsShow: 1},
		{Name: "广东邮政", Code: "guangdongyouzhengwuliu", Sort: 0, IsShow: 1},
		{Name: "广通速递", Code: "guangtong", Sort: 0, IsShow: 1},
		{Name: "国通快递", Code: "guotong", Sort: 0, IsShow: 1},
		{Name: "文捷航空速递", Code: "GZWENJIE", Sort: 0, IsShow: 1},
		{Name: "山东海红快递", Code: "haihong", Sort: 0, IsShow: 1},
		{Name: "海盟速递", Code: "haimeng", Sort: 0, IsShow: 1},
		{Name: "海外环球", Code: "haiwaihuanqiu", Sort: 0, IsShow: 1},
		{Name: "航宇快递", Code: "hangyu", Sort: 0, IsShow: 1},
		{Name: "韩润物流", Code: "hanrun", Sort: 0, IsShow: 1},
		{Name: "好来运快递", Code: "haolaiyun", Sort: 0, IsShow: 1},
		{Name: "昊盛物流", Code: "haosheng", Sort: 0, IsShow: 1},
		{Name: "好又快物流", Code: "haoyoukuai", Sort: 0, IsShow: 1},
		{Name: "河北建华物流", Code: "hebeijianhua", Sort: 0, IsShow: 1},
		{Name: "恒诚物流", Code: "HENGCHENGWL", Sort: 0, IsShow: 1},
		{Name: "恒丰物流", Code: "HENGFENGWL", Sort: 0, IsShow: 1},
		{Name: "恒路物流", Code: "henglu", Sort: 0, IsShow: 1},
		{Name: "恒宇运通", Code: "hengyu", Sort: 0, IsShow: 1},
		{Name: "和丰同城", Code: "hfwuxi", Sort: 0, IsShow: 1},
		{Name: "黑狗物流", Code: "higo", Sort: 0, IsShow: 1},
		{Name: "海派通", Code: "hipito", Sort: 0, IsShow: 1},
		{Name: "猴急送", Code: "hjs", Sort: 0, IsShow: 1},
		{Name: "香港邮政", Code: "hkpost", Sort: 0, IsShow: 1},
		{Name: "飞鹰物流", Code: "hnfy", Sort: 0, IsShow: 1},
		{Name: "宏捷国际物流", Code: "hongjie", Sort: 0, IsShow: 1},
		{Name: "鸿讯物流", Code: "hongxun", Sort: 0, IsShow: 1},
		{Name: "环球通达", Code: "hqtd", Sort: 0, IsShow: 1},
		{Name: "汇通天下物流", Code: "httx56", Sort: 0, IsShow: 1},
		{Name: "华通务达物流", Code: "htwd", Sort: 0, IsShow: 1},
		{Name: "华诚物流", Code: "huacheng", Sort: 0, IsShow: 1},
		{Name: "华达快运", Code: "huada", Sort: 0, IsShow: 1},
		{Name: "华翰物流", Code: "huahan", Sort: 0, IsShow: 1},
		{Name: "华航快递", Code: "huahang", Sort: 0, IsShow: 1},
		{Name: "黄马甲快递", Code: "huangmajia", Sort: 0, IsShow: 1},
		{Name: "环球速运", Code: "huanqiu", Sort: 0, IsShow: 1},
		{Name: "华企快运", Code: "huaqi", Sort: 0, IsShow: 1},
		{Name: "华夏龙物流", Code: "huaxialong", Sort: 0, IsShow: 1},
		{Name: "天地华宇物流", Code: "huayu", Sort: 0, IsShow: 1},
		{Name: "辉联物流", Code: "huilian", Sort: 0, IsShow: 1},
		{Name: "汇强快递", Code: "huiqiang", Sort: 0, IsShow: 1},
		{Name: "百世快递", Code: "huitong", Sort: 0, IsShow: 1},
		{Name: "汇文配送", Code: "huiwen", Sort: 0, IsShow: 1},
		{Name: "伙伴物流", Code: "huoban", Sort: 0, IsShow: 1},
		{Name: "户通物流", Code: "hutongwuliu", Sort: 0, IsShow: 1},
		{Name: "百成大达物流", Code: "idada", Sort: 0, IsShow: 1},
		{Name: "中国邮政", Code: "intmail", Sort: 0, IsShow: 1},
		{Name: "京东快递", Code: "jd", Sort: 0, IsShow: 1},
		{Name: "景光物流", Code: "jgwl", Sort: 0, IsShow: 1},
		{Name: "佳惠尔快递", Code: "jiahuier", Sort: 0, IsShow: 1},
		{Name: "佳吉物流", Code: "jiaji", Sort: 0, IsShow: 1},
		{Name: "佳家通", Code: "jiajiatong56", Sort: 0, IsShow: 1},
		{Name: "佳怡物流", Code: "jiayi", Sort: 0, IsShow: 1},
		{Name: "佳宇物流", Code: "JIAYU", Sort: 0, IsShow: 1},
		{Name: "加运美物流", Code: "jiayunmei", Sort: 0, IsShow: 1},
		{Name: "捷特快递", Code: "jiete", Sort: 0, IsShow: 1},
		{Name: "锦程国际物流", Code: "jinchengwuliu", Sort: 0, IsShow: 1},
		{Name: "金大物流", Code: "jindawuliu", Sort: 0, IsShow: 1},
		{Name: "京世物流", Code: "jingshi", Sort: 0, IsShow: 1},
		{Name: "京广速递快件", Code: "jinguangsudikuaijian", Sort: 0, IsShow: 1},
		{Name: "晋越快递", Code: "jinyue", Sort: 0, IsShow: 1},
		{Name: "九曳供应链", Code: "jiuye", Sort: 0, IsShow: 1},
		{Name: "久易快递", Code: "jiuyi", Sort: 0, IsShow: 1},
		{Name: "急先达物流", Code: "jixianda", Sort: 0, IsShow: 1},
		{Name: "嘉里大通", Code: "jldt", Sort: 0, IsShow: 1},
		{Name: "金马甲", Code: "jmjss", Sort: 0, IsShow: 1},
		{Name: "日本邮政", Code: "jppost", Sort: 0, IsShow: 1},
		{Name: "吉日优派", Code: "jrypex", Sort: 0, IsShow: 1},
		{Name: "骏川物流", Code: "JUNCHUANWL", Sort: 0, IsShow: 1},
		{Name: "骏丰国际速递", Code: "junfengguoji", Sort: 0, IsShow: 1},
		{Name: "吉祥邮", Code: "jxy", Sort: 0, IsShow: 1},
		{Name: "康力物流", Code: "klwl", Sort: 0, IsShow: 1},
		{Name: "直邮易", Code: "kuachangwuliu", Sort: 0, IsShow: 1},
		{Name: "快捷速递", Code: "kuaijie", Sort: 0, IsShow: 1},
		{Name: "快淘快递", Code: "kuaitao", Sort: 0, IsShow: 1},
		{Name: "快优达速递", Code: "kuaiyouda", Sort: 0, IsShow: 1},
		{Name: "宽容物流", Code: "kuanrong", Sort: 0, IsShow: 1},
		{Name: "跨越快递", Code: "kuayue", Sort: 0, IsShow: 1},
		{Name: "蓝镖快递", Code: "lanbiao", Sort: 0, IsShow: 1},
		{Name: "蓝弧快递", Code: "lanhu", Sort: 0, IsShow: 1},
		{Name: "宝凯物流", Code: "lbbk", Sort: 0, IsShow: 1},
		{Name: "联邦物流", Code: "LBWL", Sort: 0, IsShow: 1},
		{Name: "乐递供应链", Code: "ledii", Sort: 0, IsShow: 1},
		{Name: "乐捷递", Code: "lejiedi", Sort: 0, IsShow: 1},
		{Name: "云豹国际货运", Code: "leopard", Sort: 0, IsShow: 1},
		{Name: "联昊通快递", Code: "lianhaotong", Sort: 0, IsShow: 1},
		{Name: "成都立即送快递", Code: "lijisong", Sort: 0, IsShow: 1},
		{Name: "利民物流", Code: "LIMINWL", Sort: 0, IsShow: 1},
		{Name: "一号线", Code: "lineone", Sort: 0, IsShow: 1},
		{Name: "龙邦快运", Code: "longbang", Sort: 0, IsShow: 1},
		{Name: "隆浪快递", Code: "longlangkuaidi", Sort: 0, IsShow: 1},
		{Name: "龙胜物流", Code: "LONGSHENWL", Sort: 0, IsShow: 1},
		{Name: "恒通快递", Code: "lqht", Sort: 0, IsShow: 1},
		{Name: "乐天速递", Code: "ltexp", Sort: 0, IsShow: 1},
		{Name: "论道国际物流", Code: "lundao", Sort: 0, IsShow: 1},
		{Name: "鲁通快运", Code: "lutong", Sort: 0, IsShow: 1},
		{Name: "麦力快递", Code: "mailikuaidi", Sort: 0, IsShow: 1},
		{Name: "木春货运", Code: "mchy", Sort: 0, IsShow: 1},
		{Name: "美国快递", Code: "meiguo", Sort: 0, IsShow: 1},
		{Name: "美龙快递", Code: "meilong", Sort: 0, IsShow: 1},
		{Name: "美快国际物流", Code: "meiquick", Sort: 0, IsShow: 1},
		{Name: "美西快递", Code: "meixi", Sort: 0, IsShow: 1},
		{Name: "门对门", Code: "menduimen", Sort: 0, IsShow: 1},
		{Name: "蒙速快递", Code: "mengsu", Sort: 0, IsShow: 1},
		{Name: "民邦速递", Code: "minbang", Sort: 0, IsShow: 1},
		{Name: "明亮物流", Code: "mingliang", Sort: 0, IsShow: 1},
		{Name: "民航快递", Code: "minhang", Sort: 0, IsShow: 1},
		{Name: "闽盛物流", Code: "minsheng", Sort: 0, IsShow: 1},
		{Name: "南北快递", Code: "nanbei", Sort: 0, IsShow: 1},
		{Name: "中国南方航空股份有限公司", Code: "NANHANG", Sort: 0, IsShow: 1},
		{Name: "红马速递", Code: "nedahm", Sort: 0, IsShow: 1},
		{Name: "港中能达快递", Code: "nengda", Sort: 0, IsShow: 1},
		{Name: "新蛋奥硕物流", Code: "neweggozzo", Sort: 0, IsShow: 1},
		{Name: "华赫物流", Code: "nmhuahe", Sort: 0, IsShow: 1},
		{Name: "腾达速递", Code: "nntengda", Sort: 0, IsShow: 1},
		{Name: "偌亚奥国际", Code: "nuoyaao", Sort: 0, IsShow: 1},
		{Name: "OCS国际快递", Code: "ocs", Sort: 0, IsShow: 1},
		{Name: "一号仓", Code: "onehcang", Sort: 0, IsShow: 1},
		{Name: "onTrac", Code: "ontrac", Sort: 0, IsShow: 1},
		{Name: "中欧快运", Code: "otobv", Sort: 0, IsShow: 1},
		{Name: "澳大利亚PCA快递", Code: "pca", Sort: 0, IsShow: 1},
		{Name: "配思货运", Code: "PEISI", Sort: 0, IsShow: 1},
		{Name: "陪行物流", Code: "peixing", Sort: 0, IsShow: 1},
		{Name: "彪记快递", Code: "PEWKEE", Sort: 0, IsShow: 1},
		{Name: "皇家物流", Code: "pfcexpress", Sort: 0, IsShow: 1},
		{Name: "凤凰快递", Code: "PHOENIXEXP", Sort: 0, IsShow: 1},
		{Name: "平安达快递", Code: "pinganda", Sort: 0, IsShow: 1},
		{Name: "平安达腾飞", Code: "pingandatengfei", Sort: 0, IsShow: 1},
		{Name: "小包", Code: "pingyou", Sort: 0, IsShow: 1},
		{Name: "品骏快递", Code: "pjbest", Sort: 0, IsShow: 1},
		{Name: "贝邮宝", Code: "ppbyb", Sort: 0, IsShow: 1},
		{Name: "急顺通", Code: "pzhjst", Sort: 0, IsShow: 1},
		{Name: "秦邦快运", Code: "qbexpress", Sort: 0, IsShow: 1},
		{Name: "启辰国际物流", Code: "qichen", Sort: 0, IsShow: 1},
		{Name: "秦远物流", Code: "qinyuan", Sort: 0, IsShow: 1},
		{Name: "千顺快递", Code: "qskdyxgs", Sort: 0, IsShow: 1},
		{Name: "全晨快递", Code: "quanchen", Sort: 0, IsShow: 1},
		{Name: "全峰快递", Code: "quanfeng", Sort: 0, IsShow: 1},
		{Name: "全际通快递", Code: "quanjitong", Sort: 0, IsShow: 1},
		{Name: "全日通快递", Code: "quanritong", Sort: 0, IsShow: 1},
		{Name: "全速快运", Code: "quansu", Sort: 0, IsShow: 1},
		{Name: "全速通国际快递", Code: "quansutong", Sort: 0, IsShow: 1},
		{Name: "全信通快递", Code: "quanxintong", Sort: 0, IsShow: 1},
		{Name: "全一快递", Code: "quanyi", Sort: 0, IsShow: 1},
		{Name: "全之鑫物流", Code: "qzx56", Sort: 0, IsShow: 1},
		{Name: "日日顺物流", Code: "ririshun", Sort: 0, IsShow: 1},
		{Name: "日昱物流", Code: "riyu", Sort: 0, IsShow: 1},
		{Name: "荣庆物流", Code: "rongqing", Sort: 0, IsShow: 1},
		{Name: "RPX保时达", Code: "rpx", Sort: 0, IsShow: 1},
		{Name: "捷网俄全通", Code: "ruexp", Sort: 0, IsShow: 1},
		{Name: "如风达快递", Code: "rufeng", Sort: 0, IsShow: 1},
		{Name: "凡客如风达", Code: "rufengda", Sort: 0, IsShow: 1},
		{Name: "瑞达国际速递", Code: "ruidaex", Sort: 0, IsShow: 1},
		{Name: "瑞丰速递", Code: "ruifeng", Sort: 0, IsShow: 1},
		{Name: "全时速运", Code: "runhengfeng", Sort: 0, IsShow: 1},
		{Name: "日益通速递", Code: "rytsd", Sort: 0, IsShow: 1},
		{Name: "赛澳递", Code: "saiaodi", Sort: 0, IsShow: 1},
		{Name: "三态速递", Code: "santai", Sort: 0, IsShow: 1},
		{Name: "丰程物流", Code: "sccod", Sort: 0, IsShow: 1},
		{Name: "泰国138", Code: "sd138", Sort: 0, IsShow: 1},
		{Name: "优配速运", Code: "sdyoupei", Sort: 0, IsShow: 1},
		{Name: "速递中国", Code: "sendtochina", Sort: 0, IsShow: 1},
		{Name: "七天连锁", Code: "sevendays", Sort: 0, IsShow: 1},
		{Name: "十方通物流", Code: "sfift", Sort: 0, IsShow: 1},
		{Name: "圣安物流", Code: "shengan", Sort: 0, IsShow: 1},
		{Name: "晟邦物流", Code: "shengbang", Sort: 0, IsShow: 1},
		{Name: "盛丰物流", Code: "shengfeng", Sort: 0, IsShow: 1},
		{Name: "盛辉物流", Code: "shenghui", Sort: 0, IsShow: 1},
		{Name: "申通快递", Code: "shentong", Sort: 0, IsShow: 1},
		{Name: "昊昕物流", Code: "SHHX", Sort: 0, IsShow: 1},
		{Name: "世运快递", Code: "shiyun", Sort: 0, IsShow: 1},
		{Name: "上海林道货运", Code: "shlindao", Sort: 0, IsShow: 1},
		{Name: "顺发物流", Code: "SHUNFAWL", Sort: 0, IsShow: 1},
		{Name: "顺丰速运", Code: "shunfeng", Sort: 0, IsShow: 1},
		{Name: "顺捷丰达", Code: "shunjiefengda", Sort: 0, IsShow: 1},
		{Name: "四海快递", Code: "sihaiet", Sort: 0, IsShow: 1},
		{Name: "思迈快递", Code: "simai", Sort: 0, IsShow: 1},
		{Name: "信联通", Code: "sinatone", Sort: 0, IsShow: 1},
		{Name: "新加坡邮政", Code: "singpost", Sort: 0, IsShow: 1},
		{Name: "宋军物流", Code: "SJWL", Sort: 0, IsShow: 1},
		{Name: "荷兰", Code: "Sky", Sort: 0, IsShow: 1},
		{Name: "春风物流", Code: "spring", Sort: 0, IsShow: 1},
		{Name: "星晨急便", Code: "STARS", Sort: 0, IsShow: 1},
		{Name: "顺通快递", Code: "stkd", Sort: 0, IsShow: 1},
		{Name: "速必达物流", Code: "subida", Sort: 0, IsShow: 1},
		{Name: "速呈宅配", Code: "suchengzhaipei", Sort: 0, IsShow: 1},
		{Name: "穗佳物流", Code: "suijia", Sort: 0, IsShow: 1},
		{Name: "郑州速捷", Code: "sujievip", Sort: 0, IsShow: 1},
		{Name: "上大物流", Code: "SUNDAPOST", Sort: 0, IsShow: 1},
		{Name: "苏宁快递", Code: "suning", Sort: 0, IsShow: 1},
		{Name: "新杰物流", Code: "sunjex", Sort: 0, IsShow: 1},
		{Name: "新速航", Code: "sunspeedy", Sort: 0, IsShow: 1},
		{Name: "速尔物流", Code: "sure", Sort: 0, IsShow: 1},
		{Name: "速腾快递", Code: "suteng", Sort: 0, IsShow: 1},
		{Name: "速通物流", Code: "sutong", Sort: 0, IsShow: 1},
		{Name: "苏粤货运", Code: "SUYUE", Sort: 0, IsShow: 1},
		{Name: "盛旺货运", Code: "SWHY", Sort: 0, IsShow: 1},
		{Name: "山西红马甲", Code: "sxhongmajia", Sort: 0, IsShow: 1},
		{Name: "沈阳佳惠尔", Code: "syjiahuier", Sort: 0, IsShow: 1},
		{Name: "华宇物流", Code: "tiandihuayu", Sort: 0, IsShow: 1},
		{Name: "天河物流", Code: "TIANHEWL", Sort: 0, IsShow: 1},
		{Name: "天天快递", Code: "tiantian", Sort: 0, IsShow: 1},
		{Name: "天纵物流", Code: "tianzong", Sort: 0, IsShow: 1},
		{Name: "万家通", Code: "timedg", Sort: 0, IsShow: 1},
		{Name: "天联快运", Code: "tlky", Sort: 0, IsShow: 1},
		{Name: "TNT快递", Code: "tnt", Sort: 0, IsShow: 1},
		{Name: "通成物流", Code: "tongcheng", Sort: 0, IsShow: 1},
		{Name: "通达兴物流", Code: "tongdaxing", Sort: 0, IsShow: 1},
		{Name: "通和天下物流", Code: "tonghe", Sort: 0, IsShow: 1},
		{Name: "中运全速", Code: "topspeedex", Sort: 0, IsShow: 1},
		{Name: "汤氏物流", Code: "TSWL", Sort: 0, IsShow: 1},
		{Name: "合众速递", Code: "ucs", Sort: 0, IsShow: 1},
		{Name: "UEQ快递", Code: "ueq", Sort: 0, IsShow: 1},
		{Name: "UEX", Code: "uex", Sort: 0, IsShow: 1},
		{Name: "UPS快递", Code: "ups", Sort: 0, IsShow: 1},
		{Name: "USPS快递", Code: "usps", Sort: 0, IsShow: 1},
		{Name: "美通快递", Code: "valueway", Sort: 0, IsShow: 1},
		{Name: "鹰运国际速递", Code: "vipexpress", Sort: 0, IsShow: 1},
		{Name: "万博快递", Code: "wanbo", Sort: 0, IsShow: 1},
		{Name: "万家物流", Code: "wanjia", Sort: 0, IsShow: 1},
		{Name: "万象物流", Code: "wanxiang", Sort: 0, IsShow: 1},
		{Name: "微特派快递", Code: "weitepai", Sort: 0, IsShow: 1},
		{Name: "渥途国际速运", Code: "wotu", Sort: 0, IsShow: 1},
		{Name: "威时沛运", Code: "wtdchina", Sort: 0, IsShow: 1},
		{Name: "五环速递", Code: "wuhuan", Sort: 0, IsShow: 1},
		{Name: "微转运", Code: "wzhaunyun", Sort: 0, IsShow: 1},
		{Name: "西安胜峰", Code: "xaetc", Sort: 0, IsShow: 1},
		{Name: "鑫飞鸿物流快递", Code: "XFHONG", Sort: 0, IsShow: 1},
		{Name: "西安城联速递", Code: "xianchenglian", Sort: 0, IsShow: 1},
		{Name: "先锋快递", Code: "xianfeng", Sort: 0, IsShow: 1},
		{Name: "北青小红帽", Code: "xiaohongmao", Sort: 0, IsShow: 1},
		{Name: "喜来快递", Code: "xilaikd", Sort: 0, IsShow: 1},
		{Name: "新邦物流", Code: "xinbang", Sort: 0, IsShow: 1},
		{Name: "新蛋物流", Code: "xindan", Sort: 0, IsShow: 1},
		{Name: "信丰物流", Code: "xinfeng", Sort: 0, IsShow: 1},
		{Name: "星程宅配", Code: "xingchengzhaipei", Sort: 0, IsShow: 1},
		{Name: "鑫天顺物流", Code: "XINTIAN", Sort: 0, IsShow: 1},
		{Name: "信天捷快递", Code: "xintianjie", Sort: 0, IsShow: 1},
		{Name: "西邮寄", Code: "xipost", Sort: 0, IsShow: 1},
		{Name: "希优特快递", Code: "xiyoute", Sort: 0, IsShow: 1},
		{Name: "祥龙运通", Code: "xlyt", Sort: 0, IsShow: 1},
		{Name: "鑫世锐达", Code: "xsrd", Sort: 0, IsShow: 1},
		{Name: "鑫通宝物流", Code: "xtb", Sort: 0, IsShow: 1},
		{Name: "源安达快递", Code: "yad", Sort: 0, IsShow: 1},
		{Name: "亚风速递", Code: "yafeng", Sort: 0, IsShow: 1},
		{Name: "亚马逊物流", Code: "yamaxunwuliu", Sort: 0, IsShow: 1},
		{Name: "燕文物流", Code: "yanwen", Sort: 0, IsShow: 1},
		{Name: "邮联物流", Code: "YBWL", Sort: 0, IsShow: 1},
		{Name: "远成快运", Code: "ycgky", Sort: 0, IsShow: 1},
		{Name: "一邦快递", Code: "yibang", Sort: 0, IsShow: 1},
		{Name: "易达通快递", Code: "yidatong", Sort: 0, IsShow: 1},
		{Name: "亿领速运", Code: "yilingsuyun", Sort: 0, IsShow: 1},
		{Name: "英超物流", Code: "yingchao", Sort: 0, IsShow: 1},
		{Name: "顺捷丰达", Code: "yinjie", Sort: 0, IsShow: 1},
		{Name: "音速速运", Code: "yinsu", Sort: 0, IsShow: 1},
		{Name: "一柒物流", Code: "yiqiguojiwuliu", Sort: 0, IsShow: 1},
		{Name: "亿顺航", Code: "yishunhang", Sort: 0, IsShow: 1},
		{Name: "易通达", Code: "yitongda", Sort: 0, IsShow: 1},
		{Name: "永昌物流", Code: "yongchang", Sort: 0, IsShow: 1},
		{Name: "永旺达快递", Code: "yongwangda", Sort: 0, IsShow: 1},
		{Name: "邮必佳", Code: "youbijia", Sort: 0, IsShow: 1},
		{Name: "UC优速快递", Code: "youshuwuliu", Sort: 0, IsShow: 1},
		{Name: "优速快递", Code: "yousu", Sort: 0, IsShow: 1},
		{Name: "优速通达", Code: "yousutongda", Sort: 0, IsShow: 1},
		{Name: "挂号信", Code: "youzhengguonei", Sort: 0, IsShow: 1},
		{Name: "壹品速递", Code: "ypsd", Sort: 0, IsShow: 1},
		{Name: "一统飞鸿快递", Code: "ytfh", Sort: 0, IsShow: 1},
		{Name: "远成物流", Code: "yuancheng", Sort: 0, IsShow: 1},
		{Name: "圆通速递", Code: "yuantong", Sort: 0, IsShow: 1},
		{Name: "圆圆物流", Code: "YUANYUANWL", Sort: 0, IsShow: 1},
		{Name: "元智捷诚快递", Code: "yuanzhi", Sort: 0, IsShow: 1},
		{Name: "越丰快递", Code: "yuefeng", Sort: 0, IsShow: 1},
		{Name: "御风速运", Code: "yufeng", Sort: 0, IsShow: 1},
		{Name: "煜嘉物流", Code: "yujiawuliu", Sort: 0, IsShow: 1},
		{Name: "誉美捷快递", Code: "yumeijie", Sort: 0, IsShow: 1},
		{Name: "韵达快运", Code: "yunda", Sort: 0, IsShow: 1},
		{Name: "韵达美国件", Code: "yundaexus", Sort: 0, IsShow: 1},
		{Name: "运通快递", Code: "yuntong", Sort: 0, IsShow: 1},
		{Name: "云物流", Code: "yunwuliu", Sort: 0, IsShow: 1},
		{Name: "宇鑫物流", Code: "yuxin", Sort: 0, IsShow: 1},
		{Name: "源伟丰快递", Code: "ywfex", Sort: 0, IsShow: 1},
		{Name: "宇鑫物流", Code: "yxwl", Sort: 0, IsShow: 1},
		{Name: "远洋国际", Code: "yyexpress", Sort: 0, IsShow: 1},
		{Name: "一运全成物流", Code: "yyqc56", Sort: 0, IsShow: 1},
		{Name: "增益快递", Code: "zengyi", Sort: 0, IsShow: 1},
		{Name: "增益速递", Code: "zengyisudi", Sort: 0, IsShow: 1},
		{Name: "振刚物流", Code: "ZGWL", Sort: 0, IsShow: 1},
		{Name: "宅急送", Code: "zhaijisong", Sort: 0, IsShow: 1},
		{Name: "众辉达物流", Code: "zhdwl", Sort: 0, IsShow: 1},
		{Name: "至诚通达快递", Code: "zhichengtongda", Sort: 0, IsShow: 1},
		{Name: "芝麻开门", Code: "zhima", Sort: 0, IsShow: 1},
		{Name: "中睿速递", Code: "zhongruisudi", Sort: 0, IsShow: 1},
		{Name: "中速快件", Code: "zhongsukuaidi", Sort: 0, IsShow: 1},
		{Name: "中天快运", Code: "zhongtian", Sort: 0, IsShow: 1},
		{Name: "中铁快运", Code: "zhongtie", Sort: 0, IsShow: 1},
		{Name: "中通快递", Code: "zhongtong", Sort: 0, IsShow: 1},
		{Name: "中外运速递", Code: "zhongwaiyun", Sort: 0, IsShow: 1},
		{Name: "中信达快递", Code: "zhongxinda", Sort: 0, IsShow: 1},
		{Name: "中邮物流", Code: "zhongyou", Sort: 0, IsShow: 1},
		{Name: "纵行物流", Code: "zongxing", Sort: 0, IsShow: 1},
		{Name: "准实快运", Code: "zsky123", Sort: 0, IsShow: 1},
		{Name: "中铁物流", Code: "ZTKY", Sort: 0, IsShow: 1},
		{Name: "智通物流", Code: "ztong", Sort: 0, IsShow: 1},
		{Name: "中天万运快递", Code: "ztwy", Sort: 0, IsShow: 1},
		{Name: "佐川急便", Code: "zuochuan", Sort: 0, IsShow: 1},
		{Name: "中外速运", Code: "zwsy", Sort: 0, IsShow: 1},
		{Name: "郑州建华快递", Code: "zzjh", Sort: 0, IsShow: 1},
	}
}

func (s *Service) DeleteExpress(ctx context.Context, id uint) error {
	if _, err := s.store.GetExpress(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return s.store.SoftDeleteExpress(ctx, id)
}

func (s *Service) ListCity(ctx context.Context, parentID *uint) ([]City, error) {
	return s.store.ListCity(ctx, parentID)
}

func (s *Service) ListTemplate(ctx context.Context, merID uint, page, limit int) (*PageResult[ShippingTemplate], error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	page, limit = normalize(page, limit)
	list, total, err := s.store.ListTemplate(ctx, merID, page, limit)
	if err != nil {
		return nil, err
	}
	return &PageResult[ShippingTemplate]{List: list, Total: total, Page: page, Limit: limit}, nil
}

func (s *Service) GetTemplate(ctx context.Context, merID, id uint) (*ShippingTemplate, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	row, err := s.store.GetTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	regions, err := s.store.ListRegions(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

func (s *Service) CreateTemplate(ctx context.Context, merID uint, in TemplateInput) (*ShippingTemplate, error) {
	if merID == 0 {
		return nil, ErrBadParam
	}
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, ErrBadParam
	}
	typ := in.Type
	if typ == 0 {
		typ = 1
	}
	row := &ShippingTemplate{
		MerID: merID, Name: name, Type: typ, Appoint: in.Appoint,
		Sort: in.Sort, CreateTime: time.Now(),
	}
	regions := toRegions(0, in.Regions)
	if err := s.store.CreateTemplate(ctx, row, regions); err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, merID, id uint, in TemplateInput) (*ShippingTemplate, error) {
	row, err := s.store.GetTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if name := strings.TrimSpace(in.Name); name != "" {
		row.Name = name
	}
	if in.Type > 0 {
		row.Type = in.Type
	}
	row.Appoint = in.Appoint
	row.Sort = in.Sort
	regions := toRegions(id, in.Regions)
	if err := s.store.UpdateTemplate(ctx, row, regions); err != nil {
		return nil, err
	}
	row.Regions = regions
	return row, nil
}

// SetDefaultTemplate 在同一商户范围内保持默认模板唯一；不改变订单计价或商品模板绑定。
func (s *Service) SetDefaultTemplate(ctx context.Context, merID, id uint) (*ShippingTemplate, error) {
	if merID == 0 || id == 0 {
		return nil, ErrBadParam
	}
	if err := s.store.SetDefaultTemplate(ctx, merID, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return s.GetTemplate(ctx, merID, id)
}

func (s *Service) DeleteTemplate(ctx context.Context, merID, id uint) error {
	row, err := s.store.GetTemplate(ctx, merID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if row.IsDefault == 1 {
		return ErrDefaultTemplate
	}
	return s.store.SoftDeleteTemplate(ctx, merID, id)
}

func toRegions(templateID uint, in []RegionInput) []Region {
	if len(in) == 0 {
		return []Region{{
			TemplateID: templateID, First: 1, FirstPrice: 0, Continue: 1, ContinuePrice: 0,
		}}
	}
	out := make([]Region, 0, len(in))
	for _, r := range in {
		first := r.First
		if first <= 0 {
			first = 1
		}
		cont := r.Continue
		if cont <= 0 {
			cont = 1
		}
		out = append(out, Region{
			TemplateID: templateID, CityIDs: strings.TrimSpace(r.CityIDs),
			First: first, FirstPrice: r.FirstPrice, Continue: cont, ContinuePrice: r.ContinuePrice,
		})
	}
	return out
}

func normalize(page, limit int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return page, limit
}
