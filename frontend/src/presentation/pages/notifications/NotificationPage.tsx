import React, { useState } from 'react';
import {
    Box,
    Typography,
    Container,
    Paper,
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    Avatar,
    IconButton,
    Button,
    Chip,
    Tabs,
    Tab,
    Divider,
    alpha,
    Tooltip,
    Stack,
    Fade
} from '@mui/material';
import NotificationsIcon from '@mui/icons-material/Notifications';
import InfoIcon from '@mui/icons-material/Info';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import WarningIcon from '@mui/icons-material/Warning';
import ErrorIcon from '@mui/icons-material/Error';
import DeleteIcon from '@mui/icons-material/Delete';
import DoneAllIcon from '@mui/icons-material/DoneAll';
import MoreVertIcon from '@mui/icons-material/MoreVert';

type NotificationType = 'INFO' | 'SUCCESS' | 'WARNING' | 'ERROR';

interface Notification {
    id: string;
    type: NotificationType;
    title: string;
    message: string;
    time: string;
    read: boolean;
}

const MOCK_NOTIFICATIONS: Notification[] = [
    {
        id: '1',
        type: 'SUCCESS',
        title: 'System Update Successful',
        message: 'The latest system update has been successfully installed. Check the release notes for new features.',
        time: '2 hours ago',
        read: false,
    },
    {
        id: '2',
        type: 'INFO',
        title: 'New User Registered',
        message: 'A new user has just signed up for the application. Refresh your user list to see the update.',
        time: '5 hours ago',
        read: false,
    },
    {
        id: '3',
        type: 'WARNING',
        title: 'Storage Capacity Alert',
        message: 'Your current storage is at 85% capacity. Consider upgrading your plan or cleaning up old files.',
        time: 'Yesterday',
        read: true,
    },
    {
        id: '4',
        type: 'ERROR',
        title: 'Connection Lost',
        message: 'The connection to the primary database was lost briefly at 3:00 AM. System recovered automatically.',
        time: '2 days ago',
        read: true,
    },
];

const getTypeColor = (type: NotificationType) => {
    switch (type) {
        case 'SUCCESS': return 'success';
        case 'WARNING': return 'warning';
        case 'ERROR': return 'error';
        default: return 'primary';
    }
};

const getTypeIcon = (type: NotificationType) => {
    switch (type) {
        case 'SUCCESS': return <CheckCircleIcon />;
        case 'WARNING': return <WarningIcon />;
        case 'ERROR': return <ErrorIcon />;
        default: return <InfoIcon />;
    }
};

export const NotificationPage: React.FC = () => {
    const [notifications, setNotifications] = useState<Notification[]>(MOCK_NOTIFICATIONS);
    const [tabValue, setTabValue] = useState(0);

    const handleTabChange = (_: React.SyntheticEvent, newValue: number) => {
        setTabValue(newValue);
    };

    const filteredNotifications = notifications.filter(n => {
        if (tabValue === 1) return !n.read;
        return true;
    });

    const markAllAsRead = () => {
        setNotifications(notifications.map(n => ({ ...n, read: true })));
    };

    const clearAll = () => {
        setNotifications([]);
    };

    const deleteNotification = (id: string) => {
        setNotifications(notifications.filter(n => n.id !== id));
    };

    const toggleRead = (id: string) => {
        setNotifications(notifications.map(n =>
            n.id === id ? { ...n, read: !n.read } : n
        ));
    };

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Box sx={{ mb: 4, display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', flexWrap: 'wrap', gap: 2 }}>
                <Box>
                    <Typography variant="h4" sx={{ fontWeight: 800, mb: 1, color: 'text.primary' }}>
                        Notifications
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Stay updated with the latest activity and system alerts.
                    </Typography>
                </Box>
                <Stack direction="row" spacing={1}>
                    <Button
                        startIcon={<DoneAllIcon />}
                        onClick={markAllAsRead}
                        disabled={!notifications.some(n => !n.read)}
                        variant="outlined"
                        size="small"
                        sx={{ borderRadius: 2 }}
                    >
                        Mark all as read
                    </Button>
                    <Button
                        startIcon={<DeleteIcon />}
                        onClick={clearAll}
                        disabled={notifications.length === 0}
                        color="error"
                        variant="outlined"
                        size="small"
                        sx={{ borderRadius: 2 }}
                    >
                        Clear all
                    </Button>
                </Stack>
            </Box>

            <Paper
                elevation={0}
                sx={{
                    borderRadius: 4,
                    overflow: 'hidden',
                    border: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                    bgcolor: 'background.paper'
                }}
            >
                <Tabs
                    value={tabValue}
                    onChange={handleTabChange}
                    sx={{
                        px: 2,
                        pt: 1,
                        borderBottom: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                        '& .MuiTabs-indicator': { height: 3, borderRadius: '3px 3px 0 0' }
                    }}
                >
                    <Tab
                        label={
                            <Stack direction="row" spacing={1} alignItems="center">
                                <span>All</span>
                                <Chip label={notifications.length} size="small" sx={{ height: 20, fontSize: '0.75rem' }} />
                            </Stack>
                        }
                    />
                    <Tab
                        label={
                            <Stack direction="row" spacing={1} alignItems="center">
                                <span>Unread</span>
                                {notifications.filter(n => !n.read).length > 0 && (
                                    <Chip
                                        label={notifications.filter(n => !n.read).length}
                                        size="small"
                                        color="primary"
                                        sx={{ height: 20, fontSize: '0.75rem' }}
                                    />
                                )}
                            </Stack>
                        }
                    />
                </Tabs>

                <List disablePadding>
                    {filteredNotifications.length > 0 ? (
                        filteredNotifications.map((notif, index) => (
                            <Fade in={true} key={notif.id} timeout={300 + index * 100}>
                                <Box>
                                    <ListItem
                                        sx={{
                                            py: 3,
                                            px: { xs: 2, md: 4 },
                                            bgcolor: notif.read ? 'transparent' : (theme) => alpha(theme.palette.primary.main, 0.02),
                                            transition: 'all 0.2s',
                                            '&:hover': {
                                                bgcolor: (theme) => alpha(theme.palette.action.hover, 0.04),
                                            },
                                            position: 'relative',
                                        }}
                                        secondaryAction={
                                            <Stack direction="row" spacing={1}>
                                                <Tooltip title="Delete">
                                                    <IconButton edge="end" onClick={() => deleteNotification(notif.id)}>
                                                        <DeleteIcon fontSize="small" />
                                                    </IconButton>
                                                </Tooltip>
                                                <IconButton edge="end">
                                                    <MoreVertIcon fontSize="small" />
                                                </IconButton>
                                            </Stack>
                                        }
                                    >
                                        <ListItemAvatar sx={{ minWidth: 64 }}>
                                            <Avatar
                                                sx={{
                                                    bgcolor: (theme) => alpha(theme.palette[getTypeColor(notif.type)].main, 0.1),
                                                    color: (theme) => theme.palette[getTypeColor(notif.type)].main,
                                                    width: 48,
                                                    height: 48
                                                }}
                                            >
                                                {getTypeIcon(notif.type)}
                                            </Avatar>
                                        </ListItemAvatar>
                                        <ListItemText
                                            primary={
                                                <Stack direction="row" spacing={1} alignItems="center" sx={{ mb: 0.5 }}>
                                                    <Typography
                                                        variant="subtitle1"
                                                        sx={{
                                                            fontWeight: notif.read ? 600 : 800,
                                                            color: notif.read ? 'text.primary' : 'primary.main'
                                                        }}
                                                    >
                                                        {notif.title}
                                                    </Typography>
                                                    {!notif.read && (
                                                        <Box sx={{ width: 8, height: 8, borderRadius: '50%', bgcolor: 'primary.main' }} />
                                                    )}
                                                </Stack>
                                            }
                                            secondary={
                                                <Box>
                                                    <Typography variant="body2" color="text.secondary" sx={{ mb: 1, maxWidth: '90%' }}>
                                                        {notif.message}
                                                    </Typography>
                                                    <Typography variant="caption" color="text.disabled" sx={{ fontWeight: 500 }}>
                                                        {notif.time}
                                                    </Typography>
                                                </Box>
                                            }
                                            onClick={() => toggleRead(notif.id)}
                                            style={{ cursor: 'pointer' }}
                                        />
                                    </ListItem>
                                    {index < filteredNotifications.length - 1 && <Divider component="li" sx={{ opacity: 0.5 }} />}
                                </Box>
                            </Fade>
                        ))
                    ) : (
                        <Box sx={{ py: 10, textAlign: 'center' }}>
                            <Avatar sx={{ width: 80, height: 80, bgcolor: 'action.hover', mx: 'auto', mb: 2 }}>
                                <NotificationsIcon sx={{ fontSize: 40, color: 'text.disabled' }} />
                            </Avatar>
                            <Typography variant="h6" color="text.secondary">
                                No notifications found
                            </Typography>
                            <Typography variant="body2" color="text.disabled">
                                You're all caught up! Check back later for new updates.
                            </Typography>
                        </Box>
                    )}
                </List>
            </Paper>
        </Container>
    );
};
