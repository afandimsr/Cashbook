import type { SplitMethod } from '../../../../domain/entities/SharedExpense';

export interface CalcParticipant {
  baseShare?: number; // EXACT method
  percentage?: number; // PERCENTAGE method
}

export interface CalcItem {
  price: number;
  quantity: number;
  participantIndexes: number[];
}

export interface CalcResult {
  subtotal: number;
  total: number;
  baseShares: number[];
  finalShares: number[];
  error: string | null;
}

export interface CalcParams {
  method: SplitMethod;
  subtotal: number; // ignored for ITEM (recomputed from items)
  taxAmount: number;
  serviceCharge: number;
  otherCharge: number;
  discountAmount: number;
  participants: CalcParticipant[];
  items: CalcItem[];
  payerIndex: number;
}

/**
 * computeSplit mirrors the backend share logic (whole rupiah, payer absorbs the
 * rounding remainder) so the modal can show an accurate live preview.
 */
export function computeSplit(params: CalcParams): CalcResult {
  const { method, taxAmount, serviceCharge, otherCharge, discountAmount, participants, items, payerIndex } = params;
  const n = participants.length;
  const base = new Array<number>(n).fill(0);
  let subtotal = params.subtotal;
  let error: string | null = null;

  const result = (): CalcResult => {
    const net = Math.round(subtotal - discountAmount + taxAmount + serviceCharge + otherCharge);
    const finalShares = new Array<number>(n).fill(0);
    if (subtotal > 0 && !error) {
      const pIdx = payerIndex >= 0 && payerIndex < n ? payerIndex : 0;
      let allocated = 0;
      for (let i = 0; i < n; i++) {
        if (i === pIdx) continue;
        finalShares[i] = Math.round((base[i] / subtotal) * net);
        allocated += finalShares[i];
      }
      finalShares[pIdx] = Math.round(net - allocated);
    }
    return { subtotal, total: subtotal > 0 ? net : 0, baseShares: base, finalShares, error };
  };

  if (n === 0) {
    error = 'Add at least one participant';
    return result();
  }

  if (method === 'EQUAL') {
    if (subtotal <= 0) {
      error = 'Enter a subtotal greater than zero';
      return result();
    }
    const b = Math.floor(subtotal / n);
    for (let i = 0; i < n; i++) base[i] = b;
  } else if (method === 'EXACT') {
    if (subtotal <= 0) {
      error = 'Enter a subtotal greater than zero';
      return result();
    }
    let sum = 0;
    for (let i = 0; i < n; i++) {
      base[i] = participants[i].baseShare || 0;
      sum += base[i];
    }
    if (Math.round(sum) !== Math.round(subtotal)) {
      error = `Shares (${Math.round(sum)}) must equal subtotal (${Math.round(subtotal)})`;
    }
  } else if (method === 'PERCENTAGE') {
    if (subtotal <= 0) {
      error = 'Enter a subtotal greater than zero';
      return result();
    }
    let pct = 0;
    for (let i = 0; i < n; i++) pct += participants[i].percentage || 0;
    if (Math.abs(pct - 100) > 0.01) {
      error = `Percentages must sum to 100 (got ${pct.toFixed(2)}%)`;
    }
    for (let i = 0; i < n; i++) {
      base[i] = Math.floor((subtotal * (participants[i].percentage || 0)) / 100);
    }
  } else if (method === 'ITEM') {
    if (items.length === 0) {
      error = 'Add at least one item';
      return result();
    }
    subtotal = 0;
    for (const item of items) {
      const qty = item.quantity > 0 ? item.quantity : 1;
      const itemTotal = item.price * qty;
      subtotal += itemTotal;
      const idxs = item.participantIndexes;
      if (idxs.length === 0) {
        error = error || 'Every item needs at least one participant';
        continue;
      }
      const per = Math.floor(itemTotal / idxs.length);
      const rem = Math.round(itemTotal - per * idxs.length);
      idxs.forEach((pi, k) => {
        if (pi < 0 || pi >= n) return;
        base[pi] += per + (k === 0 ? rem : 0);
      });
    }
  }

  if (taxAmount < 0 || serviceCharge < 0 || otherCharge < 0 || discountAmount < 0) {
    error = error || 'Charges and discount must not be negative';
  }
  if (discountAmount > subtotal) error = error || 'Discount cannot exceed subtotal';

  return result();
}
