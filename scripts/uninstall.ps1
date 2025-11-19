# Lish Windows 卸载脚本
# 使用方法: powershell -ExecutionPolicy Bypass -File uninstall.ps1

param(
    [string]$InstallPath = "$env:LOCALAPPDATA\Lish"
)

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "    Lish 卸载程序" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# 确认卸载
$confirmation = Read-Host "确定要卸载 Lish 吗？(y/N)"
if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
    Write-Host "已取消卸载" -ForegroundColor Yellow
    exit 0
}

Write-Host ""
Write-Host "[1/3] 删除安装文件..." -ForegroundColor Green
if (Test-Path $InstallPath) {
    Remove-Item -Path $InstallPath -Recurse -Force
    Write-Host "      完成: 已删除 $InstallPath" -ForegroundColor Gray
} else {
    Write-Host "      跳过: 未找到安装目录" -ForegroundColor Gray
}

Write-Host "[2/3] 从 PATH 中移除..." -ForegroundColor Green
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -like "*$InstallPath*") {
    $newPath = ($currentPath -split ';' | Where-Object { $_ -ne $InstallPath }) -join ';'
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "      完成: 已从 PATH 中移除" -ForegroundColor Gray
} else {
    Write-Host "      跳过: PATH 中不存在" -ForegroundColor Gray
}

Write-Host "[3/3] 清理历史记录..." -ForegroundColor Green
$historyFile = "$env:USERPROFILE\.lish_history"
if (Test-Path $historyFile) {
    $deleteHistory = Read-Host "是否删除历史记录文件？(y/N)"
    if ($deleteHistory -eq 'y' -or $deleteHistory -eq 'Y') {
        Remove-Item -Path $historyFile -Force
        Write-Host "      完成: 已删除历史记录" -ForegroundColor Gray
    } else {
        Write-Host "      跳过: 保留历史记录" -ForegroundColor Gray
    }
} else {
    Write-Host "      跳过: 未找到历史记录文件" -ForegroundColor Gray
}

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "    卸载完成！" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "提示: 请重新打开终端使更改生效" -ForegroundColor Yellow
Write-Host ""
Write-Host "感谢使用 Lish！" -ForegroundColor Cyan
Write-Host ""

