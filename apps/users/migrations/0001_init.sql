CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE,
	username TEXT UNIQUE NOT NULL,
	active BOOLEAN ,
	otphash TEXT ,
	otpexpiry INTEGER ,
	online BOOLEAN ,
	roles TEXT 
);

CREATE TABLE IF NOT EXISTS roles (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT ,
	permissions TEXT ,
	createdat TIMESTAMP ,
	updatedat TIMESTAMP 
);

CREATE TABLE IF NOT EXISTS permissions (
	id SERIAL PRIMARY KEY,
	resource TEXT UNIQUE NOT NULL,
	action TEXT UNIQUE NOT NULL,
	description TEXT ,
	createdat TIMESTAMP ,
	updatedat TIMESTAMP 
);