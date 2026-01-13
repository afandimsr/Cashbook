import React, { useState } from 'react';
import {
    AppBar,
    Toolbar,
    IconButton,
    Typography,
    Box,
    Avatar,
    Menu,
    MenuItem,
    ListItemIcon,
    Divider,
    Tooltip,
    alpha,
    InputBase,
} from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import ChevronLeftIcon from '@mui/icons-material/ChevronLeft';
import LogoutIcon from '@mui/icons-material/Logout';
import PersonIcon from '@mui/icons-material/Person';
import SettingsIcon from '@mui/icons-material/Settings';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import SearchIcon from '@mui/icons-material/Search';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../state/authStore';
import { useThemeStore } from '../../../state/themeStore';

const drawerWidth = 260;
const miniDrawerWidth = 80;

interface TopbarProps {
    sidebarOpen: boolean;
    onMobileDrawerToggle: () => void;
    onSidebarToggle: () => void;
}

export const Topbar: React.FC<TopbarProps> = ({
    sidebarOpen,
    onMobileDrawerToggle,
    onSidebarToggle
}) => {
    const navigate = useNavigate();
    const { user, logout } = useAuthStore();
    const { mode, toggleMode } = useThemeStore();
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
    const [searchQuery, setSearchQuery] = useState('');
    const open = Boolean(anchorEl);

    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (searchQuery.trim()) {
            navigate(`/dashboard/transactions?q=${encodeURIComponent(searchQuery)}`);
        }
    };

    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };

    const handleClose = () => {
        setAnchorEl(null);
    };

    const handleLogout = () => {
        handleClose();
        logout();
    };

    return (
        <AppBar
            position="fixed"
            sx={{
                width: {
                    xs: '100%',
                    sm: sidebarOpen ? `calc(100% - ${drawerWidth}px)` : `calc(100% - ${miniDrawerWidth}px)`
                },
                ml: {
                    xs: 0,
                    sm: sidebarOpen ? `${drawerWidth}px` : `${miniDrawerWidth}px`
                },
                transition: (theme) => theme.transitions.create(['width', 'margin', 'background-color'], {
                    easing: theme.transitions.easing.sharp,
                    duration: theme.transitions.duration.leavingScreen,
                }),
                bgcolor: (theme) => alpha(theme.palette.background.paper, 0.8),
                backdropFilter: 'blur(8px)',
                color: 'text.primary',
                boxShadow: 'none',
                borderBottom: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                zIndex: (theme) => theme.zIndex.drawer + 1,
            }}
        >
            <Toolbar sx={{ px: [2, 3] }}>
                <IconButton
                    color="inherit"
                    aria-label="open drawer"
                    edge="start"
                    onClick={onMobileDrawerToggle}
                    sx={{ mr: 2, display: { sm: 'none' } }}
                >
                    <MenuIcon />
                </IconButton>
                <IconButton
                    color="inherit"
                    aria-label="toggle sidebar"
                    edge="start"
                    onClick={onSidebarToggle}
                    sx={{
                        mr: 2,
                        display: { xs: 'none', sm: 'inline-flex' },
                        bgcolor: (theme) => alpha(theme.palette.action.hover, 0.05),
                        '&:hover': { bgcolor: (theme) => alpha(theme.palette.action.hover, 0.1) }
                    }}
                >
                    {sidebarOpen ? <ChevronLeftIcon /> : <MenuIcon />}
                </IconButton>
                <Typography variant="h6" noWrap component="div" sx={{
                    fontWeight: 800,
                    fontSize: '1.25rem',
                    color: 'primary.main',
                    letterSpacing: '-0.5px',
                    display: { xs: 'none', lg: 'block' },
                    mr: 4
                }}>
                    CashBook
                </Typography>

                <Box
                    component="form"
                    onSubmit={handleSearch}
                    sx={{
                        flexGrow: 1,
                        maxWidth: 400,
                        position: 'relative',
                        display: { xs: 'none', md: 'block' }
                    }}
                >
                    <Box sx={{
                        position: 'absolute',
                        height: '100%',
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        px: 2,
                        pointerEvents: 'none',
                        color: 'text.secondary'
                    }}>
                        <SearchIcon fontSize="small" />
                    </Box>
                    <InputBase
                        placeholder="Search transactions..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        sx={{
                            width: '100%',
                            '& .MuiInputBase-input': {
                                py: 1.2,
                                pl: 6,
                                pr: 2,
                                borderRadius: 2,
                                bgcolor: (theme) => alpha(theme.palette.action.hover, 0.05),
                                border: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`,
                                transition: 'all 0.2s',
                                '&:focus': {
                                    bgcolor: (theme) => alpha(theme.palette.background.paper, 1),
                                    boxShadow: (theme) => `0 4px 12px ${alpha(theme.palette.primary.main, 0.08)}`,
                                    borderColor: 'primary.main',
                                }
                            }
                        }}
                    />
                </Box>

                <Box sx={{ flexGrow: 1, display: { md: 'none' } }} />

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, ml: 'auto', }}>
                    <Tooltip title={`Switch to ${mode === 'light' ? 'dark' : 'light'} mode`}>
                        <IconButton
                            onClick={toggleMode}
                            sx={{
                                color: mode === 'light' ? 'warning.main' : 'primary.light',
                                bgcolor: (theme) => alpha(theme.palette.action.hover, 0.05),
                                '&:hover': {
                                    bgcolor: (theme) => alpha(theme.palette.action.hover, 0.1),
                                    transform: 'scale(1.1) rotate(10deg)',
                                },
                                transition: 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)',
                                mr: 1
                            }}
                        >
                            {mode === 'light' ? <DarkModeIcon /> : <LightModeIcon />}
                        </IconButton>
                    </Tooltip>

                    <Box sx={{
                        display: { xs: 'none', md: 'flex' },
                        flexDirection: 'column',
                        alignItems: 'flex-end',
                        mr: 1
                    }}>
                        <Typography variant="body2" sx={{ fontWeight: 700, lineHeight: 1 }}>
                            {user?.name}
                        </Typography>
                        <Typography variant="caption" sx={{ color: 'text.secondary', fontWeight: 500 }}>
                            {user?.role === 'ADMIN' ? 'Administrator' : 'User'}
                        </Typography>
                    </Box>
                    <Tooltip title="Account settings">
                        <IconButton
                            onClick={handleClick}
                            size="small"
                            sx={{
                                p: 0.5,
                                border: (theme) => `1px solid ${alpha(theme.palette.divider, 0.1)}`
                            }}
                            aria-controls={open ? 'account-menu' : undefined}
                            aria-haspopup="true"
                            aria-expanded={open ? 'true' : undefined}
                        >
                            <Avatar sx={{
                                width: 34,
                                height: 34,
                                bgcolor: 'primary.main',
                                fontSize: '0.875rem',
                                fontWeight: 600,
                                boxShadow: (theme) => `0 2px 8px ${alpha(theme.palette.primary.main, 0.3)}`
                            }}>
                                {user?.username?.charAt(0).toUpperCase() || 'U'}
                            </Avatar>
                        </IconButton>
                    </Tooltip>
                </Box>
                <Menu
                    anchorEl={anchorEl}
                    id="account-menu"
                    open={open}
                    onClose={handleClose}
                    onClick={handleClose}
                    transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                    anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
                    PaperProps={{
                        elevation: 0,
                        sx: {
                            overflow: 'visible',
                            filter: 'drop-shadow(0px 4px 20px rgba(0,0,0,0.1))',
                            borderRadius: 3,
                            mt: 1.5,
                            minWidth: 200,
                            border: (theme) => `1px solid ${alpha(theme.palette.divider, 0.08)}`,
                            '& .MuiAvatar-root': {
                                width: 32,
                                height: 32,
                                ml: -0.5,
                                mr: 1,
                            },
                        },
                    }}
                >
                    <Box sx={{ px: 2, py: 1.5 }}>
                        <Typography variant="subtitle2" sx={{ fontWeight: 700 }}>My Account</Typography>
                        <Typography variant="caption" color="text.secondary">{user?.email}</Typography>
                    </Box>
                    <Divider sx={{ my: 0.5, opacity: 0.6 }} />
                    <MenuItem onClick={() => { handleClose(); navigate('/dashboard/profile'); }} sx={{ py: 1, px: 2 }}>
                        <ListItemIcon>
                            <PersonIcon fontSize="small" />
                        </ListItemIcon>
                        Profile
                    </MenuItem>
                    <MenuItem onClick={handleClose} sx={{ py: 1, px: 2 }}>
                        <ListItemIcon>
                            <SettingsIcon fontSize="small" />
                        </ListItemIcon>
                        Settings
                    </MenuItem>
                    <Divider sx={{ my: 0.5, opacity: 0.6 }} />
                    <MenuItem onClick={handleLogout} sx={{ color: 'error.main', py: 1, px: 2 }}>
                        <ListItemIcon>
                            <LogoutIcon fontSize="small" color="error" />
                        </ListItemIcon>
                        Logout
                    </MenuItem>
                </Menu>
            </Toolbar>
        </AppBar>
    );
};
