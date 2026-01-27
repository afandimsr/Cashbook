import React, { useState } from 'react';
import {
    Container,
    Box,
    Typography,
    Grid,
    Card,
    LinearProgress,
    Stack,
    IconButton,
    Tooltip,
    Snackbar,
    Alert
} from '@mui/material';
import { Edit as EditIcon } from '@mui/icons-material';
import { useBudgets } from '../../../../application/hooks/useBudgets';
import { formatIDR } from '../../../utils/formatCurrency';
import { BudgetDialog } from './components/BudgetDialog';
import type { Category } from '../../../../domain/entities/Category';
import type { Budget } from '../../../../domain/entities/Budget';

export const BudgetPage: React.FC = () => {
    const now = new Date();
    const currentMonth = now.getMonth() + 1;
    const currentYear = now.getFullYear();

    const {
        budgetAnalysis,
        handleUpdateBudget,
    } = useBudgets(currentMonth, currentYear);

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [selectedCategory, setSelectedCategory] = useState<Category | null>(null);
    const [selectedBudget, setSelectedBudget] = useState<Budget | null>(null);
    const [snack, setSnack] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' });

    const handleEditBudget = (category: Category, budget: Budget | null) => {
        setSelectedCategory(category);
        setSelectedBudget(budget);
        setIsDialogOpen(true);
    };

    const onSaveBudget = async (data: { category_id: number; amount: number; month: number; year: number }) => {
        try {
            await handleUpdateBudget(data);
            setSnack({ open: true, message: 'Budget updated successfully', severity: 'success' });
            setIsDialogOpen(false);
        } catch (error: any) {
            setSnack({ open: true, message: error?.message || 'Failed to update budget', severity: 'error' });
        }
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Box>
                    <Typography variant="h3" sx={{ fontWeight: 800, mb: 1 }}>Budgets</Typography>
                    <Typography variant="h6" color="text.secondary">
                        Plan your spending for {now.toLocaleString('default', { month: 'long' })} {currentYear}
                    </Typography>
                </Box>
            </Box>

            <Grid container spacing={3}>
                {(budgetAnalysis ?? []).map(({ category, spending, limit, percentage, isOver, budget }) => (
                    <Grid size={{ xs: 12, md: 6, lg: 4 }} key={category?.id ?? `${category?.name ?? 'cat'}-${Math.random()}`}>
                        <Card sx={{ p: 3, borderRadius: 3, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                            <Stack direction="row" justifyContent="space-between" alignItems="flex-start" mb={2}>
                                <Box>
                                    <Typography variant="h6" sx={{ fontWeight: 700 }}>{category?.name ?? '—'}</Typography>
                                    <Typography variant="body2" color="text.secondary">
                                        Budget: <strong>{formatIDR(limit)}</strong>
                                    </Typography>
                                </Box>
                                <Tooltip title="Set Budget">
                                    <IconButton size="small" onClick={() => category && handleEditBudget(category, budget || null)}>
                                        <EditIcon fontSize="small" />
                                    </IconButton>
                                </Tooltip>
                            </Stack>

                            <Box sx={{ mb: 1 }}>
                                <Stack direction="row" justifyContent="space-between" mb={1}>
                                    <Typography variant="body2" color="text.secondary">
                                        Spent: {formatIDR(spending)}
                                    </Typography>
                                    <Typography variant="body2" sx={{ fontWeight: 700, color: isOver ? 'error.main' : 'text.primary' }}>
                                        {percentage.toFixed(0)}%
                                    </Typography>
                                </Stack>
                                <LinearProgress
                                    variant="determinate"
                                    value={Math.min(percentage, 100)}
                                    color={isOver ? 'error' : (percentage > 80 ? 'warning' : 'primary')}
                                    sx={{ height: 8, borderRadius: 4 }}
                                />
                            </Box>
                        </Card>
                    </Grid>
                ))}
            </Grid>
            <BudgetDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSave={onSaveBudget}
                category={selectedCategory}
                budget={selectedBudget}
                month={currentMonth}
                year={currentYear}
            />

            <Snackbar open={snack.open} autoHideDuration={4000} onClose={() => setSnack({ ...snack, open: false })}>
                <Alert onClose={() => setSnack({ ...snack, open: false })} severity={snack.severity} sx={{ width: '100%' }}>
                    {snack.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};
