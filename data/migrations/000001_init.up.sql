CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION citext;
CREATE DOMAIN domain_email AS citext
CHECK(
   VALUE ~ '^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$'
);

CREATE TABLE accounts (
  address bytea,
  name varchar(60) NOT NULL,
  created_at timestamp DEFAULT now(),
  email domain_email NOT NULL,
  admin boolean DEFAULT FALSE,
  PRIMARY KEY (address)
);

CREATE TABLE organization (
    id uuid DEFAULT uuid_generate_v4 (),
    name varchar(60) NOT NULL,
    created_at timestamp NOT NULL,
    -- will connect a wallet later
    -- wallet_address bytea,
    owner_address bytea NOT NULL,
    token bytea NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (owner_address) REFERENCES accounts(address)
);

CREATE TABLE organization_members (
    organization_id uuid NOT NULL,
    member_address bytea NOT NULL,
    date_joined timestamp NOT NULL,
    role varchar(200) NOT NULL,
    is_admin boolean DEFAULT FALSE,
    salary integer DEFAULT 0,
    PRIMARY KEY (organization_id, member_address),
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (member_address) REFERENCES accounts(address)
);

CREATE TABLE salary_payment (
    payment_id uuid DEFAULT uuid_generate_v4 (),
    organization_id uuid NOT NULL,
    member_address bytea NOT NULL,
    transaction_hash CHAR(66) NOT NULL,
    amount integer NOT NULL,
    token bytea NOT NULL,
    date timestamp NOT NULL,
    PRIMARY KEY (payment_id),
    FOREIGN KEY (organization_id) REFERENCES organization(id),
    FOREIGN KEY (member_address) REFERENCES accounts(address)
);
