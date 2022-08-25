-- script ran to initialize the db

CREATE DATABASE dyna WITH CONNECTION LIMIT 200;
CREATE USER dyna WITH PASSWORD 'board';
GRANT ALL ON DATABASE dyna TO dyna;
