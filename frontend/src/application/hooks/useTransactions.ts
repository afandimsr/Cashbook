import { useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTransactionStore, type TransactionFilter } from '../../state/transactionStore';
import type { Transaction } from '../../domain/entities/Transaction';

export const useTransactions = () => {
    const {
        transactions,
        total,
        totalAmount,
        isLoading,
        error,
        fetchTransactions,
        addTransaction,
        deleteTransaction
    } = useTransactionStore();

    const [searchParams, setSearchParams] = useSearchParams();
    const page = parseInt(searchParams.get('page') || '1');
    const limit = parseInt(searchParams.get('limit') || '20');

    const filters: TransactionFilter = {
        q: searchParams.get('q') || '',
        category_id: searchParams.get('category_id') || '',
        type: searchParams.get('type') || '',
        start_date: searchParams.get('start_date') || '',
        end_date: searchParams.get('end_date') || '',
    };

    const refreshTransactions = useCallback(() => {
        return fetchTransactions(page, limit, filters);
    }, [fetchTransactions, page, limit, searchParams]); // Depends on page, limit, and params

    useEffect(() => {
        refreshTransactions();
    }, [refreshTransactions]);

    const setFilters = (newFilters: Partial<TransactionFilter>) => {
        const nextParams = new URLSearchParams(searchParams);
        Object.entries(newFilters).forEach(([key, value]) => {
            if (value) {
                nextParams.set(key, value.toString());
            } else {
                nextParams.delete(key);
            }
        });
        nextParams.set('page', '1'); // Reset to page 1 on filter change
        setSearchParams(nextParams);
    };

    const setPage = (newPage: number) => {
        const nextParams = new URLSearchParams(searchParams);
        nextParams.set('page', newPage.toString());
        setSearchParams(nextParams);
    };

    const clearFilters = () => {
        setSearchParams({});
    };

    const handleAddTransaction = async (data: Omit<Transaction, 'id'>) => {
        await addTransaction(data);
    };

    const handleDeleteTransaction = async (id: number) => {
        await deleteTransaction(id);
    };

    return {
        transactions,
        total,
        totalAmount,
        isLoading,
        error,
        page,
        limit,
        setPage,
        refreshTransactions,
        handleAddTransaction,
        handleDeleteTransaction,
        filters,
        setFilters,
        clearFilters
    };
};
