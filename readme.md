# 🛡️ Zero Trust ERP

A **high-performance Modular Monolithic ERP system** built in **Go**.
Inspired by the extensibility of **Odoo** and the structured routing of **Django**, this project is designed with a **Security-First (Zero Trust)** architecture.

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

* ⚙️ **Routing**: Central router for all apps
* 🔐 **Middleware**: Security, logging, request validation
* ⚠️ **Important**: Avoid modifying core logic unless necessary

Register your app routes in:

```
core/urls.go
```

---

### 2. Apps (`/apps`)

Each app is a **self-contained module**.

* **Controllers** → Request handling logic
* **Models** → Database schema (PostgreSQL, etc.)
* **Views** → HTML / templates
* **Security** → Zero Trust access control per module

---

### 3. Static Files (`/static`)

* Centralized storage for public assets
* Served securely via the Core router
* Protects against **path traversal attacks**

---

## 🚀 Getting Started

### 1. Environment Setup

Create a `.env` file in the root directory:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=zerotrust_erp
```

---

### 2. Create a New App

1. Create a folder inside `apps/`
2. Add your handlers in `controllers/`
3. Define routes in `urls.go`
4. Register routes in `core/urls.go`

---

### 3. Run the Server

```
go run main.go
```

Server will start at:

```
http://localhost:8080
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
  Every action requires explicit authentication

* 🧩 **Modular Isolation**
  Security policies are enforced per app

* 🍪 **Secure Authentication**
  Uses server-side secure cookies instead of vulnerable client-side tokens

---

## 💡 Vision

To build a **secure, scalable, and developer-friendly ERP framework** that combines:

* ⚡ Performance of Go
* 🧱 Modularity of Odoo
* 🧭 Structure of Django
* 🔐 Zero Trust Security

### CLI ###

* go run main.go migrate -> to migrate your model struct to postgresql database  
