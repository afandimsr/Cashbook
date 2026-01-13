import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box } from '@mui/material';

export const AuthLayout: React.FC = () => {
    return (
        <Box
            sx={{
                minHeight: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                bgcolor: 'background.default',
                width: '100%',
                position: 'relative',
                overflow: 'hidden'
            }}
        >
            <Outlet />
        </Box>
    );
};
