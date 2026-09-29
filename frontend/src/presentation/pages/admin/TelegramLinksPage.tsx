import React, { useEffect, useMemo, useState } from 'react';
import {
    Box,
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Chip,
    Typography,
    alpha,
    Button,
    CircularProgress,
} from '@mui/material';
import TelegramIcon from '@mui/icons-material/Telegram';
import { useUserStore } from '../../../state/userStore';
import { TelegramManageDialog } from '../users/TelegramManageDialog';
import type { User } from '../../../domain/entities/User';

export const TelegramLinksPage: React.FC = () => {
    const {
        users, fetchUsers,
        telegramLinks, isLoading, fetchTelegramLinks,
        telegramStatus, fetchTelegramStatus, generateTelegramLinkCode, unlinkTelegram,
    } = useUserStore();
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [openManage, setOpenManage] = useState(false);

    useEffect(() => {
        fetchUsers();
        fetchTelegramLinks();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    const rows = useMemo(() => {
        return telegramLinks
            .filter((link) => link.is_active)
            .map((link) => ({
                link,
                user: users.find((u) => u.id === link.user_id) || null,
            }));
    }, [telegramLinks, users]);

    const handleManage = (user: User | null) => {
        if (!user) return;
        setSelectedUser(user);
        setOpenManage(true);
    };

    const handleCloseManage = () => {
        setOpenManage(false);
        fetchTelegramLinks();
    };

    return (
        <Box sx={{ width: '100%', p: { xs: 2, md: 4 } }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 3 }}>
                <TelegramIcon color="primary" sx={{ fontSize: 32 }} />
                <Box>
                    <Typography variant="h5" sx={{ fontWeight: 700 }}>Telegram Links</Typography>
                    <Typography variant="body2" color="text.secondary">
                        Every user currently connected to the CashBook Telegram bot.
                    </Typography>
                </Box>
            </Box>

            <TableContainer
                component={Paper}
                sx={{
                    width: '100%',
                    borderRadius: 3,
                    overflow: 'hidden',
                    overflowX: 'auto',
                    boxShadow: '0 4px 20px rgba(0,0,0,0.05)',
                    border: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`
                }}
            >
                <Table sx={{ minWidth: 650 }}>
                    <TableHead sx={{ bgcolor: (theme) => alpha(theme.palette.action.hover, 0.5) }}>
                        <TableRow>
                            <TableCell sx={{ fontWeight: 600 }}>User</TableCell>
                            <TableCell sx={{ fontWeight: 600 }}>Telegram Username</TableCell>
                            <TableCell sx={{ fontWeight: 600 }}>Linked At</TableCell>
                            <TableCell sx={{ fontWeight: 600 }}>Status</TableCell>
                            <TableCell align="right" sx={{ fontWeight: 600 }}>Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {isLoading && rows.length === 0 && (
                            <TableRow>
                                <TableCell colSpan={5} align="center" sx={{ py: 6 }}>
                                    <CircularProgress size={24} />
                                </TableCell>
                            </TableRow>
                        )}
                        {!isLoading && rows.map(({ link, user }) => (
                            <TableRow
                                key={link.id}
                                sx={{
                                    '&:hover': { bgcolor: (theme) => alpha(theme.palette.primary.main, 0.02) },
                                    transition: 'background-color 0.2s'
                                }}
                            >
                                <TableCell>
                                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                        {user ? user.name : `User #${link.user_id}`}
                                    </Typography>
                                    {user && (
                                        <Typography variant="caption" color="text.secondary">
                                            {user.email}
                                        </Typography>
                                    )}
                                </TableCell>
                                <TableCell>
                                    {link.telegram_username ? `@${link.telegram_username}` : '—'}
                                </TableCell>
                                <TableCell>{new Date(link.linked_at).toLocaleString()}</TableCell>
                                <TableCell>
                                    <Chip
                                        label="Connected"
                                        color="success"
                                        size="small"
                                        sx={{
                                            borderRadius: 1.5,
                                            fontSize: '0.75rem',
                                            fontWeight: 600,
                                            bgcolor: (theme) => alpha(theme.palette.success.main, 0.1),
                                            color: (theme) => theme.palette.success.main,
                                            border: 'none'
                                        }}
                                    />
                                </TableCell>
                                <TableCell align="right">
                                    <Button
                                        size="small"
                                        variant="outlined"
                                        disabled={!user}
                                        onClick={() => handleManage(user)}
                                        sx={{ textTransform: 'none', borderRadius: '10px' }}
                                    >
                                        Manage
                                    </Button>
                                </TableCell>
                            </TableRow>
                        ))}
                        {!isLoading && rows.length === 0 && (
                            <TableRow>
                                <TableCell colSpan={5} align="center" sx={{ py: 8 }}>
                                    <Typography variant="body1" color="text.secondary">
                                        No one has connected Telegram yet.
                                    </Typography>
                                </TableCell>
                            </TableRow>
                        )}
                    </TableBody>
                </Table>
            </TableContainer>

            <TelegramManageDialog
                open={openManage}
                user={selectedUser}
                status={telegramStatus}
                isLoading={isLoading}
                onClose={handleCloseManage}
                onLoadStatus={fetchTelegramStatus}
                onGenerateCode={generateTelegramLinkCode}
                onUnlink={unlinkTelegram}
            />
        </Box>
    );
};
