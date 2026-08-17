export type SplitBillStatus = 'OPEN' | 'PARTIALLY_SETTLED' | 'SETTLED';

export type SplitMethod = 'EQUAL' | 'EXACT' | 'PERCENTAGE' | 'ITEM';

export interface SplitParticipant {
  id: string;
  split_bill_id: string;
  user_id: number | null;
  shadow_name: string;
  user_name: string;
  base_share: number;
  share_amount: number;
  is_paid: boolean;
  settled_at: string | null;
  transaction_id: number | null;
  created_at: string;
}

export interface SplitItem {
  id: string;
  split_bill_id: string;
  name: string;
  price: number;
  quantity: number;
  category_id: number | null;
  participant_ids?: string[];
  created_at: string;
}

export interface SplitBill {
  id: string;
  creator_id: number;
  payer_id: number | null;
  category_id: number;
  title: string;
  subtotal: number;
  tax_amount: number;
  service_charge: number;
  other_charge: number;
  discount_amount: number;
  total_amount: number;
  date: string;
  status: SplitBillStatus;
  split_method: SplitMethod;
  participants?: SplitParticipant[];
  items?: SplitItem[];
  created_at: string;
  updated_at: string;
}

export interface SplitSummary {
  total_owed_to_me: number;
  total_i_owe: number;
}

export interface CreateSplitParticipant {
  user_id: number | null;
  shadow_name: string;
  base_share?: number;
  percentage?: number;
}

export interface CreateSplitItem {
  name: string;
  price: number;
  quantity: number;
  category_id?: number | null;
  participant_indexes: number[];
}

export interface CreateSplitRequest {
  title: string;
  subtotal: number;
  tax_amount: number;
  service_charge: number;
  other_charge: number;
  discount_amount: number;
  category_id: number;
  date: string;
  payer_id: number;
  split_method: SplitMethod;
  participants: CreateSplitParticipant[];
  items?: CreateSplitItem[];
}
