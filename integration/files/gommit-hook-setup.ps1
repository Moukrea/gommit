# Determine architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# Download gommit
$gommitUrl = "https://github.com/Moukrea/gommit/releases/download/latest/gommit-windows-$arch"
Invoke-WebRequest -Uri $gommitUrl -OutFile ".gommit\gommit.exe"

# Prepare commit-msg hook content for gommit
$gommitHookContent = @"

# Gommit commit message validation
./.gommit/gommit.exe `$1
exit `$LASTEXITCODE
"@

# Handle commit-msg hook
$destFile = ".git\hooks\commit-msg"

if (Test-Path $destFile) {
    Write-Host "Existing commit-msg hook found."
    $content = Get-Content $destFile -Raw
    if ($content -match "/.gommit/gommit") {
        Write-Host "Gommit hook already exists in commit-msg. No changes needed."
    } else {
        $choice = Read-Host "Choose action (overwrite/append/skip)"
        switch ($choice.ToLower()) {
            "overwrite" {
                Set-Content -Path $destFile -Value "#!/bin/sh$gommitHookContent"
                Write-Host "Overwrote commit-msg hook with Gommit."
            }
            "append" {
                if ($content -match "exit\s") {
                    $content = $content -replace "exit\s", "${gommitHookContent}`nexit "
                } else {
                    $content += $gommitHookContent
                }
                Set-Content -Path $destFile -Value $content
                Write-Host "Appended Gommit to existing commit-msg hook."
            }
            "skip" {
                Write-Host "Skipped modifying commit-msg hook."
            }
            default {
                Write-Host "Invalid choice. Skipping commit-msg hook modification."
            }
        }
    }
} else {
    Set-Content -Path $destFile -Value "#!/bin/sh$gommitHookContent"
    Write-Host "Created new commit-msg hook with Gommit."
}

Write-Host "Gommit has been set up successfully."