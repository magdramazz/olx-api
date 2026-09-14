import os
import re
import uuid
from datetime import datetime

import requests

BASE_URL = (globals().get("TARGET_URL") or os.environ.get("TARGET_URL") or "https://olx-api-hdtz.onrender.com").rstrip("/")
TIMEOUT = 90
REQUIRED_FIELDS = {"id", "title", "description", "price", "city", "created_at"}


def parse_rfc3339(value):
    # Go emits RFC3339Nano (variable-length fraction, "Z" suffix); normalise for fromisoformat.
    value = value.replace("Z", "+00:00")
    m = re.match(r"^(.*T\d{2}:\d{2}:\d{2})(?:\.(\d+))?(.*)$", value)
    assert m, f"created_at is not RFC3339: {value!r}"
    head, frac, tz = m.groups()
    frac = ((frac or "") + "000000")[:6]
    return datetime.fromisoformat(f"{head}.{frac}{tz}")


def test_list_listings_returns_well_formed_array_newest_first():
    r = requests.get(f"{BASE_URL}/listings", timeout=TIMEOUT)
    assert r.status_code == 200, f"expected 200, got {r.status_code}: {r.text[:200]}"
    assert r.headers.get("Content-Type", "").startswith("application/json"), r.headers.get("Content-Type")

    body = r.json()
    assert isinstance(body, list), f"expected a JSON array, got {type(body).__name__}: {r.text[:200]}"
    assert len(body) <= 100, f"list must be capped at 100 items, got {len(body)}"

    ids = []
    timestamps = []
    for item in body:
        missing = REQUIRED_FIELDS - set(item)
        assert not missing, f"listing {item.get('id')!r} is missing fields {sorted(missing)}"
        uuid.UUID(item["id"])
        assert isinstance(item["title"], str) and item["title"].strip(), f"empty title: {item!r}"
        assert isinstance(item["city"], str) and item["city"].strip(), f"empty city: {item!r}"
        assert int(item["price"]) >= 0, f"price must be a non-negative integer: {item['price']!r}"
        ids.append(item["id"])
        timestamps.append(parse_rfc3339(item["created_at"]))

    assert len(ids) == len(set(ids)), "listing ids must be unique"
    assert timestamps == sorted(timestamps, reverse=True), "listings must be ordered by created_at, newest first"


test_list_listings_returns_well_formed_array_newest_first()
