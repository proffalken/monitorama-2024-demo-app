from django.http import HttpResponse, JsonResponse

from opentelemetry import trace

tracer = trace.get_tracer(__name__)

import os
import requests

from .models import Space

WEATHER_SERVICE_URL = os.getenv("WEATHER_SERVICE_URL", "http://localhost:8889/weather/")

# Create your views here.
def get_space_status(request):
    spaces = Space.objects.all()
    html = f"<html><head><title>Hackspace Details</title></head><body><table><tr><th>Space Name</th><th>Is Open?</th></tr>"
    for space in spaces:
        with tracer.start_as_current_span("space.lookup", kind=trace.SpanKind.INTERNAL) as span:
            span.set_attribute("space.name", space.name)
            space_data = requests.get(space.api_uri, timeout=10).json()
            if space_data["state"]["open"]:
                current_space_state = "YES"
            else:
                current_space_state = "NO"
            html = html + f"<tr><td>{space.name}</td><td>{current_space_state}</td></tr>"
    html = html + "</table></body></html>"
    return HttpResponse(html)

def get_space_status_json(request):
    spaces = Space.objects.all()
    space_return_data = []
    for space in spaces:
        with tracer.start_as_current_span("space.lookup", kind=trace.SpanKind.INTERNAL) as span:
            span.set_attribute("space.name", space.name)
            space_data = requests.get(space.api_uri, timeout=10).json()
            if space_data["state"]["open"]:
                current_space_state = "YES"
            else:
                current_space_state = "NO"
            with tracer.start_as_current_span("weather.lookup", kind=trace.SpanKind.INTERNAL):
                weather_data = requests.get(
                    f"{WEATHER_SERVICE_URL}?lat={space_data['location']['lat']}&lng={space_data['location']['lon']}",
                    timeout=10,
                ).json()
            space_return_data.append(
                    {"name": space.name, 
                     "is_open": current_space_state,
                     "weather": {
                         "temperature": weather_data["current"]["temp_c"],
                         "condition": weather_data["current"]["condition"]["text"],
                         "humidity": weather_data["current"]["humidity"]
                         }
                    })

    results = {"data": space_return_data}
    return JsonResponse(results)
