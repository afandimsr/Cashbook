import React, { useState, useEffect } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Typography,
    Box,
    CircularProgress,
    IconButton,
} from '@mui/material';
import TelegramIcon from '@mui/icons-material/Telegram';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import CheckIcon from '@mui/icons-material/Check';
import type { User } from '../../../domain/entities/User';
import type { TelegramStatus } from '../../../domain/entities/Telegram';
import { config } from '../../../app/config';

interface TelegramManageDialogProps {
    open: boolean;
    user: User | null;
    status: TelegramStatus | null;
    isLoading: boolean;
    onClose: () => void;
    onLoadStatus: (id: number) => Promise<void>;
    onGenerateCode: (id: number) => Promise<{ code: string; expires_at: string }>;
    onUnlink: (id: number) => Promise<void>;
}

export const TelegramManageDialog: React.FC<TelegramManageDialogProps> = ({
    open,
    user,
    status,
    isLoading,
    onClose,
    onLoadStatus,
    onGenerateCode,
    onUnlink,
}) => {
    const [code, setCode] = useState<{ code: string; expires_at: string } | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [copied, setCopied] = useState(false);
    const [confirmingUnlink, setConfirmingUnlink] = useState(false);

    useEffect(() => {
        if (open && user) {
            setCode(null);
            setError(null);
            setConfirmingUnlink(false);
            onLoadStatus(user.id).catch((err: any) => setError(err.message || 'Failed to load status'));
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [open, user]);

    const handleGenerateCode = async () => {
        if (!user) return;
        setError(null);
        try {
            const result = await onGenerateCode(user.id);
            setCode(result);
        } catch (err: any) {
            setError(err.message || 'Failed to generate link code');
        }
    };

    const handleCopy = () => {
        if (!code) return;
        navigator.clipboard.writeText(`/link ${code.code}`);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const handleUnlink = async () => {
        if (!user) return;
        setError(null);
        try {
            await onUnlink(user.id);
            setConfirmingUnlink(false);
        } catch (err: any) {
            setError(err.message || 'Failed to disconnect Telegram');
        }
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
            <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <TelegramIcon color="primary" />
                Manage Telegram
            </DialogTitle>
            <DialogContent dividers>
                <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                    Managing Telegram connection for <strong>{user?.name}</strong> ({user?.email})
                </Typography>

                {isLoading && !status && (
                    <Box sx={{ display: 'flex', justifyContent: 'center', py: 2 }}>
                        <CircularProgress size={24} />
                    </Box>
                )}

                {status && status.linked && (
                    <Box sx={{ mb: 2 }}>
                        <Typography variant="body2" sx={{ mb: 2 }}>
                            {status.telegram_username ? `Connected as @${status.telegram_username}` : 'Telegram is connected'}
                            {status.linked_at ? ` since ${new Date(status.linked_at).toLocaleDateString()}` : ''}.
                        </Typography>

                        {!confirmingUnlink ? (
                            <Button
                                variant="outlined"
                                color="error"
                                fullWidth
                                onClick={() => setConfirmingUnlink(true)}
                            >
                                Disconnect Telegram
                            </Button>
                        ) : (
                            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
                                <Typography variant="caption" color="error">
                                    Are you sure? This user will need a new code to reconnect.
                                </Typography>
                                <Box sx={{ display: 'flex', gap: 1 }}>
                                    <Button variant="text" onClick={() => setConfirmingUnlink(false)} fullWidth>
                                        Cancel
                                    </Button>
                                    <Button variant="contained" color="error" onClick={handleUnlink} fullWidth>
                                        Confirm Disconnect
                                    </Button>
                                </Box>
                            </Box>
                        )}
                    </Box>
                )}

                {status && !status.linked && !code && (
                    <Button
                        variant="contained"
                        fullWidth
                        onClick={handleGenerateCode}
                        startIcon={<TelegramIcon />}
                    >
                        Generate Link Code
                    </Button>
                )}

                {code && (
                    <Box>
                        <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                            Have the user open Telegram, message @{config.TELEGRAM_BOT_USERNAME}, and send:
                        </Typography>
                        <Box
                            sx={{
                                p: 1.5,
                                background: 'rgba(0,0,0,0.04)',
                                borderRadius: '8px',
                                fontFamily: 'monospace',
                                fontSize: '15px',
                                mb: 1,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'space-between',
                            }}
                        >
                            <Box sx={{ userSelect: 'all' }}>/link {code.code}</Box>
                            <IconButton size="small" onClick={handleCopy}>
                                {copied ? <CheckIcon fontSize="small" /> : <ContentCopyIcon fontSize="small" />}
                            </IconButton>
                        </Box>
                        <Typography variant="caption" color="text.secondary">
                            Expires at {new Date(code.expires_at).toLocaleTimeString()}.
                        </Typography>
                    </Box>
                )}

                {error && (
                    <Typography color="error" variant="body2" sx={{ mt: 2 }}>
                        {error}
                    </Typography>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={onClose} color="inherit">
                    Close
                </Button>
            </DialogActions>
        </Dialog>
    );
};
