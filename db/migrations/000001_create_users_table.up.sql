CREATE TABLE IF NOT EXISTS users (
   id serial PRIMARY KEY,
   username VARCHAR (50) NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL,
   site VARCHAR (300) NOT NULL
);
