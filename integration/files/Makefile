# Gommit setup
.PHONY: setup-gommit

setup-gommit:
	@echo "Setting up Gommit..."
ifeq ($(OS),Windows_NT)
	@powershell -ExecutionPolicy Bypass -File .gommit/setup.ps1
else
	@sh .gommit/setup.sh
endif

# Include this in your existing Makefile or use it as a standalone Makefile