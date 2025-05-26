package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PDOK/betterstack-exporter/internal/betterstack"
	"github.com/PDOK/betterstack-exporter/internal/metrics"
	"github.com/iancoleman/strcase"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

const (
	APITokenFlag    = "api-token"
	BindAddressFlag = "bind-address"
	PageSizeFlag    = "page-size"
)

var (
	cliFlags = []cli.Flag{
		&cli.StringFlag{
			Name:     APITokenFlag,
			Usage:    "The API token to authenticate with Better Stack.",
			EnvVars:  []string{strcase.ToScreamingSnake(APITokenFlag)},
			Required: true,
		},
		&cli.StringFlag{
			Name:    BindAddressFlag,
			Usage:   "The TCP network address addr that is listened on.",
			Value:   ":8080",
			EnvVars: []string{strcase.ToScreamingSnake(BindAddressFlag)},
		},
		&cli.IntFlag{
			Name:    PageSizeFlag,
			Usage:   "The number of monitors to request per page (max 250).",
			Value:   50,
			EnvVars: []string{strcase.ToScreamingSnake(PageSizeFlag)},
		},
	}
)

func main() {
	app := cli.NewApp()
	app.HelpName = "Better Stack Exporter"
	app.Name = "betterstack-exporter"
	app.Usage = "Collects Better Stack uptime statuses and exports as Prometheus metrics"
	app.Flags = cliFlags
	app.Action = func(c *cli.Context) error {
		config := betterstack.Config{
			APIToken: c.String(APITokenFlag),
			PageSize: c.Int(PageSizeFlag),
		}
		client := betterstack.NewClient(config)
		metricsUpdater := metrics.NewUpdater(client)
		prometheus.MustRegister(metricsUpdater)

		bindAddress := c.String(BindAddressFlag)
		http.Handle("/metrics", promhttp.Handler())
		server := &http.Server{
			Addr:              bindAddress,
			ReadHeaderTimeout: 10 * time.Second,
		}
		log.Printf("listening on %s", bindAddress)
		return server.ListenAndServe()
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
