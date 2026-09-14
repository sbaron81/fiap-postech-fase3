import os
from unittest.mock import MagicMock, patch

# app.py le essas variaveis e conecta de verdade no PostgreSQL no import
# (nivel de modulo), entao precisam existir e o pool precisa estar mockado
# antes de qualquer teste importar app.
os.environ.setdefault("DATABASE_URL", "postgresql://user:pass@localhost:5432/test")
os.environ.setdefault("AUTH_SERVICE_URL", "http://auth-service:8001")

patch("psycopg2.pool.SimpleConnectionPool", return_value=MagicMock()).start()
