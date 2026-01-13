import React, { useEffect, useMemo } from 'react';
import { BrowserRouter } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { getTheme } from '../theme';
import { AppRoutes } from './routes';
import { useThemeStore } from '../state/themeStore';
import { ScrollToTop } from '../presentation/components/utils/ScrollToTop';

interface AppProps {
    initialTitle: string;
}

const App: React.FC<AppProps> = ({ initialTitle }) => {
    const { mode } = useThemeStore();
    useEffect(() => {
        document.title = initialTitle;
    }, [initialTitle]);

    const theme = useMemo(() => getTheme(mode), [mode]);

    return (
        <ThemeProvider theme={theme}>
            <CssBaseline />
            <BrowserRouter>
                <ScrollToTop />
                <AppRoutes />
            </BrowserRouter>
        </ThemeProvider>
    );
};

export default App;
