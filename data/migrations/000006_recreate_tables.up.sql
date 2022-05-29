CREATE TABLE accounts (
  address VARCHAR(64),
  name varchar(60) NOT NULL,
  created_at timestamp DEFAULT now(),
  email domain_email NOT NULL,
  admin boolean DEFAULT FALSE,
  PRIMARY KEY (address)
);

CREATE TABLE organizations (
    id uuid DEFAULT uuid_generate_v4 (),
    name varchar(60) NOT NULL,
    created_at timestamp NOT NULL,
    -- will connect a wallet later
    -- wallet_address VARCHAR(64),
    owner_address VARCHAR(64) NOT NULL,
    token VARCHAR(64) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_address) REFERENCES accounts(address)
);

CREATE TABLE organization_members (
    organization_id uuid NOT NULL,
    member_address VARCHAR(64) NOT NULL,
    date_joined timestamp NOT NULL,
    role varchar(200) NOT NULL,
    is_admin boolean DEFAULT FALSE,
    salary NUMERIC DEFAULT 0,
    PRIMARY KEY (organization_id, member_address),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (member_address) REFERENCES accounts(address)
);

CREATE TABLE salary_payments (
    payment_id uuid DEFAULT uuid_generate_v4 (),
    organization_id uuid NOT NULL,
    member_address VARCHAR(64) NOT NULL,
    transaction_hash CHAR(66) NOT NULL,
    amount NUMERIC NOT NULL,
    token VARCHAR(64) NOT NULL,
    date timestamp NOT NULL,
    PRIMARY KEY (payment_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(id),
    FOREIGN KEY (member_address) REFERENCES accounts(address)
);
