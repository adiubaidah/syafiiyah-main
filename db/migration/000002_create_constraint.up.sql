ALTER TABLE smart_card
ADD CONSTRAINT check_santri_or_employee
CHECK (
  (santri_id IS NOT NULL AND employee_id IS NULL) OR
  (santri_id IS NULL AND employee_id IS NOT NULL) OR
  (santri_id IS NULL AND employee_id IS NULL)
);

ALTER TABLE "santri_presence"
ADD COLUMN "created_date" DATE GENERATED ALWAYS AS (("created_at" AT TIME ZONE 'UTC')::DATE) STORED;

ALTER TABLE "santri_presence"
ADD CONSTRAINT unique_santri_schedule_date
UNIQUE ("santri_id", "schedule_id", "created_date");