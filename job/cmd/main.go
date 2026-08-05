package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crmlive/pte-live-ecrm/job/internal/business/ordertimeout"
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

	dsn, err := cfg.DSNFor(config.DatabaseBusiness)
	if err != nil {
		log.Fatalf("business database: %v", err)
	}
	gdb, err := db.OpenMySQL(dsn, cfg.Server.Mode == "debug")
	if err != nil {
		log.Fatalf("business mysql: %v", err)
	}

	ttl := time.Duration(cfg.Job.UnpaidTTLMinutes) * time.Minute
	batch := cfg.Job.UnpaidBatch
	tick := time.Duration(cfg.Job.TickSeconds) * time.Second

	log.Printf(
		"crm_live_job started (business database); unpaid_ttl=%s batch=%d tick=%s; rollback=presell+assist+combination+reservation+coupon",
		ttl, batch, tick,
	)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	ticker := time.NewTicker(tick)
	defer ticker.Stop()

	runClose := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		n, err := ordertimeout.ExpireUnpaid(ctx, gdb, time.Now(), ttl, batch)
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
			log.Printf("crm_live_job shutting down")
			return
		case <-ticker.C:
			runClose()
		}
	}
}
