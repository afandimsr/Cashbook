import React from 'react';
import {
    Box,
    Typography,
    Card,
    Stack,
    IconButton,
    Chip
} from '@mui/material';
import { Delete as DeleteIcon } from '@mui/icons-material';
import type { RecurringTransaction } from '../../../../../domain/entities/RecurringTransaction';
import type { Category } from '../../../../../domain/entities/Category';
import { formatIDR } from '../../../../utils/formatCurrency';

interface RecurringCardProps {
    recurring: RecurringTransaction;
    category?: Category;
    onDelete: (id: number) => void;
}

export const RecurringCard: React.FC<RecurringCardProps> = ({ recurring, category, onDelete }) => {
    return (
        <Card sx={{ p: { xs: 2, md: 3 }, borderRadius: 3 }}>
            <Stack direction="row" justifyContent="space-between" alignItems="start" >
                <Box>
                    <Stack direction="row" spacing={1} mb={1} width="100%" overflow="auto">
                        <Chip
                            label={recurring.frequency.toUpperCase()}
                            size="small"
                            color="primary"
                            variant="outlined"
                        />
                        <Chip
                            label={recurring.type.toUpperCase()}
                            size="small"
                            color={recurring.type === 'income' ? 'success' : 'error'}
                        />
                    </Stack>
                    <Typography variant="h6" sx={{ fontWeight: 700 }}>{recurring.note || 'No Note'}</Typography>
                    <Typography color="text.secondary" variant="body2">
                        {category?.name || 'Unknown Category'}
                    </Typography>
                    <Typography variant="h6" sx={{ fontWeight: 800 }}>
                        {formatIDR(recurring.amount)}
                    </Typography>
                </Box>

            </Stack>
            <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="caption" color="text.secondary">
                    Last processed: {recurring.last_processed ? new Date(recurring.last_processed).toLocaleDateString() : 'Never'}
                </Typography>
                <IconButton color="error" onClick={() => onDelete(recurring.id)}>
                    <DeleteIcon />
                </IconButton>
            </Box>
        </Card>
    );
};
