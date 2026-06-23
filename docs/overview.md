# Grafana Overview

## What is Grafana?

Grafana is an open-source analytics and interactive visualization platform. It connects to a wide range of data sources — including time-series databases, relational databases, cloud monitoring services, and log aggregation systems — and lets you explore, query, and visualize data through customizable dashboards. Originally focused on metrics from tools like Graphite and InfluxDB, Grafana has grown into a full observability platform that handles metrics, logs, and traces side by side.

Grafana is used by thousands of companies worldwide to monitor everything from Kubernetes clusters and microservices to business KPIs and IoT sensor data.

---

## Key Features

| Feature | Description |
|---|---|
| **Unified dashboards** | Combine panels from multiple data sources into a single, cohesive view. |
| **Rich visualization library** | Time-series graphs, heatmaps, bar charts, stat panels, geo maps, and more — all configurable without writing front-end code. |
| **Template variables** | Drive dynamic, reusable dashboards with drop-down variables that filter queries at runtime. |
| **Alerting** | Define alert rules visually and route notifications to Slack, PagerDuty, OpsGenie, email, and dozens of other channels. |
| **Explore mode** | Run ad-hoc queries against any data source without saving a dashboard — ideal for incident investigation. |
| **Plugin ecosystem** | Extend Grafana with community and enterprise plugins for additional data sources, panels, and apps. |
| **Access control** | Fine-grained role-based access control (RBAC) for teams, folders, and individual dashboards. |
| **Provisioning** | Manage dashboards, data sources, and alert rules as code using YAML files or the HTTP API. |

---

## Architecture

Grafana is a single Go binary (`grafana-server`) that serves a React-based web application. At runtime it:

1. Reads configuration from `conf/grafana.ini` (or environment variables).
2. Connects to a backend database (SQLite by default; MySQL or PostgreSQL in production) to store dashboards, users, and settings.
3. Proxies data-source requests — the browser sends a query to Grafana, which forwards it to the real data source (e.g., Prometheus) and streams the response back.

This architecture means **Grafana never stores your metrics** — it only stores metadata (dashboard JSON, user accounts, alert rules). Your data stays in your existing storage systems.

---

## Basic Usage Example

### 1 — Run Grafana with Docker

```bash
docker run -d \
  --name grafana \
  -p 3000:3000 \
  grafana/grafana:latest
```

Open `http://localhost:3000` and log in with `admin` / `admin`.

### 2 — Add a Prometheus data source

```bash
# Run a local Prometheus instance (optional, for testing)
docker run -d --name prometheus -p 9090:9090 prom/prometheus
```

In the Grafana UI:

1. Go to **Configuration (⚙️) → Data Sources → Add data source**.
2. Choose **Prometheus**.
3. Set the URL to `http://prometheus:9090` (or `http://localhost:9090` if not using Docker networking).
4. Click **Save & Test**.

### 3 — Import a community dashboard

1. Go to **+ → Import**.
2. Enter a dashboard ID from [grafana.com/dashboards](https://grafana.com/grafana/dashboards) (e.g., `1860` for the Node Exporter Full dashboard).
3. Select your Prometheus data source and click **Import**.

You now have a fully populated dashboard showing host metrics!

---

## Further Reading

- **Official documentation:** <https://grafana.com/docs/grafana/latest/>
- **Community dashboards:** <https://grafana.com/grafana/dashboards>
- **Plugin catalogue:** <https://grafana.com/grafana/plugins>
- **Community forum:** <https://community.grafana.com/>
- **Contributing guide:** [CONTRIBUTING.md](../CONTRIBUTING.md)
