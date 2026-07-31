package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	adminauth "github.com/qixi-live/qixi-live-mergers/api-platform/internal/admin/auth"
	admincustomerservice "github.com/qixi-live/qixi-live-mergers/api-platform/internal/admin/customerservice"
	admindashboard "github.com/qixi-live/qixi-live-mergers/api-platform/internal/admin/dashboard"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/article"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/assist"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/attachment"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/broadcast"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/chat"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/circle"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/cloudconfig"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/combination"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/community"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/cs"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/finance"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/logistics"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/merchant"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/presell"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/productmeta"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/seckill"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/trade"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/domain/usertag"
	merchantapplicationevent "github.com/qixi-live/qixi-live-mergers/api-platform/internal/event/merchantapplication"
	platformdiyevent "github.com/qixi-live/qixi-live-mergers/api-platform/internal/event/platformdiy"
	articlepersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/article"
	assistpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/assist"
	attachmentpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/attachment"
	broadcastpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/broadcast"
	cartpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/cart"
	chatpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/chat"
	circlepersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/circle"
	cloudconfigpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/combination"
	communitypersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/community"
	contentpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/content"
	cspersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/cs"
	diypersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/diy"
	financepersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/finance"
	identitypersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/identity"
	logisticspersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/logistics"
	merchantpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/merchant"
	presellpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/presell"
	productmetapersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/productmeta"
	promotionpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/promotion"
	reservationpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/reservation"
	seckillpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/seckill"
	tradepersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/trade"
	usertagpersist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/infra/persist/usertag"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/config"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/db"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/response"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/pkg/upload"
	platformarticle "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/article"
	platformassist "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/assist"
	platformattachment "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/attachment"
	platformbroadcast "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/broadcast"
	platformcircle "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/circle"
	platformcloudconfig "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/cloudconfig"
	platformcombination "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/combination"
	platformcommunity "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/community"
	platformcontent "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/content"
	platformcoupon "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/coupon"
	platformdiy "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/diy"
	platformfinance "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/finance"
	platformlogistics "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/logistics"
	platformmerchant "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/merchant"
	nativecatalog "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/nativecatalog"
	platformnativeorder "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/nativeorder"
	platformnativerefund "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/nativerefund"
	platformpresell "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/presell"
	platformproductmeta "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/productmeta"
	platformseckill "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/seckill"
	platformsvip "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/svip"
	platformusertag "github.com/qixi-live/qixi-live-mergers/api-platform/internal/platform/usertag"
	"github.com/qixi-live/qixi-live-mergers/api-platform/internal/serviceportal"
)

func main() {
	cfgPath := flag.String("config", "conf/app.yaml", "path to app.yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn, err := cfg.DSNFor(config.DatabaseAdmin)
	if err != nil {
		log.Fatalf("platform database config: %v", err)
	}
	gdb, err := db.OpenMySQL(dsn, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	platformdiyevent.Start(context.Background(), gdb, cfg.NATS.URL)
	merchantApplicationProjection, err := merchantapplicationevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant application projection subscriber unavailable: %v", err)
	}
	if merchantApplicationProjection != nil {
		defer merchantApplicationProjection.Close()
	}
	businessDSN, err := cfg.DSNFor(config.DatabaseBusiness)
	if err != nil {
		log.Fatalf("business payment projection database config: %v", err)
	}
	businessDB, err := db.OpenMySQL(businessDSN, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("business payment projection mysql: %v", err)
	}
	merchantDSN, err := cfg.DSNFor(config.DatabaseMerchant)
	if err != nil {
		log.Fatalf("merchant product database config: %v", err)
	}
	merchantDB, err := db.OpenMySQL(merchantDSN, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("merchant product mysql: %v", err)
	}

	jwtMgr := authjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL(), cfg.JWT.RefreshTTL())
	idSvc := identity.NewService(identitypersist.NewRepo(gdb))
	merSvc := merchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(gdb)))
	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	circleSvc := circle.NewService(circlepersist.NewRepo(gdb))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
	financeSvc := finance.NewService(financepersist.NewStoreAdapter(financepersist.NewRepo(gdb)))
	contentSvc := content.NewService(contentpersist.NewRepo(gdb))
	cloudConfigSvc, err := cloudconfig.NewService(cloudconfigpersist.NewRepo(gdb), cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("cloud config crypto: %v", err)
	}
	paymentConfigStore, err := paymentconfig.NewStore(businessDB, cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("payment config projection crypto: %v", err)
	}
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
		Mode: cfg.IM.Mode, APIBase: cfg.IM.APIBase, APIPublicURL: cfg.IM.APIPublicURL, WSPublicURL: cfg.IM.WSPublicURL,
		AppID: cfg.IM.AppID, Token: cfg.IM.IntegrationToken, Secret: cfg.JWT.Secret,
	}, cloudConfigSvc)
	logisticsSvc := logistics.NewService(logisticspersist.NewRepo(gdb))
	productMetaSvc := productmeta.NewService(productmetapersist.NewRepo(gdb))
	articleSvc := article.NewService(articlepersist.NewRepo(gdb))
	userTagSvc := usertag.NewService(usertagpersist.NewRepo(gdb))
	cosStore := upload.COS{
		Bucket: cfg.COS.Bucket, Region: cfg.COS.Region,
		SecretID: cfg.COS.SecretID, SecretKey: cfg.COS.SecretKey,
		BaseURL: cfg.COS.BaseURL, KeyPrefix: cfg.COS.KeyPrefix,
	}
	var yamlFileUp upload.Store
	if cfg.COS.Enabled && cosStore.Configured() {
		yamlFileUp = cosStore
	} else {
		yamlFileUp = upload.Local{Dir: cfg.Upload.Dir, PublicBase: cfg.Upload.PublicBase}
	}
	// 后台数据库的 COS 开关优先；未启用时保持 app.yaml 的既有上传行为。
	fileUp := upload.DatabaseCOS{Resolver: cloudConfigSvc, Fallback: yamlFileUp}
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))

	platformAuthH := adminauth.NewHandler(gdb, jwtMgr)
	platformCustomerServiceH := admincustomerservice.NewHandler(gdb, businessDB)
	platformDashboardH := admindashboard.NewHandler(gdb, businessDB)
	platformMerH := platformmerchant.NewHandler(merSvc, idSvc)
	platformCircleH := platformcircle.NewHandler(circleSvc, idSvc)
	platformCatH := nativecatalog.NewHandler(gdb, merchantDB, businessDB, idSvc)
	platformOrderH := platformnativeorder.NewHandler(businessDB, merchantDB, idSvc)
	platformRefundH := platformnativerefund.NewHandler(businessDB, merchantDB, idSvc)
	platformFinanceH := platformfinance.NewHandler(financeSvc, idSvc)
	platformCouponH := platformcoupon.NewHandler(promoSvc, idSvc)
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
	platformCloudConfigH := platformcloudconfig.NewHandler(cloudConfigSvc, idSvc, paymentConfigStore)
	platformLogisticsH := platformlogistics.NewHandler(logisticsSvc)
	platformProductMetaH := platformproductmeta.NewHandler(productMetaSvc)
	platformArticleH := platformarticle.NewHandler(articleSvc)
	platformUserTagH := platformusertag.NewHandler(userTagSvc, idSvc)
	serviceH := serviceportal.NewHandler(idSvc, jwtMgr, tradeSvc, csSvc, chatSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.MaxMultipartMemory = 128 << 20
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"service": "qixi_mergers_api_platform", "ok": true})
	})
	r.Static(cfg.Upload.PublicBase, cfg.Upload.Dir)

	platformPublic := r.Group("/api/platform/v1")
	platformAuthed := r.Group("/api/platform/v1")
	platformAuthed.Use(
		middleware.JWTRequired(jwtMgr, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RestrictRegionConsole(),
	)
	platformAuthH.Register(platformPublic, platformAuthed)
	platformAuthH.RegisterSettings(platformAuthed)
	platformDashboardH.Register(platformAuthed)
	platformMerH.Register(platformAuthed)
	platformCircleH.Register(platformAuthed)
	platformCatH.Register(platformAuthed)
	platformOrderH.Register(platformAuthed)
	platformRefundH.Register(platformAuthed)
	platformFinanceH.Register(platformAuthed)
	platformCouponH.Register(platformAuthed)
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
	platformCloudConfigH.Register(platformAuthed)
	platformLogisticsH.Register(platformAuthed)
	platformProductMetaH.Register(platformAuthed)
	platformArticleH.Register(platformAuthed)
	platformUserTagH.Register(platformAuthed)
	platformCustomerServiceH.Register(platformAuthed)

	servicePublic := r.Group("/api/service/v1")
	serviceAuthed := r.Group("/api/service/v1")
	serviceAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalPlatform))
	serviceH.Register(servicePublic, serviceAuthed)

	r.GET("/swagger/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "api-platform swagger shell",
			"prefixes": []string{
				"/api/platform/v1", "/api/service/v1",
			},
		})
	})

	log.Printf("qixi_mergers_api_platform listening on %s", cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-Requested-With,Form-type,Referer,Connection,Content-Length,Host,Origin,Authori-zation,Accept,Accept-Encoding,Accept-Language,X-AppId")
		c.Header("Vary", "Origin, Access-Control-Request-Headers")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
