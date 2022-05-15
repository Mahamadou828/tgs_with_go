CREATE TABLE IF NOT EXISTS "public".provider
(
    id                          UUID         NOT NULL,
    email                       varchar(255) NOT NULL,
    billing_address_street      text         NOT NULL,
    billing_address_postal_code text         NOT NULL,
    billing_address_city        text         NOT NULL,
    billing_address_country     text         NOT NULL,
    phone_number                varchar(255) NOT NULL,
    siren                       varchar(255) NOT NULL,
    send_bill                   boolean      NOT NULL,
    generate_bill               boolean      NOT NULL,
    payment_by_tgs              boolean      NOT NULL,
    active                      boolean      NOT NULL,
    name                        varchar(255) NOT NULL,
    provider_rating             int          NOT NULL,
    api_token                   text         NOT NULL,
    created_at                  timestamp    NOT NULL,
    updated_at                  timestamp    NOT NULL,
    deleted_at                  timestamp,
    CONSTRAINT PK_168 PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS "public".fare
(
    id                                  UUID           NOT NULL,
    price_excluding_taxes               decimal(10, 2) NOT NULL,
    price_including_axes                decimal(10, 2) NOT NULL,
    price_excluding_taxes_in_eur        decimal(10, 2) NOT NULL,
    price_including_taxes_in_eur        decimal(10, 2) NOT NULL,
    final_price_excluding_taxes         decimal(10, 2) NOT NULL,
    final_price_including_taxes         decimal(10, 2) NOT NULL,
    final_price_excluding_taxes_in_eur  decimal(10, 2) NOT NULL,
    final_price_including_taxes_in_eur  decimal(10, 2) NOT NULL,
    service_fees_including_taxes        decimal(10, 2) NOT NULL,
    service_fees_excluding_taxes        decimal(10, 2) NOT NULL,
    service_fees_including_taxes_in_eur decimal(10, 2) NOT NULL,
    service_fees_excluding_taxes_in_eur decimal(10, 2) NOT NULL,
    commission                          decimal(10, 2) NOT NULL,
    commission_in_eur                   decimal(10, 2) NOT NULL,
    majoration                          decimal(10, 2) NOT NULL,
    majoration_in_eur                   decimal(10, 2) NOT NULL,
    taxes                               decimal(10, 2) NOT NULL,
    taxes_in_eur                        decimal(10, 2) NOT NULL,
    service_fees_taxes                  decimal(10, 2) NOT NULL,
    service_fees_taxes_in_euro          decimal(10, 2) NOT NULL,
    CONSTRAINT PK_114 PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS "public".country
(
    id            UUID         NOT NULL,
    name          varchar(255) NOT NULL,
    currency_name varchar(255) NOT NULL,
    currency_code varchar(255) NOT NULL,
    alpha_code2   varchar(255) NOT NULL,
    alpha_code3   varchar(255) NOT NULL,
    vat           decimal      NOT NULL,
    CONSTRAINT PK_93 PRIMARY KEY (id),
    CONSTRAINT name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS "public".aggregator
(
    id               UUID         NOT NULL,
    name             varchar(255) NOT NULL,
    code             varchar(255) NOT NULL,
    api_key          varchar(255) NOT NULL,
    provider_timeout int          NOT NULL,
    active           boolean      NOT NULL,
    type             varchar(255) NOT NULL,
    payment_by_tgs   boolean      NOT NULL,
    logo_url         text         NOT NULL,
    created_at       timestamp    NOT NULL,
    updated_at       timestamp    NOT NULL,
    deleted_at       timestamp,
    CONSTRAINT PK_23 PRIMARY KEY (id),
    CONSTRAINT Index_231 UNIQUE (api_key)
);

COMMENT ON TABLE "public".aggregator IS 'An aggregator is a client that integrate the api and redistribute the offer';

COMMENT ON COLUMN "public".aggregator.type IS 'There''s two type of aggregator:
- Enterprise: represent an aggregator that redistribute the offer under the enterprise logic
- Default: represent an aggregator that redistribute the offer under the normal workflow logic';

CREATE TABLE IF NOT EXISTS "public"."user"
(
    id                UUID         NOT NULL,
    aggregator_id     UUID         NOT NULL,
    email             varchar(255) NOT NULL,
    phone_number      varchar(255) NOT NULL,
    name              varchar(255) NOT NULL,
    stripe_id         text         NOT NULL,
    api_key           varchar(255) NOT NULL,
    active            boolean      NOT NULL,
    cognito_id        text         NOT NULL,
    is_monthly_active boolean      NOT NULL,
    is_cgu_accepted   boolean      NOT NULL,
    role              varchar(255) NOT NULL,
    updated_at        timestamp    NOT NULL,
    deleted_at        timestamp,
    created_at        timestamp    NOT NULL,
    CONSTRAINT PK_7 PRIMARY KEY (id),
    CONSTRAINT cognitoId UNIQUE (cognito_id),
    CONSTRAINT FK_50 FOREIGN KEY (aggregator_id) REFERENCES "public".aggregator (id)
);

CREATE TABLE IF NOT EXISTS "public".search
(
    id                 UUID         NOT NULL,
    user_id            UUID         NOT NULL,
    aggregator_id      UUID         NOT NULL,
    country_id         UUID         NOT NULL,
    start_date         timestamp    NOT NULL,
    is_planned         boolean      NOT NULL,
    start_full_address text         NOT NULL,
    start_longitude    decimal      NOT NULL,
    start_latitude     decimal      NOT NULL,
    start_postcode     varchar(255) NOT NULL,
    start_country      varchar(255) NOT NULL,
    start_region       varchar(255) NOT NULL,
    end_full_address   text         NOT NULL,
    end_longitude      decimal      NOT NULL,
    end_latitude       decimal      NOT NULL,
    end_postcode       varchar(255) NOT NULL,
    end_country        varchar(255) NOT NULL,
    end_region         varchar(255) NOT NULL,
    distance           int          NOT NULL,
    nbr_of_passenger   int          NOT NULL,
    flight_number      varchar(255) NOT NULL,
    train_number       varchar(255) NOT NULL,
    nbr_of_baggage     int          NOT NULL,
    current_rates      decimal      NOT NULL,
    CONSTRAINT PK_60 PRIMARY KEY (id),
    CONSTRAINT FK_109 FOREIGN KEY (country_id) REFERENCES "public".country (id),
    CONSTRAINT FK_204 FOREIGN KEY (aggregator_id) REFERENCES "public".aggregator (id),
    CONSTRAINT FK_86 FOREIGN KEY (user_id) REFERENCES "public"."user" (id)
);

CREATE INDEX FK_111 ON "public".search
    (
     country_id
        );

CREATE INDEX FK_206 ON "public".search
    (
     aggregator_id
        );

CREATE INDEX FK_88 ON "public".search
    (
     user_id
        );

CREATE TABLE IF NOT EXISTS "public".payment_method
(
    id                  UUID         NOT NULL,
    name                varchar(255) NOT NULL,
    user_id             UUID         NOT NULL,
    display_credit_card varchar(255) NOT NULL,
    stripe_id           text         NOT NULL,
    type                varchar(255) NOT NULL,
    is_favorite         boolean      NOT NULL,
    created_at          timestamp    NOT NULL,
    updated_at          timestamp    NOT NULL,
    deleted_at          timestamp,
    CONSTRAINT PK_209 PRIMARY KEY (id),
    CONSTRAINT FK_224 FOREIGN KEY (user_id) REFERENCES "public"."user" (id)
);

CREATE INDEX FK_226 ON "public".payment_method
    (
     user_id
        );

CREATE TABLE IF NOT EXISTS "public".offer
(
    id                    UUID           NOT NULL,
    fare_id               UUID           NOT NULL,
    search_id             UUID           NOT NULL,
    provider_id           UUID           NOT NULL,
    display_price         varchar(255)   NOT NULL,
    display_price_numeric decimal(10, 2) NOT NULL,
    display_provider_name varchar(255)   NOT NULL,
    co2_emission          decimal(10, 2) NOT NULL,
    provider_offer_id     text           NOT NULL,
    eta                   int            NOT NULL,
    start_date            timestamp      NOT NULL,
    logo                  text           NOT NULL,
    redirection_link      text           NOT NULL,
    term_and_condition    varchar(255)   NOT NULL,
    provider_offer_name   varchar(255)   NOT NULL,
    vehicle_type          varchar(255)   NOT NULL,
    description           varchar(255)   NOT NULL,
    createdAt             timestamp      NOT NULL,
    updatedAt             timestamp      NOT NULL,
    deletedAt             timestamp,
    CONSTRAINT PK_145 PRIMARY KEY (id),
    CONSTRAINT FK_163 FOREIGN KEY (fare_id) REFERENCES "public".fare (id),
    CONSTRAINT FK_189 FOREIGN KEY (provider_id) REFERENCES "public".provider (id),
    CONSTRAINT FK_192 FOREIGN KEY (search_id) REFERENCES "public".search (id)
);

CREATE INDEX FK_165 ON "public".offer
    (
     fare_id
        );

CREATE INDEX FK_191 ON "public".offer
    (
     provider_id
        );

CREATE INDEX FK_194 ON "public".offer
    (
     search_id
        );

CREATE INDEX Index_162 ON "public".offer
    (
     start_date
        );
