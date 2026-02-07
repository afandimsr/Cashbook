import React from 'react';
import { Container, Box, Typography, Grid, Paper, Card, Stack } from '@mui/material';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';
import { useReports } from '../../../../application/hooks/useReports';
import { formatIDR } from '../../../utils/formatCurrency';

export const ReportsPage: React.FC = () => {
    const now = new Date();
    const currentMonth = now.getMonth() + 1;
    const currentYear = now.getFullYear();

    const { chartData, sortedSpending, isLoading } = useReports(currentMonth, currentYear);

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ mb: 4 }}>
                <Typography variant="h3" sx={{ fontWeight: 800, mb: 1 }}>Financial Reports</Typography>
                <Typography variant="h6" color="text.secondary">
                    Analysis for {now.toLocaleString('default', { month: 'long' })} {currentYear}
                </Typography>
            </Box>

            <Grid container spacing={4}>
                <Grid size={{ xs: 12, md: 6 }}>
                    <Paper sx={{ p: 3, borderRadius: 3, height: 500, display: 'flex', flexDirection: 'column' }}>
                        <Typography variant="h5" sx={{ fontWeight: 700, mb: 3 }}>Spending by Category</Typography>
                        <Box sx={{ flexGrow: 1, minHeight: 0 }}>
                            <ResponsiveContainer width="100%" height="100%">
                                <PieChart>
                                    <Pie
                                        data={chartData}
                                        cx="50%"
                                        cy="50%"
                                        innerRadius={80}
                                        outerRadius={120}
                                        paddingAngle={5}
                                        dataKey="value"
                                    >
                                        {chartData.map((entry, index) => (
                                            <Cell key={`cell-${index}`} fill={entry.color} />
                                        ))}
                                    </Pie>
                                    <Tooltip
                                        formatter={(value: number | undefined) => [formatIDR(value || 0), 'Spent']}
                                        contentStyle={{ borderRadius: 8, border: 'none', boxShadow: '0 4px 20px rgba(0,0,0,0.1)' }}
                                    />
                                    <Legend verticalAlign="bottom" height={36} />
                                </PieChart>
                            </ResponsiveContainer>
                        </Box>
                    </Paper>
                </Grid>

                <Grid size={{ xs: 12, md: 6 }}>
                    <Paper sx={{ p: 3, borderRadius: 3, height: 500, display: 'flex', flexDirection: 'column' }}>
                        <Typography variant="h5" sx={{ fontWeight: 700, mb: 3 }}>Expense Breakdown</Typography>

                        <Box sx={{
                            flexGrow: 1,
                            overflowY: 'auto',
                            mb: 2,
                            pr: 1,
                            '&::-webkit-scrollbar': { width: '8px' },
                            '&::-webkit-scrollbar-thumb': { backgroundColor: 'action.hover', borderRadius: '4px' }
                        }}>
                            <Stack spacing={2}>
                                {sortedSpending?.map(item => (
                                    <Card key={item.category_id} variant="outlined" sx={{ p: 2, borderRadius: 2 }}>
                                        <Stack direction="row" justifyContent="space-between" alignItems="center">
                                            <Stack direction="row" spacing={2} alignItems="center">
                                                <Box sx={{ width: 12, height: 12, borderRadius: '50%', bgcolor: item.color }} />
                                                <Typography sx={{ fontWeight: 600 }}>{item.category_name}</Typography>
                                            </Stack>
                                            <Typography sx={{ fontWeight: 700 }}>{formatIDR(item.total_amount)}</Typography>
                                        </Stack>
                                    </Card>
                                ))}
                                {sortedSpending?.length === 0 && !isLoading && (
                                    <Typography color="text.secondary" textAlign="center" py={4}>
                                        No spending data for this month.
                                    </Typography>
                                )}
                            </Stack>
                        </Box>

                        {sortedSpending?.length > 0 && (
                            <Card
                                variant="elevation"
                                elevation={0}
                                sx={{
                                    p: 2.5,
                                    borderRadius: 3,
                                    background: 'linear-gradient(135deg, #10b981 0%, #059669 100%)',
                                    color: 'white',
                                    boxShadow: '0 8px 16px -4px rgba(16, 185, 129, 0.4)',
                                    mt: 'auto'
                                }}
                            >
                                <Stack direction="row" justifyContent="space-between" alignItems="center">
                                    <Typography sx={{ fontWeight: 600, opacity: 0.9 }}>Total Spending</Typography>
                                    <Typography variant="h5" sx={{ fontWeight: 800 }}>
                                        {formatIDR(sortedSpending.reduce((acc, item) => acc + item.total_amount, 0))}
                                    </Typography>
                                </Stack>
                            </Card>
                        )}
                    </Paper>
                </Grid>
            </Grid>
        </Container>
    );
};
