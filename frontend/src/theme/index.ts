import { createTheme } from '@mui/material/styles';
import type { ThemeOptions } from '@mui/material/styles';

export const getTheme = (mode: 'light' | 'dark') => {
    const themeOptions: ThemeOptions = {
        palette: {
            mode,
            primary: {
                main: mode === 'light' ? '#1976d2' : '#90caf9',
            },
            secondary: {
                main: mode === 'light' ? '#dc004e' : '#f48fb1',
            },
            background: {
                default: mode === 'light' ? '#f5f7fa' : '#0f172a',
                paper: mode === 'light' ? '#ffffff' : '#1e293b',
            },
        },
        components: {
            MuiAppBar: {
                styleOverrides: {
                    root: {
                        backgroundColor: mode === 'light' ? 'rgba(255, 255, 255, 0.8)' : 'rgba(15, 23, 42, 0.8)',
                    },
                },
            },
        },
    };

    return createTheme(themeOptions);
};

const defaultTheme = getTheme('light');
export default defaultTheme;
