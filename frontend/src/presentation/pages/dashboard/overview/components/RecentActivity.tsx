import React from 'react';
import {
    Card,
    Typography,
    List,
    ListItem,
    ListItemAvatar,
    ListItemText,
    Avatar,
    alpha,
    Skeleton
} from '@mui/material';
import ShoppingCartIcon from '@mui/icons-material/ShoppingCart';
import PersonAddIcon from '@mui/icons-material/PersonAdd';
import SettingsIcon from '@mui/icons-material/Settings';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

interface ActivityItem {
    id: number;
    title: string;
    time: string;
    icon: React.ReactNode;
    color: string;
}

const activities: ActivityItem[] = [
    { id: 1, title: 'New order received from John Doe', time: '5 minutes ago', icon: <ShoppingCartIcon fontSize="small" />, color: '#1976d2' },
    { id: 2, title: 'User "Alice" registered', time: '1 hour ago', icon: <PersonAddIcon fontSize="small" />, color: '#2e7d32' },
    { id: 3, title: 'System settings updated by admin', time: '3 hours ago', icon: <SettingsIcon fontSize="small" />, color: '#757575' },
    { id: 4, title: 'Server storage limit warning', time: '5 hours ago', icon: <ErrorOutlineIcon fontSize="small" />, color: '#d32f2f' },
];

export const RecentActivity: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    return (
        <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)', height: '100%' }}>
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>
                Recent Activities
            </Typography>

            <List disablePadding>
                {(isLoading ? [...Array(4)] : activities).map((activity, index) => (
                    <ListItem
                        key={activity?.id || index}
                        disableGutters
                        sx={{ mb: 2, '&:last-child': { mb: 0 } }}
                    >
                        <ListItemAvatar>
                            {isLoading ? (
                                <Skeleton variant="circular" width={40} height={40} />
                            ) : (
                                <Avatar sx={{ bgcolor: alpha(activity.color, 0.1), color: activity.color }}>
                                    {activity.icon}
                                </Avatar>
                            )}
                        </ListItemAvatar>
                        <ListItemText
                            primary={isLoading ? <Skeleton variant="text" width="80%" /> : activity.title}
                            secondary={isLoading ? <Skeleton variant="text" width="40%" /> : activity.time}
                            primaryTypographyProps={{ variant: 'body2', fontWeight: 600 }}
                            secondaryTypographyProps={{ variant: 'caption' }}
                        />
                    </ListItem>
                ))}
            </List>
        </Card>
    );
};
