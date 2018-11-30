ALTER TABLE account
  ADD remote_school_id INT
    DEFAULT 60
    NOT NULL;

UPDATE account
SET remote_school_id = 60
WHERE remote_school_id = 0;
