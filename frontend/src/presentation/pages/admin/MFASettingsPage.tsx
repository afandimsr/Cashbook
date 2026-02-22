import React, { useState, useEffect } from 'react';
import {
    Box,
    Typography,
    Switch,
    Alert,
    CircularProgress,
    Paper,
    FormControlLabel,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { apiClient } from '../../../infrastructure/apiClient';
import type { MFASettings } from '../../../domain/entities/User';

const PRIMARY_COLOR = '#10b981';

const SettingsCard = styled(Paper)({
    padding: '32px',
    borderRadius: '16px',
    maxWidth: '560px',
    margin: '0 auto',
});

export const MFASettingsPage: React.FC = () => {
    const [settings, setSettings] = useState<MFASettings | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [isSaving, setIsSaving] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        fetchSettings();
    }, []);

    const fetchSettings = async () => {
        try {
            const data = await apiClient.get<MFASettings>('/user/mfa-settings');
            setSettings(data);
        } catch (err: any) {
            setError(err.message || 'Failed to load MFA settings');
        } finally {
            setIsLoading(false);
        }
    };

    const handleToggle = async (enforce2fa: boolean) => {
        setIsSaving(true);
        setError(null);
        setSuccess(null);
        try {
            await apiClient.put('/user/mfa-settings', { enforce_2fa: enforce2fa });
            setSettings((prev) => prev ? { ...prev, enforce_2fa: enforce2fa } : prev);
            setSuccess(`2FA enforcement ${enforce2fa ? 'enabled' : 'disabled'} system-wide.`);
        } catch (err: any) {
            setError(err.message || 'Failed to update MFA settings');
        } finally {
            setIsSaving(false);
        }
    };

    if (isLoading) {
        return (
            <Box sx={{ display: 'flex', justifyContent: 'center', p: 8 }}>
                <CircularProgress sx={{ color: PRIMARY_COLOR }} />
            </Box>
        );
    }

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                MFA Security Policy
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
                Configure system-wide multi-factor authentication enforcement.
            </Typography>

            {error && <Alert severity="error" sx={{ mb: 2, borderRadius: '12px' }}>{error}</Alert>}
            {success && <Alert severity="success" sx={{ mb: 2, borderRadius: '12px' }}>{success}</Alert>}

            <SettingsCard elevation={2}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                    <span className="material-symbols-outlined" style={{ fontSize: 40, color: PRIMARY_COLOR }}>
                        admin_panel_settings
                    </span>
                    <Box>
                        <Typography variant="h6" sx={{ fontWeight: 600 }}>Enforce 2FA</Typography>
                        <Typography variant="body2" color="text.secondary">
                            When enabled, all users will be required to set up two-factor authentication.
                        </Typography>
                    </Box>
                </Box>

                <FormControlLabel
                    control={
                        <Switch
                            checked={settings?.enforce_2fa || false}
                            onChange={(e) => handleToggle(e.target.checked)}
                            disabled={isSaving}
                            sx={{
                                '& .MuiSwitch-switchBase.Mui-checked': {
                                    color: PRIMARY_COLOR,
                                },
                                '& .MuiSwitch-switchBase.Mui-checked + .MuiSwitch-track': {
                                    backgroundColor: PRIMARY_COLOR,
                                },
                            }}
                        />
                    }
                    label={
                        <Typography variant="body1" sx={{ fontWeight: 500 }}>
                            {settings?.enforce_2fa ? 'Enforced' : 'Not enforced'}
                        </Typography>
                    }
                />

                {isSaving && (
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 2 }}>
                        <CircularProgress size={16} sx={{ color: PRIMARY_COLOR }} />
                        <Typography variant="caption" color="text.secondary">Saving...</Typography>
                    </Box>
                )}

                <Alert severity="info" sx={{ mt: 3, borderRadius: '12px' }}>
                    Users who haven't set up 2FA will still be able to log in but will be encouraged to set it up.
                </Alert>
            </SettingsCard>
        </Box>
    );
};
