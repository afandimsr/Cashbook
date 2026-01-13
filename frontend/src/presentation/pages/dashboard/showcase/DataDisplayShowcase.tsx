import React from 'react';
import {
    Grid,
    Card,
    Typography,
    Badge,
    Tooltip,
    Chip,
    Avatar,
    Stack,
    Box,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
    Divider,
    Skeleton,
    alpha
} from '@mui/material';
import MailIcon from '@mui/icons-material/Mail';
import NotificationsIcon from '@mui/icons-material/Notifications';

export const DataDisplayShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(4)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12, md: 6 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="50%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={150} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Chips & Badges */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Badges & Chips</Typography>
                    <Box sx={{ mb: 4 }}>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Badges</Typography>
                        <Stack direction="row" spacing={4}>
                            <Badge badgeContent={4} color="primary">
                                <MailIcon color="action" />
                            </Badge>
                            <Badge badgeContent={10} color="secondary">
                                <NotificationsIcon color="action" />
                            </Badge>
                            <Badge variant="dot" color="error">
                                <MailIcon color="action" />
                            </Badge>
                        </Stack>
                    </Box>
                    <Divider sx={{ my: 3 }} />
                    <Box>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Chips</Typography>
                        <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap>
                            <Chip label="Primary" color="primary" />
                            <Chip label="Success" color="success" variant="outlined" />
                            <Chip label="Warning" color="warning" />
                            <Chip label="Error" color="error" clickable />
                            <Chip label="Delete" onDelete={() => { }} />
                        </Stack>
                    </Box>
                </Card>
            </Grid>

            {/* Avatars & Tooltips */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Avatars & Tooltips</Typography>
                    <Box sx={{ mb: 4 }}>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Avatars</Typography>
                        <Stack direction="row" spacing={2}>
                            <Avatar alt="Remy Sharp" src="https://mui.com/static/images/avatar/1.jpg" />
                            <Avatar alt="Travis Howard" src="https://mui.com/static/images/avatar/2.jpg" />
                            <Avatar sx={{ bgcolor: alpha('#1976d2', 0.1), color: '#1976d2' }}>JD</Avatar>
                            <Avatar variant="rounded" sx={{ bgcolor: 'secondary.main' }}>N</Avatar>
                        </Stack>
                    </Box>
                    <Divider sx={{ my: 3 }} />
                    <Box>
                        <Typography variant="subtitle2" sx={{ mb: 2 }}>Tooltips</Typography>
                        <Stack direction="row" spacing={2}>
                            <Tooltip title="Delete Action" arrow>
                                <Chip label="Hover me (Top)" sx={{ cursor: 'pointer' }} />
                            </Tooltip>
                            <Tooltip title="Custom Tooltip" placement="right">
                                <Chip label="Right Placement" variant="outlined" sx={{ cursor: 'pointer' }} />
                            </Tooltip>
                        </Stack>
                    </Box>
                </Card>
            </Grid>

            {/* Typography & Stats Table */}
            <Grid size={{ xs: 12 }}>
                <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Data Table & Typography</Typography>
                    <Grid container spacing={4}>
                        <Grid size={{ xs: 12, md: 4 }}>
                            <Typography variant="h3">H3 Heading</Typography>
                            <Typography variant="h5" color="secondary">H5 Secondary</Typography>
                            <Typography variant="body1">Body 1 text. Lorem ipsum dolor sit amet, consectetur adipiscing elit.</Typography>
                            <Typography variant="caption" color="text.secondary">Caption text for small notes.</Typography>
                        </Grid>
                        <Grid size={{ xs: 12, md: 8 }}>
                            <TableContainer component={Paper} elevation={0} sx={{ border: (theme) => `1px solid ${theme.palette.divider}`, borderRadius: 2 }}>
                                <Table size="small">
                                    <TableHead sx={{ bgcolor: (theme) => alpha(theme.palette.primary.main, 0.05) }}>
                                        <TableRow>
                                            <TableCell sx={{ fontWeight: 700 }}>Item</TableCell>
                                            <TableCell align="right" sx={{ fontWeight: 700 }}>Category</TableCell>
                                            <TableCell align="right" sx={{ fontWeight: 700 }}>Status</TableCell>
                                        </TableRow>
                                    </TableHead>
                                    <TableBody>
                                        {[
                                            { name: 'Apple UI', cat: 'Design', status: 'Live' },
                                            { name: 'Node Backend', cat: 'Dev', status: 'Pending' },
                                            { name: 'Marketing PR', cat: 'Sales', status: 'Review' },
                                        ].map((row) => (
                                            <TableRow key={row.name}>
                                                <TableCell>{row.name}</TableCell>
                                                <TableCell align="right">{row.cat}</TableCell>
                                                <TableCell align="right">
                                                    <Chip label={row.status} size="small" variant="outlined" />
                                                </TableCell>
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </TableContainer>
                        </Grid>
                    </Grid>
                </Card>
            </Grid>
        </Grid>
    );
};
