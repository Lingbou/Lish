# Lish Windows 
# powershell -ExecutionPolicy Bypass -File install.ps1

param(
    [string]$InstallPath = "$env:LOCALAPPDATA\Lish"
)

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "    Lish 安装程序" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# 检查是否以管理员身份运行（可选）
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "提示: 未以管理员身份运行" -ForegroundColor Yellow
    Write-Host "      安装到用户目录: $InstallPath" -ForegroundColor Yellow
    Write-Host ""
}

# 检查是否存在 lish.exe
if (-not (Test-Path "lish.exe")) {
    Write-Host "错误: 找不到 lish.exe" -ForegroundColor Red
    Write-Host "      请先编译: go build -ldflags='-s -w' -o lish.exe cmd/lish/main.go" -ForegroundColor Yellow
    exit 1
}

Write-Host "[1/4] 创建安装目录..." -ForegroundColor Green
if (-not (Test-Path $InstallPath)) {
    New-Item -ItemType Directory -Path $InstallPath -Force | Out-Null
}
Write-Host "      完成: $InstallPath" -ForegroundColor Gray

Write-Host "[2/4] 复制文件..." -ForegroundColor Green
Copy-Item "lish.exe" -Destination "$InstallPath\lish.exe" -Force
Write-Host "      完成: lish.exe" -ForegroundColor Gray

Write-Host "[3/4] 添加到 PATH..." -ForegroundColor Green
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$currentPath;$InstallPath", "User")
    Write-Host "      完成: 已添加到用户 PATH" -ForegroundColor Gray
    Write-Host "      注意: 请重新打开终端使 PATH 生效" -ForegroundColor Yellow
} else {
    Write-Host "      跳过: 已在 PATH 中" -ForegroundColor Gray
}

Write-Host "[4/4] 配置 Windows Terminal..." -ForegroundColor Green

# 查找 Windows Terminal 配置文件
$wtSettingsPath = "$env:LOCALAPPDATA\Packages\Microsoft.WindowsTerminal_8wekyb3d8bbwe\LocalState\settings.json"

if (Test-Path $wtSettingsPath) {
    Write-Host "      找到 Windows Terminal 配置文件" -ForegroundColor Gray
    Write-Host ""
    Write-Host "      请手动添加以下配置到 Windows Terminal:" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "      {" -ForegroundColor Cyan
    Write-Host "          `"name`": `"Lish`"," -ForegroundColor Cyan
    Write-Host "          `"commandline`": `"$InstallPath\\lish.exe`"," -ForegroundColor Cyan
    Write-Host "          `"startingDirectory`": `"%USERPROFILE%`"," -ForegroundColor Cyan
    Write-Host "          `"icon`": `"🐚`"," -ForegroundColor Cyan
    Write-Host "          `"colorScheme`": `"One Half Dark`"" -ForegroundColor Cyan
    Write-Host "      }" -ForegroundColor Cyan
    Write-Host ""
} else {
    Write-Host "      未找到 Windows Terminal" -ForegroundColor Gray
    Write-Host "      如需集成，请参考 README.md" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "    安装完成！" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "安装位置: $InstallPath" -ForegroundColor Gray
Write-Host ""
Write-Host "快速开始:" -ForegroundColor Yellow
Write-Host "  1. 重新打开终端" -ForegroundColor White
Write-Host "  2. 输入: lish" -ForegroundColor White
Write-Host "  3. 输入: help" -ForegroundColor White
Write-Host ""
Write-Host "更多信息请查看: https://github.com/Lingbou/Lish" -ForegroundColor Cyan
Write-Host ""

