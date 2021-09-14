CREATE TABLE IF NOT EXISTS "account"  (
   "id" bigserial PRIMARY KEY,
   "owner" varchar NOT NULL,
   "balance" bigint NOT NULL,
   "currency" varchar NOT NULL,
   "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "account" ("owner");