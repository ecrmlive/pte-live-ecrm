package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	apparticle "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/article"
	appassist "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/assist"
	appcallback "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/callback"
	appchat "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/chat"
	appcombination "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/combination"
	appcommunity "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/community"
	appcontent "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/content"
	appcoupon "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/coupon"
	appdiy "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/diy"
	appinvoice "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/invoice"
	apppoints "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/points"
	apppresell "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/presell"
	appreservation "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/reservation"
	appseckill "github.com/qixi-live/qixi-live-mergers/api-business/internal/app/seckill"
	businessaddress "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/address"
	businessauth "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/auth"
	businesscart "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/cart"
	businesscatalog "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/catalog"
	businesslive "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/live"
	businessorder "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/order"
	businessrefund "github.com/qixi-live/qixi-live-mergers/api-business/internal/business/refund"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/article"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/assist"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/chat"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/cloudconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/combination"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/community"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/invoice"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/presell"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/seckill"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/domain/trade"
	merchantimevent "github.com/qixi-live/qixi-live-mergers/api-business/internal/event/merchantim"
	articlepersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/article"
	assistpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/assist"
	cartpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/cart"
	catalogpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/catalog"
	chatpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/chat"
	cloudconfigpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/combination"
	communitypersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/community"
	contentpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/content"
	diypersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/diy"
	invoicepersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/invoice"
	presellpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/presell"
	promotionpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/promotion"
	reservationpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/reservation"
	seckillpersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/seckill"
	tradepersist "github.com/qixi-live/qixi-live-mergers/api-business/internal/infra/persist/trade"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/paymentconfig"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/config"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/db"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api-business/internal/pkg/response"
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
	contentSvc := content.NewService(contentpersist.NewRepo(gdb))
	diySvc := diy.NewService(diypersist.NewRepo(gdb))
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

	appH := businessauth.NewHandler(businessauth.NewService(gdb), jwtMgr)
	appCatH := businesscatalog.NewHandler(gdb)
	appAddrH := businessaddress.NewHandler(gdb)
	appCartH := businesscart.NewHandler(gdb)
	appOrderH := businessorder.NewHandler(gdb, paymentConfigStore, cfg.Payment.Sandbox)
	appOrderCallbackH := businessorder.NewCallbackHandler(gdb, cfg.Payment.Sandbox)
	appPointsH := apppoints.NewHandler(tradeSvc, catSvc)
	appRefundH := businessrefund.NewHandler(gdb)
	appCouponH := appcoupon.NewHandler(promoSvc, cartSvc)
	appContentH := appcontent.NewHandler(contentSvc)
	appDiyH := appdiy.NewHandler(diySvc)
	appSeckillH := appseckill.NewHandler(seckillSvc)
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
		response.OK(c, gin.H{"service": "qixi_mergers_api_business", "ok": true})
	})

	appPublic := r.Group("/api/app/v1")
	appAuthed := r.Group("/api/app/v1")
	appAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalApp))
	appH.Register(appPublic, appAuthed)
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
	appAddrH.Register(appAuthed)
	appCartH.Register(appAuthed)
	appOrderH.Register(appAuthed)
	appComboH.RegisterAuthed(appAuthed)
	appPresellH.RegisterAuthed(appAuthed)
	appReserveH.RegisterAuthed(appAuthed)
	appAssistH.RegisterAuthed(appAuthed)
	appPointsH.RegisterAuthed(appAuthed)
	appRefundH.Register(appAuthed)
	appCouponH.Register(appAuthed)
	appCommunityH.RegisterAuthed(appAuthed)
	appChatH.Register(appAuthed)
	appInvoiceH.Register(appAuthed)
	appLiveH.RegisterAuthed(appAuthed)

	cb := r.Group("/api/callback/v1")
	cb.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"prefix": "/api/callback/v1"})
	})
	appOrderCallbackH.RegisterMock(cb)
	appcallback.NewPayHandler(tradeSvc).Register(cb)

	r.GET("/swagger/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "api-business swagger shell",
			"prefixes": []string{"/api/app/v1", "/api/callback/v1"},
		})
	})

	log.Printf("qixi_mergers_api_business listening on %s", cfg.Server.Addr)
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
