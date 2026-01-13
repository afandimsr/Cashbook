import React, { useState } from 'react';
import {
    Grid,
    Card,
    Typography,
    Alert,
    AlertTitle,
    Button,
    Snackbar,
    CircularProgress,
    LinearProgress,
    Box,
    Stack,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Skeleton
} from '@mui/material';

export const FeedbackShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    const [snackOpen, setSnackOpen] = useState(false);
    const [dialogOpen, setDialogOpen] = useState(false);

    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(3)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="30%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={80} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Alerts */}
            <Grid size={{ xs: 12 }}>
                <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Alerts</Typography>
                    <Stack spacing={2}>
                        <Alert severity="error">
                            <AlertTitle>Error</AlertTitle>
                            This is an error alert — <strong>check it out!</strong>
                        </Alert>
                        <Alert severity="warning">This is a warning alert.</Alert>
                        <Alert severity="info">This is an info alert.</Alert>
                        <Alert severity="success">This is a success alert.</Alert>
                    </Stack>
                </Card>
            </Grid>

            {/* Progress & Interactivity */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Progress Indicators</Typography>
                    <Box sx={{ mb: 4 }}>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Circular</Typography>
                        <Stack direction="row" spacing={3} alignItems="center">
                            <CircularProgress size={24} />
                            <CircularProgress color="secondary" />
                            <CircularProgress variant="determinate" value={75} />
                        </Stack>
                    </Box>
                    <Box>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Linear</Typography>
                        <LinearProgress sx={{ mb: 2, height: 8, borderRadius: 4 }} />
                        <LinearProgress color="secondary" variant="buffer" value={50} valueBuffer={70} sx={{ height: 8, borderRadius: 4 }} />
                    </Box>
                </Card>
            </Grid>

            {/* Triggers */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Interactive Triggers</Typography>
                    <Stack spacing={3}>
                        <Box>
                            <Typography variant="subtitle2" sx={{ mb: 2 }}>Snackbar</Typography>
                            <Button variant="outlined" onClick={() => setSnackOpen(true)}>Open Snackbar</Button>
                        </Box>
                        <Box>
                            <Typography variant="subtitle2" sx={{ mb: 2 }}>Dialog</Typography>
                            <Button variant="contained" onClick={() => setDialogOpen(true)}>Open Dialog</Button>
                        </Box>
                    </Stack>
                </Card>
            </Grid>

            <Snackbar
                open={snackOpen}
                autoHideDuration={3000}
                onClose={() => setSnackOpen(false)}
                message="Note archived successfully!"
            />

            <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)}>
                <DialogTitle>{"Use Google's location service?"}</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Let Google help apps determine location. This means sending anonymous
                        location data to Google, even when no apps are running.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setDialogOpen(false)}>Disagree</Button>
                    <Button onClick={() => setDialogOpen(false)} autoFocus>Agree</Button>
                </DialogActions>
            </Dialog>
        </Grid>
    );
};
