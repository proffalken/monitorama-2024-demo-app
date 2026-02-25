from django.http import JsonResponse

from opentelemetry import trace

tracer = trace.get_tracer(__name__)

import os
import requests

WEATHER_API_URL = os.getenv("WEATHER_API_URL", "https://api.weatherapi.com/v1/current.json")
WEATHER_API_KEY = os.getenv("WEATHER_API_KEY", "50ffe0f2fd1b46ca9b1105115240502")

# Create your views here.
def get_weather(request):
    # https://api.weatherapi.com/v1/current.json?key=50ffe0f2fd1b46ca9b1105115240502&q=51.804536,-2.698146&aqi=no
    with tracer.start_as_current_span("weather.fetch", kind=trace.SpanKind.INTERNAL):
        weather_info = requests.get(
            f"{WEATHER_API_URL}?key={WEATHER_API_KEY}&q={request.GET.get('lat')},{request.GET.get('lng')}&aqi=no",
            timeout=10,
        ).json()

        return JsonResponse(weather_info)
