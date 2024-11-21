# Create .gommit directory
New-Item -ItemType Directory -Force -Path .gommit | Out-Null

# Download gommit-hook-setup.sh
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/moukrea/gommit/main/integration/files/gommit-hook-setup.sh" -OutFile ".gommit\gommit-hook-setup.sh"

# Download gommit-hook-setup.ps1
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/moukrea/gommit/main/integration/files/gommit-hook-setup.ps1" -OutFile ".gommit\gommit-hook-setup.ps1"

# Download Makefile
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/moukrea/gommit/main/integration/files/Makefile" -OutFile "Makefile"

# Add .gommit/gommit and .gommit/last_update to .gitignore
$gitignoreContent = Get-Content .gitignore -ErrorAction SilentlyContinue
if (-not $gitignoreContent) {
    $gitignoreContent = @()
}
$linesToAdd = @('.gommit/gommit', '.gommit/last_update')
foreach ($line in $linesToAdd) {
    if ($gitignoreContent -notcontains $line) {
        Add-Content .gitignore "`n$line"
    }
}

Write-Host "Gommit integration files have been set up successfully."
Write-Host "The .gitignore file has been updated with .gommit/gommit and .gommit/last_update."
Write-Host "You can now run 'make gommit-setup' to complete the setup."
