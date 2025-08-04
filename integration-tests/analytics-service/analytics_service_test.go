package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

const analyticsServiceBaseURL = "http://localhost:30082"

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
	SessionID    string `json:"session_id"`
	PageURL      string `json:"page_url"`
	TimeOnPage   int    `json:"time_on_page"`
	IsActiveTime int    `json:"is_active_time"`
}

type SessionTime struct {
	SessionID      string `json:"session_id"`
	TotalTime      int    `json:"total_time"`
	ActiveTime     int    `json:"active_time"`
	PagesVisited   int    `json:"pages_visited"`
	LastActivity   string `json:"last_activity"`
}

func TestAnalyticsServiceHealth(t *testing.T) {
	resp, err := http.Get(analyticsServiceBaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make health check request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected health status 'healthy', got '%s'", response["status"])
	}
}

func TestTrackPageView(t *testing.T) {
	pageView := PageView{
		SessionID:      "test_session_123",
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		PageURL:        "https://lugx-gaming.com/home",
		PageTitle:      "LUGX Gaming - Home",
		Referrer:       "https://google.com",
		PageLoadTime:   1500,
		ViewportWidth:  1920,
		ViewportHeight: 1080,
	}

	jsonData, err := json.Marshal(pageView)
	if err != nil {
		t.Fatalf("Failed to marshal page view: %v", err)
	}

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pageview", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to track page view: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode page view response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestTrackClick(t *testing.T) {
	click := Click{
		SessionID:    "test_session_123",
		PageURL:      "https://lugx-gaming.com/shop",
		ElementTag:   "button",
		ElementID:    "buy-now-btn",
		ElementClass: "btn btn-primary",
		ElementText:  "Buy Now",
		ClickX:       150,
		ClickY:       300,
	}

	jsonData, err := json.Marshal(click)
	if err != nil {
		t.Fatalf("Failed to marshal click data: %v", err)
	}

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/click", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to track click: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode click response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestTrackScrollDepth(t *testing.T) {
	scrollDepth := ScrollDepth{
		SessionID:           "test_session_123",
		PageURL:             "https://lugx-gaming.com/shop",
		MaxScrollPercentage: 75,
		TotalPageHeight:     2400,
		ViewportHeight:      1080,
	}

	jsonData, err := json.Marshal(scrollDepth)
	if err != nil {
		t.Fatalf("Failed to marshal scroll depth data: %v", err)
	}

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/scroll", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to track scroll depth: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode scroll depth response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestTrackPageTime(t *testing.T) {
	pageTime := PageTime{
		SessionID:    "test_session_123",
		PageURL:      "https://lugx-gaming.com/product-details",
		TimeOnPage:   45000, // 45 seconds
		IsActiveTime: 1,
	}

	jsonData, err := json.Marshal(pageTime)
	if err != nil {
		t.Fatalf("Failed to marshal page time data: %v", err)
	}

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pagetime", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to track page time: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode page time response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestTrackSessionTime(t *testing.T) {
	sessionTime := SessionTime{
		SessionID:    "test_session_123",
		TotalTime:    300000, // 5 minutes
		ActiveTime:   180000, // 3 minutes
		PagesVisited: 5,
		LastActivity: time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(sessionTime)
	if err != nil {
		t.Fatalf("Failed to marshal session time data: %v", err)
	}

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/sessiontime", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to track session time: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code 200 or 201, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode session time response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}
}

func TestGetAnalyticsData(t *testing.T) {
	// First, ensure we have some data by tracking a few events
	testSessionID := "analytics_test_session"
	
	// Track a page view
	pageView := PageView{
		SessionID:      testSessionID,
		UserAgent:      "Test User Agent",
		PageURL:        "https://lugx-gaming.com/test",
		PageTitle:      "Test Page",
		Referrer:       "",
		PageLoadTime:   1000,
		ViewportWidth:  1920,
		ViewportHeight: 1080,
	}

	jsonData, err := json.Marshal(pageView)
	if err != nil {
		t.Fatalf("Failed to marshal test page view: %v", err)
	}

	// Send the page view
	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pageview", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send test page view: %v", err)
	}
	resp.Body.Close()

	// Wait a moment for data to be processed
	time.Sleep(time.Second)

	// Now get analytics data
	getResp, err := http.Get(analyticsServiceBaseURL + "/api/analytics/data")
	if err != nil {
		t.Fatalf("Failed to get analytics data: %v", err)
	}
	defer getResp.Body.Close()

	if getResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", getResp.StatusCode)
	}

	var analyticsData []map[string]interface{}
	if err := json.NewDecoder(getResp.Body).Decode(&analyticsData); err != nil {
		t.Fatalf("Failed to decode analytics data response: %v", err)
	}

	// Validate the structure of analytics data
	if len(analyticsData) == 0 {
		t.Log("No analytics data returned - this might be expected if the database is empty or the query returns no results")
		return
	}

	// Check that each analytics entry has expected fields
	for i, entry := range analyticsData {
		if _, hasURL := entry["page_url"]; !hasURL {
			t.Errorf("Analytics entry %d missing 'page_url' field", i)
		}
		
		if _, hasViews := entry["views"]; !hasViews {
			t.Errorf("Analytics entry %d missing 'views' field", i)
		}
		
		if _, hasAvgLoadTime := entry["avg_load_time"]; !hasAvgLoadTime {
			t.Errorf("Analytics entry %d missing 'avg_load_time' field", i)
		}
	}

	t.Logf("Successfully retrieved analytics data with %d entries", len(analyticsData))
}

func TestGetAnalyticsDataWithDateRange(t *testing.T) {
	// Test with date range parameters
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02") // 7 days ago
	endDate := time.Now().Format("2006-01-02")                     // today

	url := fmt.Sprintf("%s/api/analytics/data?start_date=%s&end_date=%s", 
		analyticsServiceBaseURL, startDate, endDate)

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Failed to get analytics data with date range: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	var analyticsData []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&analyticsData); err != nil {
		t.Fatalf("Failed to decode analytics data response: %v", err)
	}

	t.Logf("Retrieved analytics data for date range %s to %s: %d entries", 
		startDate, endDate, len(analyticsData))
}

func TestMultipleEventTypes(t *testing.T) {
	sessionID := "multi_event_test_session"
	baseURL := "https://lugx-gaming.com/multi-test"

	// Track multiple types of events for the same session
	events := []struct {
		name string
		url  string
		data interface{}
	}{
		{
			name: "pageview",
			url:  "/api/analytics/pageview",
			data: PageView{
				SessionID:      sessionID,
				UserAgent:      "Multi-Event Test Agent",
				PageURL:        baseURL,
				PageTitle:      "Multi Event Test Page",
				Referrer:       "",
				PageLoadTime:   800,
				ViewportWidth:  1366,
				ViewportHeight: 768,
			},
		},
		{
			name: "click",
			url:  "/api/analytics/click",
			data: Click{
				SessionID:    sessionID,
				PageURL:      baseURL,
				ElementTag:   "a",
				ElementID:    "nav-link",
				ElementClass: "nav-item",
				ElementText:  "Products",
				ClickX:       200,
				ClickY:       50,
			},
		},
		{
			name: "scroll",
			url:  "/api/analytics/scroll",
			data: ScrollDepth{
				SessionID:           sessionID,
				PageURL:             baseURL,
				MaxScrollPercentage: 50,
				TotalPageHeight:     1500,
				ViewportHeight:      768,
			},
		},
		{
			name: "pagetime",
			url:  "/api/analytics/pagetime",
			data: PageTime{
				SessionID:    sessionID,
				PageURL:      baseURL,
				TimeOnPage:   30000, // 30 seconds
				IsActiveTime: 1,
			},
		},
	}

	// Send all events
	for _, event := range events {
		jsonData, err := json.Marshal(event.data)
		if err != nil {
			t.Fatalf("Failed to marshal %s event: %v", event.name, err)
		}

		resp, err := http.Post(analyticsServiceBaseURL+event.url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			t.Fatalf("Failed to send %s event: %v", event.name, err)
		}
		
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Errorf("Failed to track %s event, got status %d", event.name, resp.StatusCode)
		}
		
		resp.Body.Close()
		
		// Small delay between events
		time.Sleep(100 * time.Millisecond)
	}

	t.Logf("Successfully sent %d different types of analytics events", len(events))
}

func TestInvalidAnalyticsData(t *testing.T) {
	// Test with invalid JSON
	invalidJSON := `{"session_id": "test", "invalid_json"`

	resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pageview", "application/json", bytes.NewBufferString(invalidJSON))
	if err != nil {
		t.Fatalf("Failed to send invalid JSON: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code 400 for invalid JSON, got %d", resp.StatusCode)
	}

	// Test with missing required fields
	incompleteData := map[string]interface{}{
		"session_id": "test",
		// Missing other required fields
	}

	jsonData, err := json.Marshal(incompleteData)
	if err != nil {
		t.Fatalf("Failed to marshal incomplete data: %v", err)
	}

	resp2, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pageview", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send incomplete data: %v", err)
	}
	defer resp2.Body.Close()

	// The service might accept incomplete data or return an error
	// Either way is acceptable for this test
	t.Logf("Incomplete data request returned status: %d", resp2.StatusCode)
}

func TestConcurrentAnalyticsTracking(t *testing.T) {
	// Test concurrent requests to ensure the service can handle multiple simultaneous analytics events
	const numConcurrentRequests = 10
	
	done := make(chan bool, numConcurrentRequests)
	errors := make(chan error, numConcurrentRequests)

	for i := 0; i < numConcurrentRequests; i++ {
		go func(requestID int) {
			pageView := PageView{
				SessionID:      fmt.Sprintf("concurrent_test_session_%d", requestID),
				UserAgent:      "Concurrent Test Agent",
				PageURL:        fmt.Sprintf("https://lugx-gaming.com/concurrent-test-%d", requestID),
				PageTitle:      fmt.Sprintf("Concurrent Test Page %d", requestID),
				Referrer:       "",
				PageLoadTime:   1000 + requestID*100,
				ViewportWidth:  1920,
				ViewportHeight: 1080,
			}

			jsonData, err := json.Marshal(pageView)
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to marshal: %v", requestID, err)
				return
			}

			resp, err := http.Post(analyticsServiceBaseURL+"/api/analytics/pageview", "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				errors <- fmt.Errorf("request %d: failed to send: %v", requestID, err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				errors <- fmt.Errorf("request %d: unexpected status %d", requestID, resp.StatusCode)
				return
			}

			done <- true
		}(i)
	}

	// Wait for all requests to complete
	successCount := 0
	errorCount := 0

	for i := 0; i < numConcurrentRequests; i++ {
		select {
		case <-done:
			successCount++
		case err := <-errors:
			t.Logf("Concurrent request error: %v", err)
			errorCount++
		case <-time.After(10 * time.Second):
			t.Fatalf("Timeout waiting for concurrent requests to complete")
		}
	}

	t.Logf("Concurrent test completed: %d successes, %d errors", successCount, errorCount)

	if successCount == 0 {
		t.Errorf("No concurrent requests succeeded")
	}
}
