import React, { useState } from 'react';
import {
    Box,
    Typography,
    Container,
    Button,
    CircularProgress,
    Stack,
    Snackbar,
    Alert,
} from '@mui/material';
import {
    Add as AddIcon,
    FilterList as FilterIcon,
} from '@mui/icons-material';
import { useTransactions } from '../../../../application/hooks/useTransactions';
import { useCategories } from '../../../../application/hooks/useCategories';
import { TransactionList } from './components/TransactionList';
import { TransactionDialog } from './components/TransactionDialog';

export const TransactionPage: React.FC = () => {
    const {
        transactions,
        isLoading,
        handleAddTransaction,
        handleDeleteTransaction
    } = useTransactions();

    const { categories } = useCategories();
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [snack, setSnack] = useState<{ open: boolean; message: string; severity: 'success' | 'error' | 'info' | 'warning' }>({ open: false, message: '', severity: 'success' });

    const handleSaveTransaction = async (data: any) => {
        try {
            await handleAddTransaction(data);
            setSnack({ open: true, message: 'Transaction created', severity: 'success' });
            setIsDialogOpen(false);
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to create transaction', severity: 'error' });
            throw err;
        }
    };

    const handleDelete = async (id: number) => {
        try {
            // useTransactions already asks for confirmation, but ensure snackbar
            await handleDeleteTransaction(id);
            setSnack({ open: true, message: 'Transaction deleted', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to delete transaction', severity: 'error' });
            throw err;
        }
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Box>
                    <Typography variant="h4" sx={{ fontWeight: 800, color: 'text.primary' }}>
                        Transactions
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Keep track of your spending and earnings
                    </Typography>
                </Box>
                <Stack direction="row" spacing={2}>
                    <Button variant="outlined" startIcon={<FilterIcon />} sx={{ borderRadius: 2 }}>
                        Filters
                    </Button>
                    <Button
                        variant="contained"
                        startIcon={<AddIcon />}
                        onClick={() => setIsDialogOpen(true)}
                        sx={{ borderRadius: 2, px: 3 }}
                    >
                        Add Transaction
                    </Button>
                </Stack>
            </Box>

            {isLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
                    <CircularProgress />
                </Box>
            ) : (
                <TransactionList
                    transactions={transactions ?? []}
                    categories={categories ?? []}
                    onDelete={handleDelete}
                />
            )}

            <TransactionDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSave={handleSaveTransaction}
                categories={categories ?? []}
            />

            <Snackbar open={snack.open} autoHideDuration={4000} onClose={() => setSnack({ ...snack, open: false })}>
                <Alert onClose={() => setSnack({ ...snack, open: false })} severity={snack.severity} sx={{ width: '100%' }}>
                    {snack.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};
