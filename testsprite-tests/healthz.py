import os

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90


def test_healthz_returns_ok():
    r = requests.get(f"{BASE_URL}/healthz", timeout=TIMEOUT)
    assert r.status_code == 200, f"expected 200, got {r.status_code}: {r.text[:200]}"
    assert r.headers.get("Content-Type", "").startswith("application/json"), r.headers.get("Content-Type")
    body = r.json()
    assert isinstance(body, dict), f"expected a JSON object, got {body!r}"
    assert str(body.get("status", "")).startswith("ok"), f"unexpected status field: {body!r}"


test_healthz_returns_ok()
