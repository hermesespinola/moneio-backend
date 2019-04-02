CREATE TABLE IF NOT EXISTS bills(
    serialCode CHAR(8) PRIMARY KEY NOT NULL,
    denomination INT NOT NULL
        CHECK(denomination = ANY('{20, 50, 100, 500, 1000}'::int[]))
);

-- Maybe add image path
CREATE TABLE IF NOT EXISTS billEntries(
    id SERIAL UNIQUE,
    creationDate DATE NOT NULL DEFAULT CURRENT_DATE,
    serialCode CHAR(8) NOT NULL REFERENCES bills(serialCode),
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    notes VARCHAR(255)
);
