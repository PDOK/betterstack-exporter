package metrics

import (
	"log"

	"github.com/PDOK/betterstack-exporter/internal/betterstack"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	statusCodes = map[string]float64{
		"down":        0,
		"maintenance": 1,
		"up":          2,
		"paused":      3,
		"pending":     4,
		"validating":  5,
	}
)

type Updater struct {
	client                   betterstack.Client
	BetterStackMonitorStatus *prometheus.GaugeVec
}

func NewUpdater(client betterstack.Client) *Updater {
	betterStackMonitorStatus := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "betterstack",
			Name:      "monitor_status",
			Help:      "The current status of the check (0: down, 1: maintenance, 2: up, 3: paused, 4: pending, 5: validating)",
		},
		[]string{
			"id",
			"pronounceable_name",
			"url",
		},
	)
	return &Updater{
		client:                   client,
		BetterStackMonitorStatus: betterStackMonitorStatus,
	}
}

func (u *Updater) Describe(ch chan<- *prometheus.Desc) {
	u.BetterStackMonitorStatus.Describe(ch)
}

func (u *Updater) Collect(ch chan<- prometheus.Metric) {
	u.UpdatePromMetrics()
	u.BetterStackMonitorStatus.Collect(ch)
}

func (u *Updater) UpdatePromMetrics() {
	log.Println("start updating uptime metrics")
	monitors, err := u.client.ListMonitors()
	if err != nil {
		log.Fatal(err)
	}
	for _, monitor := range monitors {
		labels := map[string]string{
			"id":                 monitor.ID,
			"pronounceable_name": monitor.PronounceableName,
			"url":                monitor.URL,
		}
		u.BetterStackMonitorStatus.With(labels).Set(statusCodes[monitor.Status])
	}
	log.Println("finished updating uptime metrics")
}
