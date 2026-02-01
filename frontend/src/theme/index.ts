import { createTheme } from '@mui/material/styles';
import type { ThemeOptions } from '@mui/material/styles';

export const getTheme = (mode: 'light' | 'dark') => {
    const themeOptions: ThemeOptions = {
        palette: {
            mode,
            primary: {
                main: '#10b981', // Emerald branding
                dark: '#059669',
                light: '#34d399',
            },
            secondary: {
                main: mode === 'light' ? '#334155' : '#94a3b8',
            },
            background: {
                default: mode === 'light' ? '#F8FAFC' : '#022c22', // Emerald dark background
                paper: mode === 'light' ? '#FFFFFF' : '#064e3b',
            },
            text: {
                primary: mode === 'light' ? 'rgba(0, 0, 0, 0.87)' : '#f8fafc',
                secondary: mode === 'light' ? 'rgba(0, 0, 0, 0.6)' : 'rgba(248, 250, 252, 0.7)',
            },
        },
        typography: {
            fontFamily: "'Roboto', sans-serif",
            h1: { fontWeight: 700 },
            h2: { fontWeight: 700 },
            h3: { fontWeight: 700 },
            h4: { fontWeight: 700 },
            h5: { fontWeight: 600 },
            h6: { fontWeight: 600 },
        },
        shape: {
            borderRadius: 8,
        },
        components: {
            MuiButton: {
                styleOverrides: {
                    root: {
                        textTransform: 'none',
                        fontWeight: 600,
                    },
                },
            },
            MuiAppBar: {
                styleOverrides: {
                    root: {
                        backgroundColor: mode === 'light' ? 'rgba(255, 255, 255, 0.8)' : 'rgba(2, 44, 34, 0.8)',
                        backdropFilter: 'blur(8px)',
                        color: mode === 'light' ? 'inherit' : '#fff',
                    },
                },
            },
        },
    };

    return createTheme(themeOptions);
};

const defaultTheme = getTheme('light');
export default defaultTheme;
