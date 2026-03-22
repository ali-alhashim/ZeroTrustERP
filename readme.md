Zero Trust ERP 🛡️
A high-performance, Modular Monolithic ERP system built in Go. Inspired by the extensibility of Odoo and the structured routing of Django, this project is designed with a "Security-First" (Zero Trust) architecture.

🏗️ Project Structure
The system is divided into a protected Core and extensible Apps:

Plaintext
├── main.go            # Entry point
├── .env               # Database & Environment secrets (ignored by git)
├── core/              # The Engine (Internal logic, Middleware, Global Router)
│   └── routes.go      # Register your app URLs here
├── static/            # Global assets (CSS, JS, Images)
└── apps/              # Business Modules (The "App" layer)
    └── [app_name]/
        ├── controllers/ # HTTP Handlers & Logic
        ├── models/      # Data structures & Database schema
        ├── views/       # HTML Templates
        └── security/    # App-specific ACLs and Permissions
🛠️ Key Components
1. The Core (/core)
The Core handles the lifecycle of the application. To maintain system stability, you should generally not modify the core logic.

Registration: Use core/routes.go to hook your app's router into the main system.

2. Apps (/apps)
Each folder inside apps/ is a self-contained module.

Controllers: Contains the Go logic for processing requests.

Models: Defines your PostgreSQL/Database structs.

Views: Pure HTML/Templates for the frontend.

Security: Implements Zero Trust principles (Identity-aware access) specific to that module.

3. Static Files (/static)
Centralized storage for all public-facing assets. Served securely via the Core router to prevent path traversal.

🚀 Getting Started
1. Environment Setup
Create a .env file in the root directory:

Bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=zerotrust_erp
2. Creating a New App
Create a new folder in apps/.

Define your handlers in controllers/.

Register the path in core/routes.go.

3. Running the Server
Bash
go run main.go
The server will start at http://localhost:8080.

🔒 Zero Trust Philosophy
Unlike traditional ERPs that trust any user on the internal network, Zero Trust ERP assumes every request is a potential threat.

Identity-Defined: Every action requires explicit authentication.

Modular Isolation: Security policies are defined at the app level.

Secure Auth: Defaults to server-side secure cookies rather than vulnerable client tokens.