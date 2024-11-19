# Integrating Gommit Setup

This guide explains how to integrate the Gommit into your repository.

## Files

The Gommit setup consists of three files:

1. `Makefile` (in the root of your repository)
2. `.gommit/gommit-hook-setup.sh` (for Unix-like systems: macOS and Linux)
3. `.gommit/gommit-hook-setup.ps1` (for Windows)

## Integration Steps

1. Create a `.gommit` directory in your project's root:

   ```
   mkdir .gommit
   ```

2. Copy the provided [`integration/files/gommit-hook-setup.sh`](https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.sh) and [`integration/files/gommit-hook-setup.ps1`](https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/gommit-hook-setup.ps1)  into the `.gommit` directory.

3. Make the setup scripts executable (Unix-like systems only):

   ```
   chmod +x .gommit/gommit-hook-setup.sh
   ```

4. Integrate the Gommit setup into your project:

   a. If you don't have an existing Makefile:
      - Copy the provided [`integration/files/Makefile`](https://raw.githubusercontent.com/Moukrea/gommit/main/integration/files/Makefile) to your project root.

   b. If you have an existing Makefile:
      - Add the following content to your existing Makefile:
        ```makefile
        # Gommit setup
        .PHONY: gommit-setup

        gommit-setup:
        	@echo "Setting up Gommit..."
        ifeq ($(OS),Windows_NT)
        	@powershell -ExecutionPolicy Bypass -File .gommit/gommit-hook-setup.ps1
        else
        	@sh .gommit/gommit-hook-setup.sh
        endif
        ```

   c. If you have an existing setup script:
      - For bash scripts, add the following line to your setup script:
        ```bash
        sh .gommit/gommit-hook-setup.sh
        ```
      - For PowerShell scripts, add the following line to your setup script:
        ```powershell
        & .\.gommit\gommit-hook-setup.ps1
        ```

## Usage

After integration, users can set up Gommit using one of the following methods:

1. Using Make (on systems with Make installed):
   ```
   make gommit-setup
   ```

2. Directly running the setup script:
   - On Unix-like systems (macOS and Linux):
     ```
     sh .gommit/gommit-hook-setup.sh
     ```
   - On Windows (PowerShell):
     ```
     powershell -ExecutionPolicy Bypass -File .gommit/gommit-hook-setup.ps1
     ```

3. If integrated into an existing setup script, users can run your project's normal setup process.

## Make Installation

The Makefile approach requires Make to be installed on the user's system. If it's not already installed, here are instructions for common platforms:

### Linux
- Ubuntu/Debian: `sudo apt-get install make`
- Fedora: `sudo dnf install make`
- Arch Linux: `sudo pacman -S make`

### macOS
- Using Homebrew: `brew install make`
- Using Xcode Command Line Tools: `xcode-select --install`

### Windows
- Using Chocolatey: `choco install make`
- Using Scoop: `scoop install make`
- Alternatively, install Windows Subsystem for Linux (WSL) and follow the Linux instructions.

## Alternative Approaches

If users don't want to install Make or prefer not to use the Makefile:

1. They can use the setup scripts directly:
   - For macOS and Linux: `sh .gommit/gommit-hook-setup.sh`
   - For Windows: `powershell -ExecutionPolicy Bypass -File .gommit/gommit-hook-setup.ps1`

2. You can create custom scripts that call the appropriate setup script based on the detected operating system.

3. Users can integrate the setup commands into their own build or setup processes.

## Notes

- The setup process will download the appropriate Gommit binary for the user's system.
- If a `commit-msg` hook already exists, the user will be prompted to overwrite, append, or skip the Gommit integration.
- The setup ensures that Gommit is not added multiple times to an existing hook.
- The Gommit hook is set up to exit with Gommit's exit code, ensuring proper integration with Git's commit process.
- The `.sh` script is designed for macOS and Linux systems, while the `.ps1` script is for Windows systems.

By following these steps and choosing the appropriate method for their environment, users can easily integrate the Gommit setup process into your project, providing a seamless experience across different platforms.