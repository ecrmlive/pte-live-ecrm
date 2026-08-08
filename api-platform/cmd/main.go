package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	adminauth "github.com/crmlive/pte-live-ecrm/api-platform/internal/admin/auth"
	admincustomerservice "github.com/crmlive/pte-live-ecrm/api-platform/internal/admin/customerservice"
	admindashboard "github.com/crmlive/pte-live-ecrm/api-platform/internal/admin/dashboard"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/article"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/attachment"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/broadcast"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cart"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/chat"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/circle"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/combination"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/community"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/content"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/cs"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/diy"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/identity"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/logistics"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/merchant"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/presell"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/promotion"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/reservation"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/seckill"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/trade"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/domain/usertag"
	commentmoderationevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/commentmoderation"
	feedbackmoderationevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/feedbackmoderation"
	merchantapplicationevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantapplication"
	merchantonboardingevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantonboarding"
	merchantsettlementevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/merchantsettlement"
	platformcityevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/platformcity"
	platformdiyevent "github.com/crmlive/pte-live-ecrm/api-platform/internal/event/platformdiy"
	articlepersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/article"
	assistpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/assist"
	attachmentpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/attachment"
	broadcastpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/broadcast"
	cartpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/cart"
	chatpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/chat"
	circlepersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/circle"
	cloudconfigpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/combination"
	communitypersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/community"
	contentpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/content"
	cspersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/cs"
	diypersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/diy"
	identitypersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/identity"
	logisticspersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/logistics"
	merchantpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/merchant"
	presellpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/presell"
	promotionpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/promotion"
	reservationpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/reservation"
	seckillpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/seckill"
	tradepersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/trade"
	usertagpersist "github.com/crmlive/pte-live-ecrm/api-platform/internal/infra/persist/usertag"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/config"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/db"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/pkg/upload"
	platformarticle "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/article"
	platformassist "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/assist"
	platformattachment "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/attachment"
	platformbroadcast "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/broadcast"
	platformcircle "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/circle"
	platformcloudconfig "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/cloudconfig"
	platformcombination "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/combination"
	platformcomment "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/comment"
	platformcommunity "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/community"
	platformcontent "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/content"
	platformcoupon "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/coupon"
	platformdiy "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/diy"
	platformfeedback "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/feedback"
	platforminvoice "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/invoice"
	platformlogistics "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/logistics"
	platformmemberlevel "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/memberlevel"
	platformmerchant "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/merchant"
	platformmerchantdeposit "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/merchantdeposit"
	platformmerchanttype "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/merchanttype"
	platformstoremenu "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/storemenu"
	nativecatalog "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativecatalog"
	platformnativeconfigitem "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativeconfigitem"
	platformnativediscount "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativediscount"
	platformnativedistribution "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativedistribution"
	platformnativeledger "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativeledger"
	platformnativemarketingdecor "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativemarketingdecor"
	platformnativeorder "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativeorder"
	platformnativerefund "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativerefund"
	platformnativesettlement "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativesettlement"
	platformnativewithdraw "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/nativewithdraw"
	platformoperationlog "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/operationlog"
	platformpoints "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/points"
	platformpresell "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/presell"
	platformproductmeta "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/productmeta"
	platformprofitsharing "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/profitsharing"
	platformrecharge "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/recharge"
	platformseckill "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/seckill"
	platformstoregroup "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/storegroup"
	platformsvip "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/svip"
	platformsvipinterest "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/svipinterest"
	platformuserlist "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/userlist"
	platformusersearch "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/usersearch"
	platformusertag "github.com/crmlive/pte-live-ecrm/api-platform/internal/platform/usertag"
	"github.com/crmlive/pte-live-ecrm/api-platform/internal/serviceportal"
	"github.com/gin-gonic/gin"
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
	platformcityevent.Start(context.Background(), gdb, cfg.NATS.URL)
	merchantapplicationevent.StartReviewOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	merchantApplicationProjection, err := merchantapplicationevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant application projection subscriber unavailable: %v", err)
	}
	if merchantApplicationProjection != nil {
		defer merchantApplicationProjection.Close()
	}
	merchantSettlementProjection, err := merchantsettlementevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant settlement projection subscriber unavailable: %v", err)
	}
	if merchantSettlementProjection != nil {
		defer merchantSettlementProjection.Close()
	}
	merchantSettlementCommands, err := merchantsettlementevent.NewCommandClient(cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant settlement command client unavailable: %v", err)
	}
	if merchantSettlementCommands != nil {
		defer merchantSettlementCommands.Close()
	}
	productCommentCommands, err := commentmoderationevent.New(cfg.NATS.URL)
	if err != nil {
		log.Printf("product comment moderation command client unavailable: %v", err)
	}
	if productCommentCommands != nil {
		defer productCommentCommands.Close()
	}
	feedbackCommands, err := feedbackmoderationevent.New(cfg.NATS.URL)
	if err != nil {
		log.Printf("feedback moderation command client unavailable: %v", err)
	}
	if feedbackCommands != nil {
		defer feedbackCommands.Close()
	}
	merchantOnboarding, err := merchantonboardingevent.New(cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant onboarding command client unavailable: %v", err)
	}
	if merchantOnboarding != nil {
		defer merchantOnboarding.Close()
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
	storeJWTMgr := authjwt.NewManager(cfg.JWT.StoreConsoleSecret(), cfg.JWT.AccessTTL(), cfg.JWT.RefreshTTL())
	idSvc := identity.NewService(identitypersist.NewRepo(gdb))
	merSvc := merchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(gdb)))
	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	circleSvc := circle.NewService(circlepersist.NewRepo(gdb))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
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
	seckillSvc := seckill.NewService(seckillpersist.NewRepo(businessDB))
	comboSvc := combination.NewService(combinationpersist.NewRepo(businessDB))
	presellSvc := presell.NewService(presellpersist.NewRepo(businessDB))
	reserveSvc := reservation.NewService(reservationpersist.NewRepo(gdb))
	broadcastSvc := broadcast.NewService(broadcastpersist.NewRepo(businessDB))
	communitySvc := community.NewService(communitypersist.NewRepo(businessDB))
	assistSvc := assist.NewService(assistpersist.NewRepo(businessDB))
	attachSvc := attachment.NewService(attachmentpersist.NewRepo(gdb))
	csSvc := cs.NewService(cspersist.NewRepo(gdb))
	chatSvc := chat.NewService(chatpersist.NewRepo(gdb), chat.IMSettings{
		Mode: cfg.IM.Mode, APIBase: cfg.IM.APIBase, APIPublicURL: cfg.IM.APIPublicURL, WSPublicURL: cfg.IM.WSPublicURL,
		AppID: cfg.IM.AppID, Token: cfg.IM.IntegrationToken, Secret: cfg.JWT.Secret,
	}, cloudConfigSvc)
	logisticsSvc := logistics.NewService(logisticspersist.NewRepo(gdb))
	articleSvc := article.NewService(articlepersist.NewRepo(gdb))
	userTagSvc := usertag.NewService(usertagpersist.NewRepo(businessDB))
	// 上传只使用后台加密配置；未启用或未补齐密钥时明确拒绝上传。
	fileUp := upload.DatabaseCOS{Resolver: cloudConfigSvc}
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))

	platformAuthH := adminauth.NewHandler(gdb, jwtMgr)
	platformCustomerServiceH := admincustomerservice.NewHandler(gdb, businessDB)
	platformDashboardH := admindashboard.NewHandler(gdb, businessDB, merchantDB)
	platformMerH := platformmerchant.NewHandler(merSvc, idSvc, gdb, merchantOnboarding).WithStoreLogin(merchantDB, storeJWTMgr)
	platformMemberLevelH := platformmemberlevel.New(businessDB, gdb)
	platformMerchantTypeH := platformmerchanttype.NewHandler(gdb)
	platformMerchantDepositH := platformmerchantdeposit.NewHandler(gdb)
	platformStoreMenuH := platformstoremenu.NewHandler(merchantDB)
	platformProfitsharingH := platformprofitsharing.NewHandler(gdb)
	platformStoreGroupH := platformstoregroup.NewHandler(gdb)
	platformCircleH := platformcircle.NewHandler(circleSvc, gdb, businessDB)
	platformCatH := nativecatalog.NewHandler(gdb, merchantDB, businessDB)
	platformCatH.StartAuditOutboxDispatcher(context.Background())
	platformCatH.StartProjectionDispatcher(context.Background())
	platformOrderH := platformnativeorder.NewHandler(businessDB, merchantDB, gdb)
	platformRefundH := platformnativerefund.NewHandler(businessDB, merchantDB, gdb)
	platformFinanceH := platformnativewithdraw.NewHandler(businessDB, gdb)
	platformLedgerH := platformnativeledger.NewHandler(businessDB, gdb)
	platformDistributionH := platformnativedistribution.NewHandler(businessDB, gdb)
	platformSettlementH := platformnativesettlement.NewHandler(gdb, merchantSettlementCommands)
	platformDiscountH := platformnativediscount.NewHandler(businessDB, gdb)
	platformMarketingDecorH := platformnativemarketingdecor.NewHandler(gdb)
	platformConfigItemH := platformnativeconfigitem.NewHandler(gdb)
	platformCouponH := platformcoupon.NewHandler(promoSvc, gdb)
	platformContentH := platformcontent.NewHandler(contentSvc, gdb)
	platformDiyH := platformdiy.NewHandler(diySvc, gdb)
	platformSeckillH := platformseckill.NewHandler(seckillSvc, gdb)
	platformComboH := platformcombination.NewHandler(comboSvc, gdb)
	platformPresellH := platformpresell.NewHandler(presellSvc, gdb)
	platformBroadcastH := platformbroadcast.NewHandler(broadcastSvc, gdb)
	platformCommunityH := platformcommunity.NewHandler(communitySvc, gdb)
	platformAssistH := platformassist.NewHandler(assistSvc, gdb)
	platformAttachH := platformattachment.NewHandler(attachSvc, gdb, fileUp)
	platformSvipH := platformsvip.NewHandler(idSvc, gdb, businessDB)
	platformSvipInterestH := platformsvipinterest.New(businessDB, gdb)
	platformCloudConfigH := platformcloudconfig.NewHandler(cloudConfigSvc, idSvc, paymentConfigStore)
	platformLogisticsH := platformlogistics.NewHandler(logisticsSvc, gdb)
	platformProductMetaH := platformproductmeta.NewHandler(gdb)
	platformCommentH := platformcomment.NewHandler(businessDB, gdb, productCommentCommands)
	platformFeedbackH := platformfeedback.New(businessDB, gdb, feedbackCommands)
	platformInvoiceH := platforminvoice.New(businessDB, gdb)
	platformArticleH := platformarticle.NewHandler(articleSvc, gdb)
	platformUserTagH := platformusertag.NewHandler(userTagSvc, idSvc)
	platformUserListH := platformuserlist.New(businessDB, gdb)
	platformUserSearchH := platformusersearch.New(businessDB, gdb)
	platformOperationLogH := platformoperationlog.New(gdb)
	platformPointsH := platformpoints.NewHandler(businessDB, gdb)
	platformRechargeH := platformrecharge.NewHandler(businessDB, gdb)
	serviceH := serviceportal.NewHandler(idSvc, jwtMgr, tradeSvc, csSvc, chatSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.MaxMultipartMemory = 128 << 20
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"service": "pte_live_ecrm_api_platform", "ok": true})
	})
	r.Static(cfg.Upload.PublicBase, cfg.Upload.Dir)

	platformPublic := r.Group("/api/platform/v1")
	platformAuthed := r.Group("/api/platform/v1")
	platformAuthed.Use(
		middleware.JWTRequired(jwtMgr, authjwt.PortalPlatform),
		middleware.RequireAdminConsole(),
		middleware.RequireAdminSession(gdb),
		middleware.RestrictRoleConsole(),
		middleware.RestrictRegionConsole(),
		middleware.AuditAdminMutation(gdb),
	)
	platformAuthH.Register(platformPublic, platformAuthed)
	platformAuthH.RegisterSettings(platformAuthed)
	platformDashboardH.Register(platformAuthed)
	platformMerH.Register(platformAuthed)
	platformMemberLevelH.Register(platformAuthed)
	platformMerchantTypeH.Register(platformAuthed)
	platformMerchantDepositH.Register(platformAuthed)
	platformStoreMenuH.Register(platformAuthed)
	platformProfitsharingH.Register(platformAuthed)
	platformStoreGroupH.Register(platformAuthed)
	platformCircleH.Register(platformAuthed)
	platformCatH.Register(platformAuthed)
	platformOrderH.Register(platformAuthed)
	platformRefundH.Register(platformAuthed)
	platformFinanceH.Register(platformAuthed)
	platformLedgerH.Register(platformAuthed)
	platformDistributionH.Register(platformAuthed)
	platformSettlementH.Register(platformAuthed)
	platformDiscountH.Register(platformAuthed)
	platformMarketingDecorH.Register(platformAuthed)
	platformConfigItemH.Register(platformAuthed)
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
	platformSvipInterestH.Register(platformAuthed)
	platformCloudConfigH.Register(platformAuthed)
	platformLogisticsH.Register(platformAuthed)
	platformProductMetaH.Register(platformAuthed)
	platformCommentH.Register(platformAuthed)
	platformFeedbackH.Register(platformAuthed)
	platformInvoiceH.Register(platformAuthed)
	platformArticleH.Register(platformAuthed)
	platformUserTagH.Register(platformAuthed)
	platformUserListH.Register(platformAuthed)
	platformUserSearchH.Register(platformAuthed)
	platformOperationLogH.Register(platformAuthed)
	platformPointsH.Register(platformAuthed)
	platformRechargeH.Register(platformAuthed)
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

	log.Printf("pte_live_ecrm_api_platform listening on %s", cfg.Server.Addr)
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
