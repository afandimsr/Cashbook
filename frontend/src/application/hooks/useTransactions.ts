import { useEffect, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
import { useTransactionStore } from '../../state/transactionStore';
import type { Transaction } from '../../domain/entities/Transaction';

export const useTransactions = () => {
    const {
        transactions,
        isLoading,
        error,
        fetchTransactions,
        addTransaction,
        deleteTransaction
    } = useTransactionStore();

    const [searchParams] = useSearchParams();
    const query = searchParams.get('q') || '';

    const refreshTransactions = useCallback((page = 1, limit = 20) => {
        return fetchTransactions(page, limit, query);
    }, [fetchTransactions, query]);

    useEffect(() => {
        refreshTransactions();
    }, [refreshTransactions]);

    const handleAddTransaction = async (data: Omit<Transaction, 'id'>) => {
        await addTransaction(data);
    };

    const handleDeleteTransaction = async (id: number) => {
        if (window.confirm('Are you sure you want to delete this transaction?')) {
            await deleteTransaction(id);
        }
    };

    return {
        transactions,
        isLoading,
        error,
        refreshTransactions,
        handleAddTransaction,
        handleDeleteTransaction,
        query
    };
};
