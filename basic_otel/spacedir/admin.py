from django.contrib import admin

# Register your models here.
from .models import Space


@admin.register(Space)
class SpaceAdmin(admin.ModelAdmin):
    pass
