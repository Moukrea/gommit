# Create .gommit directory
New-Item -ItemType Directory -Force -Path .gommit | Out-Null

# Download setup.sh
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/setup.sh" -OutFile ".gommit\setup.sh"

# Download setup.ps1
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/setup.ps1" -OutFile ".gommit\setup.ps1"

# Download Makefile
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Moukrea/gommit/main/integration/Makefile" -OutFile "Makefile"

Write-Host "Gommit integration files have been set up successfully."
Write-Host "You can now run 'make setup-gommit' to complete the setup."