package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	apparticle "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/article"
	appassist "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/assist"
	appchat "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/chat"
	appcombination "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/combination"
	appcommunity "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/community"
	appcoupon "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/coupon"
	appinvoice "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/invoice"
	apppoints "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/points"
	apppresell "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/presell"
	appreservation "github.com/crmlive/qixi-live-ecrm/api-business/internal/app/reservation"
	businessaccount "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/account"
	businessaddress "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/address"
	businessauth "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/auth"
	businesscart "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/cart"
	businesscatalog "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/catalog"
	businesscontent "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/contentview"
	businesscoupon "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/coupon"
	businessdiyview "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/diyview"
	businessfavorite "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/favorite"
	businesslive "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/live"
	businessmarketing "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/marketing"
	businessmerchantapply "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/merchantapply"
	businessorder "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/order"
	businessrefund "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/refund"
	businessupload "github.com/crmlive/qixi-live-ecrm/api-business/internal/business/upload"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/article"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/assist"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/cart"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/catalog"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/chat"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/combination"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/community"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/invoice"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/presell"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/promotion"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/reservation"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/seckill"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/domain/trade"
	merchantapplicationevent "github.com/crmlive/qixi-live-ecrm/api-business/internal/event/merchantapplication"
	merchantdiyevent "github.com/crmlive/qixi-live-ecrm/api-business/internal/event/merchantdiy"
	merchantimevent "github.com/crmlive/qixi-live-ecrm/api-business/internal/event/merchantim"
	platformdiyevent "github.com/crmlive/qixi-live-ecrm/api-business/internal/event/platformdiy"
	articlepersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/article"
	assistpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/assist"
	cartpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/cart"
	catalogpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/catalog"
	chatpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/chat"
	cloudconfigpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/combination"
	communitypersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/community"
	invoicepersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/invoice"
	presellpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/presell"
	promotionpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/promotion"
	reservationpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/reservation"
	seckillpersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/seckill"
	tradepersist "github.com/crmlive/qixi-live-ecrm/api-business/internal/infra/persist/trade"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/captchaclient"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/config"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/db"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/qixi-live-ecrm/api-business/internal/pkg/upload"
)

func main() {
	cfgPath := flag.String("config", "conf/app.yaml", "path to app.yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn, err := cfg.DSNFor(config.DatabaseBusiness)
	if err != nil {
		log.Fatalf("business database config: %v", err)
	}
	gdb, err := db.OpenMySQL(dsn, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("mysql: %v", err)
	}
	imProjection, err := merchantimevent.StartBusinessProjection(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant IM projection subscriber unavailable: %v", err)
	}
	if imProjection != nil {
		defer imProjection.Close()
	}
	diyProjection, err := merchantdiyevent.StartBusinessProjection(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant DIY projection subscriber unavailable: %v", err)
	}
	if diyProjection != nil {
		defer diyProjection.Close()
	}
	platformDIYProjection, err := platformdiyevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("platform DIY projection subscriber unavailable: %v", err)
	}
	if platformDIYProjection != nil {
		defer platformDIYProjection.Close()
	}
	merchantapplicationevent.StartOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)

	jwtMgr := authjwt.NewManager(cfg.JWT.Secret, cfg.JWT.AccessTTL(), cfg.JWT.RefreshTTL())
	paymentConfigStore, err := paymentconfig.NewStore(gdb, cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("payment runtime config crypto: %v", err)
	}
	cloudConfigSvc, err := cloudconfig.NewService(cloudconfigpersist.NewRepo(gdb), cfg.JWT.Secret)
	if err != nil {
		log.Fatalf("cloud config crypto: %v", err)
	}
	catSvc := catalog.NewService(catalogpersist.NewStoreAdapter(catalogpersist.NewRepo(gdb)))
	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
	tradeSvc.SetPayment(trade.PaymentSettings{
		Sandbox:      cfg.Payment.Sandbox,
		NotifySecret: cfg.Payment.NotifySecret,
		Wechat:       cfg.Payment.Wechat,
		Alipay:       cfg.Payment.Alipay,
	})
	seckillSvc := seckill.NewService(seckillpersist.NewRepo(gdb))
	comboSvc := combination.NewService(combinationpersist.NewRepo(gdb))
	presellSvc := presell.NewService(presellpersist.NewRepo(gdb))
	reserveSvc := reservation.NewService(reservationpersist.NewRepo(gdb))
	communitySvc := community.NewService(communitypersist.NewRepo(gdb))
	assistSvc := assist.NewService(assistpersist.NewRepo(gdb))
	chatSvc := chat.NewService(chatpersist.NewRepo(gdb), chat.IMSettings{
		Mode: cfg.IM.Mode, APIBase: cfg.IM.APIBase, APIPublicURL: cfg.IM.APIPublicURL, WSPublicURL: cfg.IM.WSPublicURL,
		AppID: cfg.IM.AppID, Token: cfg.IM.IntegrationToken, Secret: cfg.JWT.Secret,
	}, cloudConfigSvc)
	invoiceSvc := invoice.NewService(invoicepersist.NewRepo(gdb))
	articleSvc := article.NewService(articlepersist.NewRepo(gdb))
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))

	captchaClient, captchaErr := captchaclient.New(cfg.Captcha)
	if captchaErr != nil {
		log.Printf("pte-tools-captcha is unavailable: %v", captchaErr)
	}
	appH := businessauth.NewHandler(businessauth.NewService(gdb), jwtMgr, captchaClient)
	appCatH := businesscatalog.NewHandler(gdb)
	appMerchantApplyH := businessmerchantapply.NewHandler(gdb)
	appUploadH := businessupload.NewHandler(gdb, upload.DatabaseCOS{Resolver: cloudConfigSvc})
	appAddrH := businessaddress.NewHandler(gdb)
	appAccountH := businessaccount.NewHandler(gdb)
	appCartH := businesscart.NewHandler(gdb)
	appOrderH := businessorder.NewHandler(gdb, paymentConfigStore, cfg.Payment.Sandbox, cloudConfigSvc)
	appOrderCallbackH := businessorder.NewCallbackHandler(gdb, paymentConfigStore, cfg.Payment.Sandbox, cloudConfigSvc)
	appPointsH := apppoints.NewHandler(tradeSvc, catSvc)
	appRefundH := businessrefund.NewHandler(gdb)
	appCouponH := appcoupon.NewHandler(promoSvc, cartSvc)
	appBusinessCouponH := businesscoupon.NewHandler(gdb)
	appContentH := businesscontent.NewHandler(gdb)
	appDiyH := businessdiyview.NewHandler(gdb)
	appFavoriteH := businessfavorite.NewHandler(gdb)
	appSeckillH := businessmarketing.NewSeckillHandler(gdb)
	appComboH := appcombination.NewHandler(comboSvc, tradeSvc)
	appPresellH := apppresell.NewHandler(presellSvc, tradeSvc)
	appReserveH := appreservation.NewHandler(reserveSvc, tradeSvc)
	appLiveH := businesslive.NewHandler(gdb)
	appCommunityH := appcommunity.NewHandler(communitySvc)
	appAssistH := appassist.NewHandler(assistSvc, tradeSvc)
	appChatH := appchat.NewHandler(chatSvc)
	appInvoiceH := appinvoice.NewHandler(invoiceSvc)
	appArticleH := apparticle.NewHandler(articleSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"service": "qixi_live_ecrm_api_business", "ok": true})
	})

	appPublic := r.Group("/api/app/v1")
	appOptional := r.Group("/api/app/v1")
	appOptional.Use(middleware.JWTOptional(jwtMgr, authjwt.PortalApp))
	appAuthed := r.Group("/api/app/v1")
	appAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalApp))
	appH.Register(appPublic, appAuthed)
	appH.RegisterCaptchaGateway(r)
	appCatH.Register(appPublic)
	appContentH.Register(appPublic)
	appArticleH.Register(appPublic)
	appDiyH.Register(appPublic)
	appSeckillH.Register(appPublic)
	appComboH.RegisterPublic(appPublic)
	appPresellH.RegisterPublic(appPublic)
	appReserveH.RegisterPublic(appPublic)
	appLiveH.RegisterPublic(appPublic)
	appCommunityH.RegisterPublic(appPublic)
	appAssistH.RegisterPublic(appPublic)
	appPointsH.RegisterPublic(appPublic)
	appBusinessCouponH.RegisterPublic(appOptional)
	appAddrH.Register(appAuthed)
	appAccountH.Register(appAuthed)
	appCartH.Register(appAuthed)
	appOrderH.Register(appAuthed)
	appComboH.RegisterAuthed(appAuthed)
	appPresellH.RegisterAuthed(appAuthed)
	appReserveH.RegisterAuthed(appAuthed)
	appAssistH.RegisterAuthed(appAuthed)
	appPointsH.RegisterAuthed(appAuthed)
	appRefundH.Register(appAuthed)
	appCouponH.RegisterSpread(appAuthed)
	appBusinessCouponH.Register(appAuthed)
	appCommunityH.RegisterAuthed(appAuthed)
	appChatH.Register(appAuthed)
	appInvoiceH.Register(appAuthed)
	appLiveH.RegisterAuthed(appAuthed)
	appMerchantApplyH.Register(appAuthed)
	appUploadH.Register(appAuthed)
	appFavoriteH.Register(appAuthed)

	cb := r.Group("/api/callback/v1")
	cb.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"prefix": "/api/callback/v1"})
	})
	appOrderCallbackH.Register(cb)

	r.GET("/swagger/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "api-business swagger shell",
			"prefixes": []string{"/api/app/v1", "/api/callback/v1"},
		})
	})

	log.Printf("qixi_live_ecrm_api_business listening on %s", cfg.Server.Addr)
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
