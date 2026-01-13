import React, { useState } from 'react';
import {
    Box,
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
}

export const CategoryDialog: React.FC<CategoryDialogProps> = ({ open, onClose, onSave }) => {
    const [formData, setFormData] = useState<Omit<Category, 'id'>>({
        name: '',
        type: 'expense',
        color: '#3498db',
        icon: 'category'
    });

    const handleSubmit = async () => {
        if (!formData.name) return;
        await onSave(formData);
        setFormData({ name: '', type: 'expense', color: '#3498db', icon: 'category' });
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
            <DialogTitle sx={{ fontWeight: 800 }}>New Category</DialogTitle>
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
                    Create Category
                </Button>
            </DialogActions>
        </Dialog>
    );
};
