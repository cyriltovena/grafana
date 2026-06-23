![Grafana](docs/logo-horizontal.png)

The open-source platform for monitoring and observability.

[![License](https://img.shields.io/github/license/grafana/grafana)](LICENSE)
[![Circle CI](https://img.shields.io/circleci/build/gh/grafana/grafana)](https://circleci.com/gh/grafana/grafana)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafana/grafana)](https://goreportcard.com/report/github.com/grafana/grafana)

Grafana allows you to query, visualize, alert on and understand your metrics no matter where they are stored. Create, explore, and share dashboards with your team and foster a data driven culture:

- **Visualize:** Fast and flexible client side graphs with a multitude of options. Panel plugins offer many different ways to visualize metrics and logs.
- **Dynamic Dashboards:** Create dynamic & reusable dashboards with template variables that appear as dropdowns at the top of the dashboard.
- **Explore Metrics:** Explore your data through ad-hoc queries and dynamic drilldown. Split view and compare different time ranges, queries and data sources side by side.
- **Explore Logs:** Experience the magic of switching from metrics to logs with preserved label filters. Quickly search through all your logs or streaming them live.
- **Alerting:** Visually define alert rules for your most important metrics. Grafana will continuously evaluate and send notifications to systems like Slack, PagerDuty, VictorOps, OpsGenie.
- **Mixed Data Sources:** Mix different data sources in the same graph! You can specify a data source on a per-query basis. This works for even custom datasources.

## Get started

- [Get Grafana](https://grafana.com/get)
- [Installation guides](http://docs.grafana.org/installation/)

Unsure if Grafana is for you? Watch Grafana in action on [play.grafana.org](https://play.grafana.org/)!

## Documentation

The Grafana documentation is available at [grafana.com/docs](https://grafana.com/docs/).

## Contributing

If you're interested in contributing to the Grafana project:

- Start by reading the [Contributing guide](/CONTRIBUTING.md).
- Learn how to set up your local environment, in our [Developer guide](/contribute/developer-guide.md).
- Explore our [beginner-friendly issues](https://github.com/grafana/grafana/issues?q=is%3Aopen+is%3Aissue+label%3A%22beginner+friendly%22).

## Get involved

- Follow [@grafana on Twitter](https://twitter.com/grafana/).
- Read and subscribe to the [Grafana blog](https://grafana.com/blog/).
- If you have a specific question, check out our [discussion forums](https://community.grafana.com/).
- For general discussions, join us on the [official Slack](http://slack.raintank.io/) team.

## License

Grafana is distributed under the [Apache 2.0 License](https://github.com/grafana/grafana/blob/master/LICENSE).

---

## About Grafana

Grafana is the leading open-source platform for monitoring, observability, and data visualization. It enables teams to query, visualize, and alert on metrics, logs, and traces from virtually any data source — all in one place.

### What is Grafana?

Grafana connects to dozens of data sources (Prometheus, Loki, Elasticsearch, InfluxDB, PostgreSQL, and many more) and lets you build beautiful, interactive dashboards without writing any front-end code. It is used by organizations of all sizes to gain insight into their infrastructure, applications, and business metrics.

### Quick Start

1. **Install Grafana** — download the latest release from [grafana.com/get](https://grafana.com/get) or run it with Docker:
   ```bash
   docker run -d -p 3000:3000 grafana/grafana
   ```
2. **Open the UI** — navigate to `http://localhost:3000` in your browser (default credentials: `admin` / `admin`).
3. **Add a data source** — go to **Configuration → Data Sources** and connect Grafana to your metrics or log backend.
4. **Build a dashboard** — click **+ → Dashboard** and start adding panels to visualize your data.

For a full walkthrough see the [Getting Started guide](https://grafana.com/docs/grafana/latest/getting-started/).
