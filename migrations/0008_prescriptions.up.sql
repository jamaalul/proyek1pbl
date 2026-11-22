CREATE TABLE prescriptions (
  id                BIGSERIAL PRIMARY KEY,
  medical_record_id BIGINT NOT NULL REFERENCES medical_records(id) ON DELETE CASCADE,
  medication_name   VARCHAR(200) NOT NULL,
  dosage            VARCHAR(100),
  instructions      TEXT,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX prescriptions_record_idx ON prescriptions (medical_record_id);
