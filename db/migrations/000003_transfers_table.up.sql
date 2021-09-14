CREATE TABLE IF NOT EXISTS "transfers" (
     "id" bigserial PRIMARY KEY,
     "from_account" bigint NOT NULL,
     "to_account" bigint NOT NULL,
     "amount" bigint NOT NULL,
     "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account") REFERENCES "account" ("id");
ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account") REFERENCES "account" ("id");

CREATE INDEX ON "transfers" ("from_account");
CREATE INDEX ON "transfers" ("to_account");
CREATE INDEX ON "transfers" ("from_account", "to_account");

COMMENT ON COLUMN "transfers"."amount" IS 'must be possitive';