from django.shortcuts import render
from django.http import HttpResponse, JsonResponse

from opentelemetry import trace

tracer = trace.get_tracer(__name__)

import json
import requests

# Create your views here.
def get_weather(request):
    # https://api.weatherapi.com/v1/current.json?key=50ffe0f2fd1b46ca9b1105115240502&q=51.804536,-2.698146&aqi=no
    uri = "https://api.weatherapi.com/v1/current.json"
    key = "50ffe0f2fd1b46ca9b1105115240502"
    with tracer.start_as_current_span("get_weather", kind=trace.SpanKind.CLIENT) as span:
        span.set_attribute("peer.service", "weatherapi.com")
        weather_info = requests.get(f"{uri}?key={key}&q={request.GET.get('lat')},{request.GET.get('lng')}&aqi=no").json()

        return JsonResponse(weather_info)

