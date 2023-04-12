CREATE TABLE IF NOT EXISTS domains (
    id INT GENERATED ALWAYS AS IDENTITY,
    fqdn TEXT NOT NULL UNIQUE,
    update_at TIMESTAMP NOT NULL,
    update_delay INTERVAL NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO domains (
    fqdn, 
    update_at, 
    update_delay
)
VALUES
    ('yandex.ru.', '2023-04-12 10:10:10', '1W'),
    ('ozon.ru.', '2023-03-01 00:00:01', '1M'),
    ('mail.ru.', '2023-04-04 12:12:12', '1Y');

CREATE TABLE IF NOT EXISTS ipv4_addresses (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    ip VARCHAR(20) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO ipv4_addresses (
    domain_id,
    ip
)
VALUES
    (1, '5.255.255.70'),
    (1, '77.88.55.88'),
    (1, '5.255.255.77'),
    (1, '77.88.55.60'),
    (2, '162.159.128.64');

CREATE TABLE IF NOT EXISTS ipv6_addresses (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    ip VARCHAR(50) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO ipv6_addresses (
    domain_id,
    ip
)
VALUES
    (1, '2a02:6b8:a::a'),
    (3, '2a00:1148:db00:0:b0b0::1'),
    (3, '3b00:1148:da00:0:b0b0::2');

CREATE TABLE IF NOT EXISTS canonical_names (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    canonical_name TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO canonical_names (
    domain_id,
    canonical_name
)
VALUES
    (1, 'kontora.ru');

CREATE TABLE IF NOT EXISTS mail_exchangers (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    host TEXT NOT NULL,
    pref INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO mail_exchangers (
    domain_id,
    host,
    pref
)
VALUES
    (1, 'mx.yandex.ru', 10),
    (2, 'mxs.mail.ru', 20);

CREATE TABLE IF NOT EXISTS name_servers (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    name_server TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO name_servers (
    domain_id,
    name_server
)
VALUES
    (1, 'ns1.yandex.ru'),
    (1, 'ns2.yandex.ru'),
    (2, 'ns1.mail.ru'),
    (2, 'ns2.mail.ru'),
    (3, 'ns8-l2.nic.ru');

CREATE TABLE IF NOT EXISTS server_selections (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    target TEXT NOT NULL,
    port INT NOT NULL,
    priority INT NOT NULL,
    weight INT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO server_selections (
    domain_id,
    target,
    port,
    priority,
    weight
)
VALUES
    (1, 'sipserver.yandex.ru', 1234, 0, 5);

CREATE TABLE IF NOT EXISTS text_strings (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    text TEXT NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO text_strings (
    domain_id,
    text
)
VALUES
    (1, 'MS=ms75457885'),
    (2, 'yandex-verification: 2bccd09858554e85'),
    (3, 'mailru-verification: 43e81e646c1675e5');


CREATE TABLE IF NOT EXISTS registrations (
    id INT GENERATED ALWAYS AS IDENTITY,
    domain_id INT NOT NULL,
    created TIMESTAMP NOT NULL,
    paid_till TIMESTAMP NOT NULL,

    PRIMARY KEY (id),
    FOREIGN KEY (domain_id)
        REFERENCES domains (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

INSERT INTO registrations (
    domain_id, 
    created, 
    paid_till
)
VALUES
    (1, '2000-04-12 10:10:10', '2023-12-12 23:00:00'),
    (2, '2001-03-01 00:00:01', '2023-08-15 20:00:00'),
    (3, '1999-04-04 12:12:12', '2023-06-20 12:30:00');
