import React, { useState } from 'react';
import { useAuthStore } from '../../../state/authStore';
import { useNavigate } from 'react-router-dom';
import {
    Box,
    Button,
    TextField,
    Typography,
    Alert,
    Stack,
    Card as MuiCard,
    FormControl,
    FormLabel,
    FormControlLabel,
    Checkbox,
    InputAdornment,
    IconButton,
    CircularProgress
} from '@mui/material';
import { styled, alpha } from '@mui/material/styles';
import { SitemarkIcon } from '../../components/common/CustomIcons';
import { Visibility, VisibilityOff, Google as GoogleIcon } from '@mui/icons-material';
import loginBg from '../../../assets/images/login_bg.png';
import { Divider } from '@mui/material';

const StyledCard = styled(MuiCard)(({ theme }) => ({
    display: 'flex',
    flexDirection: 'column',
    alignSelf: 'center',
    width: '100%',
    padding: theme.spacing(4),
    gap: theme.spacing(2),
    margin: 'auto',
    borderRadius: theme.spacing(2),
    backdropFilter: 'blur(10px)',
    backgroundColor: alpha(theme.palette.background.paper, 0.8),
    border: `1px solid ${alpha(theme.palette.divider, 0.1)}`,
    boxShadow: `0 8px 32px 0 ${alpha(theme.palette.common.black, 0.2)}`,
    [theme.breakpoints.up('sm')]: {
        maxWidth: '450px',
    },
}));

const PageContainer = styled(Stack)(({ theme }) => ({
    height: '100vh',
    width: '100vw',
    padding: theme.spacing(2),
    backgroundImage: `linear-gradient(rgba(0,0,0,0.3), rgba(0,0,0,0.3)), url(${loginBg})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    position: 'absolute',
    top: 0,
    left: 0,
}));

export const LoginPage: React.FC = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [showPassword, setShowPassword] = useState(false);
    const { login, isLoading, error } = useAuthStore();
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login(username, password);
            navigate('/dashboard');
        } catch (err) {
            // Error handled by store state
        }
    };

    return (
        <PageContainer direction="column">
            <StyledCard variant="outlined">
                <Box sx={{ display: 'flex', justifyContent: 'center', mb: 1 }}>
                    <SitemarkIcon />
                </Box>
                <Typography
                    component="h1"
                    variant="h4"
                    sx={{ width: '100%', fontSize: 'clamp(1.5rem, 10vw, 1.8rem)', fontWeight: 700, textAlign: 'center' }}
                >
                    Sign In
                </Typography>
                <Typography color="text.secondary" align="center" variant="body2" sx={{ mb: 2 }}>
                    Enter your credentials to access your account
                </Typography>

                {error && (
                    <Alert severity="error" sx={{ mb: 1, borderRadius: 2 }}>
                        {error}
                    </Alert>
                )}

                <Box
                    component="form"
                    onSubmit={handleSubmit}
                    noValidate
                    sx={{
                        display: 'flex',
                        flexDirection: 'column',
                        width: '100%',
                        gap: 2,
                    }}
                >
                    <FormControl>
                        <FormLabel htmlFor="username" sx={{ mb: 0.5, fontWeight: 500 }}>Username</FormLabel>
                        <TextField
                            id="username"
                            name="username"
                            placeholder="Type your username"
                            autoComplete="username"
                            autoFocus
                            required
                            fullWidth
                            variant="outlined"
                            size="small"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                    </FormControl>
                    <FormControl>
                        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 0.5 }}>
                            <FormLabel htmlFor="password" sx={{ fontWeight: 500 }}>Password</FormLabel>
                        </Box>
                        <TextField
                            name="password"
                            placeholder="••••••"
                            type={showPassword ? 'text' : 'password'}
                            id="password"
                            autoComplete="current-password"
                            required
                            fullWidth
                            variant="outlined"
                            size="small"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            InputProps={{
                                endAdornment: (
                                    <InputAdornment position="end">
                                        <IconButton
                                            onClick={() => setShowPassword(!showPassword)}
                                            edge="end"
                                            size="small"
                                        >
                                            {showPassword ? <VisibilityOff /> : <Visibility />}
                                        </IconButton>
                                    </InputAdornment>
                                ),
                            }}
                        />
                    </FormControl>
                    <FormControlLabel
                        control={<Checkbox value="remember" color="primary" size="small" />}
                        label={<Typography variant="body2">Remember me</Typography>}
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        disabled={isLoading}
                        sx={{
                            py: 1.2,
                            mt: 1,
                            borderRadius: '8px',
                            textTransform: 'none',
                            fontSize: '1rem',
                            fontWeight: 600,
                            boxShadow: '0 4px 12px rgba(0,0,0,0.1)'
                        }}
                    >
                        {isLoading ? <CircularProgress size={24} color="inherit" /> : 'Sign In'}
                    </Button>
                </Box>

                <Divider sx={{ my: 1 }}>
                    <Typography variant="body2" color="text.secondary">or</Typography>
                </Divider>

                <Button
                    fullWidth
                    variant="outlined"
                    onClick={() => window.location.href = 'http://localhost:8181/api/v1/auth/google/login'}
                    startIcon={<GoogleIcon />}
                    sx={{
                        py: 1,
                        borderRadius: '8px',
                        textTransform: 'none',
                        borderColor: alpha('#000', 0.1),
                        color: 'text.primary',
                        fontWeight: 500
                    }}
                >
                    Sign in with Google
                </Button>

                {/* <Divider sx={{ my: 1 }}>
                    <Typography variant="body2" color="text.secondary">or</Typography>
                </Divider> */}

                {/* <Stack spacing={2}>
                    <Button
                        fullWidth
                        variant="outlined"
                        onClick={() => alert('Sign in with Google')}
                        startIcon={<GoogleIcon />}
                        sx={{ borderRadius: '8px', textTransform: 'none', borderColor: alpha('#000', 0.1) }}
                    >
                        Sign in with Google
                    </Button>
                    <Button
                        fullWidth
                        variant="outlined"
                        onClick={() => alert('Sign in with Facebook')}
                        startIcon={<FacebookIcon />}
                        sx={{ borderRadius: '8px', textTransform: 'none', borderColor: alpha('#000', 0.1) }}
                    >
                        Sign in with Facebook
                    </Button>
                </Stack>

                <Typography sx={{ textAlign: 'center', mt: 2 }} variant="body2">
                    Don't have an account?{' '}
                    <Link
                        href="#"
                        variant="body2"
                        sx={{ fontWeight: 600 }}
                    >
                        Sign up
                    </Link>
                </Typography> */}

                <Typography variant="caption" align="center" sx={{ mt: 2, color: 'text.secondary', fontStyle: 'italic' }}>
                    Demo credentials: admin / admin
                </Typography>
            </StyledCard>
        </PageContainer>
    );
};

