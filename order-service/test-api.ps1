# PowerShell script to test Order Service API
param(
    [string]$BaseUrl = "http://localhost:8081"
)

Write-Host "Testing Order Service API at $BaseUrl" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# Function to make HTTP requests
function Invoke-ApiRequest {
    param(
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null,
        [string]$Description
    )
    
    Write-Host "`n$Description" -ForegroundColor Yellow
    Write-Host "$Method $Endpoint" -ForegroundColor Cyan
    
    try {
        $params = @{
            Uri = "$BaseUrl$Endpoint"
            Method = $Method
            ContentType = "application/json"
        }
        
        if ($Body) {
            $params.Body = $Body | ConvertTo-Json -Depth 10
            Write-Host "Request Body: $($params.Body)" -ForegroundColor Gray
        }
        
        $response = Invoke-RestMethod @params
        Write-Host "✅ Success: $($response | ConvertTo-Json -Depth 10)" -ForegroundColor Green
        return $response
    }
    catch {
        Write-Host "❌ Error: $($_.Exception.Message)" -ForegroundColor Red
        if ($_.Exception.Response) {
            $errorDetails = $_.Exception.Response | ConvertFrom-Json -ErrorAction SilentlyContinue
            if ($errorDetails) {
                Write-Host "Error Details: $($errorDetails | ConvertTo-Json)" -ForegroundColor Red
            }
        }
        return $null
    }
}

# Test 1: Health Check
$healthResponse = Invoke-ApiRequest -Method "GET" -Endpoint "/health" -Description "1. Health Check"

if (-not $healthResponse) {
    Write-Host "`n❌ Service is not running. Please start the order service first." -ForegroundColor Red
    exit 1
}

# Test 2: Create Order
$orderData = @{
    customer_id = "customer123"
    items = @(
        @{
            game_id = 1
            game_name = "Cyberpunk 2077"
            price = 59.99
            quantity = 1
        },
        @{
            game_id = 2
            game_name = "The Witcher 3"
            price = 39.99
            quantity = 2
        }
    )
}

$createdOrder = Invoke-ApiRequest -Method "POST" -Endpoint "/api/v1/orders" -Body $orderData -Description "2. Create Order"

if ($createdOrder) {
    $orderId = $createdOrder.order.id
    Write-Host "Created Order ID: $orderId" -ForegroundColor Magenta
    
    # Test 3: Get Order by ID
    Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders/$orderId" -Description "3. Get Order by ID"
    
    # Test 4: Update Order Status
    $statusUpdate = @{
        status = "confirmed"
    }
    Invoke-ApiRequest -Method "PUT" -Endpoint "/api/v1/orders/$orderId/status" -Body $statusUpdate -Description "4. Update Order Status"
    
    # Test 5: Get Updated Order
    Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders/$orderId" -Description "5. Get Updated Order"
}

# Test 6: Get All Orders
Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders" -Description "6. Get All Orders"

# Test 7: Get Orders with Pagination
Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders?page=1&page_size=5" -Description "7. Get Orders with Pagination"

# Test 8: Get Orders by Customer
Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders/customer/customer123" -Description "8. Get Orders by Customer"

# Test 9: Get Order Statistics
Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders/stats" -Description "9. Get Order Statistics"

# Test 10: Create Another Order (different customer)
$orderData2 = @{
    customer_id = "customer456"
    items = @(
        @{
            game_id = 3
            game_name = "Red Dead Redemption 2"
            price = 49.99
            quantity = 1
        }
    )
}

$createdOrder2 = Invoke-ApiRequest -Method "POST" -Endpoint "/api/v1/orders" -Body $orderData2 -Description "10. Create Another Order"

# Test 11: Test Error Cases
Write-Host "`n🧪 Testing Error Cases" -ForegroundColor Yellow

# Invalid order (missing customer_id)
$invalidOrder = @{
    items = @(
        @{
            game_id = 1
            game_name = "Test Game"
            price = 10.00
            quantity = 1
        }
    )
}
Invoke-ApiRequest -Method "POST" -Endpoint "/api/v1/orders" -Body $invalidOrder -Description "11a. Invalid Order (missing customer_id)"

# Invalid status update
if ($orderId) {
    $invalidStatus = @{
        status = "invalid_status"
    }
    Invoke-ApiRequest -Method "PUT" -Endpoint "/api/v1/orders/$orderId/status" -Body $invalidStatus -Description "11b. Invalid Status Update"
}

# Non-existent order
Invoke-ApiRequest -Method "GET" -Endpoint "/api/v1/orders/non-existent-id" -Description "11c. Get Non-existent Order"

Write-Host "`n✅ API Testing Complete!" -ForegroundColor Green
Write-Host "===========================================" -ForegroundColor Green

# Optional: Clean up (delete test orders)
$cleanup = Read-Host "`nDo you want to delete the test orders? (y/N)"
if ($cleanup -eq "y" -or $cleanup -eq "Y") {
    if ($orderId) {
        Invoke-ApiRequest -Method "DELETE" -Endpoint "/api/v1/orders/$orderId" -Description "Cleanup: Delete First Test Order"
    }
    if ($createdOrder2) {
        $orderId2 = $createdOrder2.order.id
        Invoke-ApiRequest -Method "DELETE" -Endpoint "/api/v1/orders/$orderId2" -Description "Cleanup: Delete Second Test Order"
    }
    Write-Host "Cleanup complete!" -ForegroundColor Green
}

Write-Host "`n📝 Test Summary:" -ForegroundColor Cyan
Write-Host "- Health check: Verified service is running"
Write-Host "- Order CRUD: Create, Read, Update operations tested"
Write-Host "- Pagination: Tested with query parameters"
Write-Host "- Customer filtering: Tested customer-specific orders"
Write-Host "- Statistics: Tested order analytics endpoint"
Write-Host "- Error handling: Tested invalid requests"
Write-Host "- Data validation: Tested input validation"
