from django.shortcuts import render
from django.http import HttpResponse, JsonResponse

from opentelemetry import trace

tracer = trace.get_tracer(__name__)

import json
import requests

from .models import Space

# Create your views here.
def get_space_status(request):
    spaces = Space.objects.all()
    html = f"<html><head><title>Hackspace Details</title></head><body><table><tr><th>Space Name</th><th>Is Open?</th></tr>"
    for space in spaces:
        with tracer.start_as_current_span("get_space_data", kind=trace.SpanKind.CLIENT) as span:
            span.set_attribute("peer.service", f"{space.name}-space-api")
            space_data = requests.get(space.api_uri).json()
            if space_data['state']['open'] == True:
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
        with tracer.start_as_current_span("get_space_data", kind=trace.SpanKind.CLIENT) as span:
            span.set_attribute("peer.service", f"{space.name}-space-api")
            space_data = requests.get(space.api_uri).json()
            if space_data['state']['open'] == True:
                current_space_state = "YES"
            else:
                current_space_state = "NO"
            space_return_data.append({"name": space.name, "is_open": current_space_state})

    results = {"data": space_return_data}
    return JsonResponse(results)
