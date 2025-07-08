# DevSync

**DevSync** is a distributed developer analytics and incident alerting platform.

It collects activity from developer tools (like GitHub, CI/CD, deployment pipelines), processes events, tracks engineering activity trends, and alerts teams when risky patterns or anomalies occur — such as frequent deploys before an incident or sudden inactivity in a critical repo.

---

## 🔧 Tech Stack

- **Go** – High-performance event ingestion & processing
- **NestJS (TypeScript)** – Admin APIs
- **PostgreSQL** – Persistent event & alert storage
- **Kafka** – Event pipeline between services
- **Redis** – Cache hot data & rate-limit alerts
- **Docker Compose** – Local infrastructure
- **Prometheus + Grafana** – Observability (planned)

---

## 📁 Project Structure

devsync/
│
├── infra/ # Infrastructure config (Kafka, Redis, Postgres, etc.)
│ ├── docker-compose.yml
│ └── .env
│
├── services/
│ ├── event-ingestor/ # Go: Receives webhooks, pushes to Kafka
│ ├── event-processor/ # Go: Parses, filters, stores events
│ └── admin-api/ # NestJS: APIs for dashboard, configs, alerts
│
├── README.md
└── .gitignore


---

## 🚀 Getting Started

1. **Clone the repo**

```bash
git clone https://github.com/thoraf20/devsync.git
cd devsync
```

2. **Start infrastructure**

cd infra
docker-compose up -d

3. Create .env file

Set common environment variables in infra/.env. Examples:

POSTGRES_USER=devsync
POSTGRES_PASSWORD=devsyncpass
POSTGRES_DB=devsyncdb
KAFKA_BROKER=kafka:9092
REDIS_URL=redis://redis:6379

4. Run each service

Each service runs independently. From its folder:

# example: services/event-ingestor
go run main.go


🛣️ Roadmap
 GitHub / GitLab webhook handling

 CI/CD event parsing

 Kafka consumer group for analytics

 Admin API for viewing alerts, activity heatmaps

 Slack/email/webhook notifications

 Prometheus metrics + Grafana dashboards

 CLI tool for DevSync debugging

 👤 Author
Toheeb Rauf
Backend Engineer — Go, TypeScript, Distributed Systems
https://www.github.com/thoraf20 | https://www.linkedin.com/in/toheeb-rauf-678534102/


