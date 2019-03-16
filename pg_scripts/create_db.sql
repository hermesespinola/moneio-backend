CREATE TABLE bills(
    serialCode CHAR(8) PRIMARY KEY NOT NULL
);

CREATE TABLE billEntry(
    id INT PRIMARY KEY NOT NULL,
    serialCode CHAR(8) NOT NULL
        REFERENCES bills(serialCode),
    denomination INT NOT NULL
        CHECK(denomination = ANY('{20, 50, 100, 500, 1000}'::int[])),
    latitude FLOAT NOT NULL,
    longitude FLOAT NOT NULL,
    notes VARCHAR(255)
);
