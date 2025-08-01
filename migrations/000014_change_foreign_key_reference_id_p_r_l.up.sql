
ALTER TABLE password_reset_link
DROP CONSTRAINT  fk_id_auth;

ALTER TABLE password_reset_link
ADD CONSTRAINT fk_id_auth
FOREIGN KEY (user_id) REFERENCES auth(id) ON DELETE CASCADE;