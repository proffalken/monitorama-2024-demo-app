#!/bin/bash
source ./bin/activate
cd space_weather
export DJANGO_SETTINGS_MODULE=space_weather.settings
export OTEL_RESOURCE_ATTRIBUTES="deployment.environment.name=production,service.namespace=otel_demos,service.version=1.1.0,service.instance.id=$(hostname)-$$"
opentelemetry-instrument --traces_exporter otlp,console \
    --metrics_exporter otlp --service_name space_weather \
    --exporter_otlp_endpoint http://127.0.0.1:4318 \
    --exporter_otlp_protocol http/protobuf \
    python manage.py runserver  0.0.0.0:8889 --noreload
