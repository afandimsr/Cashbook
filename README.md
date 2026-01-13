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
- **Recharts**: Modular charting components.

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

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
