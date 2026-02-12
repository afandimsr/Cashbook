import React, { useState, useEffect } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Typography,
    Box,
    LinearProgress,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
} from '@mui/material';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import LockResetIcon from '@mui/icons-material/LockReset';
import type { User } from '../../../domain/entities/User';

interface ResetPasswordDialogProps {
    open: boolean;
    user: User | null;
    onClose: () => void;
    onReset: (password: string) => Promise<void>;
}

export const ResetPasswordDialog: React.FC<ResetPasswordDialogProps> = ({
    open,
    user,
    onClose,
    onReset,
}) => {
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const validations = {
        length: password.length >= 8,
        uppercase: /[A-Z]/.test(password),
        lowercase: /[a-z]/.test(password),
        number: /[0-9]/.test(password),
        special: /[!@#$%^&*(),.?":{}|<>-]/.test(password),
        match: password === confirmPassword && password !== '',
    };

    const strength = Object.values(validations).filter(v => v).length - (validations.match ? 1 : 0);
    const strengthPercentage = (strength / 5) * 100;

    const isValid = Object.values(validations).every(v => v);

    useEffect(() => {
        if (open) {
            setPassword('');
            setConfirmPassword('');
            setError(null);
        }
    }, [open]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!isValid) return;

        setIsLoading(true);
        setError(null);
        try {
            await onReset(password);
            onClose();
        } catch (err: any) {
            setError(err.message || 'Failed to reset password');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
            <form onSubmit={handleSubmit}>
                <DialogTitle sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <LockResetIcon color="primary" />
                    Reset Password
                </DialogTitle>
                <DialogContent dividers>
                    <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
                        Resetting password for <strong>{user?.name}</strong> ({user?.email})
                    </Typography>

                    <TextField
                        fullWidth
                        label="New Password"
                        type="password"
                        margin="normal"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                        disabled={isLoading}
                    />

                    <Box sx={{ mt: 1, mb: 2 }}>
                        <LinearProgress
                            variant="determinate"
                            value={strengthPercentage}
                            color={strength <= 2 ? 'error' : strength <= 4 ? 'warning' : 'success'}
                            sx={{ height: 6, borderRadius: 3 }}
                        />
                        <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
                            Password Strength: {strength <= 2 ? 'Weak' : strength <= 4 ? 'Medium' : 'Strong'}
                        </Typography>
                    </Box>

                    <TextField
                        fullWidth
                        label="Confirm New Password"
                        type="password"
                        margin="normal"
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        required
                        disabled={isLoading}
                        error={confirmPassword !== '' && !validations.match}
                        helperText={confirmPassword !== '' && !validations.match ? 'Passwords do not match' : ''}
                    />

                    <Box sx={{ mt: 2, bgcolor: 'action.hover', p: 1.5, borderRadius: 2 }}>
                        <Typography variant="subtitle2" gutterBottom>Requirements:</Typography>
                        <List dense sx={{ py: 0 }}>
                            <ValidationItem label="At least 8 characters" met={validations.length} />
                            <ValidationItem label="At least one uppercase letter" met={validations.uppercase} />
                            <ValidationItem label="At least one lowercase letter" met={validations.lowercase} />
                            <ValidationItem label="At least one number" met={validations.number} />
                            <ValidationItem label="At least one special character" met={validations.special} />
                        </List>
                    </Box>

                    {error && (
                        <Typography color="error" variant="body2" sx={{ mt: 2 }}>
                            {error}
                        </Typography>
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={onClose} color="inherit" disabled={isLoading}>
                        Cancel
                    </Button>
                    <Button
                        type="submit"
                        variant="contained"
                        disabled={!isValid || isLoading}
                        startIcon={isLoading ? null : <LockResetIcon />}
                    >
                        {isLoading ? 'Resetting...' : 'Reset Password'}
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};

const ValidationItem: React.FC<{ label: string; met: boolean }> = ({ label, met }) => (
    <ListItem sx={{ py: 0.25, px: 0 }}>
        <ListItemIcon sx={{ minWidth: 30 }}>
            {met ? (
                <CheckCircleOutlineIcon color="success" sx={{ fontSize: 18 }} />
            ) : (
                <ErrorOutlineIcon color="disabled" sx={{ fontSize: 18 }} />
            )}
        </ListItemIcon>
        <ListItemText
            primary={label}
            primaryTypographyProps={{
                variant: 'caption',
                color: met ? 'success.main' : 'text.secondary',
                fontWeight: met ? 600 : 400
            }}
        />
    </ListItem>
);
