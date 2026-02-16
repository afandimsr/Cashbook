import React, { useState } from 'react';
import {
    Box,
    Button,
    IconButton,
    InputAdornment,
    TextField,
    Typography,
    Alert,
    CircularProgress,
    Divider
} from '@mui/material';
import { styled } from '@mui/material/styles';
import { useAuthStore } from '../../../state/authStore';
import { useNavigate } from 'react-router-dom';

// Design Constants from login.html
const PRIMARY_COLOR = '#10b981';
const PRIMARY_DARK = '#059669';
const LEFT_BG = '#062d24';
const TEXT_PRIMARY = 'rgba(0, 0, 0, 0.87)';
const TEXT_SECONDARY = 'rgba(0, 0, 0, 0.6)';
const DIVIDER_COLOR = 'rgba(0, 0, 0, 0.12)';

// Styled Components
const PageContainer = styled(Box)({
    minHeight: '100vh',
    display: 'flex',
    fontFamily: "'Roboto', sans-serif",
});

const LeftPanel = styled(Box)(({ theme }) => ({
    display: 'none',
    position: 'relative',
    flexDirection: 'column',
    justifyContent: 'space-between',
    padding: '64px',
    backgroundColor: LEFT_BG,
    overflow: 'hidden',
    width: '50%',
    [theme.breakpoints.up('lg')]: {
        display: 'flex',
    },
}));

const GlowEffect = styled(Box)({
    position: 'absolute',
    top: '-10%',
    left: '-10%',
    width: '70%',
    height: '70%',
    borderRadius: '50%',
    backgroundColor: PRIMARY_COLOR,
    filter: 'blur(120px)',
    opacity: 0.2,
    zIndex: 0,
});

const RightPanel = styled(Box)(({ theme }) => ({
    width: '100%',
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    alignItems: 'center',
    padding: '24px',
    backgroundColor: '#ffffff',
    [theme.breakpoints.up('lg')]: {
        width: '50%',
    },
}));

const CustomTextField = styled(TextField)({
    '& .MuiOutlinedInput-root': {
        color: TEXT_PRIMARY,
        '& fieldset': {
            borderColor: DIVIDER_COLOR,
            transition: 'border-color 200ms cubic-bezier(0.4, 0, 0.2, 1)',
        },
        '&:hover fieldset': {
            borderColor: 'rgba(0, 0, 0, 0.87)',
        },
        '&.Mui-focused fieldset': {
            borderColor: PRIMARY_COLOR,
            borderWidth: 2,
        },
        '& input': {
            padding: '16.5px 14px',
            fontSize: '1rem',
            lineHeight: '1.4375em',
            color: TEXT_PRIMARY, // Ensure text is dark and visible
        }
    },
    '& .MuiInputLabel-root': {
        color: TEXT_SECONDARY,
        transform: 'translate(14px, 16px) scale(1)',
        backgroundColor: 'transparent',
    },
    '& .MuiInputLabel-root.Mui-focused, & .MuiInputLabel-root.MuiFormLabel-filled': {
        transform: 'translate(10px, -9px) scale(0.75)',
        color: PRIMARY_COLOR,
        backgroundColor: 'white',
        padding: '0 4px',
    },
});

const SignInButton = styled(Button)({
    backgroundColor: PRIMARY_COLOR,
    color: 'white',
    padding: '8px 22px',
    fontSize: '0.875rem',
    fontWeight: 500,
    textTransform: 'uppercase',
    borderRadius: 4,
    boxShadow: '0px 3px 1px -2px rgba(0,0,0,0.2), 0px 2px 2px 0px rgba(0,0,0,0.14), 0px 1px 5px 0px rgba(0,0,0,0.12)',
    '&:hover': {
        backgroundColor: PRIMARY_DARK,
        boxShadow: '0px 2px 4px -1px rgba(0,0,0,0.2), 0px 4px 5px 0px rgba(0,0,0,0.14), 0px 1px 10px 0px rgba(0,0,0,0.12)',
    },
});

const GoogleButton = styled(Button)({
    border: '1px solid rgba(0, 0, 0, 0.23)',
    backgroundColor: 'transparent',
    color: TEXT_PRIMARY,
    padding: '8px 22px',
    fontSize: '0.875rem',
    fontWeight: 500,
    textTransform: 'uppercase',
    borderRadius: 4,
    '&:hover': {
        backgroundColor: 'rgba(0, 0, 0, 0.04)',
        borderColor: TEXT_PRIMARY,
    },
});

export const LoginPageNew: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);

    const { login, isLoading, error } = useAuthStore();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login(email, password);
            navigate('/dashboard');
        } catch (err) {
            // Error managed by store
        }
    };

    const handleGoogleLogin = () => {
        window.location.href = '/api/v1/auth/google/login';
    };

    return (
        <PageContainer>
            {/* Left Panel */}
            <LeftPanel>
                <GlowEffect />
                <Box sx={{ position: 'relative', zIndex: 10, display: 'flex', alignItems: 'center', gap: 2 }}>
                    <Box sx={{
                        height: 48,
                        width: 48,
                        borderRadius: 1,
                        bgcolor: PRIMARY_COLOR,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                        color: LEFT_BG,
                        boxShadow: '0px 3px 5px -1px rgba(0,0,0,0.2),0px 6px 10px 0px rgba(0,0,0,0.14),0px 1px 18px 0px rgba(0,0,0,0.12)'
                    }}>
                        <span className="material-symbols-outlined" style={{ fontSize: 32 }}>account_balance_wallet</span>
                    </Box>
                    <Typography variant="h4" sx={{ color: 'white', fontWeight: 700, letterSpacing: '-0.3px', m: 0 }}>
                        CashBook
                    </Typography>
                </Box>

                <Box sx={{ position: 'relative', zIndex: 10, display: 'flex', flexDirection: 'column', gap: 3, maxWidth: 512 }}>
                    <Typography variant="h2" sx={{
                        color: 'white',
                        fontWeight: 700,
                        lineHeight: 1.1,
                        fontSize: '3rem',
                        m: 0
                    }}>
                        Smart wealth <span style={{ color: PRIMARY_COLOR }}>management</span> for your future.
                    </Typography>
                    <Typography variant="body1" sx={{
                        color: 'rgba(209, 250, 229, 0.7)',
                        fontSize: '1.125rem',
                        lineHeight: 1.6,
                        m: 0
                    }}>
                        Take control of your financial journey with CashBook. Professional tools designed to help you track, save, and grow your assets.
                    </Typography>
                </Box>

                <Typography sx={{
                    position: 'relative',
                    zIndex: 10,
                    color: 'rgba(16, 185, 129, 0.5)',
                    fontSize: '0.75rem',
                    fontWeight: 500,
                    textTransform: 'uppercase',
                    letterSpacing: '0.2em',
                    m: 0
                }}>
                    © 2026 CashBook
                </Typography>
            </LeftPanel>

            {/* Right Panel */}
            <RightPanel>
                <Box sx={{
                    width: '100%',
                    maxWidth: 448,
                    display: 'flex',
                    flexDirection: 'column',
                    padding: { xs: '32px', sm: '48px' }
                }}>
                    {/* Mobile Logo */}
                    <Box sx={{ display: { xs: 'flex', lg: 'none' }, alignItems: 'center', gap: 1.5, mb: 4 }}>
                        <span className="material-symbols-outlined" style={{ color: PRIMARY_COLOR, fontSize: 36 }}>account_balance_wallet</span>
                        <Typography variant="h5" sx={{ fontWeight: 700, color: '#111827', m: 0 }}>
                            CashBook
                        </Typography>
                    </Box>

                    <Box sx={{ mb: 4 }}>
                        <Typography variant="h5" sx={{ fontWeight: 500, fontSize: '1.5rem', letterSpacing: 'normal', color: TEXT_PRIMARY, m: 0 }}>
                            Sign In
                        </Typography>
                        <Typography variant="body2" sx={{ color: TEXT_SECONDARY, mt: 1, m: 0 }}>
                            Enter your email and password to continue
                        </Typography>
                    </Box>

                    {error && (
                        <Alert severity="error" sx={{ mb: 3, borderRadius: 1 }}>
                            {error}
                        </Alert>
                    )}

                    <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
                        <CustomTextField
                            fullWidth
                            label="Email Address"
                            type="email"
                            variant="outlined"
                            autoComplete="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                        />

                        <CustomTextField
                            fullWidth
                            label="Password"
                            type={showPassword ? 'text' : 'password'}
                            variant="outlined"
                            autoComplete="current-password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton onClick={() => setShowPassword(!showPassword)} edge="end" sx={{ color: 'rgba(0, 0, 0, 0.5)' }}>
                                            <span className="material-symbols-outlined" style={{ fontSize: 20 }}>
                                                {showPassword ? 'visibility_off' : 'visibility'}
                                            </span>
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                        />

                        <SignInButton
                            type="submit"
                            fullWidth
                            disabled={isLoading}
                        >
                            {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Sign In'}
                        </SignInButton>

                        <Box sx={{ display: 'flex', alignItems: 'center', my: 1 }}>
                            <Divider sx={{ flexGrow: 1, bgcolor: DIVIDER_COLOR }} />
                            <Typography variant="caption" sx={{ px: 2, color: TEXT_SECONDARY, textTransform: 'uppercase', fontSize: '0.75rem' }}>
                                OR
                            </Typography>
                            <Divider sx={{ flexGrow: 1, bgcolor: DIVIDER_COLOR }} />
                        </Box>

                        <GoogleButton
                            type="button"
                            fullWidth
                            startIcon={
                                <svg width="20" height="20" viewBox="0 0 24 24">
                                    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4" />
                                    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853" />
                                    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05" />
                                    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.66l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335" />
                                </svg>
                            }
                            onClick={handleGoogleLogin}
                        >
                            <Typography variant="inherit" sx={{ ml: 1 }}>Sign in with Google</Typography>
                        </GoogleButton>
                    </Box>
                </Box>
            </RightPanel>
        </PageContainer>
    );
};

export default LoginPageNew;
