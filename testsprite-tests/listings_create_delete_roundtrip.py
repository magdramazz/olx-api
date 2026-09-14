import os
import uuid

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90


def list_listings():
    r = requests.get(f"{BASE_URL}/listings", timeout=TIMEOUT)
    assert r.status_code == 200, f"GET /listings: expected 200, got {r.status_code}: {r.text[:200]}"
    return r.json()


def test_create_list_delete_roundtrip():
    title = f"TestSprite fixture {uuid.uuid4()}"
    payload = {
        "title": title,
        "description": "Created by a TestSprite test and deleted at the end of it.",
        "price": "4200",
        "city": "Testville",
    }

    r = requests.post(f"{BASE_URL}/listings", json=payload, timeout=TIMEOUT)
    assert r.status_code == 201, f"POST /listings: expected 201, got {r.status_code}: {r.text[:200]}"
    assert r.headers.get("Content-Type", "").startswith("application/json"), r.headers.get("Content-Type")
    created = r.json()  # the whole body must be a single valid JSON document
    listing_id = created.get("id")
    uuid.UUID(listing_id)

    deleted = False
    try:
        assert created.get("title") == title, f"created listing title mismatch: {created!r}"
        assert created.get("city") == "Testville", f"created listing city mismatch: {created!r}"
        assert int(created.get("price")) == 4200, f"created listing price mismatch: {created!r}"
        assert created.get("created_at"), f"created listing has no created_at: {created!r}"

        match = next((item for item in list_listings() if item["id"] == listing_id), None)
        assert match is not None, f"created listing {listing_id} not returned by GET /listings"
        assert match["title"] == title and match["city"] == "Testville" and int(match["price"]) == 4200, match

        r = requests.delete(f"{BASE_URL}/listings/{listing_id}", timeout=TIMEOUT)
        assert r.status_code == 204, f"DELETE: expected 204, got {r.status_code}: {r.text[:200]}"
        assert r.text == "", f"204 response must have an empty body, got {r.text[:200]!r}"
        deleted = True

        assert all(item["id"] != listing_id for item in list_listings()), "deleted listing is still returned by GET /listings"

        r = requests.delete(f"{BASE_URL}/listings/{listing_id}", timeout=TIMEOUT)
        assert r.status_code == 404, f"second DELETE of the same id: expected 404, got {r.status_code}: {r.text[:200]}"
    finally:
        if not deleted:
            requests.delete(f"{BASE_URL}/listings/{listing_id}", timeout=TIMEOUT)


test_create_list_delete_roundtrip()
