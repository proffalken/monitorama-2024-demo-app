from django.shortcuts import render
from django.http import HttpResponse

import requests

from .models import Space

# Create your views here.
def get_space_status(request):
    spaces = Space.objects.all()
    html = f"<html><head><title>Hackspace Details</title></head><body><table><tr><th>Space Name</th><th>Is Open?</th></tr>"
    for space in spaces:
        space_data = requests.get(space.api_uri).json()
        if space_data['state']['open'] == True:
            current_space_state = "YES"
        else:
            current_space_state = "NO"
        html = html + f"<tr><td>{space.name}</td><td>{current_space_state}</td></tr>"
    html = html + "</table></body></html>"
    return HttpResponse(html)
