ALTER TABLE organizations ALTER COLUMN created_at SET DEFAULT now();
ALTER TABLE organization_members ALTER COLUMN date_joined SET DEFAULT now();
