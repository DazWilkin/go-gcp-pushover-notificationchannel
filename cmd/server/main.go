package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	pushover_notificationchannel "github.com/DazWilkin/go-gcp-pushover-notificationchannel"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace string = "pushover"
	subsystem string = "notificationchannel"
)

var (
	// BuildTime is the time that this binary was built represented as a UNIX epoch
	BuildTime string
	// GitCommit is the git commit value and is expected to be set during build
	GitCommit string
	// GoVersion is the Golang runtime version
	GoVersion = runtime.Version()
	// OSVersion is the OS version (uname --kernel-release) and is expected to be set during build
	OSVersion string
	// StartTime is the start time of the exporter represented as a UNIX epoch
	StartTime = time.Now().Unix()
)
var (
	counterBuildTime = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "build_info",
			Namespace: namespace,
			Subsystem: subsystem,
			Help:      "A metric with a constant '1' value labels by build|start time, git commit, OS and Go versions",
		}, []string{
			"build_time",
			"git_commit",
			"os_version",
			"go_version",
			"start_time",
		},
	)
)

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("ok")); err != nil {
		log.Printf("unable to write response: %v", err)
	}
}
func main() {
	if GitCommit == "" {
		log.Println("[main] GitCommit value unchanged: expected to be set during build")
	}
	if OSVersion == "" {
		log.Println("[main] OSVersion value unchanged: expected to be set during build")
	}

	counterBuildTime.With(prometheus.Labels{
		"build_time": BuildTime,
		"git_commit": GitCommit,
		"os_version": OSVersion,
		"go_version": GoVersion,
		"start_time": strconv.FormatInt(StartTime, 10),
	}).Inc()

	// Environment variable(s)
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("Unable to obtain `PORT` from the environment; setting as default (8080)")
		port = "8080"
	}

	http.HandleFunc("/", pushover_notificationchannel.Webhook)
	http.HandleFunc("/healthz", healthz)

	// Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening [0.0.0.0:%s]", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), nil))
}
