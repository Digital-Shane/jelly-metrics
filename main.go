package main

import (
	"github.com/Digital-Shane/jelly-metrics/jellyfin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	jellyfinHost = loadEnvVarWithDefault("JELLYFIN_HOST", "http://localhost:8096")
	port         = loadEnvVarWithDefault("PORT", "8097")

	// Create the metrics that will be updated in the script
	metricMediaCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_media_count",
			Help: "Media count by type",
		},
		[]string{"type"},
	)

	metricClientCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_connected_clients_count",
			Help: "Connected clients by username",
		},
		[]string{"username"},
	)

	metricStreamCount = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "jellyfin_stream_count",
			Help: "Actively playing streams by username",
		},
		[]string{"username"},
	)
)

func init() {
	// Disable default go collectors. This metric collector has a tiny resource footprint
	// that does not require advanced monitoring directly. Removing these collectors will reduce
	// the speed that Prometheus data directory grows without requiring users to filter out
	// base metrics.
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	prometheus.Unregister(collectors.NewGoCollector())

	// Add the metrics we are going to gather from the Jellyfin API
	prometheus.MustRegister(metricMediaCount)
	prometheus.MustRegister(metricClientCount)
	prometheus.MustRegister(metricStreamCount)
}

func main() {
	// Load the Jellyfin API token from the environment
	jellyfinToken, ok := os.LookupEnv("JELLYFIN_TOKEN")
	if !ok {
		log.Fatalln("Failed to load required jellyfin api token. Supply JELLYFIN_TOKEN environment variable.")
	}

	// Create the Jellyfin client and validate the loaded token works
	jClient := jellyfin.NewClient(jellyfinHost, jellyfinToken)
	if err := jClient.ValidateToken(); err != nil {
		log.Fatalln("Provided jellyfin api token is invalid.")
	}

	// Create a ticker to gather metrics every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when we're done

	// Start a go routine for updating metrics
	go func() {
		// Update metric values every 15 seconds
		for range ticker.C {
			updateMetrics(jClient)
			ticker.Reset(15 * time.Second)
		}
	}()

	// Create server to handle metrics endpoint
	http.Handle("/metrics", promhttp.Handler())
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Printf("Starting server on port %s", port)
	log.Fatal(server.ListenAndServe())
}

func updateMetrics(jClient *jellyfin.Client) {
	// Connected Devices
	if counts, err := jClient.GetConnectedDevicesPerUser(); err != nil {
		slog.Error("failed to gather connected devices", slog.Any("error", err))
	} else {
		updateMetric(metricClientCount, "username", counts)
	}

	// Active Streams
	if streamCounts, err := jClient.GetActiveStreamsPerUser(); err != nil {
		slog.Error("failed to gather active streams", slog.Any("error", err))
	} else {
		updateMetric(metricStreamCount, "username", streamCounts)
	}

	// Media Counts
	if mediaCounts, err := jClient.GetMediaByType(); err != nil {
		slog.Error("failed to gather media counts", slog.Any("error", err))
	} else {
		updateMetric(metricMediaCount, "type", mediaCounts)
	}
}

func updateMetric(gauge *prometheus.GaugeVec, labelKey string, counts map[string]int) {
	gauge.Reset()
	for username, count := range counts {
		gauge.With(prometheus.Labels{
			labelKey: username,
		}).Set(float64(count))
	}
}

func loadEnvVarWithDefault(key, defaultValue string) string {
	foundValue, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return foundValue
}
