CREATE TABLE IF NOT EXISTS departments (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	code TEXT UNIQUE NOT NULL,
	manager_id INT UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS job_titles (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	code TEXT UNIQUE NOT NULL,
	description TEXT ,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS employees (
	id SERIAL PRIMARY KEY,
	badge_id TEXT UNIQUE NOT NULL,
	name TEXT NOT NULL,
	department_id INT,
	local_name TEXT ,
	job_title_id INT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE,
	username TEXT UNIQUE NOT NULL,
	active BOOLEAN DEFAULT TRUE,
	otp_hash TEXT ,
	session_token TEXT ,
	session_expiry TIMESTAMP ,
	otp_expiry TIMESTAMP ,
	last_login TIMESTAMP ,
	incorrect_otp_attempts INTEGER ,
	online BOOLEAN DEFAULT FALSE,
	related_employee_id INT UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
	id SERIAL PRIMARY KEY,
	name TEXT UNIQUE NOT NULL,
	description TEXT ,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
	id SERIAL PRIMARY KEY,
	resource TEXT NOT NULL,
	action TEXT NOT NULL,
	description TEXT ,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS logs (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	username TEXT NOT NULL,
	email TEXT NOT NULL,
	resource TEXT NOT NULL,
	action TEXT NOT NULL,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_roles (
    user_id INT,
    role_id INT,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS roles_permissions (
    role_id INT,
    permission_id INT,
    PRIMARY KEY (role_id, permission_id)
);

ALTER TABLE departments ADD CONSTRAINT fk_departments_manager_id FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE CASCADE;

ALTER TABLE employees ADD CONSTRAINT fk_employees_department_id FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE;

ALTER TABLE employees ADD CONSTRAINT fk_employees_job_title_id FOREIGN KEY (job_title_id) REFERENCES job_titles(id) ON DELETE CASCADE;

ALTER TABLE users_roles ADD CONSTRAINT fk_users_roles_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE users_roles ADD CONSTRAINT fk_users_roles_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE users ADD CONSTRAINT fk_users_related_employee_id FOREIGN KEY (related_employee_id) REFERENCES employees(id) ON DELETE CASCADE;

ALTER TABLE roles_permissions ADD CONSTRAINT fk_roles_permissions_role_id FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE;

ALTER TABLE roles_permissions ADD CONSTRAINT fk_roles_permissions_permission_id FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE;