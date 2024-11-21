ALTER TABLE smart_card
ADD CONSTRAINT check_santri_or_employee
CHECK (
  (santri_id IS NOT NULL AND employee_id IS NULL) OR
  (santri_id IS NULL AND employee_id IS NOT NULL) OR
  (santri_id IS NULL AND employee_id IS NULL)
);