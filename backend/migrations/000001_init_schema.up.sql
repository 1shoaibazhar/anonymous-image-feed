CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title TEXT NOT NULL,
  file_path TEXT NOT NULL,
  mime_type TEXT NOT NULL,
  width INT,
  height INT,
  size_bytes BIGINT,
  status TEXT NOT NULL DEFAULT 'processing',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE tags (
  id SERIAL PRIMARY KEY,
  name CITEXT UNIQUE NOT NULL
);

CREATE TABLE image_tags (
  image_id UUID REFERENCES images(id) ON DELETE CASCADE,
  tag_id INT REFERENCES tags(id) ON DELETE CASCADE,
  PRIMARY KEY (image_id, tag_id)
);

CREATE INDEX idx_images_created_at ON images (created_at DESC);
CREATE INDEX idx_image_tags_tag_id ON image_tags (tag_id);