package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crmlive/pte-live-ecrm/job/internal/domain/assist"
	"github.com/crmlive/pte-live-ecrm/job/internal/domain/cart"
	"github.com/crmlive/pte-live-ecrm/job/internal/domain/combination"
	"github.com/crmlive/pte-live-ecrm/job/internal/domain/presell"
	"github.com/crmlive/pte-live-ecrm/job/internal/domain/promotion"
	"github.com/crmlive/pte-live-ecrm/job/internal/domain/trade"
	assistpersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/assist"
	cartpersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/cart"
	combinationpersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/combination"
	presellpersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/presell"
	promotionpersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/promotion"
	tradepersist "github.com/crmlive/pte-live-ecrm/job/internal/infra/persist/trade"
	"github.com/crmlive/pte-live-ecrm/job/internal/pkg/config"
	"github.com/crmlive/pte-live-ecrm/job/internal/pkg/db"
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

	cartSvc := cart.NewService(cartpersist.NewStoreAdapter(cartpersist.NewRepo(gdb)))
	promoSvc := promotion.NewService(promotionpersist.NewStoreAdapter(promotionpersist.NewRepo(gdb)))
	tradeSvc := trade.NewService(tradepersist.NewStoreAdapter(tradepersist.NewRepo(gdb)), cartSvc, promoSvc)
	presellSvc := presell.NewService(presellpersist.NewRepo(gdb))
	assistSvc := assist.NewService(assistpersist.NewRepo(gdb))
	comboSvc := combination.NewService(combinationpersist.NewRepo(gdb))
	tradeSvc.SetPresell(presell.NewTradeBridge(presellSvc))
	tradeSvc.SetAssist(assist.NewTradeBridge(assistSvc))
	tradeSvc.SetCombination(combination.NewTradeBridge(comboSvc))

	ttl := time.Duration(cfg.Job.UnpaidTTLMinutes) * time.Minute
	batch := cfg.Job.UnpaidBatch
	tick := time.Duration(cfg.Job.TickSeconds) * time.Second

	log.Printf(
		"pte_live_ecrm_job started (nats=%s); unpaid_ttl=%s batch=%d tick=%s; rollback=presell+assist+combination",
		cfg.NATS.URL, ttl, batch, tick,
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	runClose := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		n, err := tradeSvc.CloseExpiredUnpaid(ctx, ttl, batch)
		if err != nil {
			log.Printf("close unpaid: %v", err)
			return
		}
		if n > 0 {
			log.Printf("closed unpaid group orders: %d", n)
		}
	}

	runClose()

	for {
		select {
		case <-sig:
			log.Printf("pte_live_ecrm_job shutting down")
			return
		case <-ticker.C:
			runClose()
		}
	}
}
