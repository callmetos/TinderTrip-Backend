package service

import (
	"context"
	"net/http"
	"time"

	"TinderTrip-Backend/internal/utils"
	"TinderTrip-Backend/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
)

var (
	// HTTP metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status_code"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// Database metrics
	dbConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_active",
			Help: "Number of active database connections",
		},
	)

	dbConnectionsIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_connections_idle",
			Help: "Number of idle database connections",
		},
	)

	// Redis metrics
	redisConnectionsActive = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "redis_connections_active",
			Help: "Number of active Redis connections",
		},
	)

	// Application metrics
	applicationUptime = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "application_uptime_seconds",
			Help: "Application uptime in seconds",
		},
	)

	// Custom business metrics
	userRegistrationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registrations_total",
			Help: "Total number of user registrations",
		},
	)

	eventCreationsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "event_creations_total",
			Help: "Total number of events created",
		},
	)

	activeUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_users",
			Help: "Number of currently active users",
		},
	)
)

type MonitoringService struct {
	metricsServer *http.Server
	healthServer  *http.Server
	db            *gorm.DB
	startTime     time.Time
}

func NewMonitoringService(db *gorm.DB) *MonitoringService {
	return &MonitoringService{
		db:        db,
		startTime: time.Now(),
	}
}

func (ms *MonitoringService) Start() error {
	if !config.AppConfig.Monitoring.Enabled {
		utils.Logger().Info("Monitoring is disabled")
		return nil
	}

	// Start metrics server
	go ms.startMetricsServer()

	// Start health check server
	go ms.startHealthServer()

	// Start uptime counter
	go ms.startUptimeCounter()

	// Start database metrics collection
	go ms.startDatabaseMetrics()

	utils.Logger().Info("Monitoring services started")
	return nil
}

func (ms *MonitoringService) startMetricsServer() {
	router := gin.New()
	router.Use(gin.Recovery())

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	ms.metricsServer = &http.Server{
		Addr:    ":" + config.AppConfig.Monitoring.MetricsPort,
		Handler: router,
	}

	utils.Logger().Infof("Metrics server starting on port %s", config.AppConfig.Monitoring.MetricsPort)
	if err := ms.metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Logger().Errorf("Metrics server error: %v", err)
	}
}

func (ms *MonitoringService) startHealthServer() {
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check endpoints
	router.GET("/health", ms.healthCheck)
	router.GET("/ready", ms.readinessCheck)
	router.GET("/live", ms.livenessCheck)

	ms.healthServer = &http.Server{
		Addr:    ":" + config.AppConfig.Monitoring.HealthPort,
		Handler: router,
	}

	utils.Logger().Infof("Health server starting on port %s", config.AppConfig.Monitoring.HealthPort)
	if err := ms.healthServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		utils.Logger().Errorf("Health server error: %v", err)
	}
}

func (ms *MonitoringService) startUptimeCounter() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		uptime := time.Since(ms.startTime).Seconds()
		applicationUptime.Set(uptime)
	}
}

func (ms *MonitoringService) startDatabaseMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ms.collectDatabaseMetrics()
	}
}

func (ms *MonitoringService) collectDatabaseMetrics() {
	if ms.db == nil {
		return
	}

	sqlDB, err := ms.db.DB()
	if err != nil {
		utils.Logger().Errorf("Failed to get database instance: %v", err)
		return
	}

	stats := sqlDB.Stats()
	dbConnectionsActive.Set(float64(stats.OpenConnections))
	dbConnectionsIdle.Set(float64(stats.Idle))
}

func (ms *MonitoringService) healthCheck(c *gin.Context) {
	status := gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    time.Since(ms.startTime).String(),
		"version":   "1.0.1",
	}

	c.JSON(http.StatusOK, status)
}

func (ms *MonitoringService) readinessCheck(c *gin.Context) {
	// Check database connectivity
	if ms.db != nil {
		sqlDB, err := ms.db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database connection error",
			})
			return
		}

		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
				"error":  "database ping failed",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
		},
	})
}

func (ms *MonitoringService) livenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "alive",
		"uptime": time.Since(ms.startTime).String(),
	})
}

func (ms *MonitoringService) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if ms.metricsServer != nil {
		if err = ms.metricsServer.Shutdown(ctx); err != nil {
			utils.Logger().Errorf("Metrics server shutdown error: %v", err)
		}
	}

	if ms.healthServer != nil {
		if err = ms.healthServer.Shutdown(ctx); err != nil {
			utils.Logger().Errorf("Health server shutdown error: %v", err)
		}
	}

	return err
}

// Metrics collection functions
func (ms *MonitoringService) RecordHTTPRequest(method, endpoint, statusCode string, duration float64) {
	httpRequestsTotal.WithLabelValues(method, endpoint, statusCode).Inc()
	httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
}

func (ms *MonitoringService) RecordUserRegistration() {
	userRegistrationsTotal.Inc()
}

func (ms *MonitoringService) RecordEventCreation() {
	eventCreationsTotal.Inc()
}

func (ms *MonitoringService) SetActiveUsers(count float64) {
	activeUsers.Set(count)
}

// GetPrometheusMetrics returns the Prometheus registry for custom metrics
func GetPrometheusRegistry() *prometheus.Registry {
	return prometheus.DefaultRegisterer.(*prometheus.Registry)
}
