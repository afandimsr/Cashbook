import React, { useState } from 'react';
import {
    Container,
    Box,
    Typography,
    Button,
    CircularProgress,
    Snackbar,
    Alert,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useCategories } from '../../../../application/hooks/useCategories';
import { CategoryList } from './components/CategoryList';
import { CategoryDialog } from './components/CategoryDialog';
import type { Category } from '../../../../domain/entities/Category';

export const CategoryPage: React.FC = () => {
    const {
        categories,
        isLoading,
        handleAddCategory,
        handleDeleteCategory
    } = useCategories();

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [snack, setSnack] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' });

    const handleSaveCategory = async (data: Omit<Category, 'id'>) => {
        try {
            await handleAddCategory(data);
            setSnack({ open: true, message: 'Category created', severity: 'success' });
            setIsDialogOpen(false);
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to create category', severity: 'error' });
            throw err;
        }
    };

    const handleDelete = async (id: number) => {
        try {
            if (!window.confirm('Are you sure you want to delete this category?')) return;
            await handleDeleteCategory(id);
            setSnack({ open: true, message: 'Category deleted', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to delete category', severity: 'error' });
            throw err;
        }
    };

    return (
        <Container maxWidth="xl" sx={{ py: 4 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
                <Box>
                    <Typography variant="h3" sx={{ fontWeight: 800, mb: 1 }}>Categories</Typography>
                    <Typography variant="h6" color="text.secondary">Organize your finances with custom categories</Typography>
                </Box>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={() => setIsDialogOpen(true)}
                    sx={{ borderRadius: 2 }}
                >
                    Add Category
                </Button>
            </Box>

            {isLoading ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
                    <CircularProgress />
                </Box>
            ) : (
                <CategoryList
                    categories={categories ?? []}
                    onDelete={handleDelete}
                    isLoading={isLoading}
                />
            )}

                <CategoryDialog
                    open={isDialogOpen}
                    onClose={() => setIsDialogOpen(false)}
                    onSave={handleSaveCategory}
                />

            <Snackbar open={snack.open} autoHideDuration={4000} onClose={() => setSnack({ ...snack, open: false })}>
                <Alert onClose={() => setSnack({ ...snack, open: false })} severity={snack.severity} sx={{ width: '100%' }}>
                    {snack.message}
                </Alert>
            </Snackbar>
        </Container>
    );
};
