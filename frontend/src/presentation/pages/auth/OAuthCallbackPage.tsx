import React, { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuthStore } from '../../../state/authStore';
import { Box, CircularProgress, Typography, Stack } from '@mui/material';

export const OAuthCallbackPage: React.FC = () => {
    const [searchParams] = useSearchParams();
    const { handleOAuthSuccess } = useAuthStore();
    const navigate = useNavigate();

    useEffect(() => {
        const token = searchParams.get('token');
        if (token) {
            handleOAuthSuccess(token)
                .then(() => {
                    navigate('/dashboard');
                })
                .catch((err) => {
                    console.error('OAuth processing failed', err);
                    navigate('/login?error=oauth_failed');
                });
        } else {
            navigate('/login?error=token_missing');
        }
    }, [searchParams, handleOAuthSuccess, navigate]);

    return (
        <Box
            sx={{
                height: '100vh',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                backgroundColor: 'background.default'
            }}
        >
            <Stack spacing={2} alignItems="center">
                <CircularProgress size={40} />
                <Typography variant="h6">Authenticating...</Typography>
                <Typography variant="body2" color="text.secondary">
                    Please wait while we complete your sign-in
                </Typography>
            </Stack>
        </Box>
    );
};
