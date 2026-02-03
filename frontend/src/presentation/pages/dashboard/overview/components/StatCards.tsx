import React from 'react';
import { Grid, Card, Box, Typography, alpha, Skeleton, useTheme } from '@mui/material';
import { formatIDR } from '../../../../utils/formatCurrency';
import AccountBalanceIcon from '@mui/icons-material/AccountBalance';
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward';
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import AccountBalanceWalletIcon from '@mui/icons-material/AccountBalanceWallet';

interface StatCardProps {
    title: string;
    value: string;
    trend?: string;
    icon: React.ReactNode;
    color: string;
    isLoading?: boolean;
}

const StatCard: React.FC<StatCardProps> = ({ title, value, trend, icon, color, isLoading }) => {
    const theme = useTheme();

    if (isLoading) {
        return (
            <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
                    <Skeleton variant="circular" width={48} height={48} />
                </Box>
                <Skeleton variant="text" width="60%" height={40} />
                <Skeleton variant="text" width="40%" />
            </Card>
        );
    }

    return (
        <Card sx={{
            p: { xs: 3, md: 5 },
            borderRadius: 4,
            height: '100%',
            boxShadow: '0 4px 20px rgba(0,0,0,0.05)',
            position: 'relative',
            overflow: 'hidden',
        }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
                <Box sx={{
                    p: 1.5,
                    borderRadius: 3,
                    bgcolor: alpha(color, 0.1),
                    color: color,
                    display: 'flex'
                }}>
                    {icon}
                </Box>
                {trend && (
                    <Typography variant="caption" sx={{
                        color: theme.palette.success.main,
                        fontWeight: 700,
                        bgcolor: alpha(theme.palette.success.main, 0.1),
                        px: 1,
                        py: 0.5,
                        borderRadius: 1.5
                    }}>
                        {trend}
                    </Typography>
                )}
            </Box>
            <Typography variant="h5" sx={{ fontWeight: 700, mb: 0.5 }}>
                {value}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 500 }}>
                {title}
            </Typography>
        </Card>
    );
};

export const StatCards: React.FC<{ summary: any, isLoading: boolean }> = ({ summary, isLoading }) => {
    const stats = [
        {
            title: 'Current Balance',
            value: formatIDR(summary?.balance || 0),
            icon: <AccountBalanceWalletIcon />,
            color: '#1976d2'
        },
        {
            title: 'Total Income',
            value: formatIDR(summary?.total_income || 0),
            icon: <ArrowUpwardIcon />,
            color: '#2e7d32'
        },
        {
            title: 'Total Expense',
            value: formatIDR(summary?.total_expense || 0),
            icon: <ArrowDownwardIcon />,
            color: '#d32f2f'
        },
        {
            title: 'Accounts',
            value: '1',
            icon: <AccountBalanceIcon />,
            color: '#9c27b0'
        },
    ];

    return (
        <Grid
            container
            spacing={3}
            sx={{
                width: '100%',
                margin: 0,
            }}
        >
            {stats.map((stat, index) => (
                <Grid key={index} size={{ xs: 12, sm: 6, lg: 3 }}>
                    <StatCard
                        {...stat}
                        isLoading={isLoading}
                    />
                </Grid>
            ))}
        </Grid>
    );
};
