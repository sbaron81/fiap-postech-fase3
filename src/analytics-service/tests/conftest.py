import os

# app.py le essas variaveis no import (nivel de modulo) e faz sys.exit(1) se
# alguma faltar, entao precisam existir antes de qualquer teste importar app.
os.environ.setdefault("AWS_REGION", "us-east-1")
os.environ.setdefault("AWS_SQS_URL", "https://sqs.us-east-1.amazonaws.com/123456789012/test-queue")
os.environ.setdefault("AWS_DYNAMODB_TABLE", "test-table")
