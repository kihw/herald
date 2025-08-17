# Simple PowerShell Deployment script for herald.lol
param(
    [string]$ServerIP = "51.178.17.78",
    [string]$ServerUser = "root"
)

Write-Host "🚀 Deploying LoL Match Exporter to herald.lol..." -ForegroundColor Green

# Create archive excluding certain files
Write-Host "📦 Creating deployment archive..." -ForegroundColor Yellow

# Use tar to create archive (requires Git for Windows or WSL)
$excludeList = @(
    "--exclude=*.exe",
    "--exclude=*.db", 
    "--exclude=node_modules",
    "--exclude=web/dist",
    "--exclude=web/node_modules",
    "--exclude=data",
    "--exclude=exports", 
    "--exclude=logs",
    "--exclude=.git",
    "--exclude=*.log",
    "--exclude=*.tmp"
)

$tarArgs = @("-czf", "lol-exporter-deployment.tar.gz") + $excludeList + @(".")
& tar @tarArgs

Write-Host "✅ Archive created" -ForegroundColor Green

# Transfer to server
Write-Host "📤 Transferring to server..." -ForegroundColor Yellow
& scp "lol-exporter-deployment.tar.gz" "${ServerUser}@${ServerIP}:/tmp/"

Write-Host "🎉 Files transferred! Now connect to server and run deployment." -ForegroundColor Green
Write-Host "Commands to run on server:" -ForegroundColor Cyan
Write-Host "ssh $ServerUser@$ServerIP" -ForegroundColor White
Write-Host "cd /opt && mkdir -p lol-match-exporter && cd lol-match-exporter" -ForegroundColor White
Write-Host "tar -xzf /tmp/lol-exporter-deployment.tar.gz" -ForegroundColor White
Write-Host "cp .env.herald .env" -ForegroundColor White
Write-Host "docker-compose -f docker-compose.production.yml up -d --build" -ForegroundColor White

# Clean up
Remove-Item "lol-exporter-deployment.tar.gz" -Force