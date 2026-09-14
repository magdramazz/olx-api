import os
import uuid

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90


def test_request_id_is_echoed_when_supplied():
    rid = f"testsprite-{uuid.uuid4()}"
    r = requests.get(f"{BASE_URL}/healthz", headers={"X-Request-ID": rid}, timeout=TIMEOUT)
    assert r.status_code == 200, f"expected 200, got {r.status_code}"
    assert r.headers.get("X-Request-ID") == rid, f"expected X-Request-ID {rid!r}, got {r.headers.get('X-Request-ID')!r}"


def test_request_id_is_generated_when_absent():
    first = requests.get(f"{BASE_URL}/healthz", timeout=TIMEOUT).headers.get("X-Request-ID")
    second = requests.get(f"{BASE_URL}/healthz", timeout=TIMEOUT).headers.get("X-Request-ID")
    assert first and second, f"missing X-Request-ID header: {first!r}, {second!r}"
    uuid.UUID(first)
    uuid.UUID(second)
    assert first != second, "generated request ids must be unique per request"


test_request_id_is_echoed_when_supplied()
test_request_id_is_generated_when_absent()
