
-- DROP SCHEMA `h3i_kb_db`;
-- CREATE SCHEMA `h3i_kb_db`;

CREATE TABLE IF NOT EXISTS "products" (
  "id" int,
  "code" varchar(25) UNIQUE NOT NULL,
  "name" varchar(50) NOT NULL,
  "auth_user" varchar(50),
  "auth_pass" varchar(50),
  "price" float(5) DEFAULT 0,
  "renewal_day" int DEFAULT 0,
  "url_notif_sub" varchar(75),
  "url_notif_unsub" varchar(75),
  "url_notif_renewal" varchar(75),
  "url_postback" varchar(75),
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "contents" (
  "id" int,
  "product_id" int NOT NULL,
  "name" varchar(20) NOT NULL,
  "value" varchar(250) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "adnets" (
  "id" int,
  "name" varchar(20) NOT NULL,
  "value" varchar(20) NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "schedules" (
  "id" int,
  "name" varchar(20) UNIQUE NOT NULL,
  "publish_at" timestamp,
  "unlocked_at" timestamp,
  "is_unlocked" bool DEFAULT false,
  PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "subscriptions" (
  "id" SERIAL PRIMARY KEY,
  "product_id" int NOT NULL,
  "msisdn" varchar(20) NOT NULL,
  "adnet" varchar(20),
  "latest_subject" varchar(20),
  "latest_status" varchar(20),
  "renewal_at" timestamp,
  "purge_at" timestamp,
  "unsub_at" timestamp,
  "charge_at" timestamp,
  "retry_at" timestamp,
  "success_firstpush" int DEFAULT 0,
  "success_renewal" int DEFAULT 0,
  "total_success" int DEFAULT 0,
  "total_firstpush" int DEFAULT 0,
  "total_renewal" int DEFAULT 0,
  "total_amount" float(8) DEFAULT 0,
  "is_retry" bool DEFAULT false,
  "is_purge" bool DEFAULT false,
  "is_trial" bool DEFAULT false,
  "is_suspend" bool DEFAULT false,
  "is_active" bool DEFAULT false,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "transactions" (
  "id" SERIAL PRIMARY KEY,
  "transaction_id" varchar(20),
  "product_id" int NOT NULL,
  "msisdn" varchar(20) NOT NULL,
  "adnet" varchar(20),
  "amount" float(5) DEFAULT 0,
  "subject" varchar(20),
  "status" varchar(20),
  "status_detail" varchar(20),
  "payload" text,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "histories" (
  "id" SERIAL PRIMARY KEY,
  "product_id" int,
  "msisdn" varchar(20),
  "adnet" varchar(20),
  "subject" varchar(20),
  "payload" text,
  "created_at" timestamp
);

CREATE TABLE IF NOT EXISTS "blacklists" (
  "id" SERIAL PRIMARY KEY,
  "msisdn" varchar(20) UNIQUE NOT NULL,
  "created_at" timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_msisdn" ON "blacklists" ("msisdn");
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_product_msisdn" ON "subscriptions" ("product_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_product_msisdn" ON "transactions" ("product_id", "msisdn");

ALTER TABLE "contents" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "transactions" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");
ALTER TABLE "histories" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("id");