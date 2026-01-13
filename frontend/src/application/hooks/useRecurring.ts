import { useEffect } from 'react';
import { useRecurringStore } from '../../state/recurringStore';
import type { RecurringTransaction } from '../../domain/entities/RecurringTransaction';

export const useRecurring = () => {
    const {
        recurringTransactions,
        isLoading,
        error,
        fetchRecurring,
        createRecurring,
        deleteRecurring,
        processDue
    } = useRecurringStore();

    useEffect(() => {
        fetchRecurring();
    }, [fetchRecurring]);

    const handleCreateRecurring = async (data: Partial<RecurringTransaction>) => {
        await createRecurring(data);
    };

    const handleDeleteRecurring = async (id: number) => {
        await deleteRecurring(id);
    };

    const handleProcessDue = async () => {
        await processDue();
    };

    return {
        recurringTransactions,
        isLoading,
        error,
        handleCreateRecurring,
        handleDeleteRecurring,
        handleProcessDue,
        refreshRecurring: fetchRecurring
    };
};
