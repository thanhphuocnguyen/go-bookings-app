ALTER TABLE "rooms" ALTER COLUMN "room_name" SET NOT NULL;
ALTER TABLE "rooms" ADD CONSTRAINT "unique_room_name" UNIQUE ("room_name");