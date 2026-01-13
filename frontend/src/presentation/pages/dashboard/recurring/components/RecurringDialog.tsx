import React, { useState, useEffect } from 'react';
import {
    Stack,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    MenuItem
} from '@mui/material';
import type { RecurringTransaction, Frequency } from '../../../../../domain/entities/RecurringTransaction';
import type { Category } from '../../../../../domain/entities/Category';

interface RecurringDialogProps {
    open: boolean;
    onClose: () => void;
    onSave: (data: Partial<RecurringTransaction>) => Promise<void>;
    categories: Category[];
}

interface RecurringFormData {
    category_id: string;
    amount: string;
    type: 'income' | 'expense';
    frequency: Frequency;
    note: string;
    start_date: string;
}

export const RecurringDialog: React.FC<RecurringDialogProps> = ({ open, onClose, onSave, categories }) => {
    const initialState: RecurringFormData = {
        category_id: '',
        amount: '',
        type: 'expense',
        frequency: 'monthly',
        note: '',
        start_date: new Date().toISOString().split('T')[0]
    };

    const [formData, setFormData] = useState<RecurringFormData>(initialState);

    useEffect(() => {
        if (!open) {
            setFormData(initialState);
        }
    }, [open]);

    const handleSave = async () => {
        const categoryId = parseInt(formData.category_id);
        const amount = parseFloat(formData.amount);

        if (isNaN(categoryId) || isNaN(amount)) {
            return;
        }

        await onSave({
            ...formData,
            category_id: categoryId,
            amount: amount,
            start_date: new Date(formData.start_date).toISOString()
        } as Partial<RecurringTransaction>);
        onClose();
    };

    const isFormValid = formData.category_id && formData.amount && !isNaN(parseFloat(formData.amount)) && formData.start_date;

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="xs">
            <DialogTitle sx={{ fontWeight: 800 }}>New Recurring Transaction</DialogTitle>
            <DialogContent>
                <Stack spacing={3} sx={{ pt: 1 }}>
                    <TextField
                        select
                        fullWidth
                        label="Type"
                        value={formData.type}
                        onChange={(e) => setFormData({ ...formData, type: e.target.value as 'income' | 'expense' })}
                    >
                        <MenuItem value="expense">Expense</MenuItem>
                        <MenuItem value="income">Income</MenuItem>
                    </TextField>
                    <TextField
                        select
                        fullWidth
                        label="Category"
                        value={formData.category_id}
                        onChange={(e) => setFormData({ ...formData, category_id: e.target.value })}
                    >
                        {categories?.filter(c => c.type === formData.type).map(cat => (
                            <MenuItem key={cat.id} value={(cat.id ?? '').toString()}>{cat.name}</MenuItem>
                        ))}
                    </TextField>
                    <TextField
                        select
                        fullWidth
                        label="Frequency"
                        value={formData.frequency}
                        onChange={(e) => setFormData({ ...formData, frequency: e.target.value as Frequency })}
                    >
                        <MenuItem value="daily">Daily</MenuItem>
                        <MenuItem value="weekly">Weekly</MenuItem>
                        <MenuItem value="monthly">Monthly</MenuItem>
                    </TextField>
                    <TextField
                        fullWidth
                        label="Start Date"
                        type="date"
                        value={formData.start_date}
                        onChange={(e) => setFormData({ ...formData, start_date: e.target.value })}
                        slotProps={{ inputLabel: { shrink: true } }}
                    />
                    <TextField
                        fullWidth
                        label="Amount"
                        type="number"
                        value={formData.amount}
                        onChange={(e) => setFormData({ ...formData, amount: e.target.value })}
                    />
                    <TextField
                        fullWidth
                        label="Note"
                        value={formData.note}
                        onChange={(e) => setFormData({ ...formData, note: e.target.value })}
                    />
                </Stack>
            </DialogContent>
            <DialogActions sx={{ p: 3 }}>
                <Button onClick={onClose} color="inherit">Cancel</Button>
                <Button
                    onClick={handleSave}
                    variant="contained"
                    disabled={!isFormValid}
                >
                    Create
                </Button>
            </DialogActions>
        </Dialog>
    );
};
