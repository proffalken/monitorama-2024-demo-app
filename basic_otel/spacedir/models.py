from django.db import models

# Create your models here.
class Space(models.Model):
    name = models.CharField()
    api_uri = models.CharField()

    def __str__(self):
        return self.name
