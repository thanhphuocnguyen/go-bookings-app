CREATE TABLE
    "restrictions" (
        "id" SERIAL PRIMARY KEY,
        "name" varchar UNIQUE NOT NULL,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ())
    );

CREATE TABLE
    "room_restrictions" (
        "id" SERIAL PRIMARY KEY,
        "start_date" date,
        "end_date" date,
        "room_id" integer NOT NULL,
        "restriction_id" integer NOT NULL,
        "reservation_id" integer REFERENCES "reservations" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ()),
        FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY ("restriction_id") REFERENCES "restrictions" ("id") ON DELETE CASCADE ON UPDATE CASCADE
    );