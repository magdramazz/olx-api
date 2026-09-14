import os
import uuid

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90


def assert_error(r, status, code, case):
    assert r.status_code == status, f"{case}: expected {status}, got {r.status_code}: {r.text[:200]}"
    assert r.headers.get("Content-Type", "").startswith("application/json"), f"{case}: {r.headers.get('Content-Type')}"
    err = r.json().get("error")
    assert isinstance(err, dict), f"{case}: expected an error envelope, got {r.text[:200]}"
    assert err.get("code") == code, f"{case}: expected error code {code!r}, got {err.get('code')!r}"
    assert err.get("message"), f"{case}: error message must not be empty"


def test_delete_nonexistent_listing_returns_404():
    missing_id = uuid.uuid4()
    r = requests.delete(f"{BASE_URL}/listings/{missing_id}", timeout=TIMEOUT)
    assert_error(r, 404, "not_found", f"DELETE unknown id {missing_id}")


def test_delete_malformed_id_returns_400():
    for bad_id in ["not-a-uuid", "12345", "00000000-0000-0000-0000-00000000000"]:
        r = requests.delete(f"{BASE_URL}/listings/{bad_id}", timeout=TIMEOUT)
        assert_error(r, 400, "invalid_id", f"DELETE malformed id {bad_id!r}")


test_delete_nonexistent_listing_returns_404()
test_delete_malformed_id_returns_400()
