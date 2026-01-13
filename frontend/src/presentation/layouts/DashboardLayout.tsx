import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Box, Toolbar } from '@mui/material';
import { Sidebar } from '../components/ui/Sidebar';
import { Topbar } from '../components/ui/Topbar';
import { Footer } from '../components/ui/Footer';

const drawerWidth = 260;
const miniDrawerWidth = 80;

export const DashboardLayout: React.FC = () => {
    const [mobileOpen, setMobileOpen] = useState(false);
    const [sidebarOpen, setSidebarOpen] = useState(true);

    const handleDrawerToggle = () => {
        setMobileOpen(!mobileOpen);
    };

    const handleSidebarToggle = () => {
        setSidebarOpen(!sidebarOpen);
    };

    return (
        <Box sx={{
            display: 'flex',
            minHeight: '100vh',
            bgcolor: (theme) => theme.palette.mode === 'light' ? '#F8FAFC' : 'background.default',
            transition: 'background-color 0.3s ease-in-out'
        }}>
            <Topbar
                sidebarOpen={sidebarOpen}
                onMobileDrawerToggle={handleDrawerToggle}
                onSidebarToggle={handleSidebarToggle}
            />

            <Sidebar mobileOpen={mobileOpen} sidebarOpen={sidebarOpen} onClose={handleDrawerToggle} />

            <Box
                component="main"
                sx={{
                    flexGrow: 1,
                    width: {
                        xs: '100%',
                        sm: sidebarOpen ? `calc(100% - ${drawerWidth}px)` : `calc(100% - ${miniDrawerWidth}px)`
                    },
                    transition: (theme) => theme.transitions.create(['width', 'margin', 'padding'], {
                        easing: theme.transitions.easing.sharp,
                        duration: theme.transitions.duration.leavingScreen,
                    }),
                    minHeight: '100vh',
                    display: 'flex',
                    flexDirection: 'column'
                }}
            >
                <Toolbar /> {/* Spacer for Topbar */}
                <Box sx={{
                    flexGrow: 1,
                    p: { xs: 2, sm: 3, md: 4 }
                }}>
                    <Outlet />
                </Box>
                <Footer />
            </Box>
        </Box>
    );
};
