package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/aftersale"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/assist"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/attachment"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/broadcast"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/cart"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/cloudconfig"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/combination"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/community"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/cs"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/diy"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/finance"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/fulfillment"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/identity"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/invoice"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/logistics"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/merchant"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/presell"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/productmeta"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/promotion"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/reservation"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/seckill"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/domain/trade"
	merchantdiyevent "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/event/merchantdiy"
	merchantimevent "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/event/merchantim"
	aftersalepersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/aftersale"
	assistpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/assist"
	attachmentpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/attachment"
	broadcastpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/broadcast"
	cartpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/cart"
	cloudconfigpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/combination"
	communitypersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/community"
	cspersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/cs"
	diypersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/diy"
	financepersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/finance"
	fulfillmentpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/fulfillment"
	identitypersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/identity"
	invoicepersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/invoice"
	logisticspersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/logistics"
	merchantpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/merchant"
	presellpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/presell"
	productmetapersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/productmeta"
	promotionpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/promotion"
	reservationpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/reservation"
	seckillpersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/seckill"
	tradepersist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/infra/persist/trade"
	managerauth "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/manager/auth"
	managerorder "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/manager/order"
	managerrefund "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/manager/refund"
	merchantassist "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/assist"
	merchantattachment "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/attachment"
	merchantbroadcast "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/broadcast"
	merchantcombination "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/combination"
	merchantcommunity "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/community"
	merchantcoupon "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/coupon"
	merchantcs "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/cs"
	merchantdashboard "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/dashboard"
	merchantdiy "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/diy"
	merchantfinance "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/finance"
	merchantfulfillment "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/fulfillment"
	merchantimsdk "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/imsdk"
	merchantinvoice "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/invoice"
	merchantlogistics "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/logistics"
	merchantcatalog "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/nativecatalog"
	nativeorder "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/nativeorder"
	nativerefund "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/nativerefund"
	merchantorder "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/order"
	merchantpayment "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/payment"
	merchantpresell "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/presell"
	merchantproductmeta "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/productmeta"
	merchantreservation "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/reservation"
	merchantseckill "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/seckill"
	merchantsetting "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/setting"
	merchantsvip "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/merchant/svip"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/authjwt"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/config"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/db"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/response"
	"github.com/crmlive/qixi-live-ecrm/api-merchant/internal/pkg/upload"
	storeauth "github.com/crmlive/qixi-live-ecrm/api-merchant/internal/store/auth"
)

func main() {
	cfgPath := flag.String("config", "conf/app.yaml", "path to app.yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn, err := cfg.DSNFor(config.DatabaseMerchant)
	if err != nil {
		log.Fatalf("merchant database config: %v", err)
	}
	gdb, err := db.OpenMySQL(dsn, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	merchantimevent.StartMerchantOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	merchantdiyevent.StartMerchantOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	businessDSN, err := cfg.DSNFor(config.DatabaseBusiness)
	if err != nil {
		log.Fatalf("business payment projection database config: %v", err)
	}
	businessDB, err := db.OpenMySQL(businessDSN, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("business payment projection mysql: %v", err)
	}

	jwtMgr := authjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL(), cfg.JWT.RefreshTTL())
	idSvc := identity.NewService(identitypersist.NewRepo(gdb))
	merSvc := merchant.NewService(merchantpersist.NewStoreAdapter(merchantpersist.NewRepo(gdb)))
	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
	aftersaleSvc := aftersale.NewService(aftersalepersist.NewStoreAdapter(aftersalepersist.NewRepo(gdb)))
	financeSvc := finance.NewService(financepersist.NewStoreAdapter(financepersist.NewRepo(gdb)))
	cloudConfigSvc, err := cloudconfig.NewService(cloudconfigpersist.NewRepo(gdb), cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("cloud config crypto: %v", err)
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
	fulfillSvc := fulfillment.NewService(fulfillmentpersist.NewRepo(gdb))
	invoiceSvc := invoice.NewService(invoicepersist.NewRepo(gdb))
	logisticsSvc := logistics.NewService(logisticspersist.NewRepo(gdb))
	productMetaSvc := productmeta.NewService(productmetapersist.NewRepo(gdb))
	// 上传只使用后台加密配置；未启用或未补齐密钥时明确拒绝上传。
	fileUp := upload.DatabaseCOS{Resolver: cloudConfigSvc}
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))

	storeAuthH := storeauth.NewHandler(gdb, jwtMgr)
	merchantDashboardH := merchantdashboard.NewHandler(gdb, businessDB)
	merchantCatH := merchantcatalog.NewHandler(gdb, businessDB, idSvc)
	merchantOrderH := merchantorder.NewHandler(tradeSvc, idSvc, logisticsSvc)
	merchantNativeOrderH := nativeorder.NewHandler(businessDB, idSvc)
	merchantRefundH := nativerefund.NewHandler(businessDB, idSvc)
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
	merchantPaymentH := merchantpayment.NewHandler(gdb, businessDB, cfg.JWT.Secret)
	merchantIMSDKH := merchantimsdk.NewHandler(gdb)
	merchantCsH := merchantcs.NewHandler(csSvc, idSvc)
	merchantFulfillH := merchantfulfillment.NewHandler(fulfillSvc)
	merchantInvoiceH := merchantinvoice.NewHandler(invoiceSvc)
	merchantLogisticsH := merchantlogistics.NewHandler(logisticsSvc)
	merchantProductMetaH := merchantproductmeta.NewHandler(productMetaSvc)
	merchantDiyH := merchantdiy.NewHandler(diySvc)
	managerAuthH := managerauth.NewHandler(idSvc, jwtMgr)
	managerOrderH := managerorder.NewHandler(tradeSvc, idSvc)
	managerRefundH := managerrefund.NewHandler(aftersaleSvc, tradeSvc, idSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.MaxMultipartMemory = 128 << 20
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", healthHandler)
	r.Static(cfg.Upload.PublicBase, cfg.Upload.Dir)

	merchantPublic := r.Group("/api/merchant/v1")
	merchantAuthed := r.Group("/api/merchant/v1")
	merchantAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalMerchant), middleware.RequireStoreConsole(), middleware.RequireStoreAppID())
	storeAuthH.Register(merchantPublic, merchantAuthed)
	merchantDashboardH.Register(merchantAuthed)
	merchantCatH.Register(merchantAuthed)
	merchantOrderH.RegisterVerify(merchantAuthed)
	merchantNativeOrderH.Register(merchantAuthed)
	merchantRefundH.Register(merchantAuthed)
	merchantFinanceH.Register(merchantAuthed)
	merchantCouponH.Register(merchantAuthed)
	merchantSeckillH.Register(merchantAuthed)
	merchantComboH.Register(merchantAuthed)
	merchantPresellH.Register(merchantAuthed)
	merchantSvipH.Register(merchantAuthed)
	merchantSettingH.Register(merchantAuthed)
	merchantPaymentH.Register(merchantAuthed)
	merchantIMSDKH.Register(merchantAuthed)
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

	// 兼容期内店员/履约入口沿用 /api/manager/v1；仅本服务注册。
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
			"message": "api-merchant swagger shell",
			"prefixes": []string{
				"/api/merchant/v1", "/api/manager/v1",
			},
		})
	})

	log.Printf("qixi_live_ecrm_api_merchant listening on %s", cfg.Server.Addr)
	if err := r.Run(cfg.Server.Addr); err != nil {
		log.Fatalf("run: %v", err)
	}
}

func healthHandler(c *gin.Context) {
	response.OK(c, gin.H{"service": "qixi_live_ecrm_api_merchant", "ok": true})
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
