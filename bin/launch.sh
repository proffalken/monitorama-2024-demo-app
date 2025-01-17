#!/bin/bash
source ./bin/activate
cd basic_otel
export DJANGO_SETTINGS_MODULE=basic_otel.settings
export OTEL_RESOURCE_ATTRIBUTES=deployment.environment=Production,service.namespace=otel_demos,service.version=1.0,service.instance.id=1
opentelemetry-instrument --traces_exporter otlp,console \
    --metrics_exporter otlp --service_name space_api_lookup \
    --exporter_otlp_endpoint http://127.0.0.1:4318 \
    --exporter_otlp_protocol http/protobuf \
    python manage.py runserver  0.0.0.0:8898 --noreload
