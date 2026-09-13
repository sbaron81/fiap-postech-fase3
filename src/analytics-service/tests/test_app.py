import json
from unittest.mock import patch

import app


def make_message(body, message_id="msg-1", receipt_handle="receipt-1"):
    return {
        "MessageId": message_id,
        "ReceiptHandle": receipt_handle,
        "Body": json.dumps(body),
    }


def test_process_message_saves_valid_event_and_deletes_from_queue():
    message = make_message(
        {
            "user_id": "user-123",
            "flag_name": "enable-new-dashboard",
            "result": True,
            "timestamp": "2026-01-01T00:00:00Z",
        }
    )

    with patch.object(app, "dynamodb_client") as mock_dynamodb, patch.object(app, "sqs_client") as mock_sqs:
        app.process_message(message)

    mock_dynamodb.put_item.assert_called_once()
    call_kwargs = mock_dynamodb.put_item.call_args.kwargs
    assert call_kwargs["TableName"] == app.DYNAMODB_TABLE_NAME

    item = call_kwargs["Item"]
    assert item["user_id"] == {"S": "user-123"}
    assert item["flag_name"] == {"S": "enable-new-dashboard"}
    assert item["result"] == {"BOOL": True}
    assert item["timestamp"] == {"S": "2026-01-01T00:00:00Z"}
    assert "event_id" in item

    mock_sqs.delete_message.assert_called_once_with(
        QueueUrl=app.SQS_QUEUE_URL,
        ReceiptHandle="receipt-1",
    )


def test_process_message_with_invalid_json_does_not_ack():
    message = {
        "MessageId": "msg-2",
        "ReceiptHandle": "receipt-2",
        "Body": "{not-valid-json",
    }

    with patch.object(app, "dynamodb_client") as mock_dynamodb, patch.object(app, "sqs_client") as mock_sqs:
        app.process_message(message)

    # Mensagem malformada nao deve ser gravada nem removida da fila (vira poison pill).
    mock_dynamodb.put_item.assert_not_called()
    mock_sqs.delete_message.assert_not_called()


def test_process_message_with_missing_field_does_not_ack():
    message = make_message({"user_id": "user-123"})  # falta flag_name/result/timestamp

    with patch.object(app, "dynamodb_client") as mock_dynamodb, patch.object(app, "sqs_client") as mock_sqs:
        app.process_message(message)

    # Sem os campos obrigatorios, a mensagem deve continuar na fila para nova tentativa.
    mock_dynamodb.put_item.assert_not_called()
    mock_sqs.delete_message.assert_not_called()


def test_health_endpoint_returns_ok():
    client = app.app.test_client()

    response = client.get("/health")

    assert response.status_code == 200
    assert response.get_json() == {"status": "ok"}
