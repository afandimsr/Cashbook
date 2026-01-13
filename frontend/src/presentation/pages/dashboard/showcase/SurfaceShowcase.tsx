import React from 'react';
import {
    Grid,
    Card,
    Typography,
    Accordion,
    AccordionSummary,
    AccordionDetails,
    Paper,
    Box,
    Stack,
    Skeleton,
    alpha
} from '@mui/material';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';

export const SurfaceShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(4)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12, md: 6 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="40%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={120} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Accordion */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Accordions</Typography>
                    <Box>
                        <Accordion defaultExpanded elevation={0} sx={{ border: (theme) => `1px solid ${theme.palette.divider}`, mb: 1 }}>
                            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                <Typography sx={{ fontWeight: 600 }}>Accordion Item 1</Typography>
                            </AccordionSummary>
                            <AccordionDetails>
                                <Typography variant="body2" color="text.secondary">
                                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse malesuada lacus ex, sit amet blandit leo lobortis eget.
                                </Typography>
                            </AccordionDetails>
                        </Accordion>
                        <Accordion elevation={0} sx={{ border: (theme) => `1px solid ${theme.palette.divider}` }}>
                            <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                                <Typography sx={{ fontWeight: 600 }}>Accordion Item 2</Typography>
                            </AccordionSummary>
                            <AccordionDetails>
                                <Typography variant="body2" color="text.secondary">
                                    Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse malesuada lacus ex, sit amet blandit leo lobortis eget.
                                </Typography>
                            </AccordionDetails>
                        </Accordion>
                    </Box>
                </Card>
            </Grid>

            {/* Paper Variations */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Papers & Elevation</Typography>
                    <Stack spacing={2}>
                        <Paper variant="outlined" sx={{ p: 2, textAlign: 'center' }}>Outlined Paper</Paper>
                        <Paper elevation={3} sx={{ p: 2, textAlign: 'center', bgcolor: alpha('#1976d2', 0.05) }}>Elevation 3 Paper</Paper>
                        <Paper elevation={8} sx={{ p: 2, textAlign: 'center' }}>Elevation 8 Paper</Paper>
                    </Stack>
                </Card>
            </Grid>

            {/* Complex Surfaces */}
            <Grid size={{ xs: 12 }}>
                <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Complex Surface Layout</Typography>
                    <Box sx={{ p: { xs: 2, md: 4 }, bgcolor: 'grey.50', borderRadius: 3, border: (theme) => `1px dashed ${theme.palette.divider}` }}>
                        <Grid container spacing={2}>
                            <Grid size={{ xs: 12, sm: 4 }}>
                                <Paper sx={{ p: 3, textAlign: 'center', borderRadius: 3 }}>
                                    <Typography variant="h4">12</Typography>
                                    <Typography variant="caption">Files Pending</Typography>
                                </Paper>
                            </Grid>
                            <Grid size={{ xs: 12, sm: 4 }}>
                                <Paper sx={{ p: 3, textAlign: 'center', borderRadius: 3, bgcolor: 'primary.main', color: 'primary.contrastText' }}>
                                    <Typography variant="h4">99%</Typography>
                                    <Typography variant="caption">Uptime</Typography>
                                </Paper>
                            </Grid>
                            <Grid size={{ xs: 12, sm: 4 }}>
                                <Paper sx={{ p: 3, textAlign: 'center', borderRadius: 3 }}>
                                    <Typography variant="h4">45k</Typography>
                                    <Typography variant="caption">Customers</Typography>
                                </Paper>
                            </Grid>
                        </Grid>
                    </Box>
                </Card>
            </Grid>
        </Grid>
    );
};
