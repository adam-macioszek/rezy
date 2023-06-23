CREATE TABLE "reservations" (
  "id" bigserial PRIMARY KEY NOT NULL,
  "table_size" int NOT NULL,
  "start_time" timestamp NOT NULL,
  "booked" boolean NOT NULL,
  "duration" int NOt NULL
);


