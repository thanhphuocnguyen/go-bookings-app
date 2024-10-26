CREATE TABLE
    "rooms" (
        "id" integer PRIMARY KEY,
        "room_name" varchar,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ())
    );

CREATE TABLE
    "reservation" (
        "id" integer PRIMARY KEY,
        "user_id" integer NOT NULL,
        "phone" varchar NOT NULL,
        "room_id" integer NOT NULL,
        "start_date" date,
        "end_date" date,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ()),
        FOREIGN KEY ("room_id") REFERENCES "rooms" ("id"),
        FOREIGN KEY ("user_id") REFERENCES "users" ("id")
    );