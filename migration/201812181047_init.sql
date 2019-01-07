--ENTER POINT -- FIRS script
CREATE TYPE docupload  AS ENUM ('unavailable', 'optional', 'required');
ALTER TYPE docupload OWNER TO bookings;
CREATE TYPE roomtype  AS ENUM ('single', 'double', 'triple', 'quad', 'queen', 'king', 'twin', 'studio', 'suite', 'min_suite');
ALTER TYPE roomtype OWNER TO bookings;
CREATE TYPE states AS ENUM ('draft', 'cancelled', 'booked', 'pending', 'pending_resp', 'rejected', 'completed');
ALTER TYPE states OWNER TO bookings;

CREATE TABLE  rooms(
    id                      UUID PRIMARY KEY,
    name                    TEXT NOT NULL,
    provider                UUID,
    hotel_id                UUID NOT NULL ,
    type_id                 UUID NOT NULL,
    reservation_min_time    INTERVAL,
    reservation_max_time    INTERVAL,
    available_from          TIME NOT NULL,
    available_to            TIME NOT NULL,
    reservation_lead_time   INTERVAL,
    is_shared               BOOLEAN NOT NULL DEFAULT FALSE,
    document_upload         docupload NOT NULL DEFAULT 'unavailable',    
    description             TEXT,
    UNIQUE(name, hotel_id)
);
ALTER TABLE rooms OWNER TO bookings ;

CREATE TABLE bookings(
    id          UUID PRIMARY KEY,
    room    UUID NOT NULL REFERENCES rooms,
    customer_id UUID NOT NULL,
    requestor_id UUID NOT NULL,
    requested_at TIMESTAMP NOT NULL,
    start_time TIMESTAMP NOT NULL ,
    end_time TIMESTAMP NOT NULL,
    state states NOT NULL,
    state_information TEXT,
    file_name TEXT,
    description TEXT
);
ALTER TABLE bookings OWNER TO bookings ;