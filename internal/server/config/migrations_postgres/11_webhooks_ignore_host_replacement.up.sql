ALTER TABLE webhooks
    ADD COLUMN ignore_host_replacement TEXT DEFAULT 'global' NOT NULL;

commit;
