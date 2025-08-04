# Analytics Service API Test Script
# Tests the analytics service endpoints in Kubernetes

$analyticsUrl = "http://localhost:30082/api/analytics"

Write-Host "🧪 Testing Analytics Service API" -ForegroundColor Green
Write-Host "=================================" -ForegroundColor Green

# Test health endpoint
Write-Host "`n🏥 Testing health endpoint..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:30082/health" -Method GET
    Write-Host "✅ Health check passed: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Green
} catch {
    Write-Host "❌ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test page view tracking
Write-Host "`n📄 Testing page view tracking..." -ForegroundColor Yellow
$pageViewData = @{
    session_id = "test-session-" + (Get-Date -Format "yyyyMMdd-HHmmss")
    user_agent = "PowerShell-Test/1.0"
    page_url = "http://localhost:30000/test"
    page_title = "Test Page"
    referrer = ""
    page_load_time = 1250
    viewport_width = 1920
    viewport_height = 1080
}

try {
    $response = Invoke-RestMethod -Uri "$analyticsUrl/pageview" -Method POST -Body ($pageViewData | ConvertTo-Json) -ContentType "application/json"
    Write-Host "✅ Page view tracked: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Green
} catch {
    Write-Host "❌ Page view tracking failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test click tracking
Write-Host "`n🖱️ Testing click tracking..." -ForegroundColor Yellow
$clickData = @{
    session_id = $pageViewData.session_id
    page_url = "http://localhost:30000/test"
    element_tag = "button"
    element_id = "test-button"
    element_class = "btn btn-primary"
    element_text = "Test Button"
    click_x = 150
    click_y = 75
}

try {
    $response = Invoke-RestMethod -Uri "$analyticsUrl/click" -Method POST -Body ($clickData | ConvertTo-Json) -ContentType "application/json"
    Write-Host "✅ Click tracked: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Green
} catch {
    Write-Host "❌ Click tracking failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test scroll tracking
Write-Host "`n📜 Testing scroll tracking..." -ForegroundColor Yellow
$scrollData = @{
    session_id = $pageViewData.session_id
    page_url = "http://localhost:30000/test"
    max_scroll_percentage = 75
    total_page_height = 2400
    viewport_height = 1080
}

try {
    $response = Invoke-RestMethod -Uri "$analyticsUrl/scroll" -Method POST -Body ($scrollData | ConvertTo-Json) -ContentType "application/json"
    Write-Host "✅ Scroll tracked: $($response | ConvertTo-Json -Depth 2)" -ForegroundColor Green
} catch {
    Write-Host "❌ Scroll tracking failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test analytics data retrieval
Write-Host "`n📊 Testing analytics data retrieval..." -ForegroundColor Yellow
try {
    Start-Sleep 2 # Wait a moment for data to be processed
    $response = Invoke-RestMethod -Uri "$analyticsUrl/data" -Method GET
    Write-Host "✅ Analytics data retrieved:" -ForegroundColor Green
    $response | ForEach-Object {
        Write-Host "   Page: $($_.page_url), Views: $($_.views), Avg Load Time: $($_.avg_load_time)ms" -ForegroundColor White
    }
} catch {
    Write-Host "❌ Analytics data retrieval failed: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n🎉 Analytics API testing completed!" -ForegroundColor Green
Write-Host "`n💡 Tip: You can also test the analytics in the browser by visiting:" -ForegroundColor Cyan
Write-Host "   Frontend: http://localhost:30000" -ForegroundColor White
Write-Host "   The frontend will automatically track user behavior!" -ForegroundColor White
