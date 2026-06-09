DROP INDEX IF EXISTS idx_relationships_pair_unique;

ALTER TABLE relationships
    ADD CONSTRAINT relationships_actor_partner_unique UNIQUE (actor_id, partner_id);
