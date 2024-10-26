ALTER TABLE "rooms" DROP CONSTRAINT "unique_room_name";
ALTER TABLE "rooms" ALTER COLUMN "room_name" DROP NOT NULL;