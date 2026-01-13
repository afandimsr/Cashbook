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
    Box,
    Chip
} from '@mui/material';
import { Delete as DeleteIcon } from '@mui/icons-material';
import type { Category } from '../../../../../domain/entities/Category';

interface CategoryListProps {
    categories?: Category[];
    onDelete: (id: number) => void;
    isLoading: boolean;
}

export const CategoryList: React.FC<CategoryListProps> = ({ categories, onDelete, isLoading }) => {
    const items = categories ?? [];

    return (
        <TableContainer component={Paper} sx={{ borderRadius: 3, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
            <Table>
                <TableHead sx={{ backgroundColor: 'action.hover' }}>
                    <TableRow>
                        <TableCell sx={{ fontWeight: 700 }}>Name</TableCell>
                        <TableCell sx={{ fontWeight: 700 }}>Type</TableCell>
                        <TableCell sx={{ fontWeight: 700 }}>Color</TableCell>
                        <TableCell sx={{ fontWeight: 700 }} align="center">Actions</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {items.map((cat) => (
                        <TableRow key={cat.id ?? `${cat.name}-${Math.random()}`} hover>
                            <TableCell sx={{ fontWeight: 600 }}>{cat.name ?? '—'}</TableCell>
                            <TableCell>
                                <Chip
                                    label={(cat.type ?? 'expense').toString().toUpperCase()}
                                    size="small"
                                    color={cat.type === 'income' ? 'success' : 'error'}
                                    variant="outlined"
                                />
                            </TableCell>
                            <TableCell>
                                <Stack direction="row" spacing={1} alignItems="center">
                                    <Box sx={{ width: 20, height: 20, borderRadius: 1, backgroundColor: cat.color ?? '#ccc' }} />
                                    <Typography variant="body2" color="text.secondary">{cat.color ?? '—'}</Typography>
                                </Stack>
                            </TableCell>
                            <TableCell align="center">
                                <IconButton
                                    color="error"
                                    onClick={() => {
                                        if (cat.id !== undefined && cat.id !== null) onDelete(cat.id);
                                    }}
                                    disabled={cat.id === undefined || cat.id === null}
                                >
                                    <DeleteIcon />
                                </IconButton>
                            </TableCell>
                        </TableRow>
                    ))}
                    {items.length === 0 && !isLoading && (
                        <TableRow>
                            <TableCell colSpan={4} align="center" sx={{ py: 4 }}>
                                <Typography color="text.secondary">No categories found</Typography>
                            </TableCell>
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </TableContainer>
    );
};
