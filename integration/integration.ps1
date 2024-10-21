# Create .gommit directory
New-Item -ItemType Directory -Force -Path .gommit | Out-Null

# Download gommit-hook-setup.sh
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.sh" -OutFile ".gommit\gommit-hook-setup.sh"

# Download gommit-hook-setup.ps1
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.ps1" -OutFile ".gommit\gommit-hook-setup.ps1"

# Download Makefile
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/Makefile" -OutFile "Makefile"

Write-Host "Gommit integration files have been set up successfully."
Write-Host "You can now run 'make gommit-setup' to complete the setup."