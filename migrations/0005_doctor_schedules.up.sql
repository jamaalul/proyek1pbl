CREATE TABLE doctor_schedules (
  id            BIGSERIAL PRIMARY KEY,
  doctor_id     BIGINT NOT NULL REFERENCES doctors(id),
  practice_date DATE   NOT NULL,
  start_time    TIME   NOT NULL,
  end_time      TIME   NOT NULL,
  quota         INT    NOT NULL CHECK (quota > 0),
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (end_time > start_time),
  UNIQUE (doctor_id, practice_date, start_time),
  EXCLUDE USING gist (
    doctor_id WITH =,
    tsrange(practice_date + start_time, practice_date + end_time) WITH &&
  )
);

CREATE INDEX doctor_schedules_date_idx ON doctor_schedules (practice_date);
