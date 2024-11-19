# PowerShell script to build Gommit binaries for multiple platforms

# List of components to build
$COMPONENTS = @("gommit")

# List of target OS and architectures
$TARGETS = @(
    "linux/amd64",
    "linux/arm64",
    "darwin/amd64",
    "darwin/arm64",
    "windows/amd64",
    "windows/arm64"
)

# Create a build directory
New-Item -ItemType Directory -Force -Path "build" | Out-Null

# Function to initialize Go module if not present
function Initialize-GoModule {
    param (
        [string]$componentDir
    )
    Push-Location $componentDir
    if (-not (Test-Path "go.mod")) {
        Write-Host "  Initializing Go module in $componentDir" -ForegroundColor Yellow
        go mod init "github.com/Moukrea/gommit/$componentDir"
        if ($LASTEXITCODE -ne 0) {
            Write-Host "    Failed to initialize Go module in $componentDir" -ForegroundColor Red
            Pop-Location
            return $false
        }
    }
    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "    Failed to tidy Go module in $componentDir" -ForegroundColor Red
        Pop-Location
        return $false
    }
    Pop-Location
    return $true
}

# Loop through each component
foreach ($component in $COMPONENTS) {
    Write-Host "Building $component..."
    
    # Initialize Go module if necessary
    if (-not (Initialize-GoModule $component)) {
        Write-Host "  Skipping $component due to module initialization failure" -ForegroundColor Red
        continue
    }
    
    # Loop through each target
    foreach ($target in $TARGETS) {
        # Split the target into OS and architecture
        $OS, $ARCH = $target.Split("/")
        
        # Set the output filename
        $OUTPUT = if ($OS -eq "windows") {
            "build\${component}-${OS}-${ARCH}.exe"
        } else {
            "build\${component}-${OS}-${ARCH}"
        }
        
        Write-Host "  Building for $OS/$ARCH..."
        
        # Build the binary
        Push-Location $component
        $env:GOOS = $OS
        $env:GOARCH = $ARCH
        go build -o "..\$OUTPUT" .
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "    Success: $OUTPUT" -ForegroundColor Green
        } else {
            Write-Host "    Failed: $OUTPUT" -ForegroundColor Red
        }
        Pop-Location
    }
}

Write-Host "Build process completed." -ForegroundColor Cyan