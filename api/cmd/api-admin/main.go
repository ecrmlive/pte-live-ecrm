package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/aftersale"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/article"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/assist"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/attachment"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/broadcast"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/combination"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/community"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/chat"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cs"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/distribution"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/finance"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/fulfillment"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/invoice"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/logistics"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/merchant"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/openapi"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/presell"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/productmeta"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/seckill"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/usertag"
	aftersalepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/aftersale"
	articlepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/article"
	assistpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/assist"
	attachmentpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/attachment"
	broadcastpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/broadcast"
	cartpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/cart"
	catalogpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/catalog"
	chatpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/chat"
	combinationpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/combination"
	communitypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/community"
	contentpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/content"
	cspersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/cs"
	fulfillmentpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/fulfillment"
	invoicepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/invoice"
	diypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/diy"
	distributionpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/distribution"
	financepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/finance"
	identitypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/identity"
	logisticspersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/logistics"
	merchantpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/merchant"
	openapipersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/openapi"
	presellpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/presell"
	productmetapersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/productmeta"
	promotionpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/promotion"
	reservationpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/reservation"
	seckillpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/seckill"
	tradepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/trade"
	usertagpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/usertag"
	managerauth "github.com/qixi-live/qixi-live-mergers/api/internal/manager/auth"
	managerorder "github.com/qixi-live/qixi-live-mergers/api/internal/manager/order"
	managerrefund "github.com/qixi-live/qixi-live-mergers/api/internal/manager/refund"
	merchantassist "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/assist"
	merchantattachment "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/attachment"
	merchantauth "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/auth"
	merchantbroadcast "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/broadcast"
	merchantcatalog "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/catalog"
	merchantcombination "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/combination"
	merchantcommunity "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/community"
	merchantcoupon "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/coupon"
	merchantcs "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/cs"
	merchantfinance "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/finance"
	merchantfulfillment "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/fulfillment"
	merchantinvoice "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/invoice"
	merchantlogistics "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/logistics"
	merchantorder "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/order"
	merchantpresell "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/presell"
	merchantproductmeta "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/productmeta"
	merchantrefund "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/refund"
	merchantreservation "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/reservation"
	merchantseckill "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/seckill"
	merchantsetting "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/setting"
	merchantsvip "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/svip"
	openapihttp "github.com/qixi-live/qixi-live-mergers/api/internal/open"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/config"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/db"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/upload"
	platformarticle "github.com/qixi-live/qixi-live-mergers/api/internal/platform/article"
	platformassist "github.com/qixi-live/qixi-live-mergers/api/internal/platform/assist"
	platformattachment "github.com/qixi-live/qixi-live-mergers/api/internal/platform/attachment"
	platformauth "github.com/qixi-live/qixi-live-mergers/api/internal/platform/auth"
	platformbroadcast "github.com/qixi-live/qixi-live-mergers/api/internal/platform/broadcast"
	platformcatalog "github.com/qixi-live/qixi-live-mergers/api/internal/platform/catalog"
	platformcombination "github.com/qixi-live/qixi-live-mergers/api/internal/platform/combination"
	platformcommunity "github.com/qixi-live/qixi-live-mergers/api/internal/platform/community"
	platformcontent "github.com/qixi-live/qixi-live-mergers/api/internal/platform/content"
	platformcoupon "github.com/qixi-live/qixi-live-mergers/api/internal/platform/coupon"
	merchantdiy "github.com/qixi-live/qixi-live-mergers/api/internal/merchant/diy"
	platformdiy "github.com/qixi-live/qixi-live-mergers/api/internal/platform/diy"
	platformfinance "github.com/qixi-live/qixi-live-mergers/api/internal/platform/finance"
	platformlogistics "github.com/qixi-live/qixi-live-mergers/api/internal/platform/logistics"
	platformmerchant "github.com/qixi-live/qixi-live-mergers/api/internal/platform/merchant"
	platformorder "github.com/qixi-live/qixi-live-mergers/api/internal/platform/order"
	platformpresell "github.com/qixi-live/qixi-live-mergers/api/internal/platform/presell"
	platformproductmeta "github.com/qixi-live/qixi-live-mergers/api/internal/platform/productmeta"
	platformrefund "github.com/qixi-live/qixi-live-mergers/api/internal/platform/refund"
	platformseckill "github.com/qixi-live/qixi-live-mergers/api/internal/platform/seckill"
	platformsetting "github.com/qixi-live/qixi-live-mergers/api/internal/platform/setting"
	platformspread "github.com/qixi-live/qixi-live-mergers/api/internal/platform/spread"
	platformsvip "github.com/qixi-live/qixi-live-mergers/api/internal/platform/svip"
	platformusertag "github.com/qixi-live/qixi-live-mergers/api/internal/platform/usertag"
	"github.com/qixi-live/qixi-live-mergers/api/internal/serviceportal"
)

func main() {
	cfgPath := flag.String("config", "conf/admin.yaml", "path to admin.yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	gdb, err := db.OpenMySQL(cfg.MySQL.DSN, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}

	jwtMgr := authjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL(), cfg.JWT.RefreshTTL())
	idSvc := identity.NewService(identitypersist.NewRepo(gdb))
	merSvc := merchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(gdb)))
	catSvc := catalog.NewService(catalogpersist.NewStoreAdapter(catalogpersist.NewRepo(gdb)))
	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	distSvc := distribution.NewService(distributionpersist.NewStoreAdapter(distributionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
	aftersaleSvc := aftersale.NewService(aftersalepersist.NewStoreAdapter(aftersalepersist.NewRepo(gdb)))
	financeSvc := finance.NewService(financepersist.NewStoreAdapter(financepersist.NewRepo(gdb)))
	contentSvc := content.NewService(contentpersist.NewRepo(gdb))
	diySvc := diy.NewService(diypersist.NewRepo(gdb))
	seckillSvc := seckill.NewService(seckillpersist.NewRepo(gdb))
	comboSvc := combination.NewService(combinationpersist.NewRepo(gdb))
	presellSvc := presell.NewService(presellpersist.NewRepo(gdb))
	reserveSvc := reservation.NewService(reservationpersist.NewRepo(gdb))
	broadcastSvc := broadcast.NewService(broadcastpersist.NewRepo(gdb))
	communitySvc := community.NewService(communitypersist.NewRepo(gdb))
	assistSvc := assist.NewService(assistpersist.NewRepo(gdb))
	attachSvc := attachment.NewService(attachmentpersist.NewRepo(gdb))
	csSvc := cs.NewService(cspersist.NewRepo(gdb))
	chatSvc := chat.NewService(chatpersist.NewRepo(gdb), chat.IMSettings{
		Mode: cfg.IM.Mode, APIBase: cfg.IM.APIBase, WSPublicURL: cfg.IM.WSPublicURL,
		AppID: cfg.IM.AppID, Token: cfg.IM.IntegrationToken, Secret: cfg.JWT.Secret,
	})
	fulfillSvc := fulfillment.NewService(fulfillmentpersist.NewRepo(gdb))
	invoiceSvc := invoice.NewService(invoicepersist.NewRepo(gdb))
	logisticsSvc := logistics.NewService(logisticspersist.NewRepo(gdb))
	productMetaSvc := productmeta.NewService(productmetapersist.NewRepo(gdb))
	articleSvc := article.NewService(articlepersist.NewRepo(gdb))
	userTagSvc := usertag.NewService(usertagpersist.NewRepo(gdb))
	cosStore := upload.COS{
		Bucket: cfg.COS.Bucket, Region: cfg.COS.Region,
		SecretID: cfg.COS.SecretID, SecretKey: cfg.COS.SecretKey,
		BaseURL: cfg.COS.BaseURL, KeyPrefix: cfg.COS.KeyPrefix,
	}
	var fileUp upload.Store
	if cfg.COS.Enabled && cosStore.Configured() {
		fileUp = cosStore
	} else {
		fileUp = upload.Local{Dir: cfg.Upload.Dir, PublicBase: cfg.Upload.PublicBase}
	}
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))
	openSvc := openapi.NewService(openapipersist.NewStoreAdapter(openapipersist.NewRepo(gdb)), jwtMgr)

	platformAuthH := platformauth.NewHandler(idSvc, jwtMgr)
	platformMerH := platformmerchant.NewHandler(merSvc)
	platformCatH := platformcatalog.NewHandler(catSvc)
	platformOrderH := platformorder.NewHandler(tradeSvc)
	platformRefundH := platformrefund.NewHandler(aftersaleSvc, idSvc)
	platformFinanceH := platformfinance.NewHandler(financeSvc, idSvc)
	platformCouponH := platformcoupon.NewHandler(promoSvc, idSvc)
	platformSpreadH := platformspread.NewHandler(distSvc)
	platformContentH := platformcontent.NewHandler(contentSvc, idSvc)
	platformDiyH := platformdiy.NewHandler(diySvc, idSvc)
	platformSeckillH := platformseckill.NewHandler(seckillSvc)
	platformComboH := platformcombination.NewHandler(comboSvc)
	platformPresellH := platformpresell.NewHandler(presellSvc)
	platformBroadcastH := platformbroadcast.NewHandler(broadcastSvc, idSvc)
	platformCommunityH := platformcommunity.NewHandler(communitySvc, idSvc)
	platformAssistH := platformassist.NewHandler(assistSvc)
	platformAttachH := platformattachment.NewHandler(attachSvc, idSvc, fileUp)
	platformSvipH := platformsvip.NewHandler(idSvc)
	platformSettingH := platformsetting.NewHandler(idSvc)
	platformLogisticsH := platformlogistics.NewHandler(logisticsSvc)
	platformProductMetaH := platformproductmeta.NewHandler(productMetaSvc)
	platformArticleH := platformarticle.NewHandler(articleSvc)
	platformUserTagH := platformusertag.NewHandler(userTagSvc)
	merchantAuthH := merchantauth.NewHandler(idSvc, jwtMgr)
	merchantCatH := merchantcatalog.NewHandler(catSvc, idSvc)
	merchantOrderH := merchantorder.NewHandler(tradeSvc, idSvc, logisticsSvc)
	merchantRefundH := merchantrefund.NewHandler(aftersaleSvc, idSvc)
	merchantFinanceH := merchantfinance.NewHandler(financeSvc)
	merchantCouponH := merchantcoupon.NewHandler(promoSvc, idSvc)
	merchantSeckillH := merchantseckill.NewHandler(seckillSvc, idSvc)
	merchantComboH := merchantcombination.NewHandler(comboSvc, idSvc)
	merchantPresellH := merchantpresell.NewHandler(presellSvc, idSvc)
	merchantBroadcastH := merchantbroadcast.NewHandler(broadcastSvc, idSvc)
	merchantCommunityH := merchantcommunity.NewHandler(communitySvc, idSvc)
	merchantAssistH := merchantassist.NewHandler(assistSvc, idSvc)
	merchantAttachH := merchantattachment.NewHandler(attachSvc, idSvc, fileUp)
	merchantReserveH := merchantreservation.NewHandler(reserveSvc, idSvc)
	merchantSvipH := merchantsvip.NewHandler(merSvc, idSvc)
	merchantSettingH := merchantsetting.NewHandler(idSvc, merSvc)
	merchantCsH := merchantcs.NewHandler(csSvc, idSvc)
	merchantFulfillH := merchantfulfillment.NewHandler(fulfillSvc)
	merchantInvoiceH := merchantinvoice.NewHandler(invoiceSvc)
	merchantLogisticsH := merchantlogistics.NewHandler(logisticsSvc)
	merchantProductMetaH := merchantproductmeta.NewHandler(productMetaSvc)
	merchantDiyH := merchantdiy.NewHandler(diySvc)
	openH := openapihttp.NewHandler(openSvc, catSvc, tradeSvc)
	serviceH := serviceportal.NewHandler(idSvc, jwtMgr, tradeSvc, csSvc, chatSvc)
	managerAuthH := managerauth.NewHandler(idSvc, jwtMgr)
	managerOrderH := managerorder.NewHandler(tradeSvc, idSvc)
	managerRefundH := managerrefund.NewHandler(aftersaleSvc, tradeSvc, idSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"service": "qixi_mergers_api_admin", "ok": true})
	})
	r.Static(cfg.Upload.PublicBase, cfg.Upload.Dir)

	platformPublic := r.Group("/api/platform/v1")
	platformAuthed := r.Group("/api/platform/v1")
	platformAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalPlatform))
	platformAuthH.Register(platformPublic, platformAuthed)
	platformMerH.Register(platformAuthed)
	platformCatH.Register(platformAuthed)
	platformOrderH.Register(platformAuthed)
	platformRefundH.Register(platformAuthed)
	platformFinanceH.Register(platformAuthed)
	platformCouponH.Register(platformAuthed)
	platformSpreadH.Register(platformAuthed)
	platformContentH.Register(platformAuthed)
	platformDiyH.Register(platformAuthed)
	platformSeckillH.Register(platformAuthed)
	platformComboH.Register(platformAuthed)
	platformPresellH.Register(platformAuthed)
	platformBroadcastH.Register(platformAuthed)
	platformCommunityH.Register(platformAuthed)
	platformAssistH.Register(platformAuthed)
	platformAttachH.Register(platformAuthed)
	platformSvipH.Register(platformAuthed)
	platformSettingH.Register(platformAuthed)
	platformLogisticsH.Register(platformAuthed)
	platformProductMetaH.Register(platformAuthed)
	platformArticleH.Register(platformAuthed)
	platformUserTagH.Register(platformAuthed)

	merchantPublic := r.Group("/api/merchant/v1")
	merchantAuthed := r.Group("/api/merchant/v1")
	merchantAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalMerchant))
	merchantAuthH.Register(merchantPublic, merchantAuthed)
	merchantCatH.Register(merchantAuthed)
	merchantOrderH.Register(merchantAuthed)
	merchantRefundH.Register(merchantAuthed)
	merchantFinanceH.Register(merchantAuthed)
	merchantCouponH.Register(merchantAuthed)
	merchantSeckillH.Register(merchantAuthed)
	merchantComboH.Register(merchantAuthed)
	merchantPresellH.Register(merchantAuthed)
	merchantSvipH.Register(merchantAuthed)
	merchantSettingH.Register(merchantAuthed)
	merchantCsH.Register(merchantAuthed)
	merchantFulfillH.Register(merchantAuthed)
	merchantInvoiceH.Register(merchantAuthed)
	merchantLogisticsH.Register(merchantAuthed)
	merchantProductMetaH.Register(merchantAuthed)
	merchantReserveH.Register(merchantAuthed)
	merchantBroadcastH.Register(merchantAuthed)
	merchantCommunityH.Register(merchantAuthed)
	merchantAssistH.Register(merchantAuthed)
	merchantAttachH.Register(merchantAuthed)
	merchantDiyH.Register(merchantAuthed)

	openPublic := r.Group("/api/open/v1")
	openAuthed := r.Group("/api/open/v1")
	openAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalOpen), openH.RequireAuth())
	openH.Register(openPublic, openAuthed)
	openPublic.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"prefix": "/api/open/v1"})
	})

	servicePublic := r.Group("/api/service/v1")
	serviceAuthed := r.Group("/api/service/v1")
	serviceAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalService))
	serviceH.Register(servicePublic, serviceAuthed)

	managerPublic := r.Group("/api/manager/v1")
	managerAuthed := r.Group("/api/manager/v1")
	managerAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalManager))
	managerAuthH.Register(managerPublic, managerAuthed)
	managerOrderH.Register(managerAuthed)
	managerRefundH.Register(managerAuthed)
	managerPublic.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"prefix": "/api/manager/v1"})
	})

	r.GET("/swagger/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api-admin swagger shell",
			"prefixes": []string{
				"/api/platform/v1", "/api/merchant/v1",
				"/api/manager/v1", "/api/service/v1", "/api/open/v1",
			},
		})
	})

	log.Printf("qixi_mergers_api_admin listening on %s", cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
