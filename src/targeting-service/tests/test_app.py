from unittest.mock import MagicMock, patch

import requests

import app


def make_client():
    app.app.config["TESTING"] = True
    return app.app.test_client()


def mock_auth_ok():
    """Simula o auth-service aprovando a chave (GET /validate -> 200)."""
    resp = MagicMock()
    resp.status_code = 200
    return patch.object(app.requests, "get", return_value=resp)


def test_health_returns_ok():
    client = make_client()

    resp = client.get("/health")

    assert resp.status_code == 200
    assert resp.get_json() == {"status": "ok"}


def test_create_rule_without_auth_header_returns_401():
    client = make_client()

    resp = client.post("/rules", json={"flag_name": "minha-flag", "rules": {"type": "PERCENTAGE", "value": 50}})

    assert resp.status_code == 401


def test_create_rule_when_auth_service_unreachable_returns_503():
    client = make_client()

    with patch.object(app.requests, "get", side_effect=requests.exceptions.ConnectionError):
        resp = client.post(
            "/rules",
            json={"flag_name": "minha-flag", "rules": {"type": "PERCENTAGE", "value": 50}},
            headers={"Authorization": "Bearer alguma-chave"},
        )

    assert resp.status_code == 503


def test_get_rule_not_found_returns_404():
    client = make_client()

    mock_cursor = MagicMock()
    mock_cursor.fetchone.return_value = None
    mock_conn = MagicMock()
    mock_conn.cursor.return_value = mock_cursor

    with mock_auth_ok(), patch.object(app.pool, "getconn", return_value=mock_conn):
        resp = client.get("/rules/flag-inexistente", headers={"Authorization": "Bearer chave-valida"})

    assert resp.status_code == 404


def test_get_rule_found_returns_200():
    client = make_client()

    mock_cursor = MagicMock()
    mock_cursor.fetchone.return_value = {
        "id": 1,
        "flag_name": "minha-flag",
        "is_enabled": True,
        "rules": {"type": "PERCENTAGE", "value": 50},
    }
    mock_conn = MagicMock()
    mock_conn.cursor.return_value = mock_cursor

    with mock_auth_ok(), patch.object(app.pool, "getconn", return_value=mock_conn):
        resp = client.get("/rules/minha-flag", headers={"Authorization": "Bearer chave-valida"})

    assert resp.status_code == 200
    assert resp.get_json()["flag_name"] == "minha-flag"
