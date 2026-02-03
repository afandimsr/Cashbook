import React, { useEffect } from 'react';
import { Box, Typography, Container, Grid, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Stack } from '@mui/material';

// Categorized Showcase Components
import { StatCards } from './components/StatCards';
import { formatIDR } from '../../../utils/formatCurrency';
import { useTransactionStore } from '../../../../state/transactionStore';
import { useCategoryStore } from '../../../../state/categoryStore';

export const OverviewPage: React.FC = () => {
    const { summary, transactions, isLoading, fetchSummary, fetchTransactions } = useTransactionStore();
    const { categories, fetchCategories } = useCategoryStore();

    useEffect(() => {
        fetchSummary();
        fetchTransactions(1, 5); // Just 5 for recent activity
        fetchCategories();
    }, [fetchSummary, fetchTransactions, fetchCategories]);

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ mb: { xs: 3, md: 4 } }}>
                <Typography
                    variant="h3"
                    sx={{
                        fontWeight: 800,
                        color: 'text.primary',
                        mb: 1,
                        fontSize: { xs: '2rem', md: '3rem' }
                    }}
                >
                    Dashboard
                </Typography>
                <Typography
                    variant="h6"
                    color="text.secondary"
                    sx={{
                        fontWeight: 400,
                        fontSize: { xs: '1rem', md: '1.25rem' }
                    }}
                >
                    Welcome back! Here's your financial overview.
                </Typography>
            </Box>

            {/* 1. Statistics Section */}
            <Box sx={{ mb: 4 }}>
                <StatCards summary={summary} isLoading={isLoading} />
            </Box>

            <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 7 }}>
                    <Typography variant="h5" sx={{ fontWeight: 700, mb: 2 }}>
                        Recent Transactions
                    </Typography>
                    <TableContainer
                        component={Paper}
                        sx={{
                            borderRadius: 3,
                            boxShadow: '0 4px 20px rgba(0,0,0,0.05)',
                            overflowX: 'auto'
                        }}
                    >
                        <Table>
                            <TableHead sx={{ backgroundColor: 'action.hover' }}>
                                <TableRow>
                                    <TableCell sx={{ fontWeight: 700 }}>Date</TableCell>
                                    <TableCell sx={{ fontWeight: 700 }}>Category</TableCell>
                                    <TableCell sx={{ fontWeight: 700 }} align="right">Amount</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                {transactions?.slice(0, 5).map((tx) => {
                                    const category = categories?.find(c => c.id === tx.category_id);
                                    return (
                                        <TableRow key={tx.id} hover>
                                            <TableCell>{new Date(tx.date).toLocaleDateString()}</TableCell>
                                            <TableCell>
                                                <Stack direction="row" spacing={1} alignItems="center">
                                                    <Box
                                                        sx={{
                                                            width: 8,
                                                            height: 8,
                                                            borderRadius: '50%',
                                                            backgroundColor: category?.color || '#ccc'
                                                        }}
                                                    />
                                                    <Typography variant="body2">{category?.name || 'Unknown'}</Typography>
                                                </Stack>
                                            </TableCell>
                                            <TableCell align="right">
                                                <Typography
                                                    variant="body2"
                                                    sx={{
                                                        fontWeight: 700,
                                                        color: tx.type === 'income' ? 'success.main' : 'error.main'
                                                    }}
                                                >
                                                    {tx.type === 'income' ? '+' : '-'}{formatIDR(tx.amount)}
                                                </Typography>
                                            </TableCell>
                                        </TableRow>
                                    );
                                })}
                                {transactions?.length === 0 && (
                                    <TableRow>
                                        <TableCell colSpan={3} align="center" sx={{ py: 4 }}>
                                            <Typography color="text.secondary">No recent transactions</Typography>
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </Grid>
                <Grid size={{ xs: 12, md: 5 }}>
                    <Typography variant="h5" sx={{ fontWeight: 700, mb: 2 }}>
                        Financial Health
                    </Typography>
                    <Paper sx={{ p: 3, borderRadius: 3, textAlign: 'center' }}>
                        <Typography color="text.secondary" gutterBottom>
                            Spending Tip
                        </Typography>
                        <Typography variant="body1">
                            You've spent more than last month. Consider reviewing your "Food" category!
                        </Typography>
                    </Paper>
                </Grid>
            </Grid>
        </Container>
    );
};
