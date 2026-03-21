CREATE TABLE IF NOT EXISTS cars (
    id BIGSERIAL PRIMARY KEY,
    make VARCHAR NOT NULL,
    model VARCHAR NOT NULL,
    year SMALLINT NOT NULL,
    description TEXT,
    image_url TEXT,
    gearbox TEXT CHECK (gearbox IN ('manual', 'automatic', 'DCT', 'CVT')),
    drivetrain TEXT CHECK (drivetrain IN ('FWD', 'RWD', 'AWD', '4WD')),
    horsepower int,
    fuel TEXT CHECK (fuel IN ('diesel', 'gas', 'electric', 'hybrid', 'plug-in-hybrid', 'hydrogen', 'lpg', 'cng')),
    price_new NUMERIC(12, 2),
    version int DEFAULT 1 NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL
)