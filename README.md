# CashBook - Personal Finance Manager

A modern, full-stack personal finance management application built with **Go (Golang)** and **React**. CashBook helps users track transactions, manage categories, set budgets, and visualize financial growth through insightful reports.

## 🚀 Features

- **Financial Dashboard**: Overview of current balance, income, expenses, and recent activities.
- **Transaction Management**: Comprehensive tracking of all income and expenses with search and filtering.
- **Category Management**: Organize transactions with customizable categories and visual indicators (colors/icons).
- **Budgeting System**: Set monthly spending limits per category and monitor progress in real-time.
- **Recurring Transactions**: Automate your repetitive bills and subscriptions.
- **Financial Reports**: Interactive charts and data breakdown for spending analysis (powered by Recharts).
- **Dual Authentication**: Traditional Username/Password login and Google OAuth integration.
- **Two-Factor Authentication (2FA)**: TOTP-based authentication with QR code setup, backup codes, and admin-enforced MFA.
- **Progressive Web App (PWA)**: Installable on mobile and desktop devices with offline support and fast loading.

## 🛠️ Technology Stack

### Backend
- **Go 1.21+**: Core programming language.
- **Gin Web Framework**: High-performance HTTP routing.
- **PostgreSQL**: Robust relational database.
- **Clean Architecture**: Domain-driven design with clear separation of Concerns (Entities, Usecases, Repositories).
- **JWT**: Secure token-based authentication.

### Frontend
- **React 19**: Modern UI library.
- **TypeScript**: Static typing for enhanced developer experience.
- **Material UI (MUI) v6**: Professional component library for high-end aesthetics.
- **Zustand**: Lightweight, scalable state management.
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
- [Go](https://golang.org/dl/) (1.21 or higher)
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
   go run cmd/main.go
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

## 🐳 Docker Deployment

The application is fully containerized with Docker Compose for easy deployment.

### Rebuild and Deploy
To rebuild the entire stack (including any PWA changes):
```bash
docker compose --env-file ./frontend/.env up -d --build
```

### Specific Service Update
To update only the frontend (e.g., for PWA updates):
```bash
docker compose --env-file ./frontend/.env up -d --build frontend
```

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
