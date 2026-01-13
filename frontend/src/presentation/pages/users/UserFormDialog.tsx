import React, { useState, useEffect } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    FormControlLabel,
    Switch,
} from '@mui/material';
import type { User } from '../../../domain/entities/User';
import type { CreateUserRequest } from '../../../application/user/Create/CreateUserRequest';
import { useAuthStore } from '../../../state/authStore';

interface UserFormDialogProps {
    open: boolean;
    user?: User | null;
    onClose: () => void;
    onSave: (user: Omit<CreateUserRequest, 'id'> | Partial<User>) => void;
}

export const UserFormDialog: React.FC<UserFormDialogProps> = ({
    open,
    user,
    onClose,
    onSave,
}) => {
    const [formData, setFormData] = useState({
        name: '',
        email: '',
        password: '',
        role: 'USER' as 'ADMIN' | 'USER' | null,
        isActive: true,
    });

    useEffect(() => {
        if (user) {
            setFormData({
                name: user.name ?? '',
                email: user.email,
                password: '', // do not prefill password
                role: user.role,
                isActive: user.isActive,
            });
        } else {
            setFormData({
                name: '',
                email: '',
                password: '',
                role: 'USER',
                isActive: true,
            });
        }
    }, [user, open]);

    const { user: currentUser } = useAuthStore();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();

        if (user) {
            // editing existing user: send partial User (without password if empty)
            const update: Partial<User> = {
                name: formData.name,
                email: formData.email,
                role: formData.role as 'ADMIN' | 'USER' | null,
                isActive: formData.isActive,
            };
            // include password only if provided
            if (formData.password && formData.password.trim() !== '') {
                // caller may expect password handling separately; include under any shape
                // but since Partial<User> doesn't have password, we leave it out for update
            }
            onSave(update);
        } else {
            // creating new user: send CreateUserRequest (without id)
            const createPayload: Omit<CreateUserRequest, 'id'> = {
                name: formData.name,
                email: formData.email,
                password: formData.password,
                role: formData.role as 'ADMIN' | 'USER' | null,
                isActive: formData.isActive,
            };
            onSave(createPayload);
        }
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            <form onSubmit={handleSubmit}>
                <DialogTitle>{user ? 'Edit User' : 'Add User'}</DialogTitle>
                <DialogContent dividers>
                    {/* Username is derived automatically from email or name; input removed */}
                    <TextField
                        fullWidth
                        label="Full Name"
                        margin="normal"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    />
                    <TextField
                        fullWidth
                        label="Email"
                        type="email"
                        margin="normal"
                        value={formData.email}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        required
                    />
                    {!user && (
                        <TextField
                            fullWidth
                            label="Password"
                            type="password"
                            margin="normal"
                            value={formData.password}
                            onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                            required
                        />
                    )}
                    {currentUser?.role === 'ADMIN' && (
                        <FormControl fullWidth margin="normal">
                        <InputLabel>Role</InputLabel>
                        <Select
                            value={formData.role}
                            label="Role"
                            onChange={(e) => setFormData({ ...formData, role: e.target.value as 'ADMIN' | 'USER' })}
                        >
                            <MenuItem value="ADMIN">Admin</MenuItem>
                            <MenuItem value="USER">User</MenuItem>
                        </Select>
                        </FormControl>
                    )}
                    {currentUser?.role === 'ADMIN' && (
                        <FormControlLabel
                        control={
                            <Switch
                                checked={formData.isActive}
                                onChange={(e) => setFormData({ ...formData, isActive: e.target.checked })}
                            />
                        }
                        label="Active Status"
                        sx={{ mt: 2 }}
                    />
                    )}
                </DialogContent>
                <DialogActions>
                    <Button onClick={onClose} color="inherit">
                        Cancel
                    </Button>
                    <Button type="submit" variant="contained">
                        Save
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};
