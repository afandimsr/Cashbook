import React, { useState } from 'react';
import {
    Grid,
    Card,
    Typography,
    Button,
    Menu,
    MenuItem,
    Popover,
    Modal,
    Box,
    Stack,
    Fade,
    Grow,
    CircularProgress,
    Skeleton,
    alpha
} from '@mui/material';

export const UtilsShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [popoverAnchor, setPopoverAnchor] = useState<null | HTMLElement>(null);
    const [modalOpen, setModalOpen] = useState(false);
    const [showGrow, setShowGrow] = useState(true);

    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(4)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12, md: 6 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="40%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={100} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Popups & Menus */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Popups & Menus</Typography>
                    <Stack spacing={2} direction="row">
                        <Box>
                            <Button variant="outlined" onClick={(e) => setAnchorEl(e.currentTarget)}>Standard Menu</Button>
                            <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={() => setAnchorEl(null)}>
                                <MenuItem onClick={() => setAnchorEl(null)}>Profile</MenuItem>
                                <MenuItem onClick={() => setAnchorEl(null)}>Account</MenuItem>
                                <MenuItem onClick={() => setAnchorEl(null)}>Logout</MenuItem>
                            </Menu>
                        </Box>
                        <Box>
                            <Button variant="outlined" color="secondary" onClick={(e) => setPopoverAnchor(e.currentTarget)}>Popover</Button>
                            <Popover
                                open={Boolean(popoverAnchor)}
                                anchorEl={popoverAnchor}
                                onClose={() => setPopoverAnchor(null)}
                                anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
                            >
                                <Typography sx={{ p: 2 }}>This is a Popover content.</Typography>
                            </Popover>
                        </Box>
                    </Stack>
                </Card>
            </Grid>

            {/* Modal */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Modals</Typography>
                    <Button variant="contained" onClick={() => setModalOpen(true)}>Open Basic Modal</Button>
                    <Modal open={modalOpen} onClose={() => setModalOpen(false)}>
                        <Box sx={{
                            position: 'absolute',
                            top: '50%',
                            left: '50%',
                            transform: 'translate(-50%, -50%)',
                            width: 400,
                            bgcolor: 'background.paper',
                            borderRadius: 3,
                            boxShadow: 24,
                            p: 4
                        }}>
                            <Typography variant="h6" component="h2">Modal Title</Typography>
                            <Typography sx={{ mt: 2 }}>Duis mollis, est non commodo luctus, nisi erat porttitor ligula.</Typography>
                            <Button sx={{ mt: 3 }} fullWidth variant="outlined" onClick={() => setModalOpen(false)}>Close</Button>
                        </Box>
                    </Modal>
                </Card>
            </Grid>

            {/* Transitions */}
            <Grid size={{ xs: 12 }}>
                <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Transition Effects</Typography>
                    <Button onClick={() => setShowGrow(!showGrow)} sx={{ mb: 3 }}>Toggle Transition</Button>
                    <Grid container spacing={4} justifyContent="center" alignItems="center">
                        <Grid size={{ xs: 'auto' }}>
                            <Fade in={showGrow}>
                                <Box sx={{ p: 3, bgcolor: alpha('#9c27b0', 0.1), borderRadius: 2, textAlign: 'center' }}>
                                    <CircularProgress color="secondary" />
                                    <Typography variant="caption" sx={{ display: 'block', mt: 1 }}>Fade Effect</Typography>
                                </Box>
                            </Fade>
                        </Grid>
                        <Grid size={{ xs: 'auto' }}>
                            <Grow in={showGrow}>
                                <Box sx={{ p: 3, bgcolor: alpha('#2e7d32', 0.1), borderRadius: 2, textAlign: 'center' }}>
                                    <CircularProgress color="success" />
                                    <Typography variant="caption" sx={{ display: 'block', mt: 1 }}>Grow Effect</Typography>
                                </Box>
                            </Grow>
                        </Grid>
                    </Grid>
                </Card>
            </Grid>
        </Grid>
    );
};
