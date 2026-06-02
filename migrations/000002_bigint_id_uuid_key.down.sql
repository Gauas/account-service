ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_user_id_fkey;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_user_id_fkey;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_user_id_fkey;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_key_key;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_key_key;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_key_key;
ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_key_key;

ALTER TABLE identities ADD COLUMN user_id_old UUID;
UPDATE identities i SET user_id_old = u.key FROM users u WHERE i.user_id = u.id;
ALTER TABLE identities DROP COLUMN user_id;
ALTER TABLE identities RENAME COLUMN user_id_old TO user_id;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_pkey;
ALTER TABLE identities DROP COLUMN id;
ALTER TABLE identities RENAME COLUMN key TO id;
ALTER TABLE identities ADD CONSTRAINT identities_pkey PRIMARY KEY (id);
ALTER TABLE identities ADD CONSTRAINT identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE verifications ADD COLUMN user_id_old UUID;
UPDATE verifications v SET user_id_old = u.key FROM users u WHERE v.user_id = u.id;
ALTER TABLE verifications DROP COLUMN user_id;
ALTER TABLE verifications RENAME COLUMN user_id_old TO user_id;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_pkey;
ALTER TABLE verifications DROP COLUMN id;
ALTER TABLE verifications RENAME COLUMN key TO id;
ALTER TABLE verifications ADD CONSTRAINT verifications_pkey PRIMARY KEY (id);
ALTER TABLE verifications ADD CONSTRAINT verifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE mfas ADD COLUMN user_id_old UUID;
UPDATE mfas m SET user_id_old = u.key FROM users u WHERE m.user_id = u.id;
ALTER TABLE mfas DROP COLUMN user_id;
ALTER TABLE mfas RENAME COLUMN user_id_old TO user_id;
ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_pkey;
ALTER TABLE mfas DROP COLUMN id;
ALTER TABLE mfas RENAME COLUMN key TO id;
ALTER TABLE mfas ADD CONSTRAINT mfas_pkey PRIMARY KEY (id);
ALTER TABLE mfas ADD CONSTRAINT mfas_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN key TO id;
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);
