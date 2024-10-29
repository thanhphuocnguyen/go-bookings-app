CREATE TABLE
    "rooms" (
        "id" SERIAL PRIMARY KEY,
        "name" varchar,
        "description" text,
        "slug" varchar UNIQUE NOT NULL,
        "price" decimal NOT NULL,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ())
    );

CREATE TABLE
    "reservations" (
        "id" SERIAL PRIMARY KEY,
        "user_id" integer NOT NULL,
        "phone" varchar NOT NULL,
        "email" varchar NOT NULL,
        "first_name" varchar NOT NULL,
        "last_name" varchar NOT NULL,
        "room_id" integer NOT NULL,
        "start_date" date,
        "end_date" date,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ()),
        FOREIGN KEY ("room_id") REFERENCES "rooms" ("id"),
        FOREIGN KEY ("user_id") REFERENCES "users" ("id")
    );