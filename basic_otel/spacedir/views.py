from django.shortcuts import render
from django.http import HttpResponse
import requests

# Create your views here.
def get_space_status(request):
    space_data = requests.get("https://members.makemonmouth.co.uk/api/spacedirectory/").json()
    if space_data['state']['open'] == True:
        current_space_state = "YES"
    else:
        current_space_state = "NO"
    html = f"<html><head><title>Is Make Monmouth Open?</title></head><body><h1>Is Make Monmouth Open?</h1><h2>{current_space_state}</h2></body></html>"
    return HttpResponse(html)
