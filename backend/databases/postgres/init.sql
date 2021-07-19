CREATE TYPE drone_state AS ENUM ('free', 'in-flight');

CREATE TABLE drone
(
    id          SERIAL PRIMARY KEY,
    state       drone_state      DEFAULT 'free',
    weight      DOUBLE PRECISION default 4,
    consumption DOUBLE PRECISION DEFAULT 500
);

CREATE TABLE warehouse
(
    id  SERIAL PRIMARY KEY,
    location_latitude DOUBLE PRECISION DEFAULT 48.080922,
    location_longitude DOUBLE PRECISION DEFAULT 20.766208
);

CREATE TABLE shipping_address
(
    id      SERIAL PRIMARY KEY,
    name    VARCHAR(50) NOT NULL,
    country VARCHAR(10) NOT NULL,
    region  VARCHAR(15) DEFAULT NULL,
    city    VARCHAR(50) NOT NULL,
    zip     INT         NOT NULL,
    street  VARCHAR(50) NOT NULL,
    street2 VARCHAR(50) DEFAULT NULL,
    street3 VARCHAR(50) DEFAULT NULL
);
-- jsonb gyorsabb, és keresésre ugyse használjuk csak elmentjük
-- assigned_drone opcionális, a még nem szállitasban levo csomagoknál NULL
CREATE TABLE parcel
(
    id                 SERIAL PRIMARY KEY,
    name               VARCHAR(50) NOT NULL,
    tracking_id        VARCHAR(25)               default '',
    weight             DOUBLE PRECISION          default 1,
    drop_off_latitude  DOUBLE PRECISION          DEFAULT 0,
    drop_off_longitude DOUBLE PRECISION          DEFAULT 0,
    assigned_drone     INT REFERENCES drone (id) DEFAULT NULL,
    from_address       INT REFERENCES shipping_address (id),
    to_address         INT REFERENCES shipping_address (id)
);

-- tme stamp is null bc this is not the time of insertion, this timestamp is signed by the drone so we know the order of messages
CREATE TABLE telemetry
(
    id                  SERIAL PRIMARY KEY,
    drone_id            INT REFERENCES drone (id),
    speed               DOUBLE PRECISION DEFAULT 0,
    latitude            DOUBLE PRECISION DEFAULT 0,
    longitude           DOUBLE PRECISION DEFAULT 0,
    altitude            DOUBLE PRECISION default 1,
    bearing             DOUBLE PRECISION DEFAULT 0,
    acceleration        DOUBLE PRECISION DEFAULT 0,
    battery_level       INT              DEFAULT NULL,
    battery_temperature INT              DEFAULT NULL,
    motor_temperatures  INTEGER[],
    errors              INTEGER[],
    time_stamp          timestamp        DEFAULT NULL
);
INSERT INTO warehouse (id) VALUES (1);





