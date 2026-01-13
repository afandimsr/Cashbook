import React from 'react';
import {
    Paper,
    Table,
    TableBody,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    IconButton,
    Typography,
    Stack,
    Box
} from '@mui/material';
import { Delete as DeleteIcon } from '@mui/icons-material';
import type { Transaction } from '../../../../../domain/entities/Transaction';
import type { Category } from '../../../../../domain/entities/Category';
import { formatIDR } from '../../../../utils/formatCurrency';

interface TransactionListProps {
    transactions?: Transaction[];
    categories?: Category[];
    onDelete: (id: number) => void;
}

export const TransactionList: React.FC<TransactionListProps> = ({ transactions, categories, onDelete }) => {
    const items = transactions ?? [];
    const cats = categories ?? [];

    return (
        <TableContainer component={Paper} sx={{ borderRadius: 3, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
            <Table>
                <TableHead sx={{ backgroundColor: 'action.hover' }}>
                    <TableRow>
                        <TableCell sx={{ fontWeight: 700 }}>Date</TableCell>
                        <TableCell sx={{ fontWeight: 700 }}>Category</TableCell>
                        <TableCell sx={{ fontWeight: 700 }}>Note</TableCell>
                        <TableCell sx={{ fontWeight: 700 }} align="right">Amount</TableCell>
                        <TableCell sx={{ fontWeight: 700 }} align="center">Actions</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {items.map((tx) => {
                        const category = cats.find(c => c.id === tx.category_id);
                        return (
                            <TableRow key={tx.id} hover>
                                <TableCell>{tx.date ? new Date(tx.date).toLocaleDateString() : '—'}</TableCell>
                                <TableCell>
                                    <Stack direction="row" spacing={1} alignItems="center">
                                        <Box
                                            sx={{
                                                width: 10,
                                                height: 10,
                                                borderRadius: '50%',
                                                backgroundColor: category?.color || '#ccc'
                                            }}
                                        />
                                        <Typography variant="body2">{category?.name || 'Unknown'}</Typography>
                                    </Stack>
                                </TableCell>
                                <TableCell>{tx.note || '—'}</TableCell>
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
                                <TableCell align="center">
                                    <IconButton size="small" color="error" onClick={() => { if (tx.id !== undefined && tx.id !== null) onDelete(tx.id); }} disabled={tx.id === undefined || tx.id === null}>
                                        <DeleteIcon fontSize="small" />
                                    </IconButton>
                                </TableCell>
                            </TableRow>
                        );
                    })}
                    {items.length === 0 && (
                        <TableRow>
                            <TableCell colSpan={5} align="center" sx={{ py: 4 }}>
                                <Typography color="text.secondary">No transactions found</Typography>
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    );
};
