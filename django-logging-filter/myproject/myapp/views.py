from django.shortcuts import render

# Create your views here.
# myapp/views.py
from django.http import JsonResponse
from django.views import View

class HealthCheckView(View):
    def get(self, request, *args, **kwargs):
        return JsonResponse({"status": "healthy"})
