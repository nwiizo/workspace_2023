# myapp/urls.py
from django.urls import path

from .views import HealthCheckView

urlpatterns = [
    path('healthcheck/', HealthCheckView.as_view(), name='healthcheck'),
]
