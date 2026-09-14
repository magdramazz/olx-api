import os

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90
VALID = {"title": "Validation probe", "description": "Should never be stored.", "price": "100", "city": "Testville"}


def cleanup_if_created(r):
    # If validation is broken and the API stored the listing anyway, don't leave it behind.
    if r.status_code == 201:
        try:
            requests.delete(f"{BASE_URL}/listings/{r.json()['id']}", timeout=TIMEOUT)
        except Exception:
            pass


def assert_error(r, status, code, case):
    cleanup_if_created(r)
    assert r.status_code == status, f"{case}: expected {status}, got {r.status_code}: {r.text[:200]}"
    assert r.headers.get("Content-Type", "").startswith("application/json"), f"{case}: {r.headers.get('Content-Type')}"
    err = r.json().get("error")
    assert isinstance(err, dict), f"{case}: expected an error envelope, got {r.text[:200]}"
    assert err.get("code") == code, f"{case}: expected error code {code!r}, got {err.get('code')!r}"
    assert err.get("message"), f"{case}: error message must not be empty"


def test_malformed_json_is_rejected():
    r = requests.post(
        f"{BASE_URL}/listings", data="{not json", headers={"Content-Type": "application/json"}, timeout=TIMEOUT
    )
    assert_error(r, 400, "invalid_body", "malformed JSON")


def test_invalid_fields_are_rejected():
    cases = {
        "empty title": {**VALID, "title": "   "},
        "missing description": {k: v for k, v in VALID.items() if k != "description"},
        "missing city": {k: v for k, v in VALID.items() if k != "city"},
        "missing price": {k: v for k, v in VALID.items() if k != "price"},
        "non-numeric price": {**VALID, "price": "abc"},
        "negative price": {**VALID, "price": "-5"},
    }
    for case, payload in cases.items():
        r = requests.post(f"{BASE_URL}/listings", json=payload, timeout=TIMEOUT)
        assert_error(r, 400, "validation_error", case)


test_malformed_json_is_rejected()
test_invalid_fields_are_rejected()
