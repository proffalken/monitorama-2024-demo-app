#!/bin/bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"
export OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"
export OTEL_SERVICE_NAME="otel-basic-go-client"
export OTEL_RESOURCE_ATTRIBUTES="deployment.environment.name=production,service.namespace=oteldemostack,service.version=1.1.0,service.instance.id=$(hostname)-$$"
cd goclient
go run .
