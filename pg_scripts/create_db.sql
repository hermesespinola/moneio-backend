CREATE TABLE IF NOT EXISTS bills(
    serialCode CHAR(8) PRIMARY KEY NOT NULL,
    denomination INT NOT NULL
        CHECK(denomination = ANY('{20, 50, 100, 500, 1000}'::int[]))
);

CREATE TABLE IF NOT EXISTS billEntry(
    id SERIAL UNIQUE,
    serialCode CHAR(8) NOT NULL
        REFERENCES bills(serialCode),
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    notes VARCHAR(255)
);
