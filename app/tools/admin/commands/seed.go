package commands

import (
	"context"
	"github.com/Mahamadou828/tgs_with_golang/business/data/schema"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"go.uber.org/zap"
	"time"
)

func Seed(cfg database.Config, version string, log *zap.SugaredLogger) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := schema.Seed(ctx, db, version); err != nil {
		return err
	}
	log.Infow("database seeding completed")
	return nil
}
