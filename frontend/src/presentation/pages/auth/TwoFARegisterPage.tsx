import React, { useState } from 'react';
import {
    Box,
    Button,
    TextField,
    Typography,
    Alert,
    CircularProgress,
    Paper,
    Chip,
    IconButton,
    Snackbar,
    Dialog,
    DialogContent,
    DialogActions,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { AuthRepository } from '../../../infrastructure/auth/AuthRepository';
import { useNavigate } from 'react-router-dom';
import { tokenStorage } from '../../../infrastructure/storage/tokenStorage';
import { safeDecodeTempJwt } from '../../../infrastructure/auth/jwt.service';

const PRIMARY_COLOR = '#10b981';
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

const GlassCard = styled(Paper)({
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

const SetupButton = styled(Button)({
    background: `linear-gradient(135deg, ${PRIMARY_COLOR} 0%, ${PRIMARY_DARK} 100%)`,
    color: '#fff',
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
        background: 'rgba(255, 255, 255, 0.1)',
        color: 'rgba(255, 255, 255, 0.3)',
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
    width: '100%',
    '&:hover': {
        background: `linear-gradient(135deg, ${PRIMARY_DARK} 0%, #047857 100%)`,
    },
    '&:disabled': {
        background: 'rgba(0, 0, 0, 0.1)',
        color: 'rgba(0, 0, 0, 0.3)',
    },
});

const authRepo = new AuthRepository();

export const TwoFARegisterPage: React.FC = () => {
    const navigate = useNavigate();
    const [step, setStep] = useState<'intro' | 'scan' | 'verify'>('intro');
    const [secret, setSecret] = useState('');
    const [qrCode, setQrCode] = useState('');
    const [code, setCode] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);
    const [successDialogOpen, setSuccessDialogOpen] = useState(false);

    const handleCopySecret = () => {
        navigator.clipboard.writeText(secret);
        setCopied(true);
    };

    const handleSetup = async () => {
        const tempToken = tokenStorage.getToken();
        const payload = safeDecodeTempJwt(tempToken || '');

        if (!payload || !payload.purpose) {
            navigate('/login');
            return;
        }

        setIsLoading(true);
        setError(null);
        try {
            const result = await authRepo.setup2FA();
            setSecret(result.secret);
            setQrCode(result.qr_code);
            setStep('scan');
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Failed to initiate 2FA setup';
            if (
                errorMessage.toLowerCase().includes('expired') ||
                errorMessage.toLowerCase().includes('unauthorized') ||
                errorMessage.toLowerCase().includes('token')
            ) {
                navigate('/login');
            } else {
                setError(errorMessage || 'Failed to initiate 2FA setup');
            }
        } finally {
            setIsLoading(false);
        }
    };

    const handleVerify = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);
        try {
            await authRepo.verifySetup2FA(code);
            tokenStorage.clearToken();
            setSuccessDialogOpen(true);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Invalid verification code';
            if (
                errorMessage.toLowerCase().includes('expired') ||
                errorMessage.toLowerCase().includes('invalid') ||
                errorMessage.toLowerCase().includes('unauthorized') ||
                errorMessage.toLowerCase().includes('token')
            ) {
                navigate('/login');
            } else {
                setError(errorMessage || 'Invalid verification code');
            }
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <PageContainer>
            <GlassCard elevation={0}>
                <Box sx={{ mb: 3 }}>
                    <span className="material-symbols-outlined" style={{ fontSize: 56, color: PRIMARY_COLOR }}>
                        security
                    </span>
                </Box>

                <Typography variant="h5" sx={{ color: '#000', fontWeight: 700, mb: 1 }}>
                    Setup Two-Factor Authentication
                </Typography>
                <Typography variant="body2" sx={{ color: 'rgba(0,0,0,0.6)', mb: 4 }}>
                    Your administrator requires you to enable 2FA for your account.
                </Typography>

                {error && (
                    <Alert severity="error" sx={{ mb: 3, borderRadius: '12px', bgcolor: 'rgba(211,47,47,0.1)' }}>
                        {error}
                    </Alert>
                )}

                {step === 'intro' && (
                    <>
                        <Typography variant="body2" sx={{ color: 'rgba(0,0,0,0.7)', mb: 3 }}>
                            Add an extra layer of security to your account with TOTP-based authentication.
                        </Typography>
                        <SetupButton
                            fullWidth
                            onClick={handleSetup}
                            disabled={isLoading}
                            startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                        >
                            {isLoading ? 'Setting up...' : 'Start 2FA Setup'}
                        </SetupButton>
                    </>
                )}

                {step === 'scan' && (
                    <>
                        <Chip
                            label="Step 1 of 2"
                            color="primary"
                            size="small"
                            sx={{ mb: 2, borderRadius: '8px' }}
                        />
                        <Typography variant="body2" sx={{ color: 'rgba(0,0,0,0.7)', mb: 2 }}>
                            Scan this QR code with your authenticator app:
                        </Typography>

                        {qrCode && (
                            <Box sx={{ display: 'flex', justifyContent: 'center', mb: 3 }}>
                                <Box sx={{ p: 2, background: '#fff', borderRadius: '16px' }}>
                                    <img src={`data:image/png;base64,${qrCode}`} alt="QR Code" width={180} height={180} />
                                </Box>
                            </Box>
                        )}

                        <Typography variant="caption" sx={{ display: 'block', mb: 1, color: 'rgba(0,0,0,0.5)' }}>
                            Or enter this key manually:
                        </Typography>
                        <Box sx={{
                            p: 1.5,
                            background: 'rgba(255,255,255,0.05)',
                            borderRadius: '8px',
                            fontFamily: 'monospace',
                            fontSize: '12px',
                            mb: 3,
                            wordBreak: 'break-all',
                            color: 'rgba(0,0,0,0.7)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'space-between',
                            gap: 1
                        }}>
                            <Box sx={{ flex: 1, overflow: 'hidden' }}>{secret}</Box>
                            <IconButton
                                size="small"
                                onClick={handleCopySecret}
                                sx={{ color: PRIMARY_COLOR, flexShrink: 0 }}
                            >
                                <span className="material-symbols-outlined">
                                    {copied ? 'check' : 'content_copy'}
                                </span>
                            </IconButton>
                        </Box>

                        <Chip
                            label="Step 2 of 2"
                            color="primary"
                            size="small"
                            sx={{ mb: 2, borderRadius: '8px' }}
                        />
                        <Typography variant="body2" sx={{ color: 'rgba(0,0,0,0.7)', mb: 2 }}>
                            Enter the 6-digit code to verify:
                        </Typography>
                        <form onSubmit={handleVerify}>
                            <TextField
                                fullWidth
                                placeholder="000000"
                                value={code}
                                disabled={successDialogOpen}
                                onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                                inputProps={{
                                    maxLength: 6,
                                    style: {
                                        textAlign: 'center',
                                        fontSize: '24px',
                                        letterSpacing: '8px',
                                        color: '#000'
                                    }
                                }}
                                sx={{
                                    mb: 2,
                                    '& .MuiOutlinedInput-root': {
                                        borderRadius: '12px',
                                        bgcolor: 'rgba(0,0,0,0.05)',
                                        '& fieldset': { borderColor: 'rgba(0,0,0,0.15)' },
                                        '&:hover fieldset': { borderColor: PRIMARY_COLOR },
                                        '&.Mui-focused fieldset': { borderColor: PRIMARY_COLOR },
                                    },
                                }}
                            />
                            <VerifyButton
                                type="submit"
                                disabled={code.length !== 6 || isLoading || successDialogOpen}
                                startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
                            >
                                {isLoading ? 'Verifying...' : successDialogOpen ? 'Verified!' : 'Verify & Continue'}
                            </VerifyButton>
                        </form>
                    </>
                )}
            </GlassCard>

            <Snackbar
                open={copied}
                autoHideDuration={2000}
                onClose={() => setCopied(false)}
                message="Secret copied to clipboard"
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            />

            <Dialog
                open={successDialogOpen}
                maxWidth="xs"
                fullWidth
                PaperProps={{
                    sx: { borderRadius: 3, textAlign: 'center', p: 2 }
                }}
            >
                <DialogContent>
                    <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', py: 2 }}>
                        <Box sx={{
                            width: 80,
                            height: 80,
                            borderRadius: '50%',
                            backgroundColor: 'rgba(16, 185, 129, 0.1)',
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            mb: 2
                        }}>
                            <span className="material-symbols-outlined" style={{ fontSize: 48, color: PRIMARY_COLOR }}>
                                check_circle
                            </span>
                        </Box>
                        <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                            2FA Enabled!
                        </Typography>
                        <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                            Your account is now protected with two-factor authentication.
                        </Typography>
                        <Typography variant="body2" sx={{ color: 'text.secondary', mt: 2 }}>
                            Please login again with your authenticator app.
                        </Typography>
                    </Box>
                </DialogContent>
                <DialogActions sx={{ justifyContent: 'center', pb: 2 }}>
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={() => navigate('/login')}
                        sx={{
                            background: `linear-gradient(135deg, ${PRIMARY_COLOR}, #059669)`,
                            borderRadius: '12px',
                            textTransform: 'none',
                            fontWeight: 600,
                            py: 1.5,
                            mx: 3,
                            maxWidth: 300
                        }}
                    >
                        Back to Login
                    </Button>
                </DialogActions>
            </Dialog>
        </PageContainer>
    );
};
