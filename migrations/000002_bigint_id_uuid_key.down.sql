ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_user_id_type_key;
ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_user_id_fkey;
ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_key_key;
ALTER TABLE mfas DROP CONSTRAINT IF EXISTS mfas_pkey;
ALTER TABLE mfas ADD COLUMN user_key UUID;
UPDATE mfas
SET user_key = users.key
FROM users
WHERE mfas.user_id = users.id;
ALTER TABLE mfas ALTER COLUMN user_key SET NOT NULL;
ALTER TABLE mfas DROP COLUMN user_id;
ALTER TABLE mfas RENAME COLUMN user_key TO user_id;
ALTER TABLE mfas DROP COLUMN IF EXISTS id;
ALTER TABLE mfas RENAME COLUMN key TO id;
ALTER TABLE mfas ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE mfas ADD CONSTRAINT mfas_pkey PRIMARY KEY (id);

ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_user_id_method_target_key;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_user_id_fkey;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_key_key;
ALTER TABLE verifications DROP CONSTRAINT IF EXISTS verifications_pkey;
ALTER TABLE verifications ADD COLUMN user_key UUID;
UPDATE verifications
SET user_key = users.key
FROM users
WHERE verifications.user_id = users.id;
ALTER TABLE verifications ALTER COLUMN user_key SET NOT NULL;
ALTER TABLE verifications DROP COLUMN user_id;
ALTER TABLE verifications RENAME COLUMN user_key TO user_id;
ALTER TABLE verifications DROP COLUMN IF EXISTS id;
ALTER TABLE verifications RENAME COLUMN key TO id;
ALTER TABLE verifications ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE verifications ADD CONSTRAINT verifications_pkey PRIMARY KEY (id);

ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_user_id_provider_key;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_user_id_fkey;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_key_key;
ALTER TABLE identities DROP CONSTRAINT IF EXISTS identities_pkey;
ALTER TABLE identities ADD COLUMN user_key UUID;
UPDATE identities
SET user_key = users.key
FROM users
WHERE identities.user_id = users.id;
ALTER TABLE identities ALTER COLUMN user_key SET NOT NULL;
ALTER TABLE identities DROP COLUMN user_id;
ALTER TABLE identities RENAME COLUMN user_key TO user_id;
ALTER TABLE identities DROP COLUMN IF EXISTS id;
ALTER TABLE identities RENAME COLUMN key TO id;
ALTER TABLE identities ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE identities ADD CONSTRAINT identities_pkey PRIMARY KEY (id);

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_key_key;
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE users DROP COLUMN IF EXISTS id;
ALTER TABLE users RENAME COLUMN key TO id;
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (id);

ALTER TABLE identities ADD CONSTRAINT identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE verifications ADD CONSTRAINT verifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE mfas ADD CONSTRAINT mfas_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
