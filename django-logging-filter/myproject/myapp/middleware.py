# myapp/middleware.py
import logging

class LoggingMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response
        self.logger = logging.getLogger('django')

    def __call__(self, request):
        request_info_str = self.get_request_info(request)
        self.logger.info(request_info_str)
        response = self.get_response(request)
        return response

    def get_request_info(self, request):
        return f"{request.method} {request.get_full_path()} {request.META.get('HTTP_USER_AGENT', '')}"

class SuppressHealthCheckMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        # スーパークラスのメソッドを呼び出す
        response = self.get_response(request)

        # ヘルスチェックのリクエストに対して、サーバーログを抑制する
        if request.path == "/healthcheck/":
            old_print = print

            def new_print(*args, **kwargs):
                if 'GET /healthcheck/' not in args[0]:
                    old_print(*args, **kwargs)

            print = new_print

        return response
