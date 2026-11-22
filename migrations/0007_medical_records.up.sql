CREATE TABLE medical_records (
  id             BIGSERIAL PRIMARY KEY,
  appointment_id BIGINT NOT NULL UNIQUE REFERENCES appointments(id),
  doctor_id      BIGINT NOT NULL REFERENCES doctors(id),
  complaint      TEXT,
  diagnosis      TEXT NOT NULL,
  notes          TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX medical_records_doctor_idx ON medical_records (doctor_id);
