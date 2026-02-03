import React, { useState } from 'react';
import {
    Container, Box, Typography, Button, Stack, Grid, Snackbar, Alert
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import AutoFixHighIcon from '@mui/icons-material/AutoFixHigh';
import { useRecurring } from '../../../../application/hooks/useRecurring';
import { useCategories } from '../../../../application/hooks/useCategories';
import { RecurringCard } from './components/RecurringCard';
import { RecurringDialog } from './components/RecurringDialog';
import { ConfirmDialog } from '../../../components/common/ConfirmDialog';

export const RecurringPage: React.FC = () => {
    const {
        recurringTransactions,
        handleCreateRecurring,
        handleDeleteRecurring,
        handleProcessDue
    } = useRecurring();

    const { categories, getCategoryById } = useCategories();
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [itemToDelete, setItemToDelete] = useState<number | null>(null);
    const [snack, setSnack] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' });

    const handleSaveRecurring = async (data: Partial<any>) => {
        try {
            await handleCreateRecurring(data);
            setSnack({ open: true, message: 'Recurring created', severity: 'success' });
            setIsDialogOpen(false);
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to create recurring', severity: 'error' });
            throw err;
        }
    };

    const handleDeleteClick = (id: number) => {
        setItemToDelete(id);
        setConfirmOpen(true);
    };

    const handleConfirmDelete = async () => {
        if (itemToDelete === null) return;
        try {
            await handleDeleteRecurring(itemToDelete);
            setSnack({ open: true, message: 'Recurring deleted', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to delete recurring', severity: 'error' });
        } finally {
            setConfirmOpen(false);
            setItemToDelete(null);
        }
    };

    const handleProcess = async () => {
        try {
            await handleProcessDue();
            setSnack({ open: true, message: 'Processed due recurring items', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to process due items', severity: 'error' });
            throw err;
        }
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{
                display: 'flex',
                flexDirection: { xs: 'column', sm: 'row' },
                justifyContent: 'space-between',
                alignItems: { xs: 'flex-start', sm: 'center' },
                mb: 4,
                gap: 2
            }}>
                <Box>
                    <Typography
                        variant="h3"
                        sx={{
                            fontWeight: 800,
                            mb: 1,
                            fontSize: { xs: '2rem', md: '3rem' }
                        }}
                    >
                        Recurring
                    </Typography>
                    <Typography variant="h6" color="text.secondary" sx={{ fontSize: { xs: '1rem', md: '1.25rem' } }}>
                        Automate your regular expenses and income
                    </Typography>
                </Box>
                <Stack direction="row" spacing={2} sx={{ width: { xs: '100%', sm: 'auto' } }}>
                    <Button
                        variant="outlined"
                        startIcon={<AutoFixHighIcon />}
                        onClick={handleProcess}
                        sx={{ borderRadius: 2, flex: { xs: 1, sm: 'auto' } }}
                    >
                        Process Due
                    </Button>
                    <Button
                        variant="contained"
                        startIcon={<AddIcon />}
                        onClick={() => setIsDialogOpen(true)}
                        sx={{ borderRadius: 2, flex: { xs: 1, sm: 'auto' } }}
                    >
                        New
                    </Button>
                </Stack>
            </Box>

            <Grid container spacing={3}>
                {(recurringTransactions ?? []).map((rt) => (
                    rt && (
                        <Grid key={rt.id ?? Math.random()} size={{ xs: 12, sm: 6, lg: 4 }}>
                            <RecurringCard
                                recurring={rt}
                                category={getCategoryById ? getCategoryById(rt.category_id) : undefined}
                                onDelete={handleDeleteClick}
                            />
                        </Grid>
                    )
                ))}
            </Grid>

            <RecurringDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSave={handleSaveRecurring}
                categories={categories}
            />

            <ConfirmDialog
                open={confirmOpen}
                title="Delete Recurring"
                message="Are you sure you want to delete this recurring item? This action cannot be undone."
                onConfirm={handleConfirmDelete}
                onCancel={() => setConfirmOpen(false)}
            />

            <Snackbar open={snack.open} autoHideDuration={4000} onClose={() => setSnack({ ...snack, open: false })}>
                <Alert onClose={() => setSnack({ ...snack, open: false })} severity={snack.severity} sx={{ width: '100%' }}>
                    {snack.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};
