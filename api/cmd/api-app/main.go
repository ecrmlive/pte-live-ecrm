package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	appaddress "github.com/qixi-live/qixi-live-mergers/api/internal/app/address"
	appauth "github.com/qixi-live/qixi-live-mergers/api/internal/app/auth"
	appcallback "github.com/qixi-live/qixi-live-mergers/api/internal/app/callback"
	appcart "github.com/qixi-live/qixi-live-mergers/api/internal/app/cart"
	appcatalog "github.com/qixi-live/qixi-live-mergers/api/internal/app/catalog"
	appcombination "github.com/qixi-live/qixi-live-mergers/api/internal/app/combination"
	appassist "github.com/qixi-live/qixi-live-mergers/api/internal/app/assist"
	appbroadcast "github.com/qixi-live/qixi-live-mergers/api/internal/app/broadcast"
	appcommunity "github.com/qixi-live/qixi-live-mergers/api/internal/app/community"
	apppresell "github.com/qixi-live/qixi-live-mergers/api/internal/app/presell"
	appreservation "github.com/qixi-live/qixi-live-mergers/api/internal/app/reservation"
	appcontent "github.com/qixi-live/qixi-live-mergers/api/internal/app/content"
	appcoupon "github.com/qixi-live/qixi-live-mergers/api/internal/app/coupon"
	appdiy "github.com/qixi-live/qixi-live-mergers/api/internal/app/diy"
	apporder "github.com/qixi-live/qixi-live-mergers/api/internal/app/order"
	apppoints "github.com/qixi-live/qixi-live-mergers/api/internal/app/points"
	apprefund "github.com/qixi-live/qixi-live-mergers/api/internal/app/refund"
	appseckill "github.com/qixi-live/qixi-live-mergers/api/internal/app/seckill"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/aftersale"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/assist"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/cart"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/catalog"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/combination"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/content"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/diy"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/identity"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/broadcast"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/community"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/presell"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/promotion"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/reservation"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/seckill"
	"github.com/qixi-live/qixi-live-mergers/api/internal/domain/trade"
	aftersalepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/aftersale"
	assistpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/assist"
	broadcastpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/broadcast"
	cartpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/cart"
	communitypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/community"
	catalogpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/catalog"
	combinationpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/combination"
	contentpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/content"
	diypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/diy"
	identitypersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/identity"
	presellpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/presell"
	promotionpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/promotion"
	reservationpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/reservation"
	seckillpersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/seckill"
	tradepersist "github.com/qixi-live/qixi-live-mergers/api/internal/infra/persist/trade"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/authjwt"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/config"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/db"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/middleware"
	"github.com/qixi-live/qixi-live-mergers/api/internal/pkg/response"
)

func main() {
	cfgPath := flag.String("config", "conf/app.yaml", "path to app.yaml")
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
	aftersaleSvc := aftersale.NewService(aftersalepersist.NewStoreAdapter(aftersalepersist.NewRepo(gdb)))
	contentSvc := content.NewService(contentpersist.NewRepo(gdb))
	diySvc := diy.NewService(diypersist.NewRepo(gdb))
	seckillSvc := seckill.NewService(seckillpersist.NewRepo(gdb))
	comboSvc := combination.NewService(combinationpersist.NewRepo(gdb))
	presellSvc := presell.NewService(presellpersist.NewRepo(gdb))
	reserveSvc := reservation.NewService(reservationpersist.NewRepo(gdb))
	broadcastSvc := broadcast.NewService(broadcastpersist.NewRepo(gdb))
	communitySvc := community.NewService(communitypersist.NewRepo(gdb))
	assistSvc := assist.NewService(assistpersist.NewRepo(gdb))
	tradeSvc.SetSeckill(seckillSvc)
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetReservation(reservation.NewTradeBridge(reserveSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))

	appH := appauth.NewHandler(idSvc, jwtMgr)
	appCatH := appcatalog.NewHandler(catSvc, diySvc)
	appAddrH := appaddress.NewHandler(cartSvc)
	appCartH := appcart.NewHandler(cartSvc)
	appOrderH := apporder.NewHandler(tradeSvc)
	appPointsH := apppoints.NewHandler(tradeSvc, catSvc)
	appRefundH := apprefund.NewHandler(aftersaleSvc)
	appCouponH := appcoupon.NewHandler(promoSvc, cartSvc)
	appContentH := appcontent.NewHandler(contentSvc)
	appDiyH := appdiy.NewHandler(diySvc)
	appSeckillH := appseckill.NewHandler(seckillSvc)
	appComboH := appcombination.NewHandler(comboSvc, tradeSvc)
	appPresellH := apppresell.NewHandler(presellSvc, tradeSvc)
	appReserveH := appreservation.NewHandler(reserveSvc, tradeSvc)
	appBroadcastH := appbroadcast.NewHandler(broadcastSvc)
	appCommunityH := appcommunity.NewHandler(communitySvc)
	appAssistH := appassist.NewHandler(assistSvc, tradeSvc)

	gin.SetMode(cfg.Server.Mode)
	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), corsMiddleware())

	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"service": "qixi_mergers_api_app", "ok": true})
	})

	appPublic := r.Group("/api/app/v1")
	appAuthed := r.Group("/api/app/v1")
	appAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalApp))
	appH.Register(appPublic, appAuthed)
	appCatH.Register(appPublic)
	appContentH.Register(appPublic)
	appDiyH.Register(appPublic)
	appSeckillH.Register(appPublic)
	appComboH.RegisterPublic(appPublic)
	appPresellH.RegisterPublic(appPublic)
	appReserveH.RegisterPublic(appPublic)
	appBroadcastH.Register(appPublic)
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

	cb := r.Group("/api/callback/v1")
	cb.GET("/ping", func(c *gin.Context) {
		response.OK(c, gin.H{"prefix": "/api/callback/v1"})
	})
	appcallback.NewPayHandler(tradeSvc).Register(cb)

	r.GET("/swagger/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":  "api-app swagger shell",
			"prefixes": []string{"/api/app/v1", "/api/callback/v1"},
		})
	})

	log.Printf("qixi_mergers_api_app listening on %s", cfg.Server.Addr)
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
