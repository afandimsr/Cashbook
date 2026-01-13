import React, { useState, useEffect } from 'react';
import {
    Box,
    Button,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    MenuItem
} from '@mui/material';
import type { Transaction } from '../../../../../domain/entities/Transaction';
import type { Category } from '../../../../../domain/entities/Category';

interface TransactionDialogProps {
    open: boolean;
    onClose: () => void;
    onSave: (data: Omit<Transaction, 'id'>) => Promise<void>;
    categories?: Category[];
}

export const TransactionDialog: React.FC<TransactionDialogProps> = ({ open, onClose, onSave, categories }) => {
    const [formData, setFormData] = useState<Omit<Transaction, 'id'>>({
        category_id: 0,
        amount: 0,
        note: '',
        date: new Date().toISOString().split('T')[0],
        type: 'expense',
    });

    useEffect(() => {
        if (!open) {
            setFormData({
                category_id: 0,
                amount: 0,
                note: '',
                date: new Date().toISOString().split('T')[0],
                type: 'expense',
            });
        }
    }, [open]);

    const handleSubmit = async () => {
        await onSave({ ...formData, amount: Number(formData.amount) });
        onClose();
    };

    return (
        <Dialog open={open} onClose={onClose} maxWidth="xs" fullWidth>
            <DialogTitle>New Transaction</DialogTitle>
            <DialogContent>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
                    <TextField
                        label="Amount"
                        type="number"
                        fullWidth
                        value={formData.amount}
                        onChange={(e) => setFormData({ ...formData, amount: Number(e.target.value) })}
                    />
                    <TextField
                        select
                        label="Category"
                        fullWidth
                        value={formData.category_id || ''}
                        onChange={(e) => {
                            const catId = Number(e.target.value);
                            const cat = (categories ?? []).find(c => c.id === catId);
                            setFormData({
                                ...formData,
                                category_id: catId,
                                type: cat?.type || 'expense'
                            });
                        }}
                    >
                        {(categories ?? []).map((cat) => (
                            <MenuItem key={cat.id} value={cat.id}>
                                {cat.name} ({cat.type})
                            </MenuItem>
                        ))}
                    </TextField>
                    <TextField
                        label="Date"
                        type="date"
                        fullWidth
                        InputLabelProps={{ shrink: true }}
                        value={formData.date}
                        onChange={(e) => setFormData({ ...formData, date: e.target.value })}
                    />
                    <TextField
                        label="Note"
                        fullWidth
                        multiline
                        rows={2}
                        value={formData.note}
                        onChange={(e) => setFormData({ ...formData, note: e.target.value })}
                    />
                </Box>
            </DialogContent>
            <DialogActions sx={{ p: 2, pt: 0 }}>
                <Button onClick={onClose}>Cancel</Button>
                <Button variant="contained" onClick={handleSubmit}>Save</Button>
            </DialogActions>
        </Dialog>
    );
};
