-- Create "feeds" table
CREATE TABLE "feeds" (
  "id" uuid NOT NULL,
  "url" character varying NOT NULL,
  "title" character varying NOT NULL,
  "created_at" timestamptz NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "feeds_url_key" to table: "feeds"
CREATE UNIQUE INDEX "feeds_url_key" ON "feeds" ("url");
