# SI Template - Modern React Admin Dashboard

A professional, feature-rich, and highly customizable React admin dashboard template built with modern best practices. Designed to provide a premium user experience with a focus on aesthetics and performance.

## ✨ Key Features

- **🌙 Dynamic Dark Mode**: Seamless theme switching with persistent state management (Zustand).
- **📂 Nested Sidebar Menus**: Support for hierarchical navigation structures with smooth collapse/expand transitions.
- **📱 Fully Responsive**: Optimized for Mobile, Tablet, and Desktop screens using a modern flexible layout.
- **👤 Profile Management**: Comprehensive profile page featuring loading skeletons and responsive sections.
- **🔔 Notification Center**: Categorized alert system (Success, Info, Warning, Error) with management actions.
- **🔐 Smart Routing**: Protected routes for authenticated users and automatic redirection for guest access.
- **🎨 Premium Aesthetics**: Modern design language utilizing Glassmorphism, smooth transitions, and curated color palettes.
- **🏛️ Clean Architecture**: Organized codebase following separation of concerns (Presentation, App, Infrastructure, State).

## 🛠️ Tech Stack

- **Framework**: [React 18+](https://reactjs.org/)
- **Build Tool**: [Vite](https://vitejs.dev/)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **UI Components**: [Material UI (MUI)](https://mui.com/)
- **State Management**: [Zustand](https://github.com/pmndrs/zustand)
- **Routing**: [React Router v6](https://reactrouter.com/)
- **Icons**: [MUI Icons](https://mui.com/material-ui/material-icons/)

## 🚀 Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm or yarn

### Installation

1. Clone the repository
2. Install dependencies:
```bash
npm install
```

### Development

Run the development server:
```bash
npm run dev
```

### Build

Build for production:
```bash
npm run build
```

## 📂 Project Structure

```text
src/
├── app/                # App-wide config, routes, and navigation
│   ├── providers/     # React Context providers
│   ├── App.tsx        # Root App component
│   ├── navigation.tsx # Sidebar menu configuration
│   └── routes.tsx     # Application route definitions
├── application/        # Application logic and use cases
├── domain/             # Entities, interfaces, and business logic
├── infrastructure/     # API services and external integrations
│   └── api/           # Axios/Fetch API implementations
├── presentation/       # UI Components, Layouts, and Pages
│   ├── components/    # Reusable UI atoms and molecules
│   │   ├── auth/      # Auth-related guards (ProtectedRoute, GuestRoute)
│   │   ├── ui/        # Shared UI components (Topbar, Sidebar, Footer)
│   │   └── utils/     # Utility components (ScrollToTop)
│   ├── layouts/       # Master layouts (Dashboard, Auth)
│   └── pages/         # Individual feature pages
│       ├── auth/      # Login and authentication pages
│       ├── dashboard/ # Main showcase dashboard
│       ├── notifications/ # Notification center
│       ├── profile/   # User profile and settings
│       └── users/     # User management pages
├── state/              # Global state stores (Zustand)
├── theme/              # MUI theme configuration and constants
└── shared/             # Shared utilities and constants
```

## 📝 License

This project is licensed under the Apache License 2.0.
