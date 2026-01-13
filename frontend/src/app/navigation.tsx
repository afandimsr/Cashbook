import React from 'react';
import DashboardIcon from '@mui/icons-material/Dashboard';
import PeopleIcon from '@mui/icons-material/People';
import SettingsIcon from '@mui/icons-material/Settings';
import NotificationsIcon from '@mui/icons-material/Notifications';
import PersonIcon from '@mui/icons-material/Person';
import CategoryIcon from '@mui/icons-material/Category';
import ReceiptIcon from '@mui/icons-material/Receipt';
import BarChartIcon from '@mui/icons-material/BarChart';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import AccountBalanceIcon from '@mui/icons-material/AccountBalance';

export interface NavItem {
    text: string;
    icon: React.ReactNode;
    path?: string;
    children?: NavItem[];
    roles?: string[];
}

export const menuItems: NavItem[] = [
    {
        text: 'Dashboard',
        icon: <DashboardIcon />,
        path: '/dashboard'
    },
    {
        text: 'Transactions',
        icon: <ReceiptIcon />,
        path: '/dashboard/transactions'
    },
    {
        text: 'Budgets',
        icon: <AccountBalanceIcon />,
        path: '/dashboard/budgets'
    },
    {
        text: 'Recurring',
        icon: <AutorenewIcon />,
        path: '/dashboard/recurring'
    },
    {
        text: 'Reports',
        icon: <BarChartIcon />,
        path: '/dashboard/reports'
    },
    {
        text: 'Categories',
        icon: <CategoryIcon />,
        path: '/dashboard/categories'
    },
    {
        text: 'Users',
        icon: <PeopleIcon />,
        path: '/dashboard/users',
        roles: ['ADMIN']
    },
    {
        text: 'Settings',
        icon: <SettingsIcon />,
        children: [
            {
                text: 'Profile',
                icon: <PersonIcon />,
                path: '/dashboard/profile'
            },
            {
                text: 'Notifications',
                icon: <NotificationsIcon />,
                path: '/dashboard/settings/notifications'
            }
        ]
    }
];
