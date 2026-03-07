import { Component, type ErrorInfo, type ReactNode } from 'react';
import { Box, Typography, Button, Container, Paper } from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';

interface Props {
    children?: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        // Update state so the next render will show the fallback UI.
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Uncaught error:', error, errorInfo);
    }

    private handleReset = () => {
        this.setState({ hasError: false, error: null });
        window.location.href = '/dashboard';
    };

    private handleGoHome = () => {
        window.location.href = '/login';
    };

    public render() {
        if (this.state.hasError) {
            return (
                <Container maxWidth="sm" sx={{ 
                    height: '100vh', 
                    display: 'flex', 
                    alignItems: 'center', 
                    justifyContent: 'center' 
                }}>
                    <Paper elevation={3} sx={{ 
                        p: 5, 
                        textAlign: 'center', 
                        borderRadius: 4,
                        bgcolor: 'background.paper'
                    }}>
                        <ErrorOutlineIcon sx={{ fontSize: 80, color: 'error.main', mb: 2 }} />
                        <Typography variant="h4" gutterBottom sx={{ fontWeight: 700 }}>
                            Oops! Something went wrong
                        </Typography>
                        <Typography variant="body1" color="text.secondary" sx={{ mb: 4 }}>
                            {this.state.error?.message || 'An unexpected error occurred.'}
                        </Typography>
                        <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
                            <Button 
                                variant="contained" 
                                color="primary" 
                                onClick={this.handleReset}
                                sx={{ borderRadius: 2, px: 3 }}
                            >
                                Try Again
                            </Button>
                            <Button 
                                variant="outlined" 
                                onClick={this.handleGoHome}
                                sx={{ borderRadius: 2, px: 3 }}
                            >
                                Go to Login
                            </Button>
                        </Box>
                    </Paper>
                </Container>
            );
        }

        return this.props.children;
    }
}
