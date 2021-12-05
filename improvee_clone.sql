CREATE TYPE "routine_repeat_time" AS ENUM (
  'daily',
  'weekly',
  'monthly'
);

CREATE TYPE "routine_schedule" AS ENUM (
  'anytime',
  'morning',
  'afternoon',
  'nighttime'
);

CREATE TABLE "Users" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "firstName" varchar,
  "lastName" varchar,
  "email" varchar,
  "password" varchar,
  "created_at" timestamp
);

CREATE TABLE "sessions" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "name" varchar,
  "email" varchar,
  "user_id" int,
  "created_at" timestamp
);

CREATE TABLE "Tasks" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "title" varchar,
  "completed" boolean,
  "deadline" timestamp,
  "colorTag" varchar,
  "created_at" timestamp,
  "user_id" int
);

CREATE TABLE "SubTasks" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "title" varchar,
  "completed" boolean,
  "created_at" timestamp,
  "task_id" int
);

CREATE TABLE "Routines" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "title" varchar,
  "colorTag" varchar,
  "daysDone" timestamp,
  "status" boolean,
  "repeatTime" routine_repeat_time,
  "schedule" routine_schedule,
  "created_at" timestamp,
  "user_id" int
);

CREATE TABLE "Goals" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "title" varchar,
  "colorTag" varchar,
  "target" int,
  "progress" int,
  "is_active" boolean,
  "deadline" timestamp,
  "created_at" timestamp,
  "user_id" int
);

CREATE TABLE "Dreams" (
  "id" SERIAL PRIMARY KEY,
  "uuid" varchar,
  "title" varchar,
  "description" varchar,
  "category" varchar,
  "why_important" varchar,
  "gains_achieving" varchar,
  "risks_and_obstacle" varchar,
  "image" varbinary,
  "user_id" int
);

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "sessions" ("user_id");

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "Tasks" ("user_id");

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "Routines" ("user_id");

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "Goals" ("user_id");

ALTER TABLE "Users" ADD FOREIGN KEY ("id") REFERENCES "Dreams" ("user_id");

ALTER TABLE "Tasks" ADD FOREIGN KEY ("id") REFERENCES "SubTasks" ("task_id");
