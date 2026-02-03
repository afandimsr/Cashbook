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
import { ConfirmDialog } from '../../../components/common/ConfirmDialog';
import type { Category } from '../../../../domain/entities/Category';

export const CategoryPage: React.FC = () => {
    const {
        categories,
        isLoading,
        handleAddCategory,
        handleUpdateCategory,
        handleDeleteCategory
    } = useCategories();

    const [isDialogOpen, setIsDialogOpen] = useState(false);
    const [editingCategory, setEditingCategory] = useState<Category | null>(null);
    const [confirmOpen, setConfirmOpen] = useState(false);
    const [itemToDelete, setItemToDelete] = useState<number | null>(null);
    const [snack, setSnack] = useState<{ open: boolean; message: string; severity: 'success' | 'error' }>({ open: false, message: '', severity: 'success' });

    const handleAddClick = () => {
        setEditingCategory(null);
        setIsDialogOpen(true);
    };

    const handleEditClick = (category: Category) => {
        setEditingCategory(category);
        setIsDialogOpen(true);
    };

    const handleSaveCategory = async (data: Omit<Category, 'id'>) => {
        try {
            if (editingCategory?.id) {
                await handleUpdateCategory(editingCategory.id, data);
                setSnack({ open: true, message: 'Category updated', severity: 'success' });
            } else {
                await handleAddCategory(data);
                setSnack({ open: true, message: 'Category created', severity: 'success' });
            }
            setIsDialogOpen(false);
            setEditingCategory(null);
        } catch (err: any) {
            const action = editingCategory ? 'update' : 'create';
            setSnack({ open: true, message: err?.message || `Failed to ${action} category`, severity: 'error' });
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
            await handleDeleteCategory(itemToDelete);
            setSnack({ open: true, message: 'Category deleted', severity: 'success' });
        } catch (err: any) {
            setSnack({ open: true, message: err?.message || 'Failed to delete category', severity: 'error' });
        } finally {
            setConfirmOpen(false);
            setItemToDelete(null);
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
                        Categories
                    </Typography>
                    <Typography variant="h6" color="text.secondary" sx={{ fontSize: { xs: '1rem', md: '1.25rem' } }}>
                        Organize your finances with custom categories
                    </Typography>
                </Box>
                <Button
                    variant="contained"
                    startIcon={<AddIcon />}
                    onClick={handleAddClick}
                    sx={{ borderRadius: 2, width: { xs: '100%', sm: 'auto' } }}
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
                    onEdit={handleEditClick}
                    onDelete={handleDeleteClick}
                    isLoading={isLoading}
                />
            )}

            <CategoryDialog
                open={isDialogOpen}
                onClose={() => {
                    setIsDialogOpen(false);
                    setEditingCategory(null);
                }}
                onSave={handleSaveCategory}
                initialData={editingCategory}
            />

            <ConfirmDialog
                open={confirmOpen}
                title="Delete Category"
                message="Are you sure you want to delete this category? This will not delete the transactions in this category."
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
