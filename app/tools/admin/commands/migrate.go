package commands

import (
	"context"
	"fmt"
	"github.com/Mahamadou828/tgs_with_golang/business/data/schema"
	"github.com/Mahamadou828/tgs_with_golang/business/sys/database"
	"go.uber.org/zap"
	"time"
)

func Migrate(cfg database.Config, version string, log *zap.SugaredLogger) error {
	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := schema.Migrate(ctx, db, version); err != nil {
		return fmt.Errorf("can't migrate database schema: %v", err)
	}
	log.Info("migrations complete")
	return nil
}
