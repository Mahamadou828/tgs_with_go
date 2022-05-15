CREATE TABLE IF NOT EXISTS "public".provider
(
    id                       UUID NOT NULL,
    email                    varchar(255) NOT NULL,
    billingAddressStreet     text NOT NULL,
    billingAddressPostalCode text NOT NULL,
    billingAddressCity       text NOT NULL,
    billingAddressCountry    text NOT NULL,
    phoneNumber              varchar(255) NOT NULL,
    siren                    varchar(255) NOT NULL,
    sendBill                 boolean NOT NULL,
    generateBill             boolean NOT NULL,
    paymentByTgs             boolean NOT NULL,
    active                   boolean NOT NULL,
    name                     varchar(255) NOT NULL,
    providerRating           int NOT NULL,
    apiToken                 text NOT NULL,
    createdAt                timestamp NOT NULL,
    updatedAt                timestamp NOT NULL,
    deletedAt                timestamp,
    CONSTRAINT PK_168 PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS "public".fare
(
    id                             UUID NOT NULL,
    priceExludingTaxes             decimal(10,2) NOT NULL,
    priceIncudingTaxes             decimal(10,2) NOT NULL,
    priceExludingTaxesInEur        decimal(10,2) NOT NULL,
    priceIncudingTaxesInEur        decimal(10,2) NOT NULL,
    finalPriceExcludingTaxes       decimal(10,2) NOT NULL,
    finalPriceIncludingTaxes       decimal(10,2) NOT NULL,
    finalPriceExcludingTaxesInEur  decimal(10,2) NOT NULL,
    finalPriceIncludingTaxesInEur  decimal(10,2) NOT NULL,
    serviceFeesIncludingTaxes      decimal(10,2) NOT NULL,
    serviceFeesExcludingTaxes      decimal(10,2) NOT NULL,
    serviceFeesIncludingTaxesInEur decimal(10,2) NOT NULL,
    serviceFeesExcludingTaxesInEur decimal(10,2) NOT NULL,
    commission                     decimal(10,2) NOT NULL,
    commissionInEur                decimal(10,2) NOT NULL,
    majoration                     decimal(10,2) NOT NULL,
    majorationInEur                decimal(10,2) NOT NULL,
    taxes                          decimal(10,2) NOT NULL,
    taxesInEur                     decimal(10,2) NOT NULL,
    servicefeesTaxes               decimal(10,2) NOT NULL,
    servicefeesTaxesInEuro         decimal(10,2) NOT NULL,
    CONSTRAINT PK_114 PRIMARY KEY ( id )
);

CREATE TABLE IF NOT EXISTS "public".country
(
    id           UUID NOT NULL,
    name         varchar(255) NOT NULL,
    currencyName varchar(255) NOT NULL,
    currencyCode varchar(255) NOT NULL,
    alphaCode2   varchar(255) NOT NULL,
    alphaCode3   varchar(255) NOT NULL,
    vat          decimal NOT NULL,
    CONSTRAINT PK_93 PRIMARY KEY ( id ),
    CONSTRAINT name UNIQUE ( name )
);

CREATE TABLE IF NOT EXISTS "public".aggregator
(
    id              UUID NOT NULL,
    name            varchar(255) NOT NULL,
    code            varchar(255) NOT NULL,
    apiKey          varchar(255) NOT NULL,
    providerTimeout int NOT NULL,
    active          boolean NOT NULL,
    type            varchar(255) NOT NULL,
    paymentByTgs    boolean NOT NULL,
    logoUrl         text NOT NULL,
    createdAt       timestamp NOT NULL,
    updatedAt       timestamp NOT NULL,
    deletedAt       timestamp,
    CONSTRAINT PK_23 PRIMARY KEY ( id ),
    CONSTRAINT Index_231 UNIQUE ( apiKey )
);

COMMENT ON TABLE "public".aggregator IS 'An aggregator is a client that integrate the api and redistribute the offer';

COMMENT ON COLUMN "public".aggregator.type IS 'There''s two type of aggregator:
- Enterprise: represent an aggregator that redistribute the offer under the enterprise logic
- Default: represent an aggregator that redistribute the offer under the normal workflow logic';

CREATE TABLE IF NOT EXISTS "public"."user"
(
    id              UUID NOT NULL,
    aggregatorId    UUID NOT NULL,
    email           varchar(255) NOT NULL,
    phoneNumber     varchar(255) NOT NULL,
    name            varchar(255) NOT NULL,
    stripeId        text NOT NULL,
    apiKey          varchar(255) NOT NULL,
    active          boolean NOT NULL,
    cognitoId       text NOT NULL,
    isMonthlyActive boolean NOT NULL,
    isCGUAccepted   boolean NOT NULL,
    role            varchar(255) NOT NULL,
    updatedAt       timestamp NOT NULL,
    deletedAt       timestamp,
    createdAt       timestamp NOT NULL,
    CONSTRAINT PK_7 PRIMARY KEY ( id ),
    CONSTRAINT cognitoId UNIQUE ( cognitoId ),
    CONSTRAINT FK_50 FOREIGN KEY ( aggregatorId ) REFERENCES "public".aggregator ( id )
);

CREATE TABLE IF NOT EXISTS "public".search
(
    id               UUID NOT NULL,
    userId           UUID NOT NULL,
    aggregatorId     UUID NOT NULL,
    countryId        UUID NOT NULL,
    startDate        timestamp NOT NULL,
    isPlanned        boolean NOT NULL,
    startFullAddress text NOT NULL,
    startLongitude   decimal NOT NULL,
    startLatitude    decimal NOT NULL,
    startPostcode    varchar(255) NOT NULL,
    startCountry     varchar(255) NOT NULL,
    startRegion      varchar(255) NOT NULL,
    endFullAddress   text NOT NULL,
    endLongitude     decimal NOT NULL,
    endLatitude      decimal NOT NULL,
    endPostcode      varchar(255) NOT NULL,
    endCountry       varchar(255) NOT NULL,
    endRegion        varchar(255) NOT NULL,
    distance         int NOT NULL,
    nbrOfPassager    int NOT NULL,
    flightNumber     varchar(255) NOT NULL,
    trainNumber      varchar(255) NOT NULL,
    nbrOfBagages     int NOT NULL,
    currentRates     decimal NOT NULL,
    CONSTRAINT PK_60 PRIMARY KEY ( id ),
    CONSTRAINT FK_109 FOREIGN KEY ( countryId ) REFERENCES "public".country ( id ),
    CONSTRAINT FK_204 FOREIGN KEY ( aggregatorId ) REFERENCES "public".aggregator ( id ),
    CONSTRAINT FK_86 FOREIGN KEY ( userId ) REFERENCES "public"."user" ( id )
);

CREATE INDEX FK_111 ON "public".search
    (
     countryId
        );

CREATE INDEX FK_206 ON "public".search
    (
     aggregatorId
        );

CREATE INDEX FK_88 ON "public".search
    (
     userId
        );

CREATE TABLE IF NOT EXISTS "public".payment_method
(
    id                UUID NOT NULL,
    name              varchar(255) NOT NULL,
    userId            UUID NOT NULL,
    displayCreditCard varchar(255) NOT NULL,
    stripeId          text NOT NULL,
    type              varchar(255) NOT NULL,
    isFavorite        boolean NOT NULL,
    createdAt         timestamp NOT NULL,
    updadtedAt        timestamp NOT NULL,
    deletedAt         timestamp,
    CONSTRAINT PK_209 PRIMARY KEY ( id ),
    CONSTRAINT FK_224 FOREIGN KEY ( userId ) REFERENCES "public"."user" ( id )
);

CREATE INDEX FK_226 ON "public".payment_method
    (
     userId
        );

CREATE TABLE IF NOT EXISTS "public".offer
(
    id                  UUID NOT NULL,
    fareId              UUID NOT NULL,
    searchId            UUID NOT NULL,
    providerId          UUID NOT NULL,
    displayPrice        varchar(255) NOT NULL,
    displayPriceNumeric decimal(10,2) NOT NULL,
    displayProviderName varchar(255) NOT NULL,
    co2Emission         decimal(10,2) NOT NULL,
    providerOfferId     text NOT NULL,
    eta                 int NOT NULL,
    startDate           timestamp NOT NULL,
    logo                text NOT NULL,
    redirectionLink     text NOT NULL,
    termAndCondition    varchar(255) NOT NULL,
    providerOfferName   varchar(255) NOT NULL,
    vehicleType         varchar(255) NOT NULL,
    description         varchar(255) NOT NULL,
    createdAt           timestamp NOT NULL,
    updatedAt           timestamp NOT NULL,
    deletedAt           timestamp,
    CONSTRAINT PK_145 PRIMARY KEY ( id ),
    CONSTRAINT FK_163 FOREIGN KEY ( fareId ) REFERENCES "public".fare ( id ),
    CONSTRAINT FK_189 FOREIGN KEY ( providerId ) REFERENCES "public".provider ( id ),
    CONSTRAINT FK_192 FOREIGN KEY ( searchId ) REFERENCES "public".search ( id )
);

CREATE INDEX FK_165 ON "public".offer
    (
     fareId
        );

CREATE INDEX FK_191 ON "public".offer
    (
     providerId
        );

CREATE INDEX FK_194 ON "public".offer
    (
     searchId
        );

CREATE INDEX Index_162 ON "public".offer
    (
     startDate
        );
