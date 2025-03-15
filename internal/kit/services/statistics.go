package services

import (
	"github.com/Minh2009/pv_soa/pkgs/log"
	"github.com/uptrace/bun"
)

type StatisticsSvc struct {
	db     *bun.DB
	logger *log.MultiLogger
}
