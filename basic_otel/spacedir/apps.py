from django.apps import AppConfig
from django.conf import settings

_DB_INSTRUMENTED = False


class SpacedirConfig(AppConfig):
    default_auto_field = 'django.db.models.BigAutoField'
    name = 'spacedir'

    def ready(self):
        global _DB_INSTRUMENTED
        if _DB_INSTRUMENTED:
            return

        db_engine = settings.DATABASES.get("default", {}).get("ENGINE", "")
        if "postgresql" not in db_engine:
            return

        try:
            from opentelemetry.instrumentation.psycopg2 import Psycopg2Instrumentor

            Psycopg2Instrumentor().instrument(enable_commenter=True)
            _DB_INSTRUMENTED = True
        except Exception:
            # Keep startup resilient if OTEL instrumentation isn't available.
            return
