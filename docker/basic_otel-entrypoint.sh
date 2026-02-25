#!/usr/bin/env sh
set -eu

cd /app/basic_otel
python manage.py migrate --noinput

exec opentelemetry-instrument \
  --traces_exporter otlp \
  --metrics_exporter otlp \
  --service_name "${OTEL_SERVICE_NAME:-space-api-lookup}" \
  --exporter_otlp_endpoint "${OTEL_EXPORTER_OTLP_ENDPOINT:-http://otel-collector:4318}" \
  --exporter_otlp_protocol "${OTEL_EXPORTER_OTLP_PROTOCOL:-http/protobuf}" \
  python manage.py runserver 0.0.0.0:8898 --noreload
