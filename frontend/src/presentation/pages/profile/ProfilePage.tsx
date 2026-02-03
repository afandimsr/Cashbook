import React, { useState, useEffect } from 'react';
import {
    Box,
    Container,
    Paper,
    Typography,
    Avatar,
    Grid,
    Divider,
    Button,
    alpha,
    useTheme,
    Card,
    CardContent,
    Stack,
    IconButton,
    Skeleton,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import EmailIcon from '@mui/icons-material/Email';
import BadgeIcon from '@mui/icons-material/Badge';
import SecurityIcon from '@mui/icons-material/Security';
import CalendarMonthIcon from '@mui/icons-material/CalendarMonth';
import CameraAltIcon from '@mui/icons-material/CameraAlt';
import { useAuthStore } from '../../../state/authStore';

const ProfileSkeleton = () => {
    const theme = useTheme();
    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            <Skeleton variant="rectangular" height={240} sx={{ borderRadius: 4, mb: 4 }} />
            <Box sx={{
                display: 'flex',
                flexDirection: { xs: 'column', sm: 'row' },
                alignItems: { xs: 'center', sm: 'flex-end' },
                gap: 3,
                px: { xs: 2, sm: 4 },
                mt: { xs: -10, sm: -5 },
                mb: 4
            }}>
                <Skeleton variant="circular" width={140} height={140} sx={{ border: `6px solid ${theme.palette.background.paper}`, flexShrink: 0 }} />
                <Box sx={{ flexGrow: 1, pb: 1, textAlign: { xs: 'center', sm: 'left' } }}>
                    <Skeleton variant="text" width="60%" height={60} sx={{ mx: { xs: 'auto', sm: 0 } }} />
                    <Skeleton variant="text" width="40%" sx={{ mx: { xs: 'auto', sm: 0 } }} />
                </Box>
            </Box>
            <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 8 }}>
                    <Stack spacing={3}>
                        <Skeleton variant="rectangular" height={160} sx={{ borderRadius: 4 }} />
                        <Skeleton variant="rectangular" height={240} sx={{ borderRadius: 4 }} />
                    </Stack>
                </Grid>
                <Grid size={{ xs: 12, md: 4 }}>
                    <Stack spacing={3}>
                        <Skeleton variant="rectangular" height={200} sx={{ borderRadius: 4 }} />
                        <Skeleton variant="rectangular" height={150} sx={{ borderRadius: 4 }} />
                    </Stack>
                </Grid>
            </Grid>
        </Container>
    );
};

export const ProfilePage: React.FC = () => {
    const theme = useTheme();
    const { user } = useAuthStore();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const timer = setTimeout(() => {
            setLoading(false);
        }, 1500);
        return () => clearTimeout(timer);
    }, []);

    if (loading) return <ProfileSkeleton />;

    return (
        <Container maxWidth="lg" sx={{ py: 4 }}>
            {/* Header / Cover Section */}
            <Paper
                elevation={0}
                sx={{
                    position: 'relative',
                    height: { xs: 200, sm: 240 },
                    borderRadius: 4,
                    mb: { xs: 10, sm: 4 },
                    background: theme.palette.mode === 'light'
                        ? 'linear-gradient(135deg, #1e293b 0%, #334155 100%)'
                        : 'linear-gradient(135deg, #0f172a 0%, #1e293b 100%)',
                    overflow: 'visible',
                    display: 'flex',
                    alignItems: 'flex-end',
                    px: { xs: 2, sm: 4 },
                    pb: 0,
                    boxShadow: '0 10px 30px -10px rgba(0,0,0,0.3)',
                }}
            >
                <IconButton
                    sx={{
                        position: 'absolute',
                        top: 16,
                        right: 16,
                        bgcolor: 'rgba(255,255,255,0.1)',
                        '&:hover': { bgcolor: 'rgba(255,255,255,0.2)' },
                        color: 'white'
                    }}
                >
                    <CameraAltIcon />
                </IconButton>

                <Box
                    sx={{
                        display: 'flex',
                        flexDirection: { xs: 'column', sm: 'row' },
                        alignItems: { xs: 'center', sm: 'flex-end' },
                        textAlign: { xs: 'center', sm: 'left' },
                        gap: { xs: 1, sm: 3 },
                        transform: { xs: 'translateY(80px)', sm: 'translateY(40px)' },
                        pb: 2,
                        width: '100%',
                    }}
                >
                    <Avatar
                        sx={{
                            width: { xs: 110, sm: 140 },
                            height: { xs: 110, sm: 140 },
                            border: `6px solid ${theme.palette.background.paper}`,
                            bgcolor: 'primary.main',
                            fontSize: { xs: '2.5rem', sm: '3rem' },
                            fontWeight: 700,
                            boxShadow: '0 8px 24px rgba(0,0,0,0.12)',
                        }}
                    >
                        {user?.username?.charAt(0).toUpperCase()}
                    </Avatar>
                    <Box sx={{ mb: { xs: 0, sm: 6 }, flexGrow: 1 }}>
                        <Typography variant="h4" sx={{
                            fontWeight: 800,
                            color: 'white',
                            textShadow: '0 2px 4px rgba(0,0,0,0.2)',
                            letterSpacing: '-0.02em',
                            fontSize: { xs: '1.75rem', sm: '2.125rem' }
                        }}>
                            {user?.username}
                        </Typography>
                        <Stack
                            direction={{ xs: 'column', sm: 'row' }}
                            spacing={1}
                            alignItems="center"
                            sx={{ mt: 0.5 }}
                        >
                            <Typography variant="subtitle1" sx={{ color: 'rgba(255,255,255,0.9)', fontWeight: 500 }}>
                                {user?.role === 'ADMIN' ? 'Administrator' : 'Team Member'}
                            </Typography>
                            <Box sx={{
                                width: 4,
                                height: 4,
                                borderRadius: '50%',
                                bgcolor: 'rgba(255,255,255,0.4)',
                                alignSelf: 'center',
                                display: { xs: 'none', sm: 'block' }
                            }} />
                            <Typography variant="subtitle1" sx={{ color: 'rgba(255,255,255,0.9)', fontWeight: 500 }}>
                                {user?.email}
                            </Typography>
                        </Stack>
                    </Box>
                    <Button
                        variant="contained"
                        startIcon={<EditIcon />}
                        sx={{
                            mb: { xs: 0, sm: 4 },
                            borderRadius: 2,
                            textTransform: 'none',
                            fontWeight: 600,
                            px: 3,
                            py: 1,
                            bgcolor: { xs: 'primary.main', sm: 'white' },
                            color: { xs: 'white', sm: '#1e293b' },
                            '&:hover': { bgcolor: { xs: 'primary.dark', sm: 'rgba(255,255,255,0.9)' } },
                            mt: { xs: 1, sm: 0 },
                            width: { xs: '100%', sm: 'auto' }
                        }}
                    >
                        Edit Profile
                    </Button>
                </Box>
            </Paper>

            <Grid container spacing={4} sx={{ mt: { xs: 12, sm: 4 } }}>
                {/* Left Column - Information */}
                <Grid size={{ xs: 12, md: 8 }}>
                    <Stack spacing={3}>
                        <Card sx={{ borderRadius: 4, border: `1px solid ${alpha(theme.palette.divider, 0.1)}` }}>
                            <CardContent sx={{ p: { xs: 3, sm: 4 } }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>About Me</Typography>
                                <Typography variant="body1" color="text.secondary" sx={{ lineHeight: 1.7 }}>
                                    Welcome to your profile. Here you can manage your personal information, security settings, and notification preferences. As an {user?.role?.toLowerCase()}, you have access to various system features tailored to your role.
                                </Typography>
                            </CardContent>
                        </Card>

                        <Card sx={{ borderRadius: 4, border: `1px solid ${alpha(theme.palette.divider, 0.1)}` }}>
                            <CardContent sx={{ p: { xs: 3, sm: 4 } }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>General Information</Typography>
                                <Grid container spacing={3}>
                                    {[
                                        { label: 'Full Name', value: user?.name, icon: <BadgeIcon color="primary" /> },
                                        { label: 'Email Address', value: user?.email, icon: <EmailIcon color="primary" /> },
                                        { label: 'Role', value: user?.role, icon: <SecurityIcon color="primary" /> },
                                        { label: 'Member Since', value: 'January 2026', icon: <CalendarMonthIcon color="primary" /> },
                                    ].map((info, idx) => (
                                        <Grid size={{ xs: 12, sm: 6 }} key={idx}>
                                            <Box sx={{
                                                display: 'flex',
                                                gap: 2,
                                                p: 2,
                                                borderRadius: 3,
                                                bgcolor: alpha(theme.palette.action.hover, 0.3),
                                                border: `1px solid ${alpha(theme.palette.divider, 0.05)}`,
                                                transition: 'all 0.2s',
                                                '&:hover': {
                                                    bgcolor: alpha(theme.palette.action.hover, 0.5),
                                                    transform: 'translateY(-2px)'
                                                }
                                            }}>
                                                <Box sx={{
                                                    p: 1.5,
                                                    borderRadius: 2,
                                                    bgcolor: alpha(theme.palette.primary.main, 0.1),
                                                    display: 'flex',
                                                    alignItems: 'center',
                                                    justifyContent: 'center',
                                                    height: 'fit-content'
                                                }}>
                                                    {info.icon}
                                                </Box>
                                                <Box sx={{ overflow: 'hidden' }}>
                                                    <Typography variant="caption" color="text.secondary" sx={{ fontWeight: 600, display: 'block' }}>
                                                        {info.label}
                                                    </Typography>
                                                    <Typography variant="body1" sx={{ fontWeight: 600, wordBreak: 'break-all' }}>
                                                        {info.value}
                                                    </Typography>
                                                </Box>
                                            </Box>
                                        </Grid>
                                    ))}
                                </Grid>
                            </CardContent>
                        </Card>
                    </Stack>
                </Grid>

                {/* Right Column - Secondary Actions / Stats */}
                <Grid size={{ xs: 12, md: 4 }}>
                    <Stack spacing={3}>
                        <Card sx={{
                            borderRadius: 4,
                            border: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                            bgcolor: theme.palette.mode === 'light' ? alpha(theme.palette.primary.main, 0.03) : alpha(theme.palette.primary.main, 0.05)
                        }}>
                            <CardContent sx={{ p: { xs: 3, sm: 4 } }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 2 }}>Account Status</Typography>
                                <Stack spacing={2}>
                                    <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                        <Typography variant="body2" color="text.secondary">Verification</Typography>
                                        <Box sx={{ px: 1.5, py: 0.5, borderRadius: 10, bgcolor: 'success.main', color: 'white', fontSize: '0.75rem', fontWeight: 700 }}>
                                            Verified
                                        </Box>
                                    </Box>
                                    <Divider />
                                </Stack>
                            </CardContent>
                        </Card>

                        <Card sx={{ borderRadius: 4, border: `1px solid ${alpha(theme.palette.divider, 0.1)}` }}>
                            <CardContent sx={{ p: { xs: 3, sm: 4 } }}>
                                <Typography variant="h6" sx={{ fontWeight: 700, mb: 2 }}>Security</Typography>
                                <Stack spacing={2}>
                                    <Button
                                        variant="text"
                                        fullWidth
                                        sx={{
                                            justifyContent: 'flex-start',
                                            textTransform: 'none',
                                            fontWeight: 600,
                                            color: 'text.primary',
                                            '&:hover': { bgcolor: alpha(theme.palette.action.hover, 0.05) }
                                        }}
                                    >
                                        Change Password
                                    </Button>
                                    <Button
                                        variant="text"
                                        fullWidth
                                        sx={{
                                            justifyContent: 'flex-start',
                                            textTransform: 'none',
                                            fontWeight: 600,
                                            color: 'text.primary',
                                            '&:hover': { bgcolor: alpha(theme.palette.action.hover, 0.05) }
                                        }}
                                    >
                                        Two-Factor Authentication
                                    </Button>
                                </Stack>
                            </CardContent>
                        </Card>
                    </Stack>
                </Grid>
            </Grid>
        </Container>
    );
};
