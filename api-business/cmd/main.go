package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	apparticle "github.com/crmlive/pte-live-ecrm/api-business/internal/app/article"
	appassist "github.com/crmlive/pte-live-ecrm/api-business/internal/app/assist"
	appchat "github.com/crmlive/pte-live-ecrm/api-business/internal/app/chat"
	appcombination "github.com/crmlive/pte-live-ecrm/api-business/internal/app/combination"
	appcommunity "github.com/crmlive/pte-live-ecrm/api-business/internal/app/community"
	appinvoice "github.com/crmlive/pte-live-ecrm/api-business/internal/app/invoice"
	apppresell "github.com/crmlive/pte-live-ecrm/api-business/internal/app/presell"
	appreservation "github.com/crmlive/pte-live-ecrm/api-business/internal/app/reservation"
	businessaccount "github.com/crmlive/pte-live-ecrm/api-business/internal/business/account"
	businessaddress "github.com/crmlive/pte-live-ecrm/api-business/internal/business/address"
	businessassistorder "github.com/crmlive/pte-live-ecrm/api-business/internal/business/assistorder"
	businessauth "github.com/crmlive/pte-live-ecrm/api-business/internal/business/auth"
	businesscart "github.com/crmlive/pte-live-ecrm/api-business/internal/business/cart"
	businesscatalog "github.com/crmlive/pte-live-ecrm/api-business/internal/business/catalog"
	businesscity "github.com/crmlive/pte-live-ecrm/api-business/internal/business/city"
	businesscombinationorder "github.com/crmlive/pte-live-ecrm/api-business/internal/business/combinationorder"
	businesscomment "github.com/crmlive/pte-live-ecrm/api-business/internal/business/comment"
	businesscontent "github.com/crmlive/pte-live-ecrm/api-business/internal/business/contentview"
	businesscoupon "github.com/crmlive/pte-live-ecrm/api-business/internal/business/coupon"
	businessdistribution "github.com/crmlive/pte-live-ecrm/api-business/internal/business/distribution"
	businessdiyview "github.com/crmlive/pte-live-ecrm/api-business/internal/business/diyview"
	businessfavorite "github.com/crmlive/pte-live-ecrm/api-business/internal/business/favorite"
	businessfeedback "github.com/crmlive/pte-live-ecrm/api-business/internal/business/feedback"
	businessfunding "github.com/crmlive/pte-live-ecrm/api-business/internal/business/funding"
	businesshistory "github.com/crmlive/pte-live-ecrm/api-business/internal/business/history"
	businesslive "github.com/crmlive/pte-live-ecrm/api-business/internal/business/live"
	businessmarketing "github.com/crmlive/pte-live-ecrm/api-business/internal/business/marketing"
	businessmerchantapply "github.com/crmlive/pte-live-ecrm/api-business/internal/business/merchantapply"
	businessnotification "github.com/crmlive/pte-live-ecrm/api-business/internal/business/notification"
	businessopenscreen "github.com/crmlive/pte-live-ecrm/api-business/internal/business/openscreen"
	businessorder "github.com/crmlive/pte-live-ecrm/api-business/internal/business/order"
	businesspoints "github.com/crmlive/pte-live-ecrm/api-business/internal/business/points"
	businesspresellorder "github.com/crmlive/pte-live-ecrm/api-business/internal/business/presellorder"
	businessrefund "github.com/crmlive/pte-live-ecrm/api-business/internal/business/refund"
	businessrefundprocessor "github.com/crmlive/pte-live-ecrm/api-business/internal/business/refundprocessor"
	businessreservation "github.com/crmlive/pte-live-ecrm/api-business/internal/business/reservation"
	businesssearchhistory "github.com/crmlive/pte-live-ecrm/api-business/internal/business/searchhistory"
	businessupload "github.com/crmlive/pte-live-ecrm/api-business/internal/business/upload"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/article"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cart"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/chat"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/cloudconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/combination"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/community"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/invoice"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/presell"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/promotion"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/reservation"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/seckill"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/domain/trade"
	commentmoderationevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/commentmoderation"
	feedbackmoderationevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/feedbackmoderation"
	merchantapplicationevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantapplication"
	merchantdiyevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantdiy"
	merchantimevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantim"
	merchantintegralpolicyevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantintegralpolicy"
	merchantledgerevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantledger"
	merchantstockevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/merchantstock"
	platformcityevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/platformcity"
	platformdiyevent "github.com/crmlive/pte-live-ecrm/api-business/internal/event/platformdiy"
	articlepersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/article"
	assistpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/assist"
	cartpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/cart"
	chatpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/chat"
	cloudconfigpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/cloudconfig"
	combinationpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/combination"
	communitypersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/community"
	invoicepersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/invoice"
	presellpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/presell"
	promotionpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/promotion"
	reservationpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/reservation"
	seckillpersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/seckill"
	tradepersist "github.com/crmlive/pte-live-ecrm/api-business/internal/infra/persist/trade"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/paymentconfig"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/authjwt"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/captchaclient"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/config"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/db"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/middleware"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/response"
	"github.com/crmlive/pte-live-ecrm/api-business/internal/pkg/upload"
	"github.com/gin-gonic/gin"
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
	integralPolicyProjection, err := merchantintegralpolicyevent.StartBusinessProjection(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant integral policy projection subscriber unavailable: %v", err)
	}
	if integralPolicyProjection != nil {
		defer integralPolicyProjection.Close()
	}
	diyProjection, err := merchantdiyevent.StartBusinessProjection(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant DIY projection subscriber unavailable: %v", err)
	}
	if diyProjection != nil {
		defer diyProjection.Close()
	}
	platformCityProjection, err := platformcityevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("platform city projection subscriber unavailable: %v", err)
	}
	if platformCityProjection != nil {
		defer platformCityProjection.Close()
	}

	platformDIYProjection, err := platformdiyevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("platform DIY projection subscriber unavailable: %v", err)
	}
	if platformDIYProjection != nil {
		defer platformDIYProjection.Close()
	}
	merchantApplicationReviewProjection, err := merchantapplicationevent.StartReviewProjection(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("merchant application review projection subscriber unavailable: %v", err)
	}
	if merchantApplicationReviewProjection != nil {
		defer merchantApplicationReviewProjection.Close()
	}
	merchantapplicationevent.StartOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	merchantstockevent.StartOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	merchantledgerevent.StartOutboxDispatcher(context.Background(), gdb, cfg.NATS.URL)
	productCommentCommands, err := commentmoderationevent.StartCommandSubscriber(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("product comment moderation command subscriber unavailable: %v", err)
	}
	if productCommentCommands != nil {
		defer productCommentCommands.Close()
	}
	feedbackCommands, err := feedbackmoderationevent.Start(context.Background(), gdb, cfg.NATS.URL)
	if err != nil {
		log.Printf("feedback moderation command subscriber unavailable: %v", err)
	}
	if feedbackCommands != nil {
		defer feedbackCommands.Close()
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
	businessrefundprocessor.New(gdb, paymentConfigStore, cloudConfigSvc).Start(context.Background())
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
	appH := businessauth.NewHandler(businessauth.NewService(gdb), jwtMgr, captchaClient, cloudConfigSvc)
	appCatH := businesscatalog.NewHandler(gdb)
	appCityH := businesscity.NewHandler(gdb)
	appMerchantApplyH := businessmerchantapply.NewHandler(gdb)
	appUploadH := businessupload.NewHandler(gdb, upload.DatabaseCOS{Resolver: cloudConfigSvc})
	appAddrH := businessaddress.NewHandler(gdb)
	appAccountH := businessaccount.NewHandler(gdb)
	appFundingH := businessfunding.NewHandler(gdb, cloudConfigSvc)
	appCartH := businesscart.NewHandler(gdb)
	appCommentH := businesscomment.NewHandler(gdb)
	appOrderH := businessorder.NewHandler(gdb, paymentConfigStore, cfg.Payment.Sandbox, cloudConfigSvc)
	appOrderCallbackH := businessorder.NewCallbackHandler(gdb, paymentConfigStore, cfg.Payment.Sandbox, cloudConfigSvc)
	appPointsH := businesspoints.NewHandler(gdb)
	appAssistOrderH := businessassistorder.NewHandler(gdb)
	appCombinationOrderH := businesscombinationorder.NewHandler(gdb)
	appPresellOrderH := businesspresellorder.NewHandler(gdb)
	appReservationOrderH := businessreservation.NewHandler(gdb)
	appRefundH := businessrefund.NewHandler(gdb)
	appBusinessCouponH := businesscoupon.NewHandler(gdb)
	appContentH := businesscontent.NewHandler(gdb)
	appDiyH := businessdiyview.NewHandler(gdb)
	appDistributionH := businessdistribution.NewHandler(gdb)
	appFavoriteH := businessfavorite.NewHandler(gdb)
	appFeedbackH := businessfeedback.NewHandler(gdb)
	appHistoryH := businesshistory.NewHandler(gdb)
	appSearchHistoryH := businesssearchhistory.NewHandler(gdb)
	appNotificationH := businessnotification.NewHandler(gdb)
	appOpenScreenH := businessopenscreen.NewHandler(gdb)
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
		response.OK(c, gin.H{"service": "pte_live_ecrm_api_business", "ok": true})
	})

	appPublic := r.Group("/api/app/v1")
	appOptional := r.Group("/api/app/v1")
	appOptional.Use(middleware.JWTOptional(jwtMgr, authjwt.PortalApp))
	appAuthed := r.Group("/api/app/v1")
	appAuthed.Use(middleware.JWTRequired(jwtMgr, authjwt.PortalApp), middleware.CUserSessionRequired(gdb))
	appH.Register(appPublic, appAuthed)
	appH.RegisterPassword(appAuthed)
	appH.RegisterCancellation(appAuthed)
	appH.RegisterCaptchaGateway(r)
	appCatH.Register(appPublic)
	appCityH.Register(appPublic)
	appOpenScreenH.Register(appPublic)
	appCommentH.RegisterPublic(appPublic)
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
	appFundingH.Register(appAuthed)
	appCartH.Register(appAuthed)
	appOrderH.Register(appAuthed)
	appComboH.RegisterAuthed(appAuthed)
	appCombinationOrderH.Register(appAuthed)
	appPresellH.RegisterAuthed(appAuthed)
	appPresellOrderH.Register(appAuthed)
	appReservationOrderH.Register(appAuthed)
	appAssistH.RegisterAuthed(appAuthed)
	appAssistOrderH.Register(appAuthed)
	appPointsH.RegisterAuthed(appAuthed)
	appRefundH.Register(appAuthed)
	appDistributionH.Register(appAuthed)
	appBusinessCouponH.Register(appAuthed)
	appCommunityH.RegisterAuthed(appAuthed)
	appChatH.Register(appAuthed)
	appInvoiceH.Register(appAuthed)
	appLiveH.RegisterAuthed(appAuthed)
	appMerchantApplyH.Register(appAuthed)
	appUploadH.Register(appAuthed)
	appFavoriteH.Register(appAuthed)
	appFeedbackH.Register(appAuthed)
	appHistoryH.Register(appAuthed)
	appSearchHistoryH.Register(appAuthed)
	appNotificationH.Register(appAuthed)
	appCommentH.Register(appAuthed)

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

	log.Printf("pte_live_ecrm_api_business listening on %s", cfg.Server.Addr)
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
