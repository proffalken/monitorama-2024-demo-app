#!/bin/bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
cd goclient
go run .
