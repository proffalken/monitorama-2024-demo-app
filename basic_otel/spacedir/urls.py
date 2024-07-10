from django.urls import path

from . import views

urlpatterns = [
        path("space_state", views.get_space_status),
        path("space_json", views.get_space_status_json),
        ]
