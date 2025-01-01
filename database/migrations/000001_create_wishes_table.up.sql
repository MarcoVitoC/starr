BEGIN;
CREATE TABLE IF NOT EXISTS "wishes"(
	"id" UUID PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"description" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP,
	"updated_at" TIMESTAMP
);
COMMIT;