package cron

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/lnkk-app/lnkk-api/internal/statistics"
	"github.com/lnkk-app/lnkk-api/pkg/api"
)

const (
	hourlyStats     = "HOURLY_STATS"
	dailyStats      = "DAILY_STATS"
	assetExpiration = "ASSET_EXPIRATION"
)

// HourlyTasks receives hourly cron task requests
func HourlyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()

	// usage metrics
	last := platform.GetJobTimestamp(ctx, hourlyStats)
	platform.CreateSimpleTask(ctx, api.WorkerBaseURL+"/statistics/assets", fmt.Sprintf("%s:%s:%d", statistics.HourlyAssetMetric, "-", last))
	platform.CreateSimpleTask(ctx, api.WorkerBaseURL+"/statistics/redirects", fmt.Sprintf("%s:%s:%d", statistics.HourlyRedirectMetric, "-", last))
	platform.UpdateJob(ctx, hourlyStats, now)

	c.Status(http.StatusOK)
}

// DailyTasks receives daily cron task requests
func DailyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()

	// usage metrics
	last := platform.GetJobTimestamp(ctx, dailyStats)
	platform.CreateSimpleTask(ctx, api.WorkerBaseURL+"/statistics/assets", fmt.Sprintf("%s:%s:%d", statistics.DailyAssetMetric, "-", last))
	platform.CreateSimpleTask(ctx, api.WorkerBaseURL+"/statistics/redirects", fmt.Sprintf("%s:%s:%d", statistics.DailyRedirectMetric, "-", last))
	platform.UpdateJob(ctx, dailyStats, now)

	// asset expiration
	last = platform.GetJobTimestamp(ctx, assetExpiration)
	platform.CreateSimpleTask(ctx, api.WorkerBaseURL+"/expire", fmt.Sprintf("%d", last))
	platform.UpdateJob(ctx, assetExpiration, now)

	c.Status(http.StatusOK)
}
