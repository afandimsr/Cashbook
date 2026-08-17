# CashBook - Personal Finance Manager

A modern, full-stack personal finance management application built with **Go (Golang)** and **React**. CashBook helps users track transactions, manage categories, set budgets, and visualize financial growth through insightful reports.

## 🚀 Features

- **Financial Dashboard**: Overview of current balance, income, expenses, and recent activities.
- **Transaction Management**: Comprehensive tracking of all income and expenses with search and filtering.
- **Category Management**: Organize transactions with customizable categories and visual indicators (colors/icons).
- **Budgeting System**: Set monthly spending limits per category and monitor progress in real-time.
- **Recurring Transactions**: Automate your repetitive bills and subscriptions.
- **Shared Expenses (Split Bill)**: Split a bill among registered users and ad-hoc "shadow" participants, with multiple split methods (equal, exact amounts, percentage, per-item), tax / service charge / other charge / discount, per-item categories, and per-participant settlement tracking that posts to the ledger.
- **Financial Reports**: Interactive charts and data breakdown for spending analysis (powered by Recharts).
- **Dual Authentication**: Traditional Username/Password login and Google OAuth integration.
- **Two-Factor Authentication (2FA)**: TOTP-based authentication with QR code setup, backup codes, and admin-enforced MFA.
- **Progressive Web App (PWA)**: Installable on mobile and desktop devices with offline support and fast loading.

## 🛠️ Technology Stack

### Backend
- **Go 1.25**: Core programming language.
- **Gin Web Framework**: High-performance HTTP routing.
- **PostgreSQL**: Robust relational database.
- **Clean Architecture**: Domain-driven design with clear separation of Concerns (Entities, Usecases, Repositories).
- **JWT**: Secure token-based authentication.
- **Testing**: Unit tests with `testify` (mocked repositories) and repository tests with `go-sqlmock`.

### Frontend
- **React 19**: Modern UI library.
- **TypeScript**: Static typing for enhanced developer experience.
- **Material UI (MUI) v7**: Professional component library for high-end aesthetics.
- **Zustand**: Lightweight, scalable state management.
- **React Router v7**: Client-side routing.
- **Vite**: Ultra-fast build tool and development server.
- **Vite PWA**: Professional PWA integration for installation and offline support.
- **Recharts**: Modular charting components.

## 🔐 Two-Factor Authentication (2FA)

CashBook supports TOTP-based Two-Factor Authentication for enhanced security.

### Features
- **TOTP Authentication**: Time-based one-time passwords using authenticator apps (Google Authenticator, Authy, etc.)
- **QR Code Setup**: Easy scanning of QR codes to set up 2FA
- **Backup Codes**: Generate 10 one-time backup codes for account recovery
- **Admin Enforcement**: Administrators can require all users to enable 2FA system-wide

### Login Flow with 2FA
1. User enters email/password
2. If 2FA is enabled, user is prompted to enter TOTP code or backup code
3. After verification, user gains access to dashboard

### Admin 2FA Settings
- Navigate to `/dashboard/user/mfa-settings` to enforce 2FA for all users
- Users without 2FA enabled will be prompted to set it up on next login

## 🧾 Shared Expenses (Split Bill)

Split a shared bill and track who owes what, at `/dashboard/shared-expenses`. Amounts are handled in whole Rupiah (IDR).

### Split methods
- **Equal** — divide the subtotal evenly.
- **Exact** — assign a specific amount per participant (must sum to the subtotal).
- **Percentage** — assign a percentage per participant (must sum to 100%).
- **Per-item** — add line items, assign each item to the participants who share it, and (optionally) give each item its own category.

### Charges
- **Tax**, **Service charge**, **Other charge**, and **Discount** — entered as a fixed amount (Rp) or a percentage. They are distributed proportionally across participants; the payer absorbs any rounding remainder so the shares always sum to the total.

### Ledger integration
- Creating a split records the payer's expense in the transaction ledger. With per-item categories, the expense is split into one transaction per category (allocated proportionally). The transaction note includes the bill title and item names.
- Marking a participant as **paid** records the settlement as an income for the payer and an expense for the registered participant, and flips the bill to `SETTLED` once everyone has paid.

## 🧪 Testing

Backend unit tests cover the use-case layer (mocked repositories), JWT/TOTP, HTTP handlers (`httptest`), and repositories (`go-sqlmock`).

```bash
cd backend
go test ./...          # run all tests
go test ./... -cover   # with coverage per package
```

## 📁 System Architecture

The project follows **Clean Architecture** and **Screaming Architecture** principles:

```text
CashBook/
├── backend/                # Go Gin Server
│   ├── internal/
│   │   ├── domain/         # Entities & Interfaces
│   │   ├── usecase/        # Business Logic
│   │   ├── repository/     # Data Persistence
│   │   └── delivery/       # HTTP Handlers & Middlewares
├── frontend/               # React Vite Client
│   ├── src/
│   │   ├── application/    # Custom Hooks (Logic)
│   │   ├── domain/         # Entities & Types
│   │   ├── presentation/   # Pages, Layouts & Components
│   │   └── state/          # Global State Store
```

## ⚙️ Getting Started

### Prerequisites
- [Go](https://golang.org/dl/) (1.25 or higher)
- [Node.js](https://nodejs.org/) (18 or higher)
- [PostgreSQL](https://www.postgresql.org/download/)

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Create your environment file:
   ```bash
   cp .env.example .env
   ```
3. Configure your database and Google OAuth credentials in `.env`.
4. Install dependencies:
   ```bash
   go mod tidy
   ```
5. Running factory seeder to create default Admin:
   ```bash
   go run cmd/create_admin/main.go
   ```
6. Run the server:
   ```bash
   go run cmd/api/main.go
   ```

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Run the development server:
   ```bash
   npm run dev
   ```

## 🔑 Demo Access
- **Default Admin**: `admin@example.com` / `admin123`
- **OAuth**: Click the "Sign in with Google" button (requires valid credentials in backend `.env`).

## 🐳 Deployment

The application is fully containerized with Docker Compose. We provide a helper script to automate the deployment process.

### Automated Deployment (Recommended)
The `deploy.sh` script handles network creation, building images, and running migrations automatically.

1.  **Configure Environment**:
    Make sure your root `.env` is configured (see `.env.example`).
2.  **Run Deploy Script**:
    ```bash
    chmod +x deploy.sh
    ./deploy.sh
    ```

### Manual Docker Deployment
If you prefer to run commands manually:

1.  **Build and Start Services**:
    ```bash
    docker compose --env-file ./frontend/.env up -d --build backend frontend
    ```
2.  **Run Migrations**:
    ```bash
    docker compose --env-file .env run --rm migrate
    ```

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
