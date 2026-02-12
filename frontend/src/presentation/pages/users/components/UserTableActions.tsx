import React from 'react';
import { Box, IconButton, Tooltip } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import LockResetIcon from '@mui/icons-material/LockReset';
import type { User } from '../../../../domain/entities/User';

interface UserTableActionsProps {
    user: User;
    onEdit: (user: User) => void;
    onDelete: (user: User) => void;
    onResetPassword: (user: User) => void;
}

export const UserTableActions: React.FC<UserTableActionsProps> = ({ user, onEdit, onDelete, onResetPassword }) => {
    return (
        <Box sx={{ display: 'flex', justifyContent: 'flex-end', gap: 0.5 }}>
            <Tooltip title="Reset Password">
                <IconButton
                    onClick={() => onResetPassword(user)}
                    size="small"
                    sx={{
                        color: 'warning.main',
                        '&:hover': { bgcolor: 'warning.lighter' }
                    }}
                >
                    <LockResetIcon fontSize="small" />
                </IconButton>
            </Tooltip>
            <Tooltip title="Edit User">
                <IconButton
                    onClick={() => onEdit(user)}
                    size="small"
                    sx={{
                        color: 'primary.main',
                        '&:hover': { bgcolor: 'primary.lighter' }
                    }}
                >
                    <EditIcon fontSize="small" />
                </IconButton>
            </Tooltip>
            <Tooltip title="Delete User">
                <IconButton
                    onClick={() => onDelete(user)}
                    size="small"
                    sx={{
                        color: 'error.main',
                        '&:hover': { bgcolor: 'error.lighter' }
                    }}
                >
                    <DeleteIcon fontSize="small" />
                </IconButton>
            </Tooltip>
        </Box>
    );
};
