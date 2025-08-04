package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type AnalyticsService struct {
	db *sql.DB
}

type PageView struct {
	SessionID      string `json:"session_id"`
	UserAgent      string `json:"user_agent"`
	IPAddress      string `json:"ip_address"`
	PageURL        string `json:"page_url"`
	PageTitle      string `json:"page_title"`
	Referrer       string `json:"referrer"`
	PageLoadTime   int    `json:"page_load_time"`
	ViewportWidth  int    `json:"viewport_width"`
	ViewportHeight int    `json:"viewport_height"`
}

type Click struct {
	SessionID    string `json:"session_id"`
	PageURL      string `json:"page_url"`
	ElementTag   string `json:"element_tag"`
	ElementID    string `json:"element_id"`
	ElementClass string `json:"element_class"`
	ElementText  string `json:"element_text"`
	ClickX       int    `json:"click_x"`
	ClickY       int    `json:"click_y"`
}

type ScrollDepth struct {
	SessionID           string `json:"session_id"`
	PageURL             string `json:"page_url"`
	MaxScrollPercentage int    `json:"max_scroll_percentage"`
	TotalPageHeight     int    `json:"total_page_height"`
	ViewportHeight      int    `json:"viewport_height"`
}

type PageTime struct {
	SessionID      string `json:"session_id"`
	PageURL        string `json:"page_url"`
	TimeOnPage     int    `json:"time_on_page"`
	IsActiveTime   int    `json:"is_active_time"`
}

type SessionTime struct {
	SessionID              string `json:"session_id"`
	StartTime              string `json:"start_time"`
	EndTime                string `json:"end_time"`
	TotalSessionDuration   int    `json:"total_session_duration"`
	PagesVisited           int    `json:"pages_visited"`
	TotalClicks            int    `json:"total_clicks"`
	DeviceType             string `json:"device_type"`
	Browser                string `json:"browser"`
	OperatingSystem        string `json:"operating_system"`
}

func NewAnalyticsService() *AnalyticsService {
	// ClickHouse connection string
	clickhouseHost := getEnv("CLICKHOUSE_HOST", "localhost")
	clickhousePort := getEnv("CLICKHOUSE_PORT", "9000")
	clickhouseUser := getEnv("CLICKHOUSE_USER", "analytics_user")
	clickhousePassword := getEnv("CLICKHOUSE_PASSWORD", "password")
	clickhouseDB := getEnv("CLICKHOUSE_DB", "analytics")

	var dsn string
	if clickhousePassword != "" {
		dsn = fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
			clickhouseUser, clickhousePassword, clickhouseHost, clickhousePort, clickhouseDB)
	} else {
		dsn = fmt.Sprintf("clickhouse://%s@%s:%s/%s",
			clickhouseUser, clickhouseHost, clickhousePort, clickhouseDB)
	}

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		log.Fatal("Failed to connect to ClickHouse:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping ClickHouse:", err)
	}

	log.Println("Connected to ClickHouse successfully")

	return &AnalyticsService{db: db}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (as *AnalyticsService) TrackPageView(w http.ResponseWriter, r *http.Request) {
	var pageView PageView
	if err := json.NewDecoder(r.Body).Decode(&pageView); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Get IP address from request
	pageView.IPAddress = getClientIP(r)

	query := `INSERT INTO analytics.page_views 
		(session_id, user_agent, ip_address, page_url, page_title, referrer, page_load_time, viewport_width, viewport_height) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := as.db.Exec(query, pageView.SessionID, pageView.UserAgent, pageView.IPAddress,
		pageView.PageURL, pageView.PageTitle, pageView.Referrer, pageView.PageLoadTime,
		pageView.ViewportWidth, pageView.ViewportHeight)

	if err != nil {
		log.Printf("Error inserting page view: %v", err)
		http.Error(w, "Failed to track page view", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (as *AnalyticsService) TrackClick(w http.ResponseWriter, r *http.Request) {
	var click Click
	if err := json.NewDecoder(r.Body).Decode(&click); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO analytics.clicks 
		(session_id, page_url, element_tag, element_id, element_class, element_text, click_x, click_y) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := as.db.Exec(query, click.SessionID, click.PageURL, click.ElementTag,
		click.ElementID, click.ElementClass, click.ElementText, click.ClickX, click.ClickY)

	if err != nil {
		log.Printf("Error inserting click: %v", err)
		http.Error(w, "Failed to track click", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (as *AnalyticsService) TrackScrollDepth(w http.ResponseWriter, r *http.Request) {
	var scrollDepth ScrollDepth
	if err := json.NewDecoder(r.Body).Decode(&scrollDepth); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO analytics.scroll_depth 
		(session_id, page_url, max_scroll_percentage, total_page_height, viewport_height) 
		VALUES (?, ?, ?, ?, ?)`

	_, err := as.db.Exec(query, scrollDepth.SessionID, scrollDepth.PageURL,
		scrollDepth.MaxScrollPercentage, scrollDepth.TotalPageHeight, scrollDepth.ViewportHeight)

	if err != nil {
		log.Printf("Error inserting scroll depth: %v", err)
		http.Error(w, "Failed to track scroll depth", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (as *AnalyticsService) TrackPageTime(w http.ResponseWriter, r *http.Request) {
	var pageTime PageTime
	if err := json.NewDecoder(r.Body).Decode(&pageTime); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO analytics.page_time 
		(session_id, page_url, time_on_page, is_active_time) 
		VALUES (?, ?, ?, ?)`

	_, err := as.db.Exec(query, pageTime.SessionID, pageTime.PageURL,
		pageTime.TimeOnPage, pageTime.IsActiveTime)

	if err != nil {
		log.Printf("Error inserting page time: %v", err)
		http.Error(w, "Failed to track page time", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (as *AnalyticsService) TrackSessionTime(w http.ResponseWriter, r *http.Request) {
	var sessionTime SessionTime
	if err := json.NewDecoder(r.Body).Decode(&sessionTime); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO analytics.session_time 
		(session_id, start_time, end_time, total_session_duration, pages_visited, total_clicks, device_type, browser, operating_system) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	startTime, _ := time.Parse(time.RFC3339, sessionTime.StartTime)
	endTime, _ := time.Parse(time.RFC3339, sessionTime.EndTime)

	_, err := as.db.Exec(query, sessionTime.SessionID, startTime, endTime,
		sessionTime.TotalSessionDuration, sessionTime.PagesVisited, sessionTime.TotalClicks,
		sessionTime.DeviceType, sessionTime.Browser, sessionTime.OperatingSystem)

	if err != nil {
		log.Printf("Error inserting session time: %v", err)
		http.Error(w, "Failed to track session time", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (as *AnalyticsService) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	// Sample endpoint to retrieve analytics data
	query := `SELECT 
		page_url, 
		count() as views,
		avg(page_load_time) as avg_load_time
		FROM analytics.page_views 
		WHERE timestamp >= now() - INTERVAL 24 HOUR 
		GROUP BY page_url 
		ORDER BY views DESC`

	rows, err := as.db.Query(query)
	if err != nil {
		log.Printf("Error querying analytics: %v", err)
		http.Error(w, "Failed to get analytics", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var pageURL string
		var views int
		var avgLoadTime float64

		err := rows.Scan(&pageURL, &views, &avgLoadTime)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}

		results = append(results, map[string]interface{}{
			"page_url":       pageURL,
			"views":          views,
			"avg_load_time":  avgLoadTime,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func getClientIP(r *http.Request) string {
	// Try to get IP from X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return forwarded
	}

	// Try to get IP from X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}

func main() {
	analyticsService := NewAnalyticsService()
	defer analyticsService.db.Close()

	router := mux.NewRouter()

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	// Analytics endpoints
	router.HandleFunc("/api/analytics/pageview", analyticsService.TrackPageView).Methods("POST")
	router.HandleFunc("/api/analytics/click", analyticsService.TrackClick).Methods("POST")
	router.HandleFunc("/api/analytics/scroll", analyticsService.TrackScrollDepth).Methods("POST")
	router.HandleFunc("/api/analytics/pagetime", analyticsService.TrackPageTime).Methods("POST")
	router.HandleFunc("/api/analytics/sessiontime", analyticsService.TrackSessionTime).Methods("POST")
	router.HandleFunc("/api/analytics/data", analyticsService.GetAnalytics).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	}).Methods("GET")

	handler := c.Handler(router)

	port := getEnv("PORT", "8080")
	log.Printf("Analytics service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
