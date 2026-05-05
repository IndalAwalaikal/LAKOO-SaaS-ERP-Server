# Lakoo SaaS - Server Ecosystem ⚙️

Lakoo SaaS Server is the robust backend infrastructure for a sophisticated, multi-tenant Enterprise Resource Planning (ERP) platform designed for Small and Medium Enterprises (SMEs). The system employs a distributed micro-architecture pattern, ensuring high availability, bulletproof security, and fault isolation between the core ERP transaction logic and the data-intensive AI analytics services.

---

## 🏛️ System Architecture

The backend consists of two primary services that coordinate with several highly-optimized data layers:

1. **Golang Core API**: The primary engine handling authentication, Role-Based Access Control (RBAC), multi-tenant CRUD operations, and external storage communication. Built using Domain-Driven Design (DDD) principles.
2. **Python AI Service**: A specialized Fast API microservice that consumes transaction data to generate demand projections and statistical sales insights using Pandas.
3. **Data Layer**:
   - **MySQL 8.0 (Relational)**: Stores transactional entities, isolated strictly by tenant ID.
   - **Redis 7.0 (Cache/Session)**: Manages volatile sessions, rate-limiting, and high-speed data caching.
   - **MinIO (Object Storage)**: S3-compatible, self-hosted storage for binary assets (Logos, Receipts, QRIS images).

---

## 🛠️ Tech Stack

| Component | Technology | Description |
| :--- | :--- | :--- |
| **Core API** | Golang 1.21+, Gin | High-performance compiled language for core business logic. |
| **AI Service** | Python 3.10+, FastAPI | High-speed Python framework for data-intensive processing. |
| **Analysis** | Pandas, Scikit-learn | Advanced data manipulation and statistical modelling. |
| **Database** | MySQL 8.0 | Robust relational storage for transactional integrity. |
| **Cache** | Redis 7.0 | In-memory caching for speed and rate limiting. |
| **Storage** | MinIO | Self-hosted S3-compatible object storage. |
| **Orchestration**| Docker, Docker Compose | Unified container orchestration for easy deployment. |

---

## 📂 Directory Structure

```text
server/
├── api/                     # Golang Backend (Clean Architecture)
│   ├── cmd/api/main.go      # Dependency Injection & Bootstrapper
│   ├── internal/            # Core Business Logic (Encapsulated)
│   │   ├── domain/          # Shared Interfaces & Models
│   │   ├── dto/             # Request/Response Contract Validations
│   │   ├── usecase/         # Pure Logical Orchestration
│   │   ├── repository/      # MySQL Persistence Implementation
│   │   ├── http/            # Delivery Layer (Handlers & Routes)
│   │   └── middleware/      # Security layers (RBAC, Rate-Limit, JWT)
│   ├── pkg/                 # Internal Generic Packages (Redis, MinIO)
│   └── migrations/          # SQL Versioning (Schema evolvement)
│
├── ai-service/              # Python Analytical Layer
│   ├── main.py              # Service orchestration & Routing
│   ├── services/            # Statistical computing logic
│   └── models/              # Pre-calculated data & regression structures
│
├── docker-compose.yml       # Multi-container orchestration logic
├── .env                     # Centralized environment configurations
└── .env.example             # Template for environment variables
```

---

## 🛡️ Security and Hardening Specifications

### **Brute-force Mitigation**
Powered by a Redis-based Rate Limiter. Attempts to access the authentication endpoints are tracked by IP and Email identity. 
- **Limit**: MAX 5 attempts per 15 Minutes.
- **Action**: Automatic block with `429 Too Many Requests` response.

### **Global Security Headers**
The following headers are injected into every API response via `SecurityMiddleware`:
- `X-Frame-Options: DENY`: Prevents Clickjacking.
- `Strict-Transport-Security`: Enforces HTTPS (HSTS).
- `Content-Security-Policy`: Restricts resource loading to trusted origins.
- `X-Content-Type-Options: nosniff`: Prevents MIME-type sniffing.

### **Data Protection & Multi-Tenancy**
- **Password Enforcement**: Minimum 8 characters required (Bcrypt hashed at DB).
- **Tenant Isolation**: Every SQL query and data access object is strictly scoped with a `tenant_id` to prevent cross-tenant data leaks.
- **Role-Based Access Control**: Strict middleware checking for Owner, Manager, and Cashier clearance levels before executing handler logic.

---

## 🚀 Getting Started (Docker Installation)

The absolute easiest way to run the entire backend infrastructure is through Docker Compose. This ensures you do not have to manually configure MySQL, Redis, or MinIO.

### Prerequisites
- Docker Engine
- Docker Compose

### Running the Services
1. Copy `.env.example` to `.env` and fill in the required passwords and secrets.
   ```bash
   cp .env.example .env
   ```
2. Build and start the infrastructure in detached mode:
   ```bash
   docker-compose up -d --build
   ```
3. Verify that the services are running. The following endpoints will be available:
   - **Core API**: `http://localhost:8080`
   - **AI-Service**: `http://localhost:8000`
   - **MySQL**: `localhost:3306` (mapped to `3307` externally if configured)
   - **MinIO Console**: `http://localhost:9001`
   - **Redis**: `localhost:6379`

---

## 📜 License

This project is licensed under the **MIT License**. You are free to use, modify, and distribute this software in compliance with the license terms.

---

## ✨ Closing Words & Acknowledgements

Thank you for visiting the Lakoo SaaS Server repository! Building a reliable backend is no easy feat, but we designed this architecture to be as scalable and secure as possible for enterprise use. We hope this codebase serves as a solid foundation for your operations or as a learning resource for implementing Clean Architecture in Golang alongside AI microservices.

If you find this project valuable, please consider giving it a star! ⭐ We welcome any issues, feature requests, or pull requests to make Lakoo even better.

*Developed with ❤️ by **Indal Awalaikal**.*
