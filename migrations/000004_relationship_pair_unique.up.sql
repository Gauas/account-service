ALTER TABLE relationships DROP CONSTRAINT IF EXISTS relationships_actor_partner_unique;

CREATE UNIQUE INDEX IF NOT EXISTS idx_relationships_pair_unique
    ON relationships (LEAST(actor_id, partner_id), GREATEST(actor_id, partner_id));
