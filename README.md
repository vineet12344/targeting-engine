# 🎯 Targeting Engine - Campaign Delivery Backend

A high-performance Go backend that selects and serves targeted ad campaigns using include/exclude rules for app, country, and OS — optimized with in-memory caching, Prometheus metrics, and production-ready structure.

---

## 📦 Features

- ✅ Rule-based targeting using Include/Exclude logic
- 🚀 Fast in-memory campaign matching with goroutines
- 🧠 Automatic cache refresh every 60 seconds
- 📊 Prometheus metrics & Grafana support
- 🧪 Unit tested matching logic
- 📦 Clean project structure & modular packages

---

## 🔧 Setup Instructions

### 1. Clone & Create `.env`

```bash
git clone https://github.com/your-username/targeting-engine.git
cd targeting-engine
```

Create a `.env` file:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=targeting_engine
PORT=8080
```

### 2. Start PostgreSQL

Ensure your PostgreSQL server is running and a database named `targeting_engine` is created.

### 3. Run the App

```bash
go run ./cmd/server/main.go
```

This will:

- Connect to DB
- Auto-migrate tables
- Load campaigns into in-memory cache
- Start a background goroutine for cache refresh
- Start the server on `localhost:8080`

---

## 🌐 API Endpoints

### 📤 Match Campaigns

```
GET /campaign?app=<app_id>&os=<os>&country=<country>
```

#### ✅ Example

```
GET /campaign?app=com.spotify.music&os=android&country=IN
```

#### 🔁 Response (HTTP 200)

```json
[
  {
    "ID": "cmp001",
    "Status": "ACTIVE",
    "ImageURL": "https://ads.spotifycdn.com/promo1.png",
    "CTA": "Listen Now"
  },
  {
    "ID": "cmp004",
    "Status": "ACTIVE",
    "ImageURL": "https://ads.netflix.com/new-season.png",
    "CTA": "Watch New Series"
  }
]
```

#### ❌ No Match (HTTP 204)

```json
{
  "message": "No matching campaign found"
}
```

---

## 🧠 Targeting Rules Logic

Campaigns match based on one or more targeting rules. Each rule supports:

| Rule Field     | Type   | Behavior                                   |
|----------------|--------|--------------------------------------------|
| IncludeApp     | string | Only allow listed apps (comma-separated)   |
| ExcludeApp     | string | Block listed apps                          |
| IncludeCountry | string | Only allow users from listed countries     |
| ExcludeCountry | string | Block listed countries                     |
| IncludeOS      | string | Only allow listed OS types (e.g. android)  |
| ExcludeOS      | string | Block listed OS types                      |

> ❗ If any `Include` field is filled, the corresponding `Exclude` **must be empty**.

---

## 🧱 Database Schema (Auto-migrated by GORM)

### Campaign Table

| Field     | Type    | Description             |
|-----------|---------|-------------------------|
| ID        | string  | Primary key             |
| Status    | string  | ACTIVE or INACTIVE      |
| ImageURL  | string  | Campaign banner image   |
| CTA       | string  | Call to action text     |

### TargetingRule Table

| Field           | Type   | Description                       |
|-----------------|--------|-----------------------------------|
| ID              | uint   | Primary key                       |
| CampaignID      | string | Foreign key (indexed)             |
| IncludeApp      | string | Comma-separated app IDs           |
| ExcludeApp      | string | Comma-separated app IDs           |
| IncludeCountry  | string | Comma-separated ISO country codes |
| ExcludeCountry  | string | Comma-separated ISO country codes |
| IncludeOS       | string | Comma-separated OS types          |
| ExcludeOS       | string | Comma-separated OS types          |

> ✅ Indexes improve performance for large datasets.

---

## 📊 Metrics (Prometheus Ready)

### Endpoint

```
GET /metrics
```

### Tracked Metrics

| Metric Name                      | Description                          |
|----------------------------------|--------------------------------------|
| `campaign_match_success_total`  | Count of successful matches          |
| `campaign_match_miss_total`     | Count of requests with no match      |
| `campaign_request_duration_seconds` | Histogram of request durations |

🧪 Integrated using Prometheus client with `/metrics` exposed.

---

## 📊 Grafana Integration

- Add Prometheus data source pointing to `http://localhost:8080/metrics`
- Create dashboards to visualize match counts and latency
- Works locally without Docker too

---

## 🧪 Run Test Cases

```bash
go test ./...
```

Includes:

- ✅ Unit tests for matching logic
- ✅ In-memory cache logic
- ✅ Rule validation edge cases

---

## 🐳 Docker (App Only)

```bash
docker build -t targeting-engine .
docker run -p 8080:8080 --env-file .env targeting-engine
```

> Optional: you can skip Docker Compose and run DB locally.

---

## 📁 Project Structure

```
.
├── cmd/server/             # App entrypoint
├── internal/campaign/      # Handlers, models, services
├── middleware/             # Prometheus metrics middleware
├── pkg/db/                 # DB connection + AutoMigrate
├── go.mod / go.sum         # Dependencies
├── .env                    # Environment config (not committed)
```

---

## 💬 Example Request for Multiple Matches

```http
GET /campaign?app=com.spotify.music&os=android&country=US
```

→ Will match multiple campaigns like Spotify + Netflix (if rules allow).

---

## 👨‍💻 Author

Made with ❤️ and hardwork by [Vineet Salve](https://github.com/vineet12344)
