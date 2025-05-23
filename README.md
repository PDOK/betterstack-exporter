# Better Stack exporter

[![Build](https://github.com/PDOK/betterstack-exporter/actions/workflows/build-and-publish-image.yml/badge.svg)](https://github.com/PDOK/betterstack-exporter/actions/workflows/build-and-publish-image.yml)
[![Lint (go)](https://github.com/PDOK/betterstack-exporter/actions/workflows/lint-go.yml/badge.svg)](https://github.com/PDOK/betterstack-exporter/actions/workflows/lint-go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/PDOK/betterstack-exporter)](https://goreportcard.com/report/github.com/PDOK/betterstack-exporter)
[![Coverage (go)](https://github.com/PDOK/betterstack-exporter/wiki/coverage.svg)](https://raw.githack.com/wiki/PDOK/betterstack-exporter/coverage.html)
[![GitHub license](https://img.shields.io/github/license/PDOK/betterstack-exporter)](https://github.com/PDOK/betterstack-exporter/blob/master/LICENSE)
[![Docker Pulls](https://img.shields.io/docker/pulls/pdok/betterstack-exporter.svg)](https://hub.docker.com/r/pdok/betterstack-exporter)

This Prometheus exporter exposes statistics about [Better Stack uptime monitoring](https://betterstack.com/uptime).
This data is collected and exposed as a Prometheus metrics endpoint. The goal is to expose stats regarding the status of each configured uptime monitor.

## Example metrics output

```text
# HELP betterstack_monitor_status The current status of the check (0: down, 1: maintenance, 2: up, 3: paused, 4: pending, 5: validating)
# TYPE betterstack_monitor_status gauge
betterstack_monitor_status{id="3269909",pronouncable_name="UptimeOperatorBetterStackTestCheck - Updated",url="https://service.pdok.nl/cbs/landuse/wfs/v1_0?request=GetCapabilities&service=WFS"} 2
betterstack_monitor_status{id="3269912",pronouncable_name="UptimeOperatorBetterStackTestCheck_2_Updated",url="https://service.pdok.nl/cbs/landuse/wms/v1_0?request=GetCapabilities&service=WMS"} 2
# ...
```

## Build

```shell
docker build .
```

## Run

```text
USAGE:
   betterstack-exporter [global options] command [command options] 

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --api-token value        The API token to authenticate with Better Stack. [$API_TOKEN]
   --bind-address value     The TCP network address addr that is listened on. (default: ":8080") [$BIND_ADDRESS]
   --scrape-interval value  The interval in seconds between scraping of monitor statuses. (default: 60) [$SCRAPE_INTERVAL]
   --page-size value        The number of monitors to request per page (max 250). (default: 50) [$PAGE_SIZE]
   --help, -h               show help
```

### Linting

Install [golangci-lint](https://golangci-lint.run/usage/install/) and run `golangci-lint run`
from the root.
