import React from 'react';
import {
    Grid,
    Card,
    Typography,
    Pagination,
    Stepper,
    Step,
    StepLabel,
    Breadcrumbs,
    Link,
    BottomNavigation,
    BottomNavigationAction,
    Box,
    Stack,
    Skeleton,
    Paper
} from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import RestoreIcon from '@mui/icons-material/Restore';
import FavoriteIcon from '@mui/icons-material/Favorite';
import LocationOnIcon from '@mui/icons-material/LocationOn';

export const NavigationShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    const steps = ['Step One', 'Step Two', 'Step Three'];

    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(3)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="30%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={100} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Breadcrumbs & Pagination */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Breadcrumbs & Pagination</Typography>
                    <Box sx={{ mb: 4 }}>
                        <Breadcrumbs aria-label="breadcrumb">
                            <Link underline="hover" color="inherit" href="/" sx={{ display: 'flex', alignItems: 'center' }}>
                                <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
                                MUI
                            </Link>
                            <Link underline="hover" color="inherit" href="/dashboard">Dashboard</Link>
                            <Typography color="text.primary">Navigation</Typography>
                        </Breadcrumbs>
                    </Box>
                    <Box>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Pagination Components</Typography>
                        <Stack spacing={2}>
                            <Pagination count={10} color="primary" />
                            <Pagination count={10} variant="outlined" shape="rounded" />
                        </Stack>
                    </Box>
                </Card>
            </Grid>

            {/* Stepper */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Stepper Flow</Typography>
                    <Stepper activeStep={1} orientation="vertical">
                        {steps.map((label) => (
                            <Step key={label}>
                                <StepLabel>{label}</StepLabel>
                            </Step>
                        ))}
                    </Stepper>
                </Card>
            </Grid>

            {/* Bottom Nav */}
            <Grid size={{ xs: 12 }}>
                <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Bottom Navigation</Typography>
                    <Box sx={{ width: '100%', display: 'flex', justifyContent: 'center', p: 4, bgcolor: 'grey.50', borderRadius: 3 }}>
                        <Paper elevation={3} sx={{ width: 400 }}>
                            <BottomNavigation showLabels>
                                <BottomNavigationAction label="Recents" icon={<RestoreIcon />} />
                                <BottomNavigationAction label="Favorites" icon={<FavoriteIcon />} />
                                <BottomNavigationAction label="Nearby" icon={<LocationOnIcon />} />
                            </BottomNavigation>
                        </Paper>
                    </Box>
                </Card>
            </Grid>
        </Grid>
    );
};
