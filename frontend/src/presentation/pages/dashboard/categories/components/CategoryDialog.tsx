import React, { useState, useEffect } from 'react';
import {
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    MenuItem,
    Stack
} from '@mui/material';
import type { Category } from '../../../../../domain/entities/Category';

interface CategoryDialogProps {
    open: boolean;
    onClose: () => void;
    onSave: (data: Omit<Category, 'id'>) => Promise<void>;
    initialData?: Category | null;
}

export const CategoryDialog: React.FC<CategoryDialogProps> = ({ open, onClose, onSave, initialData }) => {
    const [formData, setFormData] = useState<Omit<Category, 'id' | 'user_id'>>({
        name: '',
        type: 'expense',
        color: '#3498db',
        icon: 'category'
    });

    useEffect(() => {
        if (initialData) {
            setFormData({
                name: initialData.name,
                type: initialData.type,
                color: initialData.color,
                icon: initialData.icon
            });
        } else {
            setFormData({
                name: '',
                type: 'expense',
                color: '#3498db',
                icon: 'category'
            });
        }
    }, [initialData, open]);

    const handleSubmit = async () => {
        if (!formData.name) return;
        await onSave(formData);
        if (!initialData) {
            setFormData({ name: '', type: 'expense', color: '#3498db', icon: 'category' });
        }
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
            <DialogTitle sx={{ fontWeight: 800 }}>
                {initialData ? 'Edit Category' : 'New Category'}
            </DialogTitle>
            <DialogContent>
                <Stack spacing={3} sx={{ pt: 1 }}>
                    <TextField
                        label="Category Name"
                        fullWidth
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="e.g. Groceries, Salary"
                    />
                    <TextField
                        select
                        label="Type"
                        fullWidth
                        value={formData.type}
                        onChange={(e) => setFormData({ ...formData, type: e.target.value as 'income' | 'expense' })}
                    >
                        <MenuItem value="expense">Expense</MenuItem>
                        <MenuItem value="income">Income</MenuItem>
                    </TextField>
                    <TextField
                        label="Icon"
                        fullWidth
                        value={formData.icon}
                        onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
                        placeholder="e.g. shopping_cart, account_balance"
                        helperText="Material Icon name"
                    />
                    <TextField
                        label="Color"
                        type="color"
                        fullWidth
                        value={formData.color}
                        onChange={(e) => setFormData({ ...formData, color: e.target.value })}
                        helperText="Pick a color for the dashboard charts"
                    />
                </Stack>
            </DialogContent>
            <DialogActions sx={{ p: 3 }}>
                <Button onClick={onClose} color="inherit">Cancel</Button>
                <Button variant="contained" onClick={handleSubmit} disabled={!formData.name}>
                    {initialData ? 'Update Category' : 'Create Category'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};
