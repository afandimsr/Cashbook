export function formatIDR(amount: number | undefined | null): string {
    const value = amount ?? 0;
    try {
        return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(value);
    } catch (e) {
        // fallback
        return `Rp ${Math.round(value).toLocaleString('id-ID')}`;
    }
}
