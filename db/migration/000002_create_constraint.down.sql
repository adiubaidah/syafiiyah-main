ALTER TABLE "smart_card" DROP CONSTRAINT IF EXISTS "check_santri_or_employee";
ALTER TABLE "santri_presence" DROP COLUMN IF EXISTS "created_date";
ALTER TABLE "santri_presence" DROP CONSTRAINT IF EXISTS "unique_santri_schedule_date";