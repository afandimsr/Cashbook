import React from 'react';
import {
    Container, Box, Typography, Grid, Card, LinearProgress,
    Stack
} from '@mui/material';
import { useBudgets } from '../../../../application/hooks/useBudgets';
import { formatIDR } from '../../../utils/formatCurrency';

export const BudgetPage: React.FC = () => {
    const now = new Date();
    const currentMonth = now.getMonth() + 1;
    const currentYear = now.getFullYear();

    const {
        budgetAnalysis,
    } = useBudgets(currentMonth, currentYear);

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
                {(budgetAnalysis ?? []).map(({ category, spending, limit, percentage, isOver }) => (
                    <Grid size={{ xs: 12, md: 6, lg: 4 }} key={category?.id ?? `${category?.name ?? 'cat'}-${Math.random()}`}>
                        <Card sx={{ p: 3, borderRadius: 3, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                            <Stack direction="row" justifyContent="space-between" alignItems="center" mb={2}>
                                <Typography variant="h6" sx={{ fontWeight: 700 }}>{category?.name ?? '—'}</Typography>
                                    <Typography variant="body2" color="text.secondary">
                                    Budget: <strong>{formatIDR(limit)}</strong>
                                </Typography>
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
        </Container>
    );
};
