import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuthStore } from '../../../state/authStore';
import { Box, CircularProgress } from '@mui/material';

export const AdminRoute: React.FC = () => {
    const { isAuthenticated, isInitialized, user } = useAuthStore();

    if (!isInitialized) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                <CircularProgress />
            </Box>
        );
    }

    if (!isAuthenticated) {
        return <Navigate to="/login" replace />;
    }

    const roles = user?.roles || (user?.role ? [user.role] : []);
    if (!roles.includes('ADMIN')) {
        return <Navigate to="/dashboard" replace />;
    }

    return <Outlet />;
};
