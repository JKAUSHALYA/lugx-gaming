# Analytics Service

This service provides analytics tracking for the Lugx Gaming website using ClickHouse as the data storage backend.

## Features

- **Page Views**: Tracks page visits with metadata including load time, viewport size, user agent, referrer
- **Click Tracking**: Records user clicks with element details and coordinates
- **Scroll Depth**: Monitors how far users scroll on each page
- **Page Time**: Measures time spent on each page and active time
- **Session Tracking**: Tracks complete user sessions with device and browser information

## Architecture

### Components

1. **ClickHouse Database**: High-performance columnar database for analytics data (separate service)
2. **Analytics Service**: Go-based API service for data collection
3. **Frontend Analytics Library**: JavaScript library for automatic tracking

### Prerequisites

**Important**: This service requires ClickHouse to be running. ClickHouse is now a separate service located in the `../clickhouse/` directory.

1. Start ClickHouse first:

   ```powershell
   cd ..\clickhouse
   .\start.ps1
   ```

2. Then start the Analytics Service:
   ```powershell
   .\start.ps1
   ```

### Database Schema

#### page_views

- `id`: Unique identifier (UUID)
- `session_id`: Session identifier
- `user_agent`: Browser user agent
- `ip_address`: Client IP address
- `page_url`: Current page URL
- `page_title`: Page title
- `referrer`: Referring URL
- `timestamp`: Event timestamp
- `page_load_time`: Page load time in milliseconds
- `viewport_width/height`: Browser viewport dimensions

#### clicks

- `id`: Unique identifier (UUID)
- `session_id`: Session identifier
- `page_url`: Current page URL
- `element_tag`: HTML tag of clicked element
- `element_id`: Element ID attribute
- `element_class`: Element class attribute
- `element_text`: Element text content (truncated)
- `click_x/y`: Click coordinates
- `timestamp`: Event timestamp

#### scroll_depth

- `id`: Unique identifier (UUID)
- `session_id`: Session identifier
- `page_url`: Current page URL
- `max_scroll_percentage`: Maximum scroll percentage reached
- `total_page_height`: Total page height
- `viewport_height`: Browser viewport height
- `timestamp`: Event timestamp

#### page_time

- `id`: Unique identifier (UUID)
- `session_id`: Session identifier
- `page_url`: Current page URL
- `time_on_page`: Total time on page in seconds
- `is_active_time`: Active time (user interaction) in seconds
- `timestamp`: Event timestamp

#### session_time

- `id`: Unique identifier (UUID)
- `session_id`: Session identifier
- `start_time`: Session start timestamp
- `end_time`: Session end timestamp
- `total_session_duration`: Total session duration in seconds
- `pages_visited`: Number of pages visited
- `total_clicks`: Total clicks in session
- `device_type`: Device type (desktop/mobile/tablet)
- `browser`: Browser name
- `operating_system`: Operating system

## API Endpoints

### Analytics Data Collection

- `POST /api/analytics/pageview` - Track page views
- `POST /api/analytics/click` - Track clicks
- `POST /api/analytics/scroll` - Track scroll depth
- `POST /api/analytics/pagetime` - Track page time
- `POST /api/analytics/sessiontime` - Track session data

### Data Retrieval

- `GET /api/analytics/data` - Get analytics summary (last 24 hours)

### Health Check

- `GET /health` - Service health check

## Environment Variables

| Variable              | Default   | Description         |
| --------------------- | --------- | ------------------- |
| `CLICKHOUSE_HOST`     | localhost | ClickHouse host     |
| `CLICKHOUSE_PORT`     | 9000      | ClickHouse port     |
| `CLICKHOUSE_USER`     | default   | ClickHouse username |
| `CLICKHOUSE_PASSWORD` | password  | ClickHouse password |
| `CLICKHOUSE_DB`       | analytics | ClickHouse database |
| `PORT`                | 8080      | Service port        |

## Deployment

### Using Docker Compose

1. Start the analytics service with ClickHouse:

```bash
cd analytics-service
docker-compose up -d
```

2. The service will be available at `http://localhost:8080`
3. ClickHouse will be available at `http://localhost:8123` (HTTP) and `localhost:9000` (TCP)

### Manual Deployment

1. Start ClickHouse:

```bash
docker run -d --name clickhouse-server \
  -p 8123:8123 -p 9000:9000 \
  -e CLICKHOUSE_DB=analytics \
  clickhouse/clickhouse-server:23.8
```

2. Build and run the analytics service:

```bash
go mod tidy
go run main.go
```

## Frontend Integration

The analytics tracking is automatically enabled on all HTML pages through the `analytics.js` script. The script:

1. Generates a unique session ID stored in sessionStorage
2. Automatically tracks page views when pages load
3. Monitors user interactions (clicks, scrolls, focus/blur)
4. Sends data to the analytics service via HTTP POST requests
5. Handles session cleanup on page unload

### Configuration

To customize the analytics endpoint, modify the initialization in your HTML:

```javascript
window.lugxAnalytics = new LugxAnalytics({
  apiUrl: "http://your-analytics-service:8080/api/analytics",
});
```

### Manual Event Tracking

You can also track custom events:

```javascript
window.lugxAnalytics.trackCustomEvent("button_click", {
  button_name: "signup",
  location: "header",
});
```

## Data Analysis

### Sample Queries

Get top pages by views (last 24 hours):

```sql
SELECT
    page_url,
    count() as views,
    avg(page_load_time) as avg_load_time
FROM analytics.page_views
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY page_url
ORDER BY views DESC;
```

Get click heatmap data:

```sql
SELECT
    page_url,
    element_tag,
    count() as clicks
FROM analytics.clicks
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY page_url, element_tag
ORDER BY clicks DESC;
```

Get scroll depth analysis:

```sql
SELECT
    page_url,
    avg(max_scroll_percentage) as avg_scroll,
    count() as sessions
FROM analytics.scroll_depth
WHERE timestamp >= now() - INTERVAL 24 HOUR
GROUP BY page_url
ORDER BY avg_scroll DESC;
```

## Monitoring

### Health Checks

- Service health: `GET http://localhost:8080/health`
- ClickHouse health: `GET http://localhost:8123/ping`

### Logs

The service logs all operations and errors to stdout. Monitor for:

- Connection errors to ClickHouse
- Failed data insertions
- HTTP request errors

## Security Considerations

1. **CORS**: The service allows all origins by default. Configure for production use.
2. **Rate Limiting**: Consider implementing rate limiting for production deployments.
3. **Data Privacy**: Ensure compliance with data protection regulations (GDPR, etc.).
4. **IP Anonymization**: Consider anonymizing IP addresses for privacy.

## Performance

- ClickHouse is optimized for analytical workloads with high write throughput
- Data is partitioned by date for efficient queries
- Tables use appropriate sorting keys for common query patterns
- Consider data retention policies for large datasets

## Troubleshooting

### Common Issues

1. **Service won't start**: Check ClickHouse connection and credentials
2. **No data appearing**: Verify network connectivity and CORS settings
3. **High memory usage**: Monitor ClickHouse memory configuration
4. **Slow queries**: Check table partitioning and indexes

### Debug Mode

Enable detailed logging by checking the ClickHouse logs:

```bash
docker logs clickhouse
```

And service logs:

```bash
docker logs analytics-service
```
