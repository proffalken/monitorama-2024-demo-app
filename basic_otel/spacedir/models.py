from django.db import models

# Create your models here.
class Space(models.Model):
    name = models.CharField(max_length=255)
    api_uri = models.CharField(max_length=255)

    def __str__(self):
        return self.name
