package metrics

import (
	"github.com/PDOK/betterstack-exporter/internal/betterstack"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
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
	betterStackMonitorStatus *prometheus.GaugeVec
}

func NewUpdater(client betterstack.Client) *Updater {
	betterStackMonitorStatus := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "betterstack_monitor_status",
			Help: "The current status of the check (0: down, 1: maintenance, 2: up, 3: paused, 4: pending, 5: validating)",
		},
		[]string{
			"id",
			"pronouncable_name",
			"url",
		},
	)
	return &Updater{
		client:                   client,
		betterStackMonitorStatus: betterStackMonitorStatus,
	}
}

func (u *Updater) UpdatePromMetrics() error {
	monitors, err := u.client.ListMonitors()
	if err != nil {
		return err
	}
	for _, monitor := range monitors {
		labels := map[string]string{
			"id":                monitor.ID,
			"pronouncable_name": monitor.PronounceableName,
			"url":               monitor.URL,
		}
		u.betterStackMonitorStatus.With(labels).Set(statusCodes[monitor.Status])
	}
	return nil
}
