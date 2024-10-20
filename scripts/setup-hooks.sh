#!/bin/sh
set -e

REPO_ROOT=$(git rev-parse --show-toplevel)

# Check if hooks are already set up
if [ "$(git config --get core.hooksPath)" = "$REPO_ROOT/scripts/git-hooks" ]; then
    echo "Gommit hooks are already set up."
else
    # Set up the hooks path
    git config core.hooksPath "$REPO_ROOT/scripts/git-hooks"
    echo "Git hooks have been set up successfully."
fi

# Ensure the Gommit binary is downloaded
"$REPO_ROOT/scripts/gommit.sh"

echo "Gommit is ready to use."
