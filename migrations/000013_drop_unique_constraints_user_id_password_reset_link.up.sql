ALTER TABLE password_reset_link
DROP CONSTRAINT  password_reset_link_user_id_key;

ALTER TABLE password_reset_link
ADD CONSTRAINT fk_id_auth
FOREIGN KEY (id) REFERENCES auth(id) ON DELETE CASCADE;

