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
    TextField,
    MenuItem,
    Grid,
    InputAdornment,
    Collapse,
    Chip,
    Pagination,
} from '@mui/material';
import {
    Add as AddIcon,
    FilterList as FilterIcon,
    Search as SearchIcon,
    Clear as ClearIcon,
    CalendarToday as DateIcon,
} from '@mui/icons-material';
import { useTransactions } from '../../../../application/hooks/useTransactions';
import { useCategories } from '../../../../application/hooks/useCategories';
import { TransactionList } from './components/TransactionList';
import { TransactionDialog } from './components/TransactionDialog';
import { ConfirmDialog } from '../../../components/common/ConfirmDialog';

export const TransactionPage: React.FC = () => {
    const {
        transactions,
        total,
        totalAmount,
        isLoading,
        page,
        limit,
        setPage,
        handleAddTransaction,
        handleDeleteTransaction,
        filters,
        setFilters,
        clearFilters
    } = useTransactions();

    const { categories } = useCategories();
    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [itemToDelete, setItemToDelete] = useState<number | null>(null);
    const [showFilters, setShowFilters] = useState(false);
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

    const handleDeleteClick = (id: number) => {
        setItemToDelete(id);
        setConfirmOpen(true);
    };

    const handleConfirmDelete = async () => {
        if (itemToDelete === null) return;
        try {
            await handleDeleteTransaction(itemToDelete);
            setSnack({ open: true, message: 'Transaction deleted', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to delete transaction', severity: 'error' });
        } finally {
            setConfirmOpen(false);
            setItemToDelete(null);
        }
    };

    const hasActiveFilters = Object.values(filters).some(v => v !== '' && v !== undefined);

    return (
        <Container maxWidth="xl" sx={{ py: { xs: 2, md: 4 } }}>
            <Box sx={{
                display: 'flex',
                flexDirection: { xs: 'column', sm: 'row' },
                justifyContent: 'space-between',
                alignItems: { xs: 'flex-start', sm: 'center' },
                mb: 4,
                gap: 2
            }}>
                <Box>
                    <Typography variant="h4" sx={{ fontWeight: 800, color: 'text.primary', fontSize: { xs: '1.75rem', md: '2.125rem' } }}>
                        Transactions
                    </Typography>
                    <Typography variant="body1" color="text.secondary">
                        Keep track of your spending and earnings
                    </Typography>
                </Box>
                <Stack direction="row" spacing={2} sx={{ width: { xs: '100%', sm: 'auto' } }}>
                    <Button
                        variant={showFilters ? "contained" : "outlined"}
                        color={hasActiveFilters ? "primary" : "inherit"}
                        startIcon={<FilterIcon />}
                        onClick={() => setShowFilters(!showFilters)}
                        sx={{ borderRadius: 2, flex: { xs: 1, sm: 'auto' } }}
                    >
                        Filters {hasActiveFilters && `(Active)`}
                    </Button>
                    <Button
                        variant="contained"
                        startIcon={<AddIcon />}
                        onClick={() => setIsDialogOpen(true)}
                        sx={{ borderRadius: 2, px: 3, flex: { xs: 1, sm: 'auto' } }}
                    >
                        Add
                    </Button>
                </Stack>
            </Box>

            <Collapse in={showFilters}>
                <Box sx={{
                    p: 3,
                    mb: 4,
                    bgcolor: 'background.paper',
                    borderRadius: 3,
                    boxShadow: '0 4px 20px rgba(0,0,0,0.05)',
                    border: '1px solid',
                    borderColor: 'divider'
                }}>
                    <Grid container spacing={2} alignItems="center">
                        <Grid size={{ xs: 12, md: 4 }}>
                            <TextField
                                fullWidth
                                label="Search"
                                size="small"
                                value={filters.q}
                                onChange={(e) => setFilters({ q: e.target.value })}
                                InputProps={{
                                    startAdornment: (
                                        <InputAdornment position="start">
                                            <SearchIcon fontSize="small" />
                                        </InputAdornment>
                                    ),
                                }}
                            />
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6, md: 2 }}>
                            <TextField
                                select
                                fullWidth
                                label="Type"
                                size="small"
                                value={filters.type}
                                onChange={(e) => setFilters({ type: e.target.value })}
                            >
                                <MenuItem value="">All Types</MenuItem>
                                <MenuItem value="income">Income</MenuItem>
                                <MenuItem value="expense">Expense</MenuItem>
                            </TextField>
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6, md: 2 }}>
                            <TextField
                                select
                                fullWidth
                                label="Category"
                                size="small"
                                value={filters.category_id}
                                onChange={(e) => setFilters({ category_id: e.target.value })}
                            >
                                <MenuItem value="">All Categories</MenuItem>
                                {categories?.map((cat) => (
                                    <MenuItem key={cat.id} value={cat.id}>
                                        {cat.name}
                                    </MenuItem>
                                ))}
                            </TextField>
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6, md: 2 }}>
                            <TextField
                                fullWidth
                                label="From"
                                type="date"
                                size="small"
                                value={filters.start_date}
                                onChange={(e) => setFilters({ start_date: e.target.value })}
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                        <Grid size={{ xs: 12, sm: 6, md: 2 }}>
                            <TextField
                                fullWidth
                                label="To"
                                type="date"
                                size="small"
                                value={filters.end_date}
                                onChange={(e) => setFilters({ end_date: e.target.value })}
                                InputLabelProps={{ shrink: true }}
                            />
                        </Grid>
                    </Grid>
                    {hasActiveFilters && (
                        <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end' }}>
                            <Button
                                size="small"
                                startIcon={<ClearIcon />}
                                onClick={clearFilters}
                                color="inherit"
                            >
                                Clear All Filters
                            </Button>
                        </Box>
                    )}
                </Box>
            </Collapse>

            {hasActiveFilters && !showFilters && (
                <Stack direction="row" spacing={1} sx={{ mb: 2, flexWrap: 'wrap', gap: 1 }}>
                    {filters.q && <Chip label={`Search: ${filters.q}`} size="small" onDelete={() => setFilters({ q: '' })} />}
                    {filters.type && <Chip label={`Type: ${filters.type}`} size="small" onDelete={() => setFilters({ type: '' })} />}
                    {filters.category_id && (
                        <Chip
                            label={`Category: ${categories?.find(c => c.id?.toString() === filters.category_id?.toString())?.name || '...'}`}
                            size="small"
                            onDelete={() => setFilters({ category_id: '' })}
                        />
                    )}
                    {(filters.start_date || filters.end_date) && (
                        <Chip
                            icon={<DateIcon />}
                            label={`${filters.start_date || '...'} to ${filters.end_date || '...'}`}
                            size="small"
                            onDelete={() => setFilters({ start_date: '', end_date: '' })}
                        />
                    )}
                </Stack>
            )}

            {hasActiveFilters && (
                <Box sx={{ mb: 3, p: 2, bgcolor: 'primary.light', color: 'primary.contrastText', borderRadius: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
                        Filtered Total
                    </Typography>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>
                        {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(totalAmount)}
                    </Typography>
                </Box>
            )}

            {isLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
                    <CircularProgress />
                </Box>
            ) : (
                <TransactionList
                    transactions={transactions ?? []}
                    categories={categories ?? []}
                    onDelete={handleDeleteClick}
                />
            )}

            {!isLoading && total > limit && (
                <Box sx={{ mt: 4, display: 'flex', justifyContent: 'center' }}>
                    <Pagination
                        count={Math.ceil(total / limit)}
                        page={page}
                        onChange={(_, value) => setPage(value)}
                        color="primary"
                        size="large"
                        sx={{
                            '& .MuiPaginationItem-root': {
                                borderRadius: 2,
                                fontWeight: 600
                            }
                        }}
                    />
                </Box>
            )}

            <TransactionDialog
                open={isDialogOpen}
                onClose={() => setIsDialogOpen(false)}
                onSave={handleSaveTransaction}
                categories={categories ?? []}
            />

            <ConfirmDialog
                open={confirmOpen}
                title="Delete Transaction"
                message="Are you sure you want to delete this transaction?"
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
