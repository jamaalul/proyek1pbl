CREATE TABLE doctors (
  id             BIGSERIAL PRIMARY KEY,
  user_id        BIGINT NOT NULL UNIQUE REFERENCES users(id),
  polyclinic_id  BIGINT NOT NULL REFERENCES polyclinics(id),
  specialization VARCHAR(100),
  created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX doctors_polyclinic_idx ON doctors (polyclinic_id);
