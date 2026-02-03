import React from 'react';
import { Box, Typography, Button, TextField, InputAdornment, Stack } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import SearchIcon from '@mui/icons-material/Search';

interface UserTableHeaderProps {
    onAddClick: () => void;
    searchQuery: string;
    onSearchChange: (value: string) => void;
    canCreate?: boolean;
}
export const UserTableHeader: React.FC<UserTableHeaderProps> = ({ onAddClick, searchQuery, onSearchChange, canCreate = false }) => {
    return (
        <Box sx={{ mb: 4 }}>
            <Stack
                direction={{ xs: 'column', sm: 'row' }}
                justifyContent="space-between"
                alignItems={{ xs: 'stretch', sm: 'center' }}
                spacing={2}
            >
                <Box>
                    <Typography
                        variant="h4"
                        sx={{
                            fontWeight: 700,
                            color: 'text.primary',
                            mb: 0.5,
                            fontSize: { xs: '1.75rem', md: '2.125rem' }
                        }}
                    >
                        User Management
                    </Typography>
                    <Typography variant="body2" color="text.secondary" sx={{ fontSize: { xs: '0.875rem', md: '1rem' } }}>
                        Manage your team members and their account permissions here.
                    </Typography>
                </Box>
                {canCreate && (
                    <Button
                        variant="contained"
                        startIcon={<AddIcon />}
                        onClick={onAddClick}
                        sx={{
                            borderRadius: 2,
                            px: 3,
                            py: 1,
                            textTransform: 'none',
                            fontWeight: 600,
                            boxShadow: '0 4px 12px rgba(25, 118, 210, 0.2)'
                        }}
                    >
                        Add User
                    </Button>
                )}
            </Stack>

            <Box sx={{ mt: 3, display: 'flex', gap: 2 }}>
                <TextField
                    placeholder="Search users..."
                    size="small"
                    value={searchQuery}
                    onChange={(e) => onSearchChange(e.target.value)}
                    sx={{
                        maxWidth: { xs: '100%', sm: 400 },
                        width: '100%',
                        '& .MuiOutlinedInput-root': {
                            borderRadius: 2,
                            bgcolor: 'background.paper'
                        }
                    }}
                    InputProps={{
                        startAdornment: (
                            <InputAdornment position="start">
                                <SearchIcon fontSize="small" color="action" />
                            </InputAdornment>
                        ),
                    }}
                />
            </Box>
        </Box>
    );
};
