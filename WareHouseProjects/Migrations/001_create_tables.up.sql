CREATE TABLE "branches" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "address" varchar,
  "phone" varchar,
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp
);

CREATE TABLE "category" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "parent_id" UUID REFERENCES "category"("id"),
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMP
);

CREATE TABLE "product" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar UNIQUE NOT NULL,
  "category_id" uuid REFERENCES "category"("id"),
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp
);

CREATE TABLE "coming_table" (
  "id" uuid PRIMARY KEY,
  "coming_id" varchar NOT NULL,
  "branch_id" uuid REFERENCES "branches"("id"),
  "date_time" timestamp,
  "status" varchar DEFAULT 'in_process',
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp
);

CREATE TABLE "coming_table_product" (
  "id" uuid PRIMARY KEY,
  "category_id" uuid REFERENCES "category"("id"),
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar UNIQUE NOT NULL,
  "count" numeric NOT NULL DEFAULT 0,
  "total_price" numeric DEFAULT 0,
  "coming_table_id" uuid REFERENCES "coming_table"("id"),
  "created_at" timestamp DEFAULT current_timestamp,
  "updated_at" timestamp
);

CREATE TABLE "remaining" (
  "id" uuid PRIMARY KEY,
  "branch_id" uuid REFERENCES "branches"("id"),
  "category_id" uuid REFERENCES "category"("id"),
  "name" varchar NOT NULL,
  "price" numeric NOT NULL,
  "barcode" varchar UNIQUE NOT NULL,
  "count" numeric NOT NULL DEFAULT 0,
  "total_price" numeric DEFAULT 0,
  "created_at" timestamp DEFAULT (current_timestamp),
  "updated_at" timestamp
);

ALTER TABLE "category" ADD FOREIGN KEY ("parent_id") REFERENCES "category" ("id");

ALTER TABLE "product" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "coming_table" ADD FOREIGN KEY ("branch_id") REFERENCES "branch" ("id");

ALTER TABLE "coming_table_product" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");

ALTER TABLE "coming_table_product" ADD FOREIGN KEY ("coming_table_id") REFERENCES "coming_table" ("id");

ALTER TABLE "remaining" ADD FOREIGN KEY ("branch_id") REFERENCES "branch" ("id");

ALTER TABLE "remaining" ADD FOREIGN KEY ("category_id") REFERENCES "category" ("id");