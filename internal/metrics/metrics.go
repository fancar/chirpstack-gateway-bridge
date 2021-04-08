package metrics

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/brocaar/chirpstack-gateway-bridge/internal/config"
)

// Setup configures the metrics package.
func Setup(conf config.Config) error {
	if !conf.Metrics.Prometheus.EndpointEnabled || !conf.Metrics.Profiler.EndpointEnabled {
		return nil
	}

	r := mux.NewRouter()
	log.WithFields(log.Fields{
		"bind": conf.Metrics.Bind,
	}).Info("metrics: starting metrics http server ...")

	if conf.Metrics.Prometheus.EndpointEnabled {
		log.WithFields(log.Fields{
			"endpoint": conf.Metrics.Bind + "/metrics",
		}).Info("metrics: starting prometheus metrics server ...")
		r.Path("/metrics").Handler(promhttp.Handler())
	}

	if conf.Metrics.Profiler.EndpointEnabled {
		log.WithFields(log.Fields{
			"endpoint": conf.Metrics.Bind + "/debug/pprof/",
		}).Info("metrics: starting pprof: profiling data collector ...")
		r.PathPrefix("/debug/pprof/").HandlerFunc(pprof.Index)
	}

	server := http.Server{
		// Handler: promhttp.Handler(),
		Handler: r,
		Addr:    conf.Metrics.Bind,
	}

	go func() {
		err := server.ListenAndServe()
		log.WithError(err).Error("metrics: prometheus metrics server error")
	}()

	return nil
}
