CREATE TABLE "events"
(
    "id"          UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    "action"      TEXT        NOT NULL,
    "product"     TEXT        NOT NULL,
    "fingerprint" TEXT        NOT NULL,
    "created"     timestamptz NOT NULL DEFAULT current_timestamp
);