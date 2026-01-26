import { useEffect, useCallback } from 'react';
import { useCategoryStore } from '../../state/categoryStore';

export const useCategories = () => {
    const {
        categories,
        isLoading,
        error,
        fetchCategories,
        addCategory,
        updateCategory,
        deleteCategory
    } = useCategoryStore();

    useEffect(() => {
        fetchCategories();
    }, [fetchCategories]);

    const getCategoryById = useCallback((id: number) => {
        return categories ? categories.find(c => c.id === id) : undefined;
    }, [categories]);

    const getCategoriesByType = useCallback((type: 'income' | 'expense') => {
        return categories ? categories.filter(c => c.type === type) : [];
    }, [categories]);

    return {
        categories,
        isLoading,
        error,
        getCategoryById,
        getCategoriesByType,
        refreshCategories: fetchCategories,
        handleAddCategory: addCategory,
        handleUpdateCategory: updateCategory,
        handleDeleteCategory: deleteCategory
    };
};
