CREATE INDEX "idx_rooms_room_name" ON "rooms" ("room_name");
CREATE INDEX "idx_reservation_user_id" ON "reservation" ("user_id");
CREATE INDEX "idx_reservation_room_id" ON "reservation" ("room_id");

CREATE INDEX "idx_reservation_start_end_date" ON "reservation" ("start_date", "end_date");
CREATE INDEX "idx_room_restrictions_room_id" ON "room_restrictions" ("room_id");
CREATE INDEX "idx_room_restrictions_restriction_id" ON "room_restrictions" ("restriction_id");
CREATE INDEX "idx_restrictions_unique_name" ON "restrictions" ("name");