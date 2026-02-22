import React, { useState, useRef, useEffect } from 'react';
import {
    Box,
    Button,
    TextField,
    Typography,
    Alert,
    CircularProgress,
    Link,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuthStore } from '../../../state/authStore';
import { useNavigate } from 'react-router-dom';

const PRIMARY_COLOR = '#10b981';
// const PRIMARY_COLOR = '#ffffff';
const PRIMARY_DARK = '#059669';
const BG_COLOR = '#01281fff';

const PageContainer = styled(Box)({
    minHeight: '100vh',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    background: `linear-gradient(135deg, ${BG_COLOR} 0%, #0a3d32 50%, #062d24 100%)`,
    fontFamily: "'Roboto', sans-serif",
});

const GlassCard = styled(Box)({
    background: 'rgba(255, 255, 255, 0.9)',
    backdropFilter: 'blur(20px)',
    borderRadius: '24px',
    border: '1px solid rgba(255, 255, 255, 0.1)',
    padding: '48px',
    maxWidth: '440px',
    width: '100%',
    margin: '24px',
    textAlign: 'center',
    boxShadow: '0 32px 64px rgba(0, 0, 0, 0.3)',
});

const OTPInput = styled(TextField)({
    '& .MuiInputBase-input': {
        textAlign: 'center',
        fontSize: '28px',
        fontWeight: 600,
        letterSpacing: '12px',
        padding: '16px',
        color: '#000',
    },
    '& .MuiOutlinedInput-root': {
        borderRadius: '16px',
        backgroundColor: 'rgba(0, 0, 0, 0.06)',
        '& fieldset': { borderColor: 'rgba(0, 0, 0, 0.15)' },
        '&:hover fieldset': { borderColor: PRIMARY_COLOR },
        '&.Mui-focused fieldset': { borderColor: PRIMARY_COLOR },
    },
});

const VerifyButton = styled(Button)({
    background: `linear-gradient(135deg, ${PRIMARY_COLOR} 0%, ${PRIMARY_DARK} 100%)`,
    color: '#000',
    fontWeight: 600,
    fontSize: '16px',
    padding: '14px',
    borderRadius: '16px',
    textTransform: 'none',
    marginTop: '16px',
    '&:hover': {
        background: `linear-gradient(135deg, ${PRIMARY_DARK} 0%, #047857 100%)`,
        transform: 'translateY(-1px)',
        boxShadow: `0 8px 24px ${PRIMARY_COLOR}40`,
    },
    '&:disabled': {
        background: 'rgba(0, 0, 0, 0.1)',
        color: 'rgba(0, 0, 0, 0.3)',
    },
});

export const TwoFAVerifyPage: React.FC = () => {
    const [code, setCode] = useState('');
    const [showBackup, setShowBackup] = useState(false);
    const [backupCode, setBackupCode] = useState('');
    const { verify2FA, verifyBackupCode, isLoading, error, requires2FA, clear2FAState } = useAuthStore();
    const navigate = useNavigate();
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        if (!requires2FA) {
            navigate('/login', { replace: true });
        }
    }, [requires2FA, navigate]);

    useEffect(() => {
        inputRef.current?.focus();
    }, [showBackup]);

    const handleVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await verify2FA(code);
            navigate('/dashboard', { replace: true });
        } catch {
            // Error is handled in store
        }
    };

    const handleBackupVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await verifyBackupCode(backupCode);
            navigate('/dashboard', { replace: true });
        } catch {
            // Error is handled in store
        }
    };

    const handleCancel = () => {
        clear2FAState();
        navigate('/login', { replace: true });
    };

    return (
        <PageContainer>
            <GlassCard >
                <Box sx={{ mb: 3 }}>
                    <span className="material-symbols-outlined" style={{ fontSize: 56, color: PRIMARY_COLOR }}>
                        security
                    </span>
                </Box>

                <Typography variant="h5" sx={{ color: '#000', fontWeight: 700, mb: 1 }}>
                    Two-Factor Authentication
                </Typography>
                <Typography variant="body2" sx={{ color: 'rgba(0,0,0,0.6)', mb: 4 }}>
                    {showBackup
                        ? 'Enter one of your backup codes'
                        : 'Enter the 6-digit code from your authenticator app'}
                </Typography>

                {error && (
                    <Alert severity="error" sx={{ mb: 3, borderRadius: '12px' }}>
                        {error}
                    </Alert>
                )}

                {!showBackup ? (
                    <form onSubmit={handleVerify}>
                        <OTPInput
                            inputRef={inputRef}
                            fullWidth
                            placeholder="000000"
                            value={code}
                            onChange={(e) => {
                                const val = e.target.value.replace(/\D/g, '').slice(0, 6);
                                setCode(val);
                            }}
                            inputProps={{ maxLength: 6, autoComplete: 'one-time-code' }}
                        />
                        <VerifyButton
                            type="submit"
                            fullWidth
                            disabled={code.length !== 6 || isLoading}
                            startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                        >
                            {isLoading ? 'Verifying...' : 'Verify Code'}
                        </VerifyButton>
                    </form>
                ) : (
                    <form onSubmit={handleBackupVerify}>
                        <OTPInput
                            inputRef={inputRef}
                            fullWidth
                            placeholder="xxxx-xxxx"
                            value={backupCode}
                            onChange={(e) => setBackupCode(e.target.value)}
                            sx={{
                                '& .MuiInputBase-input': {
                                    letterSpacing: '4px',
                                    fontSize: '22px',
                                },
                            }}
                        />
                        <VerifyButton
                            type="submit"
                            fullWidth
                            disabled={backupCode.length < 5 || isLoading}
                            startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                        >
                            {isLoading ? 'Verifying...' : 'Use Backup Code'}
                        </VerifyButton>
                    </form>
                )}

                <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center', gap: 2 }}>
                    <Link
                        component="button"
                        onClick={() => {
                            setShowBackup(!showBackup);
                            setCode('');
                            setBackupCode('');
                        }}
                        sx={{
                            color: PRIMARY_COLOR,
                            fontSize: '14px',
                            textDecoration: 'none',
                            cursor: 'pointer',
                            '&:hover': { textDecoration: 'underline' },
                        }}
                    >
                        {showBackup ? 'Use authenticator app' : 'Use backup code'}
                    </Link>
                    <Typography sx={{ color: 'rgba(0,0,0,0.2)' }}>|</Typography>
                    <Link
                        component="button"
                        onClick={handleCancel}
                        sx={{
                            color: 'rgba(0,0,0,0.5)',
                            fontSize: '14px',
                            textDecoration: 'none',
                            cursor: 'pointer',
                            '&:hover': { textDecoration: 'underline', color: 'rgba(0,0,0,0.7)' },
                        }}
                    >
                        Cancel
                    </Link>
                </Box>
            </GlassCard>
        </PageContainer>
    );
};
