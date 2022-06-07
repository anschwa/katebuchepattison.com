# GNUmakefile
SHELL := bash
.ONESHELL:
.DELETE_ON_ERROR:
.SHELLFLAGS := -eu -o pipefail -c
MAKEFLAGS += --warn-undefined-variables

.PHONY: help
help: ## Display this help section
	@echo -e "Usage: make <command>\n\nAvailable commands are:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "  %-38s %s\n", $$1, $$2}' ${MAKEFILE_LIST}
.DEFAULT_GOAL := help

.PHONY: build
build: ## Build website for production
	npm run build-prod
	gutenblog -outDir=docs build .

.PHONY: serve
serve: ## Start development server with tailwindcss
	npm run build-dev
	gutenblog -outDir=docs serve .
