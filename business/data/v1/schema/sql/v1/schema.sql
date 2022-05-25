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
    id                 UUID         NOT NULL,
    aggregator_id      UUID         NOT NULL,
    enterprise_id      UUID,
    enterprise_team_id UUID,
    email              varchar(255) NOT NULL,
    phone_number       varchar(255) NOT NULL,
    name               varchar(255) NOT NULL,
    stripe_id          text         NOT NULL,
    api_key            varchar(255) NOT NULL,
    active             boolean      NOT NULL,
    cognito_id         text         NOT NULL,
    is_monthly_active  boolean      NOT NULL,
    is_cgu_accepted    boolean      NOT NULL,
    role               varchar(255) NOT NULL,
    updated_at         timestamp    NOT NULL,
    deleted_at         timestamp,
    created_at         timestamp    NOT NULL,
    CONSTRAINT PK_7 PRIMARY KEY (id),
    CONSTRAINT FK_50 FOREIGN KEY (aggregator_id) REFERENCES "public".aggregator (id),
    CONSTRAINT FK_255 FOREIGN KEY ( enterprise_id ) REFERENCES "public".enterprise ( "id" ),
    CONSTRAINT FK_301 FOREIGN KEY ( enterprise_team_id ) REFERENCES "public".enterprise_team ( "id" )
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
    CONSTRAINT FK_86 FOREIGN KEY (user_id) REFERENCES "public"."user" (id),
);

CREATE TABLE "public".enterprise
(
    "id"                        uuid NOT NULL,
    code                      varchar(255) NOT NULL,
    pack_id                   uuid NOT NULL,
    name                      varchar(255) NOT NULL,
    created_at                timestamp NOT NULL,
    updated_at                timestamp NOT NULL,
    deleted_at                timestamp NOT NULL,
    contact_email             varchar(255) NOT NULL,
    description               varchar(255) NOT NULL,
    logo_url                  varchar(255) NOT NULL,
    max_carbon_emission_offer varchar(255) NOT NULL,
    blocked_provider          text[] NOT NULL,
    blocked_product_type      text[] NOT NULL,
    active                    boolean NOT NULL,
    CONSTRAINT PK_234 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_267 FOREIGN KEY ( pack_id ) REFERENCES "public".enterprise_pack ( "id" )
);

CREATE TABLE "public".facturation_entity
(
    "id"                 uuid NOT NULL,
    denomination       varchar(255) NOT NULL,
    enterprise_id      uuid NOT NULL,
    facturation_number varchar(255) NOT NULL,
    vat                varchar(255) NOT NULL,
    street             varchar(255) NOT NULL,
    postal_code        varchar(255) NOT NULL,
    town               varchar(255) NOT NULL,
    country            varchar(255) NOT NULL,
    created_at         timestamp NOT NULL,
    updated_at         timestamp NOT NULL,
    deleted_at         timestamp NOT NULL,
    CONSTRAINT PK_318 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_329 FOREIGN KEY ( enterprise_id ) REFERENCES "public".enterprise ( "id" )
);

CREATE INDEX FK_331 ON "public".facturation_entity
    (
     enterprise_id
        );

CREATE TABLE "public".enterprise_policy
(
    "id"                  uuid NOT NULL,
    name                varchar(255) NOT NULL,
    enterprise_id       uuid NOT NULL,
    description         varchar(255) NOT NULL,
    collaborator_budget int NOT NULL,
    start_service_time  varchar(255) NOT NULL,
    end_service_time    varchar(255) NOT NULL,
    blocked_days        text[] NOT NULL,
    budget_type         varchar(255) NOT NULL,
    CONSTRAINT PK_281 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_289 FOREIGN KEY ( enterprise_id ) REFERENCES "public".enterprise ( "id" )
);

CREATE INDEX FK_291 ON "public".enterprise_policy
    (
     enterprise_id
        );

CREATE TABLE "public".enterprise_team
(
    "id"                    uuid NOT NULL,
    name                  varchar(255) NOT NULL,
    facturation_entity_id uuid NOT NULL,
    enterprise_id         uuid NOT NULL,
    description           varchar(255) NOT NULL,
    payment_method        varchar(255) NOT NULL,
    created_at            timestamp NOT NULL,
    updated_at            timestamp NOT NULL,
    deleted_at            timestamp NOT NULL,
    CONSTRAINT PK_294 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_313 FOREIGN KEY ( enterprise_id ) REFERENCES "public".enterprise ( "id" ),
    CONSTRAINT FK_335 FOREIGN KEY ( facturation_entity_id ) REFERENCES "public".facturation_entity ( "id" )
);

CREATE INDEX FK_315 ON "public".enterprise_team
    (
     enterprise_id
        );

CREATE INDEX FK_337 ON "public".enterprise_team
    (
     facturation_entity_id
        );


CREATE INDEX FK_269 ON "public".enterprise
    (
     pack_id
        );



CREATE TABLE "public".enterprise_pack
(
    "id"                                uuid NOT NULL,
    name                              varchar(255) NOT NULL,
    send_monthly_report               boolean NOT NULL,
    can_customize_report              boolean NOT NULL,
    send_expense_report               boolean NOT NULL,
    included_field_in_report          text[] NOT NULL,
    max_active_callaborator_per_month int NOT NULL,
    CONSTRAINT PK_260 PRIMARY KEY ( "id" )
);

CREATE TABLE "public".enterprise_policy
(
    "id"                uuid            NOT NULL,
    name                varchar(255)    NOT NULL,
    enterprise_id       uuid            NOT NULL,
    description         varchar(255)    NOT NULL,
    collaborator_budget int             NOT NULL,
    start_service_time  varchar(255)    NOT NULL,
    end_service_time    varchar(255)    NOT NULL,
    blocked_days        text            NOT NULL,
    budget_type         varchar(255)    NOT NULL,
    CONSTRAINT PK_281 PRIMARY KEY ( "id" ),
    CONSTRAINT FK_289 FOREIGN KEY ( enterprise_id ) REFERENCES "public".enterprise ( "id" )
);

CREATE INDEX FK_291 ON "public".enterprise_policy
    (
     enterprise_id
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

CREATE INDEX FK_257 ON "public"."user"
    (
     enterprise_id
        );
