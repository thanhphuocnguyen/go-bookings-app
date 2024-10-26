CREATE TABLE
    "restrictions" (
        "id" integer PRIMARY KEY,
        "name" varchar NOT NULL,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ())
    );

CREATE TABLE
    "room_restrictions" (
        "id" integer PRIMARY KEY,
        "start_date" date,
        "end_date" date,
        "room_id" integer NOT NULL,
        "restriction_id" integer NOT NULL,
        "created_at" timestamp DEFAULT (now ()),
        "updated_at" timestamp DEFAULT (now ()),
        FOREIGN KEY ("room_id") REFERENCES "rooms" ("id") ON DELETE CASCADE ON UPDATE CASCADE,
        FOREIGN KEY ("restriction_id") REFERENCES "restrictions" ("id") ON DELETE CASCADE ON UPDATE CASCADE
    );