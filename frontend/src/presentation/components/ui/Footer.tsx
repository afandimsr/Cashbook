import React from 'react';
import { Box, Typography, Container, Link, alpha, Stack } from '@mui/material';
import { config } from '../../../app/config';

export const Footer: React.FC = () => {
    const currentYear = new Date().getFullYear();

    return (
        <Box
            component="footer"
            sx={{
                py: 2.5,
                mt: 'auto',
                borderTop: (theme) => `1px solid ${theme.palette.divider}`,
                bgcolor: (theme) => alpha(theme.palette.background.paper, 0.4),
                backdropFilter: 'blur(8px)',
            }}
        >
            <Container maxWidth="xl">
                <Stack
                    direction={{ xs: 'column', sm: 'row' }}
                    spacing={2}
                    justifyContent="space-between"
                    alignItems="center"
                >
                    <Box sx={{ textAlign: { xs: 'center', sm: 'left' } }}>
                        <Typography variant="body2" color="text.secondary" sx={{ fontWeight: 500 }}>
                            © {currentYear} {config.APP_TITLE}. All rights reserved.
                        </Typography>
                    </Box>

                    <Stack
                        direction="row"
                        spacing={3}
                        alignItems="center"
                        sx={{
                            flexWrap: 'wrap',
                            justifyContent: 'center'
                        }}
                    >
                        <Link href="#" color="inherit" underline="hover" sx={{ fontSize: '0.875rem', color: 'text.secondary', '&:hover': { color: 'primary.main' } }}>
                            Privacy Policy
                        </Link>
                        <Link href="#" color="inherit" underline="hover" sx={{ fontSize: '0.875rem', color: 'text.secondary', '&:hover': { color: 'primary.main' } }}>
                            Terms of Service
                        </Link>
                        <Box sx={{
                            display: 'flex',
                            alignItems: 'center',
                            gap: 1,
                            px: 1.5,
                            py: 0.5,
                            borderRadius: 2,
                            bgcolor: (theme) => alpha(theme.palette.primary.main, 0.08),
                            border: (theme) => `1px solid ${alpha(theme.palette.primary.main, 0.1)}`
                        }}>
                            <Typography variant="caption" sx={{ color: 'primary.main', fontWeight: 700, letterSpacing: '0.05em' }}>
                                v{config.APP_VERSION}
                            </Typography>
                        </Box>
                    </Stack>
                </Stack>
            </Container>
        </Box>
    );
};
