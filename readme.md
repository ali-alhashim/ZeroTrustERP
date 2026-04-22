# 🛡️ Zero Trust ERP

A **high-performance Modular Monolithic ERP system** built in **Go**.
Inspired by the extensibility of **Odoo** and the structured of **Django**, this project is designed with a **Security-First (Zero Trust)** architecture.

---

## 🏗️ Project Structure

The system is divided into a protected **Core** and extensible **Apps**:

```
.
├── main.go            # Entry point
├── .env               # Database & environment secrets (ignored by Git)
├── core/              # Engine (internal logic, middleware, global router)
│   └── urls.go        # app routes here
├── static/            # Global assets (CSS, JS, images)
└── apps/              # Business modules ("App layer")
    └── [app_name]/
        ├── controllers/  # HTTP handlers & business logic
        ├── models/       # Data structures & DB schema
        ├── views/        # HTML templates
        ├── security/     # App-specific ACLs & permissions 
        └── urls.go       # App route definitions
```

---

## 🛠️ Key Components

### 1. Core (`/core`)

The **Core** manages the application lifecycle.

* ⚙️ **urls**: Central router for all apps
* 🔐 **auth.go**: Security, logging, request validation
* ⚠️ **Important**: Avoid modifying core logic unless necessary

Register your app in:

```
main.go
under -> // Import app packages to register their routes & Models
```

---

### 2. Apps (`/apps`)

Each app is a **self-contained module**.

* **Controllers** → Request handling logic
* **Models** → Database schema (PostgreSQL)
* **urls** → apps routes
* **Views** → HTML / templates
* **Security** → Zero Trust access control per module
* **init** → register routes & models **this important if you not register the app will ignore the migration of your model and the routes **

---

### 3. Static Files (`/static`)

* Centralized storage for public assets
* Served securely via the Core router (urls.go)
* Protects against **path traversal attacks**

---

## 🚀 Getting Started


### Requirements
* Go https://go.dev/
* Postgres https://www.postgresql.org/
* Linux , Mac or Windows

### 1. Environment Setup

Create a `.env` file in the root directory:

```
# PostgreSQL Connection Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=zerotrustertb
DB_SSLMODE=disable
DB_TIMEZONE=UTC


# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=Your Email Address
SMTP_PASSWORD=Your Email App Password
SMTP_FROM=Your Email Address

# Application Configuration
sessionSecret=your-session-secret-key
```

---

### 2. Create a New App / if you want to add more you can check dev documentation 

1. Create a folder inside `apps/`
2. Add your handlers in `controllers/`
3. Define routes in `urls.go`
4. Register routes & Models in `init.go`

---

### 3. Database Migration

```
go run main.go migrate
```

### 4. Run the Server


```
go run main.go
```

Server will start at:

```
http://localhost:8000
```

---

## 🔒 Zero Trust & Minimal Dependency Philosophy
this project follows a **Minimal Dependency Architecture**:

- 🚫 **No third-party frameworks or heavy libraries**
- 🧩 Uses only Go standard library + SQL driver
- ⚙️ Full control over routing, middleware, and security
- 🔍 Easier auditing and bug bounty analysis
- ⚡ Better performance and lower overhead

### Core Principles

* 🔑 **Identity-Driven Access**
  Every action requires explicit authentication * permissions

* 🧩 **Modular Isolation**
  Security policies are enforced per app

* 🍪 **Secure Authentication**
  Uses  secure cookies [ HMAC instance using SHA256 and your secret key + email] so the validation required 2 factore token and the email for that token 

---

## 💡 Vision

To build a **secure, scalable, and developer-friendly ERP framework** that combines:

* ⚡ Performance of Go
* 🧱 Modularity of Odoo
* 🧭 Structure of Django
* 🔐 Zero Trust Security

### CLI ###

* go run main.go migrate -> to migrate your model struct to postgresql database  
* go run main.go migrate <appname>



### Minimum Viable Product (MVP) ###
* Login system -> Ok
* Users CRUD -> OK
* Roles & Permission
* Employees CRUD
* IT Ticketing System
* Simple dashboard
### Phase 2 ###
* Action Tracker
* QR Assets
* Leave Management
* Warehouses and Items
### Phase 3 ###
* Payroll & WPS File
* Purchase Orders And Sales Orders
* Clearance and End of Service
### Phase 4 ###
* Sales Invoice
* Bills
* Chart of Accounts
### Phase 5 ####
* ZATCA
* CEO Dashboard
* E-commerce
* Governance and Compliance