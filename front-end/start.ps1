# PowerShell script to build and run the frontend Docker container

param(
    [string]$Action = "help"
)

$ImageName = "lugx-gaming-frontend"
$ContainerName = "lugx-gaming-frontend"
$Port = "3000"

function Show-Help {
    Write-Host "Lugx Gaming Frontend Docker Management" -ForegroundColor Green
    Write-Host "Usage: .\start.ps1 [action]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Actions:" -ForegroundColor Cyan
    Write-Host "  build     - Build the Docker image"
    Write-Host "  run       - Run the container"
    Write-Host "  stop      - Stop and remove the container"
    Write-Host "  clean     - Remove the Docker image"
    Write-Host "  up        - Build and run (full start)"
    Write-Host "  rebuild   - Stop, clean, build, and run"
    Write-Host "  logs      - Show container logs"
    Write-Host "  status    - Show container status"
    Write-Host "  help      - Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Magenta
    Write-Host "  .\start.ps1 up"
    Write-Host "  .\start.ps1 logs"
    Write-Host "  .\start.ps1 stop"
}

function Build-Image {
    Write-Host "Building Docker image: $ImageName" -ForegroundColor Green
    docker build -t $ImageName .
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Image built successfully!" -ForegroundColor Green
    } else {
        Write-Host "❌ Failed to build image" -ForegroundColor Red
        exit 1
    }
}

function Run-Container {
    Write-Host "Starting container: $ContainerName on port $Port" -ForegroundColor Green
    
    # Stop existing container if running
    $existing = docker ps -q -f name=$ContainerName
    if ($existing) {
        Write-Host "Stopping existing container..." -ForegroundColor Yellow
        docker stop $ContainerName | Out-Null
        docker rm $ContainerName | Out-Null
    }
    
    docker run -d --name $ContainerName -p "$Port:80" $ImageName
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Container started successfully!" -ForegroundColor Green
        Write-Host "🌐 Access the frontend at: http://localhost:$Port" -ForegroundColor Cyan
    } else {
        Write-Host "❌ Failed to start container" -ForegroundColor Red
        exit 1
    }
}

function Stop-Container {
    Write-Host "Stopping and removing container: $ContainerName" -ForegroundColor Yellow
    docker stop $ContainerName 2>$null
    docker rm $ContainerName 2>$null
    Write-Host "✅ Container stopped and removed" -ForegroundColor Green
}

function Clean-Image {
    Stop-Container
    Write-Host "Removing Docker image: $ImageName" -ForegroundColor Yellow
    docker rmi $ImageName 2>$null
    Write-Host "✅ Image removed" -ForegroundColor Green
}

function Show-Logs {
    Write-Host "Showing logs for container: $ContainerName" -ForegroundColor Green
    docker logs -f $ContainerName
}

function Show-Status {
    Write-Host "Container Status:" -ForegroundColor Green
    docker ps -f name=$ContainerName
    Write-Host ""
    Write-Host "Image Status:" -ForegroundColor Green
    docker images $ImageName
}

switch ($Action.ToLower()) {
    "build" { Build-Image }
    "run" { Run-Container }
    "stop" { Stop-Container }
    "clean" { Clean-Image }
    "up" { 
        Build-Image
        Run-Container
    }
    "rebuild" {
        Clean-Image
        Build-Image
        Run-Container
    }
    "logs" { Show-Logs }
    "status" { Show-Status }
    "help" { Show-Help }
    default { Show-Help }
}
