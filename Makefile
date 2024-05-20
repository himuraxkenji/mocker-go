SHELL := /bin/bash

local:
	@echo "Ejecutando la aplicación en modo desarrollo..." 
	@set -a && source ./environment/local.env && set +a && go run cmd/main.go

sample:
	@echo "Ejecutando la aplicación en modo desarrollo..."
	@set -a && source ./environment/sample.env && set +a && go run cmd/main.go
