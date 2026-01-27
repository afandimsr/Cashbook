import React, { useState, useEffect } from 'react';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    TextField,
    Typography,
    Box
} from '@mui/material';
import type { Category } from '../../../../../domain/entities/Category';
import type { Budget } from '../../../../../domain/entities/Budget';

interface BudgetDialogProps {
    open: boolean;
    onClose: () => void;
    onSave: (data: { category_id: number; amount: number; month: number; year: number }) => Promise<void>;
    category: Category | null;
    budget: Budget | null;
    month: number;
    year: number;
}

export const BudgetDialog: React.FC<BudgetDialogProps> = ({
    open,
    onClose,
    onSave,
    category,
    budget,
    month,
    year
}) => {
    const [amount, setAmount] = useState<string>('0');
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        if (budget) {
            setAmount(budget.amount.toString());
        } else {
            setAmount('0');
        }
    }, [budget, open]);

    const handleSave = async () => {
        if (!category) return;

        setIsLoading(true);
        try {
            await onSave({
                category_id: category.id!,
                amount: parseFloat(amount),
                month,
                year
            });
            onClose();
        } catch (error) {
            console.error('Failed to save budget', error);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <Dialog open={open} onClose={onClose} fullWidth maxWidth="xs">
            <DialogTitle sx={{ fontWeight: 700 }}>
                Set Budget for {category?.name}
            </DialogTitle>
            <DialogContent>
                <Box sx={{ mt: 1 }}>
                    <Typography variant="body2" color="text.secondary" gutterBottom>
                        Plan your spending for {category?.name}.
                    </Typography>
                    <TextField
                        fullWidth
                        label="Budget Amount (IDR)"
                        type="number"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                        sx={{ mt: 2 }}
                        autoFocus
                    />
                </Box>
            </DialogContent>
            <DialogActions sx={{ p: 3 }}>
                <Button onClick={onClose} color="inherit">Cancel</Button>
                <Button
                    onClick={handleSave}
                    variant="contained"
                    disabled={isLoading || !amount || parseFloat(amount) < 0}
                >
                    {isLoading ? 'Saving...' : 'Save Budget'}
                </Button>
            </DialogActions>
        </Dialog>
    );
};
