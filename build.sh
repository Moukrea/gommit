#!/bin/bash

# Set color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# List of components to build
COMPONENTS=("gommit")

# List of target OS and architectures
TARGETS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

# Create a build directory
mkdir -p build

# Function to initialize Go module if not present
initialize_go_module() {
    local component_dir=$1
    pushd "$component_dir" > /dev/null
    if [ ! -f "go.mod" ]; then
        echo -e "  ${YELLOW}Initializing Go module in $component_dir${NC}"
        if ! go mod init "github.com/Moukrea/gommit/$component_dir"; then
            echo -e "    ${RED}Failed to initialize Go module in $component_dir${NC}"
            popd > /dev/null
            return 1
        fi
    fi
    if ! go mod tidy; then
        echo -e "    ${RED}Failed to tidy Go module in $component_dir${NC}"
        popd > /dev/null
        return 1
    fi
    popd > /dev/null
    return 0
}

# Loop through each component
for component in "${COMPONENTS[@]}"; do
    echo -e "Building ${CYAN}$component${NC}..."
    
    # Initialize Go module if necessary
    if ! initialize_go_module "$component"; then
        echo -e "  ${RED}Skipping $component due to module initialization failure${NC}"
        continue
    fi
    
    # Loop through each target
    for target in "${TARGETS[@]}"; do
        # Split the target into OS and architecture
        IFS='/' read -r OS ARCH <<< "$target"
        
        # Set the output filename
        if [ "$OS" == "windows" ]; then
            OUTPUT="build/${component}-${OS}-${ARCH}.exe"
        else
            OUTPUT="build/${component}-${OS}-${ARCH}"
        fi
        
        echo -e "  Building for ${CYAN}$OS/$ARCH${NC}..."
        
        # Build the binary
        pushd "$component" > /dev/null
        GOOS=$OS GOARCH=$ARCH go build -o "../$OUTPUT" .
        
        if [ $? -eq 0 ]; then
            echo -e "    ${GREEN}Success: $OUTPUT${NC}"
        else
            echo -e "    ${RED}Failed: $OUTPUT${NC}"
        fi
        popd > /dev/null
    done
done

echo -e "${CYAN}Build process completed.${NC}"