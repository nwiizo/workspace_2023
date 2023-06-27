# myapp/logging_filters.py
import logging

class HealthCheckFilter(logging.Filter):
    def filter(self, record):
        return "/healthcheck/" not in record.getMessage()

class KubeProbeFilter(logging.Filter):
    def filter(self, record):
        return 'kube-probe' not in record.getMessage()
