ALTER TABLE 
  posts
ALTER 
    COLUMN version SET DEFAULT 0;

UPDATE posts
SET version = 0
WHERE version IS NULL;
