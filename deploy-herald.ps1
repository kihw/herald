# PowerShell Deployment script for herald.lol (51.178.17.78)
# Run this script from your local Windows machine

param(
    [string]$ServerIP = "51.178.17.78",
    [string]$ServerUser = "root"
)

$ErrorActionPreference = "Stop"

Write-Host "🚀 Deploying LoL Match Exporter to herald.lol..." -ForegroundColor Green

# 1. Create deployment archive
Write-Host "📦 Creating deployment archive..." -ForegroundColor Yellow

$excludePatterns = @(
    "*.exe",
    "*.db", 
    "node_modules",
    "web/dist",
    "web/node_modules",
    "data",
    "exports", 
    "logs",
    ".git",
    "*.log",
    "*.tmp"
)

# Create temp directory for filtered files
$tempDir = Join-Path $env:TEMP "lol-exporter-deploy"
if (Test-Path $tempDir) {
    Remove-Item $tempDir -Recurse -Force
}
New-Item -ItemType Directory -Path $tempDir | Out-Null

# Copy files excluding patterns
$sourceFiles = Get-ChildItem -Path . -Recurse | Where-Object {
    $file = $_
    $shouldExclude = $false
    foreach ($pattern in $excludePatterns) {
        if ($file.FullName -like "*$pattern*") {
            $shouldExclude = $true
            break
        }
    }
    -not $shouldExclude
}

foreach ($file in $sourceFiles) {
    $relativePath = $file.FullName.Substring((Get-Location).Path.Length + 1)
    $destPath = Join-Path $tempDir $relativePath
    $destDir = Split-Path $destPath -Parent
    
    if (-not (Test-Path $destDir)) {
        New-Item -ItemType Directory -Path $destDir -Force | Out-Null
    }
    
    if ($file.PSIsContainer -eq $false) {
        Copy-Item $file.FullName $destPath
    }
}

# Create tar.gz archive using 7zip or tar (if available)
$archivePath = "lol-exporter-deployment.tar.gz"
if (Get-Command tar -ErrorAction SilentlyContinue) {
    tar -czf $archivePath -C $tempDir .
} else {
    Write-Error "tar command not found. Please install Git for Windows or WSL."
    exit 1
}

Remove-Item $tempDir -Recurse -Force
Write-Host "✅ Archive created: $archivePath" -ForegroundColor Green

# 2. Transfer to server
Write-Host "📤 Transferring files to $ServerIP..." -ForegroundColor Yellow
scp $archivePath ${ServerUser}@${ServerIP}:/tmp/

# 3. Deploy on server
Write-Host "🏗️ Deploying on server..." -ForegroundColor Yellow

$sshScript = @"
    # Update system
    apt-get update -y
    
    # Install Docker if not present
    if ! command -v docker &> /dev/null; then
        echo "🐳 Installing Docker..."
        curl -fsSL https://get.docker.com -o get-docker.sh
        sh get-docker.sh
        systemctl enable docker
        systemctl start docker
    fi
    
    # Install Docker Compose if not present
    if ! command -v docker-compose &> /dev/null; then
        echo "🐳 Installing Docker Compose..."
        curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-`$(uname -s)-`$(uname -m)" -o /usr/local/bin/docker-compose
        chmod +x /usr/local/bin/docker-compose
    fi
    
    # Create deployment directory
    mkdir -p /opt/lol-match-exporter
    cd /opt/lol-match-exporter
    
    # Stop existing containers if running
    docker-compose -f docker-compose.production.yml down 2>/dev/null || true
    
    # Remove old files
    rm -rf * .*config 2>/dev/null || true
    
    # Extract new deployment
    tar -xzf /tmp/lol-exporter-deployment.tar.gz
    
    # Create required directories
    mkdir -p data exports logs logs/nginx
    
    # Copy environment file
    cp .env.herald .env
    
    # Set permissions
    chown -R 1000:1000 data exports logs
    
    # Build and start containers
    echo "🚀 Starting containers..."
    docker-compose -f docker-compose.production.yml up -d --build
    
    # Wait for services to start
    sleep 15
    
    # Check health
    echo "🏥 Checking health..."
    curl -f http://localhost/health || echo "⚠️ Health check failed"
    
    # Show status
    docker-compose -f docker-compose.production.yml ps
    
    # Show logs if there are issues
    docker-compose -f docker-compose.production.yml logs --tail=20
    
    # Clean up
    rm /tmp/lol-exporter-deployment.tar.gz
"@

ssh ${ServerUser}@${ServerIP} $sshScript

# 4. Test deployment
Write-Host "🧪 Testing deployment..." -ForegroundColor Yellow
Start-Sleep 5

try {
    $response = Invoke-WebRequest -Uri "http://herald.lol/health" -UseBasicParsing -TimeoutSec 10
    if ($response.StatusCode -eq 200) {
        Write-Host "✅ herald.lol is responding!" -ForegroundColor Green
    } else {
        Write-Host "❌ Unexpected response code: $($response.StatusCode)" -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Health check failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "🔍 Checking if server is accessible..." -ForegroundColor Yellow
    Test-NetConnection -ComputerName herald.lol -Port 80
}

# Clean up local archive
Remove-Item $archivePath -Force

Write-Host "🎉 Deployment completed!" -ForegroundColor Green
Write-Host "🌐 Application should be available at: http://herald.lol" -ForegroundColor Cyan
Write-Host "📊 Health endpoint: http://herald.lol/health" -ForegroundColor Cyan
Write-Host "📚 API docs: http://herald.lol/docs" -ForegroundColor Cyan

# Optional: Open browser
$openBrowser = Read-Host "Open herald.lol in browser? (y/N)"
if ($openBrowser -eq 'y' -or $openBrowser -eq 'Y') {
    Start-Process "http://herald.lol"
}