
-- DROP SCHEMA `pass_tsel`;
-- CREATE SCHEMA `pass_tsel`;

CREATE TABLE IF NOT EXISTS "services" (
  "id" SERIAL PRIMARY KEY,
  "category" varchar(20) NOT NULL,
  "code" varchar(25) UNIQUE NOT NULL,
  "name" varchar(50) NOT NULL,
  "price" float(5) DEFAULT 0,
  "program_id"  varchar(25),
  "sid" varchar(35),
  "renewal_day" int DEFAULT 0,
  "trial_day" int DEFAULT 0,
  "url_telco" varchar(85),
  "url_portal" varchar(85),
  "url_callback" varchar(85),
  "url_notif_sub" varchar(85),
  "url_notif_unsub" varchar(85),
  "url_notif_renewal" varchar(85),
  "url_postback" varchar(85)
);

CREATE TABLE IF NOT EXISTS "contents" (
  "id" SERIAL PRIMARY KEY,
  "service_id" int NOT NULL,
  "name" varchar(20) NOT NULL,
  "value" varchar(400) NOT NULL,
  "tid" varchar(5) NOT NULL
);

CREATE TABLE IF NOT EXISTS "adnets" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar(20) NOT NULL,
  "value" varchar(20) NOT NULL
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
  "service_id" int NOT NULL,
  "category" varchar(20) NOT NULL,
  "msisdn" varchar(20) NOT NULL,
  "channel" varchar(20),
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "latest_trxid" varchar(100),
  "latest_keyword" varchar(180),
  "latest_subject" varchar(20),
  "latest_status" varchar(20),
  "latest_pin" int DEFAULT 0,
  "amount" float(8) DEFAULT 0,
  "trial_at" timestamp,
  "renewal_at" timestamp,
  "unsub_at" timestamp,
  "charge_at" timestamp,
  "retry_at" timestamp,
  "success" int DEFAULT 0,
  "ip_address" varchar(50),
  "is_trial" bool DEFAULT false,
  "is_retry" bool DEFAULT false,
  "is_active" bool DEFAULT false,
  "charging_count" int DEFAULT 0,
  "charging_count_all" int DEFAULT 0,
  "total_firstpush" int DEFAULT 0,
  "total_renewal" int DEFAULT 0,
  "total_amount_firstpush" float(8) DEFAULT 0,
  "total_amount_renewal" float(8) DEFAULT 0,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "transactions" (
  "id" SERIAL PRIMARY KEY,
  "tx_id" varchar(100),
  "service_id" int NOT NULL,
  "msisdn" varchar(20) NOT NULL,
  "channel" varchar(20) NOT NULL,
  "adnet" varchar(55),
  "pub_id" varchar(55),
  "aff_sub" varchar(100),
  "keyword" varchar(180),
  "amount" float(8) DEFAULT 0,
  "pin" int DEFAULT 0,
  "status" varchar(45),
  "status_code" varchar(45),
  "status_detail" varchar(100),
  "subject" varchar(45),
  "ip_address" varchar(45),
  "payload" text,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE TABLE IF NOT EXISTS "histories" (
  "id" SERIAL PRIMARY KEY,
  "service_id" int NOT NULL,
  "msisdn" varchar(20),
  "channel" varchar(20),
  "adnet" varchar(55),
  "keyword" varchar(180),
  "subject" varchar(20),
  "status" varchar(45),
  "ip_address" varchar(45),
  "created_at" timestamp
);

CREATE TABLE IF NOT EXISTS "blacklists" (
  "id" SERIAL PRIMARY KEY,
  "msisdn" varchar(20) UNIQUE NOT NULL,
  "created_at" timestamp
);

CREATE TABLE IF NOT EXISTS "vips" (
  "id" SERIAL PRIMARY KEY,
  "msisdn" varchar(20) UNIQUE NOT NULL,
  "created_at" timestamp
);

CREATE TABLE IF NOT EXISTS "campaigns" (
  "id" int,
  "service_id" int NOT NULL,
  "adnet" varchar(20) NOT NULL,
  "limit_mo" int DEFAULT 0,
  "total_mo" int DEFAULT 0,
  "created_at" timestamp,
  "updated_at" timestamp,
  PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX IF NOT EXISTS "uidx_msisdn" ON "blacklists" ("msisdn");
CREATE UNIQUE INDEX IF NOT EXISTS "uidx_service_msisdn" ON "subscriptions" ("service_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_category_msisdn" ON "subscriptions" ("category", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_service_msisdn" ON "transactions" ("service_id", "msisdn");
CREATE INDEX IF NOT EXISTS "idx_service_adnet" ON "campaigns" ("service_id", "adnet");
CREATE INDEX IF NOT EXISTS "idx_service_name" ON "contents" ("service_id", "name");
CREATE INDEX IF NOT EXISTS "idx_name_publish_at" ON "schedules" ("name", "publish_at");

CREATE INDEX IF NOT EXISTS "idx_name_publish_at" ON "transactions_" ("created_at", "publish_at");

ALTER TABLE "contents" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "subscriptions" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "transactions" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "histories" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");
ALTER TABLE "campaigns" ADD FOREIGN KEY ("service_id") REFERENCES "services" ("id");

/**
(SELECT 'service_id', 'msisdn', 'channel', 'adnet', 'pub_id', 'aff_sub', 'latest_subject', 'latest_status', 'amount', 'trial_at', 'renewal_at', 'unsub_at', 'charge_at', 'retry_at','success', 'ip_address', 'is_trial', 'is_retry', 'is_active', 'total_firstpush', 'total_renewal', 'total_amount_firstpush', 'total_amount_renewal', 'created_at', 'updated_at')  UNION (SELECT service_id, msisdn, channel, adnet, pub_id, aff_sub, latest_subject, latest_status, amount, trial_at, renewal_at, unsub_at, charge_at, retry_at, success, ip_address, is_trial, is_retry, is_active, total_firstpush, total_renewal, total_amount_firstpush, total_amount_renewal, created_at, updated_at INTO OUTFILE '/tmp/subscriptions.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.subscriptions);
**/

/**
(SELECT 'id', 'service_id', 'name', 'value') UNION (SELECT id, service_id, name, value INTO OUTFILE '/tmp/contents.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.contents);
**/

/**
(SELECT 'id', 'service_id', 'adnet', 'limit_mo', 'total_mo', 'created_at', 'updated_at') UNION (SELECT id, service_id, adnet, limit_mo, total_mo, created_at, updated_at INTO OUTFILE '/tmp/campaigns.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.campaigns);
**/

/**
(SELECT 'id', 'msisdn', 'created_at') UNION (SELECT id, msisdn, created_at INTO OUTFILE '/tmp/blacklists.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.blacklists);
**/

/**
(SELECT 'id', 'msisdn', 'created_at') UNION (SELECT id, msisdn, created_at INTO OUTFILE '/tmp/vips.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.vips);
**/

/**
(SELECT 'id', 'name', 'publish_at', 'unlocked_at', 'is_unlocked') UNION (SELECT id, name, publish_at, locked_at, status INTO OUTFILE '/tmp/schedules.csv' FIELDS TERMINATED BY ';' ENCLOSED BY '"' LINES TERMINATED BY '\n' FROM telenor.schedules);
**/