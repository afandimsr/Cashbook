import React, { useState, useEffect } from 'react';
import {
    Box,
    Button,
    TextField,
    Typography,
    Alert,
    CircularProgress,
    Paper,
    Chip,
    Divider,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { AuthRepository } from '../../../infrastructure/auth/AuthRepository';
import { useAuthStore } from '../../../state/authStore';

const PRIMARY_COLOR = '#10b981';

const SetupCard = styled(Paper)({
    padding: '32px',
    borderRadius: '16px',
    maxWidth: '560px',
    margin: '0 auto',
});

const StepBadge = styled(Chip)({
    borderRadius: '12px',
    fontWeight: 700,
    fontSize: '12px',
    height: '28px',
});

const authRepo = new AuthRepository();

export const TwoFASetupPage: React.FC = () => {
    const { user } = useAuthStore();
    const [step, setStep] = useState<'idle' | 'setup' | 'verify' | 'backup' | 'done'>('idle');
    const [secret, setSecret] = useState('');
    const [qrCode, setQrCode] = useState('');
    const [code, setCode] = useState('');
    const [backupCodes, setBackupCodes] = useState<string[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

    useEffect(() => {
        if (user?.totpEnabled) {
            setStep('verify');
            setSuccess('2FA is already enabled!');
        }
    }, [user]);

    const handleSetup = async () => {
        setIsLoading(true);
        setError(null);
        try {
            const result = await authRepo.setup2FA();
            console.log(result);
            setSecret(result.secret);
            setQrCode(result.qr_code);
            setStep('setup');
        } catch (err: any) {
            setError(err.message || 'Failed to initiate 2FA setup');
        } finally {
            setIsLoading(false);
        }
    };

    const handleVerifySetup = async (e: React.FormEvent) => {
        e.preventDefault();
        setIsLoading(true);
        setError(null);
        try {
            await authRepo.verifySetup2FA(code);
            setStep('verify');
            setSuccess('2FA has been enabled successfully!');
        } catch (err: any) {
            setError(err.message || 'Invalid verification code');
        } finally {
            setIsLoading(false);
        }
    };

    const handleGenerateBackupCodes = async () => {
        setIsLoading(true);
        setError(null);
        try {
            const codes = await authRepo.generateBackupCodes();
            setBackupCodes(codes);
            setStep('backup');
        } catch (err: any) {
            setError(err.message || 'Failed to generate backup codes');
        } finally {
            setIsLoading(false);
        }
    };

    const handleDisable = async () => {
        setIsLoading(true);
        setError(null);
        try {
            await authRepo.disable2FA();
            setStep('idle');
            setSuccess('2FA has been disabled.');
            setSecret('');
            setQrCode('');
            setCode('');
            setBackupCodes([]);
        } catch (err: any) {
            setError(err.message || 'Failed to disable 2FA');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                Two-Factor Authentication
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
                Add an extra layer of security to your account with TOTP-based authentication.
            </Typography>

            {error && <Alert severity="error" sx={{ mb: 2, borderRadius: '12px' }}>{error}</Alert>}
            {success && <Alert severity="success" sx={{ mb: 2, borderRadius: '12px' }}>{success}</Alert>}

            {step === 'idle' && (
                <SetupCard elevation={2}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                        <span className="material-symbols-outlined" style={{ fontSize: 40, color: PRIMARY_COLOR }}>
                            shield
                        </span>
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 600 }}>Enable 2FA</Typography>
                            <Typography variant="body2" color="text.secondary">
                                Use an authenticator app like Google Authenticator or Authy
                            </Typography>
                        </Box>
                    </Box>
                    <Button
                        variant="contained"
                        onClick={handleSetup}
                        disabled={isLoading}
                        startIcon={isLoading ? <CircularProgress size={20} /> : undefined}
                        sx={{
                            background: `linear-gradient(135deg, ${PRIMARY_COLOR}, #059669)`,
                            borderRadius: '12px',
                            textTransform: 'none',
                            fontWeight: 600,
                            py: 1.5,
                            px: 4,
                        }}
                    >
                        Set Up 2FA
                    </Button>
                    <Divider sx={{ my: 3 }} />
                    <Button
                        variant="outlined"
                        color="error"
                        onClick={handleDisable}
                        disabled={isLoading}
                        sx={{ borderRadius: '12px', textTransform: 'none' }}
                    >
                        Disable 2FA
                    </Button>
                </SetupCard>
            )}

            {step === 'setup' && (
                <SetupCard elevation={2}>
                    <StepBadge label="Step 1 of 2" color="primary" size="small" sx={{ mb: 2 }} />
                    <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                        Scan QR Code
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                        Scan this QR code with your authenticator app, or manually enter the secret key.
                    </Typography>

                    {qrCode && (
                        <Box sx={{ display: 'flex', justifyContent: 'center', mb: 3 }}>
                            <Box
                                sx={{
                                    p: 2,
                                    background: '#fff',
                                    borderRadius: '16px',
                                    display: 'inline-block',
                                }}
                            >
                                <img src={`data:image/png;base64,${qrCode}`} alt="QR Code" width={200} height={200} />
                            </Box>
                        </Box>
                    )}

                    <Typography variant="caption" sx={{ display: 'block', mb: 1, color: 'text.secondary' }}>
                        Manual entry key:
                    </Typography>
                    <Box
                        sx={{
                            p: 1.5,
                            background: 'rgba(0,0,0,0.04)',
                            borderRadius: '8px',
                            fontFamily: 'monospace',
                            fontSize: '14px',
                            mb: 3,
                            wordBreak: 'break-all',
                            userSelect: 'all',
                        }}
                    >
                        {secret}
                    </Box>

                    <StepBadge label="Step 2 of 2" color="primary" size="small" sx={{ mb: 2 }} />
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        Enter the 6-digit code from your authenticator app to confirm:
                    </Typography>
                    <form onSubmit={handleVerifySetup}>
                        <TextField
                            fullWidth
                            placeholder="000000"
                            value={code}
                            onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                            inputProps={{ maxLength: 6, style: { textAlign: 'center', fontSize: '24px', letterSpacing: '8px' } }}
                            sx={{
                                mb: 2,
                                '& .MuiOutlinedInput-root': { borderRadius: '12px' },
                            }}
                        />
                        <Button
                            type="submit"
                            variant="contained"
                            fullWidth
                            disabled={code.length !== 6 || isLoading}
                            sx={{
                                background: `linear-gradient(135deg, ${PRIMARY_COLOR}, #059669)`,
                                borderRadius: '12px',
                                textTransform: 'none',
                                fontWeight: 600,
                                py: 1.5,
                            }}
                        >
                            {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Verify & Enable'}
                        </Button>
                    </form>
                </SetupCard>
            )}

            {step === 'verify' && (
                <SetupCard elevation={2}>
                    <Box sx={{ textAlign: 'center', mb: 3 }}>
                        <span className="material-symbols-outlined" style={{ fontSize: 56, color: PRIMARY_COLOR }}>
                            check_circle
                        </span>
                        <Typography variant="h6" sx={{ fontWeight: 600, mt: 1 }}>
                            2FA Enabled!
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                            Your account is now protected with two-factor authentication.
                        </Typography>
                    </Box>
                    <Button
                        variant="contained"
                        onClick={handleGenerateBackupCodes}
                        disabled={isLoading}
                        fullWidth
                        sx={{
                            background: `linear-gradient(135deg, ${PRIMARY_COLOR}, #059669)`,
                            borderRadius: '12px',
                            textTransform: 'none',
                            fontWeight: 600,
                            py: 1.5,
                        }}
                    >
                        Generate Backup Codes
                    </Button>
                </SetupCard>
            )}

            {step === 'backup' && (
                <SetupCard elevation={2}>
                    <Typography variant="h6" sx={{ fontWeight: 600, mb: 1 }}>
                        Backup Codes
                    </Typography>
                    <Alert severity="warning" sx={{ mb: 3, borderRadius: '12px' }}>
                        Save these codes in a secure place. Each code can only be used once.
                    </Alert>
                    <Box
                        sx={{
                            display: 'grid',
                            gridTemplateColumns: 'repeat(2, 1fr)',
                            gap: 1.5,
                            mb: 3,
                        }}
                    >
                        {backupCodes.map((bc, i) => (
                            <Box
                                key={i}
                                sx={{
                                    p: 1.5,
                                    background: 'rgba(0,0,0,0.04)',
                                    borderRadius: '8px',
                                    fontFamily: 'monospace',
                                    fontSize: '15px',
                                    fontWeight: 600,
                                    textAlign: 'center',
                                    userSelect: 'all',
                                }}
                            >
                                {bc}
                            </Box>
                        ))}
                    </Box>
                    <Button
                        variant="outlined"
                        onClick={() => setStep('idle')}
                        fullWidth
                        sx={{ borderRadius: '12px', textTransform: 'none', fontWeight: 600 }}
                    >
                        Done
                    </Button>
                </SetupCard>
            )}
        </Box>
    );
};
