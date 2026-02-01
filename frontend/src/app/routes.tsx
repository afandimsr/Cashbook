import React from 'react';
import { useRoutes, Navigate } from 'react-router-dom';
import { DashboardLayout } from '../presentation/layouts/DashboardLayout';
import { OverviewPage } from '../presentation/pages/dashboard/overview/OverviewPage';
import { AuthLayout } from '../presentation/layouts/AuthLayout';
import { LoginPageNew } from '../presentation/pages/auth/LoginPage';
import { ProtectedRoute } from '../presentation/components/auth/ProtectedRoute';
import { GuestRoute } from '../presentation/components/auth/GuestRoute';
import { AdminRoute } from '../presentation/components/auth/AdminRoute';

import { UserListPage } from '../presentation/pages/users';
import { ProfilePage } from '../presentation/pages/profile/ProfilePage';
import { NotificationPage } from '../presentation/pages/notifications/NotificationPage';
import { DebugAuthPage } from '../presentation/pages/debug/AuthDebug';
import { OAuthCallbackPage } from '../presentation/pages/auth/OAuthCallbackPage';
import { CategoryPage } from '../presentation/pages/dashboard/categories/CategoryPage';
import { TransactionPage } from '../presentation/pages/dashboard/transactions/TransactionPage';
import { BudgetPage } from '../presentation/pages/dashboard/budgets/BudgetPage';
import { ReportsPage } from '../presentation/pages/dashboard/reports/ReportsPage';
import { RecurringPage } from '../presentation/pages/dashboard/recurring/RecurringPage';
import config from './config';

export const AppRoutes: React.FC = () => {
    const routes = [
        {
            path: '/',
            element: <Navigate to="/dashboard" replace />
        },
        {
            element: <GuestRoute />, // Prevent authenticated users from login
            children: [
                {
                    path: '/login',
                    element: <AuthLayout />,
                    children: [
                        { index: true, element: <LoginPageNew /> }
                    ]
                },
                {
                    path: '/oauth/callback',
                    element: <OAuthCallbackPage />
                }
            ]
        },
        // Only include debug route in development
        ...(config.IS_DEV ? [
            {
                path: '/debug',
                element: <DebugAuthPage />
            }
        ] : []),
        {
            path: '/',
            element: <ProtectedRoute />, // Protect these routes
            children: [
                {
                    path: 'dashboard',
                    element: <DashboardLayout />,
                    children: [
                        { index: true, element: <OverviewPage /> },
                        { path: 'transactions', element: <TransactionPage /> },
                        { path: 'budgets', element: <BudgetPage /> },
                        { path: 'reports', element: <ReportsPage /> },
                        { path: 'recurring', element: <RecurringPage /> },
                        { path: 'categories', element: <CategoryPage /> },
                        {
                            path: 'users',
                            element: <AdminRoute />,
                            children: [
                                { index: true, element: <UserListPage /> }
                            ]
                        },
                        { path: 'profile', element: <ProfilePage /> },
                        { path: 'settings/notifications', element: <NotificationPage /> },
                    ],
                }
            ]
        },
        { path: '*', element: <Navigate to="/login" replace /> }
    ];
    const element = useRoutes(routes);
    return element;
};
