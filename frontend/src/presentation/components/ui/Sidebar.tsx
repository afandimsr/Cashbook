import React, { useState } from 'react';
import {
    Box,
    Drawer,
    List,
    ListItem,
    ListItemButton,
    ListItemIcon,
    ListItemText,
    Toolbar,
    Typography,
    alpha,
    Collapse,
} from '@mui/material';
import LogoutIcon from '@mui/icons-material/Logout';
import ExpandLess from '@mui/icons-material/ExpandLess';
import ExpandMore from '@mui/icons-material/ExpandMore';
import { useNavigate, useLocation, Link } from 'react-router-dom';
import { useAuthStore } from '../../../state/authStore';
import { menuItems } from '../../../app/navigation';
import type { NavItem } from '../../../app/navigation';

const drawerWidth = 260;
const miniDrawerWidth = 80;

interface SidebarProps {
    mobileOpen: boolean;
    sidebarOpen: boolean;
    onClose: () => void;
}

const SidebarItem: React.FC<{
    item: NavItem;
    depth?: number;
    sidebarOpen: boolean;
    onClose: () => void;
    active: boolean;
}> = ({ item, depth = 0, sidebarOpen, onClose, active }) => {
    const [open, setOpen] = useState(false);
    const hasChildren = item.children && item.children.length > 0;
    const location = useLocation();

    const handleToggle = (e: React.MouseEvent) => {
        if (hasChildren) {
            e.preventDefault();
            setOpen(!open);
        } else {
            onClose();
        }
    };

    return (
        <>
            <ListItem disablePadding sx={{ display: 'block' }}>
                <ListItemButton
                    component={item.path ? Link : 'div'}
                    {...(item.path ? { to: item.path } : {})}
                    onClick={handleToggle}
                    sx={{
                        minHeight: 48,
                        justifyContent: sidebarOpen ? 'initial' : 'center',
                        px: 2.5,
                        pl: sidebarOpen ? 2.5 + depth * 2 : 2.5,
                        borderRadius: 3,
                        mb: 0.5,
                        transition: 'all 0.2s ease-in-out',
                        bgcolor: active ? (theme) => alpha(theme.palette.primary.main, 0.12) : 'transparent',
                        color: active ? 'primary.main' : 'rgba(255, 255, 255, 0.7)',
                        '&:hover': {
                            bgcolor: active
                                ? (theme) => alpha(theme.palette.primary.main, 0.18)
                                : 'rgba(255, 255, 255, 0.05)',
                            transform: sidebarOpen ? 'translateX(4px)' : 'none',
                        },
                        '& .MuiListItemIcon-root': {
                            color: active ? 'primary.main' : 'rgba(255, 255, 255, 0.7)',
                        }
                    }}
                >
                    <ListItemIcon
                        sx={{
                            minWidth: 0,
                            mr: sidebarOpen ? 2 : 'auto',
                            justifyContent: 'center',
                            transition: 'all 0.2s',
                            transform: active ? 'scale(1.1)' : 'none'
                        }}
                    >
                        {item.icon}
                    </ListItemIcon>
                    {sidebarOpen && (
                        <>
                            <ListItemText
                                primary={item.text}
                                sx={{
                                    '& span': {
                                        fontWeight: active ? 700 : 500,
                                        fontSize: depth > 0 ? '0.85rem' : '0.9rem'
                                    }
                                }}
                            />
                            {hasChildren && (open ? <ExpandLess sx={{ fontSize: 18 }} /> : <ExpandMore sx={{ fontSize: 18 }} />)}
                        </>
                    )}
                </ListItemButton>
            </ListItem>
            {hasChildren && sidebarOpen && (
                <Collapse in={open} timeout="auto" unmountOnExit>
                    <List component="div" disablePadding sx={{ gap: 0.5, display: 'flex', flexDirection: 'column' }}>
                        {item.children?.map((child) => (
                            <SidebarItem
                                key={child.text}
                                item={child}
                                depth={depth + 1}
                                sidebarOpen={sidebarOpen}
                                onClose={onClose}
                                active={!!child.path && location.pathname === child.path}
                            />
                        ))}
                    </List>
                </Collapse>
            )}
        </>
    );
};

export const Sidebar: React.FC<SidebarProps> = ({ mobileOpen, sidebarOpen, onClose }) => {
    const navigate = useNavigate();
    const location = useLocation();
    const { logout } = useAuthStore();
    const { user } = useAuthStore();

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    const drawer = (
        <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Toolbar sx={{ px: [2, 3], mb: 2 }}>
                <Box sx={{
                    display: 'flex',
                    alignItems: 'center',
                    gap: 1.5,
                    overflow: 'hidden',
                    transition: 'all 0.3s'
                }}>
                    <Box sx={{
                        width: 40,
                        height: 40,
                        borderRadius: 2,
                        bgcolor: 'primary.main',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: '#022c22', // Match deeper emerald
                        boxShadow: (theme) => `0 4px 12px ${alpha(theme.palette.primary.main, 0.4)}`,
                        flexShrink: 0
                    }}>
                        <span className="material-symbols-outlined" style={{ fontSize: 24, fontWeight: 'bold' }}>account_balance_wallet</span>
                    </Box>
                    {sidebarOpen && (
                        <Typography variant="h6" sx={{
                            fontWeight: 800,
                            fontSize: '1.25rem',
                            letterSpacing: '-0.02em',
                            color: '#fff',
                            whiteSpace: 'nowrap'
                        }}>
                            CashBook
                        </Typography>
                    )}
                </Box>
            </Toolbar>

            <Box sx={{ px: 2, flexGrow: 1, overflowY: 'auto', overflowX: 'hidden' }}>
                <List sx={{ gap: 0.5, display: 'flex', flexDirection: 'column' }}>
                    {menuItems
                        .filter((item) => {
                            if (!item.roles || item.roles.length === 0) return true;
                            if (!user || !user.roles) return false;
                            return item.roles.some(r => user.roles?.includes(r));
                        })
                        .map((item) => {
                            const isChildActive = item.children?.some(child => !!child.path && location.pathname === child.path);
                            const active = (!!item.path && location.pathname === item.path) || !!isChildActive;
                            return (
                                <SidebarItem
                                    key={item.text}
                                    item={item}
                                    sidebarOpen={sidebarOpen}
                                    onClose={onClose}
                                    active={active}
                                />
                            );
                        })}
                </List>
            </Box>

            <Box sx={{ p: 2, borderTop: (theme) => `1px solid ${theme.palette.divider}` }}>
                <ListItemButton
                    onClick={handleLogout}
                    sx={{
                        minHeight: 48,
                        justifyContent: sidebarOpen ? 'initial' : 'center',
                        px: 2.5,
                        borderRadius: 3,
                        color: 'rgba(255, 255, 255, 0.5)',
                        '&:hover': {
                            bgcolor: (theme) => alpha(theme.palette.error.main, 0.08),
                            color: 'error.main',
                        },
                    }}
                >
                    <ListItemIcon
                        sx={{
                            minWidth: 0,
                            mr: sidebarOpen ? 2 : 'auto',
                            justifyContent: 'center',
                            color: 'inherit'
                        }}
                    >
                        <LogoutIcon />
                    </ListItemIcon>
                    {sidebarOpen && (
                        <ListItemText
                            primary="Sign Out"
                            sx={{ '& span': { fontWeight: 600, fontSize: '0.9rem' } }}
                        />
                    )}
                </ListItemButton>
            </Box>
        </Box>
    );

    return (
        <Box
            component="nav"
            sx={{
                width: { sm: sidebarOpen ? drawerWidth : miniDrawerWidth },
                flexShrink: { sm: 0 },
                transition: 'width 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
            }}
        >
            <Drawer
                variant="temporary"
                open={mobileOpen}
                onClose={onClose}
                ModalProps={{ keepMounted: true }}
                sx={{
                    display: { xs: 'block', sm: 'none' },
                    '& .MuiDrawer-paper': {
                        boxSizing: 'border-box',
                        width: drawerWidth,
                        bgcolor: '#022c22', // Darker, richer emerald
                        color: 'rgba(255, 255, 255, 0.7)',
                        backgroundImage: 'none',
                        borderRight: 'none',
                        boxShadow: '4px 0 24px rgba(0,0,0,0.1)'
                    },
                }}
            >
                {drawer}
            </Drawer>

            <Drawer
                variant="permanent"
                sx={{
                    display: { xs: 'none', sm: 'block' },
                    '& .MuiDrawer-paper': {
                        boxSizing: 'border-box',
                        width: sidebarOpen ? drawerWidth : miniDrawerWidth,
                        transition: 'width 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                        overflowX: 'hidden',
                        bgcolor: '#022c22', // Darker, richer emerald for more "professional" feel
                        color: 'rgba(255, 255, 255, 0.7)',
                        borderRight: 'none',
                        backgroundImage: 'none',
                        boxShadow: '4px 0 24px rgba(0,0,0,0.05)'
                    },
                }}
                open
            >
                {drawer}
            </Drawer>
        </Box>
    );
};
