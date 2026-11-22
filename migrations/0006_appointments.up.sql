CREATE TABLE appointments (
  id           BIGSERIAL PRIMARY KEY,
  patient_id   BIGINT NOT NULL REFERENCES users(id),
  schedule_id  BIGINT NOT NULL REFERENCES doctor_schedules(id),
  queue_number INT    NOT NULL CHECK (queue_number > 0),
  status       VARCHAR(20) NOT NULL DEFAULT 'booked'
               CHECK (status IN ('booked','arrived','in_consultation','completed','cancelled','no_show')),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (schedule_id, queue_number)
);

CREATE UNIQUE INDEX appointments_patient_schedule_active_uidx
  ON appointments (patient_id, schedule_id) WHERE status <> 'cancelled';
CREATE INDEX appointments_patient_idx  ON appointments (patient_id);
CREATE INDEX appointments_schedule_idx ON appointments (schedule_id, queue_number);
