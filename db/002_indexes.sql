CREATE INDEX idx_age2
    ON "Population".person USING btree
    (age ASC NULLS LAST)
    WITH (deduplicate_items=True)
;


CREATE INDEX idx_name2
    ON "Population".person USING hash
    (name)
;