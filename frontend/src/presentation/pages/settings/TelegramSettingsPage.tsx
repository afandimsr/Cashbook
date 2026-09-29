import React, { useState, useEffect } from 'react';
import {
    Box,
    Button,
    Typography,
    Alert,
    CircularProgress,
    Paper,
    IconButton,
    Snackbar,
    Divider,
} from '@mui/material';
import { styled } from '@mui/material/styles';
import TelegramIcon from '@mui/icons-material/Telegram';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import CheckIcon from '@mui/icons-material/Check';
import { TelegramRepository } from '../../../infrastructure/telegram/TelegramRepository';
import type { TelegramLinkCodeResponse } from '../../../domain/entities/Telegram';
import { config } from '../../../app/config';

const PRIMARY_COLOR = '#10b981';

const SetupCard = styled(Paper)({
    padding: '32px',
    borderRadius: '16px',
    maxWidth: '560px',
    margin: '0 auto',
});

const telegramRepo = new TelegramRepository();

export const TelegramSettingsPage: React.FC = () => {
    const [step, setStep] = useState<'idle' | 'loading-status' | 'not-connected' | 'code-shown' | 'connected'>('loading-status');
    const [linkCode, setLinkCode] = useState<TelegramLinkCodeResponse | null>(null);
    const [username, setUsername] = useState<string | undefined>(undefined);
    const [linkedAt, setLinkedAt] = useState<string | undefined>(undefined);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);

    const loadStatus = async () => {
        setStep('loading-status');
        setError(null);
        try {
            const status = await telegramRepo.getStatus();
            if (status.linked) {
                setUsername(status.telegram_username);
                setLinkedAt(status.linked_at);
                setStep('connected');
            } else {
                setStep('not-connected');
            }
        } catch (err: any) {
            setError(err.message || 'Failed to load Telegram status');
            setStep('not-connected');
        }
    };

    useEffect(() => {
        loadStatus();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const handleGenerateCode = async () => {
        setIsLoading(true);
        setError(null);
        try {
            const result = await telegramRepo.generateLinkCode();
            setLinkCode(result);
            setStep('code-shown');
        } catch (err: any) {
            setError(err.message || 'Failed to generate link code');
        } finally {
            setIsLoading(false);
        }
    };

    const handleCopyCode = () => {
        if (!linkCode) return;
        navigator.clipboard.writeText(`/link ${linkCode.code}`);
        setCopied(true);
    };

    const handleDisconnect = async () => {
        setIsLoading(true);
        setError(null);
        try {
            await telegramRepo.unlink();
            setLinkCode(null);
            setStep('not-connected');
        } catch (err: any) {
            setError(err.message || 'Failed to disconnect Telegram');
        } finally {
            setIsLoading(false);
        }
    };

    const expiresInMinutes = linkCode
        ? Math.max(0, Math.round((new Date(linkCode.expires_at).getTime() - Date.now()) / 60000))
        : 0;

    return (
        <Box sx={{ p: 3 }}>
            <Typography variant="h5" sx={{ fontWeight: 700, mb: 1 }}>
                Telegram
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
                Connect your Telegram account to log transactions and get reports via chat.
            </Typography>

            {error && <Alert severity="error" sx={{ mb: 2, borderRadius: '12px' }}>{error}</Alert>}

            {step === 'loading-status' && (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
                    <CircularProgress size={28} />
                </Box>
            )}

            {step === 'not-connected' && (
                <SetupCard elevation={2}>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mb: 3 }}>
                        <TelegramIcon sx={{ fontSize: 40, color: PRIMARY_COLOR }} />
                        <Box>
                            <Typography variant="h6" sx={{ fontWeight: 600 }}>Connect Telegram</Typography>
                            <Typography variant="body2" color="text.secondary">
                                Log expenses and check your balance straight from a Telegram chat.
                            </Typography>
                        </Box>
                    </Box>
                    <Button
                        variant="contained"
                        onClick={handleGenerateCode}
                        disabled={isLoading}
                        startIcon={isLoading ? <CircularProgress size={20} /> : <TelegramIcon />}
                        sx={{
                            background: `linear-gradient(135deg, ${PRIMARY_COLOR}, #059669)`,
                            borderRadius: '12px',
                            textTransform: 'none',
                            fontWeight: 600,
                            py: 1.5,
                            px: 4,
                        }}
                    >
                        Connect Telegram
                    </Button>
                </SetupCard>
            )}

            {step === 'code-shown' && linkCode && (
                <SetupCard elevation={2}>
                    <Typography variant="h6" sx={{ fontWeight: 600, mb: 2 }}>
                        Almost there
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
                        Open Telegram, message{' '}
                        <strong>@{config.TELEGRAM_BOT_USERNAME}</strong>, and send:
                    </Typography>

                    <Box
                        sx={{
                            p: 1.5,
                            background: 'rgba(0,0,0,0.04)',
                            borderRadius: '8px',
                            fontFamily: 'monospace',
                            fontSize: '16px',
                            mb: 1,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'space-between',
                            gap: 1,
                        }}
                    >
                        <Box sx={{ userSelect: 'all' }}>/link {linkCode.code}</Box>
                        <IconButton size="small" onClick={handleCopyCode} sx={{ color: PRIMARY_COLOR }}>
                            {copied ? <CheckIcon fontSize="small" /> : <ContentCopyIcon fontSize="small" />}
                        </IconButton>
                    </Box>
                    <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mb: 3 }}>
                        This code expires in about {expiresInMinutes} minute{expiresInMinutes === 1 ? '' : 's'}.
                    </Typography>

                    <Divider sx={{ my: 2 }} />

                    <Button
                        variant="outlined"
                        onClick={loadStatus}
                        fullWidth
                        sx={{ borderRadius: '12px', textTransform: 'none', fontWeight: 600 }}
                    >
                        I've sent it — check status
                    </Button>
                </SetupCard>
            )}

            {step === 'connected' && (
                <SetupCard elevation={2}>
                    <Box sx={{ textAlign: 'center', mb: 3 }}>
                        <TelegramIcon sx={{ fontSize: 56, color: PRIMARY_COLOR }} />
                        <Typography variant="h6" sx={{ fontWeight: 600, mt: 1 }}>
                            Connected
                        </Typography>
                        <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                            {username ? `Linked to @${username}` : 'Your Telegram account is linked'}
                            {linkedAt ? ` on ${new Date(linkedAt).toLocaleDateString()}` : ''}.
                        </Typography>
                    </Box>
                    <Button
                        variant="outlined"
                        color="error"
                        onClick={handleDisconnect}
                        disabled={isLoading}
                        fullWidth
                        sx={{ borderRadius: '12px', textTransform: 'none', fontWeight: 600 }}
                    >
                        {isLoading ? <CircularProgress size={20} /> : 'Disconnect'}
                    </Button>
                </SetupCard>
            )}

            <Snackbar
                open={copied}
                autoHideDuration={2000}
                onClose={() => setCopied(false)}
                message="Copied to clipboard"
                anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
            />
        </Box>
    );
};
