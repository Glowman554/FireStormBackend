CREATE TABLE packages (
    owner TEXT NOT NULL,
    package TEXT PRIMARY KEY,
    date_updated timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    package TEXT NOT NULL,
    content TEXT NOT NULL,
    PRIMARY KEY(name, version),
    CONSTRAINT fk_package FOREIGN KEY (package) REFERENCES packages (package) ON DELETE CASCADE ON UPDATE CASCADE
);
