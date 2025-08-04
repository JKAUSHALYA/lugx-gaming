# Analytics Service - Testing and Development

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for development)

### Setup

1. **Start the database service first:**

   ```bash
   # Start PostgreSQL
   cd ../db-service
   .\start.ps1

   # Or manually
   cd ../db-service
   docker network create lugx-network 2>$null
   docker-compose up -d
   ```

2. **Start the analytics services:**

   ```bash
   # Using PowerShell (Windows)
   .\start.ps1

   # Using Makefile (Linux/Mac)
   make all

   # Using Docker Compose directly
   docker-compose up -d
   ```

3. **Verify services are running:**

   ```bash
   # Check ClickHouse
   curl http://localhost:8123/ping

   # Check Analytics Service
   curl http://localhost:8080/health

   # Check PostgreSQL (from db-service)
   docker exec db-service-postgres-1 pg_isready -U postgres
   ```

4. **Access web interfaces:**
   - ClickHouse Web UI: http://localhost:8123/play
   - Analytics API: http://localhost:8080/api/analytics

### Testing the Analytics

1. **Open the frontend application:**

   ```bash
   cd ../front-end
   # Serve the files using a local web server
   python -m http.server 3000
   # or
   npx serve -p 3000
   ```

2. **Browse the website:**

   - Visit http://localhost:3000
   - Navigate between pages (index.html, shop.html, product-details.html, contact.html)
   - Click on various elements
   - Scroll up and down on pages
   - Stay on pages for different durations

3. **Check collected data:**
   Open ClickHouse Web UI (http://localhost:8123/play) and run queries:

   ```sql
   -- Check page views
   SELECT * FROM analytics.page_views ORDER BY timestamp DESC LIMIT 10;

   -- Check clicks
   SELECT * FROM analytics.clicks ORDER BY timestamp DESC LIMIT 10;

   -- Check scroll data
   SELECT * FROM analytics.scroll_depth ORDER BY timestamp DESC LIMIT 10;

   -- Check page time data
   SELECT * FROM analytics.page_time ORDER BY timestamp DESC LIMIT 10;

   -- Check session data
   SELECT * FROM analytics.session_time ORDER BY start_time DESC LIMIT 10;
   ```

### Sample Analytics Queries

```sql
-- Most popular pages (last 24 hours)
SELECT
    page_url,
    count() as views,
    avg(page_load_time) as avg_load_time_ms
FROM analytics.page_views
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY page_url
ORDER BY views DESC;

-- Click analysis by element type
SELECT
    element_tag,
    element_id,
    count() as clicks
FROM analytics.clicks
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY element_tag, element_id
ORDER BY clicks DESC;

-- Average scroll depth by page
SELECT
    page_url,
    avg(max_scroll_percentage) as avg_scroll_depth,
    count(DISTINCT session_id) as unique_sessions
FROM analytics.scroll_depth
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY page_url
ORDER BY avg_scroll_depth DESC;

-- Session analysis
SELECT
    device_type,
    browser,
    avg(total_session_duration) as avg_session_duration_sec,
    avg(pages_visited) as avg_pages_per_session,
    count() as total_sessions
FROM analytics.session_time
WHERE start_time >= now() - INTERVAL 24 HOUR
GROUP BY device_type, browser
ORDER BY total_sessions DESC;

-- Page engagement metrics
SELECT
    pv.page_url,
    count(DISTINCT pv.session_id) as unique_visitors,
    count(pv.id) as total_views,
    avg(pt.time_on_page) as avg_time_on_page_sec,
    avg(pt.is_active_time) as avg_active_time_sec,
    avg(sd.max_scroll_percentage) as avg_scroll_depth
FROM analytics.page_views pv
LEFT JOIN analytics.page_time pt ON pv.session_id = pt.session_id AND pv.page_url = pt.page_url
LEFT JOIN analytics.scroll_depth sd ON pv.session_id = sd.session_id AND pv.page_url = sd.page_url
WHERE pv.timestamp >= now() - INTERVAL 24 HOUR
GROUP BY pv.page_url
ORDER BY unique_visitors DESC;
```

### API Testing

You can also test the API directly:

```bash
# Test page view tracking
curl -X POST http://localhost:8080/api/analytics/pageview \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test_session_123",
    "user_agent": "Mozilla/5.0...",
    "page_url": "http://localhost:3000/index.html",
    "page_title": "Home Page",
    "referrer": "",
    "page_load_time": 1200,
    "viewport_width": 1920,
    "viewport_height": 1080
  }'

# Test click tracking
curl -X POST http://localhost:8080/api/analytics/click \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test_session_123",
    "page_url": "http://localhost:3000/index.html",
    "element_tag": "button",
    "element_id": "cta-button",
    "element_class": "btn btn-primary",
    "element_text": "Buy Now",
    "click_x": 100,
    "click_y": 200
  }'

# Get analytics data
curl http://localhost:8080/api/analytics/data
```

### Troubleshooting

1. **Services not starting:**

   - Check if ports 8080, 8123, 9000 are available
   - Ensure Docker is running
   - Check logs: `docker-compose logs`

2. **Frontend not tracking:**

   - Check browser console for errors
   - Verify analytics.js is loaded
   - Check CORS settings

3. **No data in ClickHouse:**

   - Verify analytics service can connect to ClickHouse
   - Check service logs: `docker logs analytics-service`
   - Verify database initialization: `docker logs clickhouse`

4. **CORS errors:**
   - Analytics service allows all origins by default
   - If modified, ensure frontend domain is allowed

### Development

For local development:

```bash
# Start only ClickHouse
docker run -d --name clickhouse-dev \
  -p 8123:8123 -p 9000:9000 \
  clickhouse/clickhouse-server:23.8

# Run analytics service locally
export CLICKHOUSE_HOST=localhost
export CLICKHOUSE_PORT=9000
export CLICKHOUSE_USER=default
export CLICKHOUSE_PASSWORD=
export CLICKHOUSE_DB=default
go run main.go
```
