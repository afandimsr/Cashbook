import React, { useEffect, useState } from 'react';
import { Box, Typography, Paper, Button } from '@mui/material';
import { useAuthStore } from '../../../state/authStore';
import { tokenStorage } from '../../../infrastructure/storage/tokenStorage';

export const DebugAuthPage: React.FC = () => {
    const authState = useAuthStore();
    const [localStorageToken, setLocalStorageToken] = useState<string | null>(null);

    useEffect(() => {
        const interval = setInterval(() => {
            setLocalStorageToken(tokenStorage.getToken());
        }, 1000);
        return () => clearInterval(interval);
    }, []);

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h4" gutterBottom>Auth Debug Panel</Typography>

            <Paper sx={{ p: 2, mb: 2 }}>
                <Typography variant="h6">Zustand State:</Typography>
                <pre>{JSON.stringify(authState, null, 2)}</pre>
            </Paper>

            <Paper sx={{ p: 2, mb: 2 }}>
                <Typography variant="h6">localStorage Token:</Typography>
                <Typography>{localStorageToken || 'null'}</Typography>
            </Paper>

            <Button
                variant="contained"
                onClick={() => {
                    console.log('Current authStore state:', useAuthStore.getState());
                    console.log('localStorage token:', tokenStorage.getToken());
                }}
            >
                Log to Console
            </Button>
        </Box>
    );
};
