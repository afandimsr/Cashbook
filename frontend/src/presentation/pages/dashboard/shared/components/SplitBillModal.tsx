import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Autocomplete,
  Box,
  Typography,
  IconButton,
  List,
  ListItem,
  ListItemText,
  Divider,
  ToggleButton,
  ToggleButtonGroup,
  MenuItem,
  Chip,
  Alert,
  Stack,
  Paper,
} from '@mui/material';
import { Delete as DeleteIcon, Add as AddIcon } from '@mui/icons-material';
import { useUserStore } from '../../../../../state/userStore';
import { useCategoryStore } from '../../../../../state/categoryStore';
import { useSplitStore } from '../../../../../state/splitStore';
import type { CreateSplitRequest, SplitMethod } from '../../../../../domain/entities/SharedExpense';
import { useAuthStore } from '../../../../../state/authStore';
import { formatIDR } from '../../../../utils/formatCurrency';
import { computeSplit } from '../splitCalc';

interface SplitBillModalProps {
  open: boolean;
  onClose: () => void;
}

interface ParticipantEntry {
  user_id: number | null;
  name: string;
  baseShare: string; // EXACT method (Rp)
  percentage: string; // PERCENTAGE method (%)
}

interface ItemEntry {
  name: string;
  price: string;
  quantity: string;
  categoryId: number | '';
  participantIndexes: number[];
}

type AmountMode = 'amount' | 'percent';

const num = (v: string): number => {
  const n = parseFloat(v);
  return Number.isFinite(n) ? n : 0;
};

export const SplitBillModal: React.FC<SplitBillModalProps> = ({ open, onClose }) => {
  const { users, fetchUsers } = useUserStore();
  const { categories, fetchCategories } = useCategoryStore();
  const { createSplit } = useSplitStore();

  const [title, setTitle] = useState('');
  const [categoryID, setCategoryID] = useState<number | ''>('');
  const [date, setDate] = useState(new Date().toISOString().split('T')[0]);
  const [method, setMethod] = useState<SplitMethod>('EQUAL');
  const [subtotalInput, setSubtotalInput] = useState('');
  const [participants, setParticipants] = useState<ParticipantEntry[]>([]);
  const [shadowName, setShadowName] = useState('');
  const [items, setItems] = useState<ItemEntry[]>([]);
  const [discount, setDiscount] = useState<{ mode: AmountMode; value: string }>({ mode: 'amount', value: '' });
  const [tax, setTax] = useState<{ mode: AmountMode; value: string }>({ mode: 'amount', value: '' });
  const [service, setService] = useState<{ mode: AmountMode; value: string }>({ mode: 'amount', value: '' });
  const [other, setOther] = useState<{ mode: AmountMode; value: string }>({ mode: 'amount', value: '' });
  const [submitError, setSubmitError] = useState<string | null>(null);

  const currentUser = useAuthStore((state) => state.user);

  const resetForm = () => {
    setTitle('');
    setCategoryID('');
    setDate(new Date().toISOString().split('T')[0]);
    setMethod('EQUAL');
    setSubtotalInput('');
    setShadowName('');
    setItems([]);
    setDiscount({ mode: 'amount', value: '' });
    setTax({ mode: 'amount', value: '' });
    setService({ mode: 'amount', value: '' });
    setOther({ mode: 'amount', value: '' });
    setSubmitError(null);
    setParticipants(
      currentUser ? [{ user_id: currentUser.id, name: `${currentUser.name} (You)`, baseShare: '', percentage: '' }] : [],
    );
  };

  useEffect(() => {
    if (open) {
      fetchUsers();
      fetchCategories();
      resetForm();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open]);

  const payerIndex = participants.findIndex((p) => p.user_id === currentUser?.id);

  const handleAddUser = (user: any) => {
    if (user && !participants.find((p) => p.user_id === user.id)) {
      setParticipants([...participants, { user_id: user.id, name: user.name, baseShare: '', percentage: '' }]);
    }
  };

  const handleAddShadow = () => {
    if (shadowName.trim()) {
      setParticipants([...participants, { user_id: null, name: shadowName.trim(), baseShare: '', percentage: '' }]);
      setShadowName('');
    }
  };

  const handleRemoveParticipant = (index: number) => {
    if (participants[index].user_id === currentUser?.id) return; // keep the payer
    setParticipants(participants.filter((_, i) => i !== index));
    // Drop the removed participant from item assignments and reindex.
    setItems((prev) =>
      prev.map((it) => ({
        ...it,
        participantIndexes: it.participantIndexes
          .filter((pi) => pi !== index)
          .map((pi) => (pi > index ? pi - 1 : pi)),
      })),
    );
  };

  const updateParticipant = (index: number, patch: Partial<ParticipantEntry>) => {
    setParticipants((prev) => prev.map((p, i) => (i === index ? { ...p, ...patch } : p)));
  };

  const addItem = () => {
    setItems([...items, { name: '', price: '', quantity: '1', categoryId: '', participantIndexes: [] }]);
  };

  const updateItem = (index: number, patch: Partial<ItemEntry>) => {
    setItems((prev) => prev.map((it, i) => (i === index ? { ...it, ...patch } : it)));
  };

  const removeItem = (index: number) => {
    setItems(items.filter((_, i) => i !== index));
  };

  const toggleItemParticipant = (itemIndex: number, participantIndex: number) => {
    setItems((prev) =>
      prev.map((it, i) => {
        if (i !== itemIndex) return it;
        const has = it.participantIndexes.includes(participantIndex);
        return {
          ...it,
          participantIndexes: has
            ? it.participantIndexes.filter((pi) => pi !== participantIndex)
            : [...it.participantIndexes, participantIndex],
        };
      }),
    );
  };

  // --- Derived amounts (shared between preview and submit) ---
  const currentSubtotal =
    method === 'ITEM'
      ? items.reduce((s, it) => s + num(it.price) * (parseInt(it.quantity, 10) || 1), 0)
      : num(subtotalInput);

  const discountAmount = discount.mode === 'percent' ? Math.round((currentSubtotal * num(discount.value)) / 100) : num(discount.value);
  const afterDiscount = currentSubtotal - discountAmount;
  // Percentage-based charges are computed on the subtotal after discount.
  const pctAmount = (c: { mode: AmountMode; value: string }) =>
    c.mode === 'percent' ? Math.round((afterDiscount * num(c.value)) / 100) : num(c.value);
  const taxAmount = pctAmount(tax);
  const serviceAmount = pctAmount(service);
  const otherAmount = pctAmount(other);

  const preview = computeSplit({
    method,
    subtotal: currentSubtotal,
    taxAmount,
    serviceCharge: serviceAmount,
    otherCharge: otherAmount,
    discountAmount,
    participants: participants.map((p) => ({ baseShare: num(p.baseShare), percentage: num(p.percentage) })),
    items: items.map((it) => ({ price: num(it.price), quantity: parseInt(it.quantity, 10) || 1, participantIndexes: it.participantIndexes })),
    payerIndex,
  });

  const handleSubmit = async () => {
    setSubmitError(null);
    if (!title || !categoryID || participants.length === 0) {
      setSubmitError('Title, category, and at least one participant are required.');
      return;
    }
    if (preview.error) {
      setSubmitError(preview.error);
      return;
    }

    const request: CreateSplitRequest = {
      title,
      subtotal: method === 'ITEM' ? 0 : currentSubtotal, // backend recomputes for ITEM
      tax_amount: taxAmount,
      service_charge: serviceAmount,
      other_charge: otherAmount,
      discount_amount: discountAmount,
      category_id: categoryID as number,
      date: new Date(date).toISOString(),
      payer_id: currentUser?.id ?? 0,
      split_method: method,
      participants: participants.map((p) => ({
        user_id: p.user_id,
        shadow_name: p.user_id ? '' : p.name,
        base_share: method === 'EXACT' ? num(p.baseShare) : undefined,
        percentage: method === 'PERCENTAGE' ? num(p.percentage) : undefined,
      })),
      items:
        method === 'ITEM'
          ? items.map((it) => ({
              name: it.name || 'Item',
              price: num(it.price),
              quantity: parseInt(it.quantity, 10) || 1,
              category_id: it.categoryId === '' ? null : it.categoryId,
              participant_indexes: it.participantIndexes,
            }))
          : undefined,
    };

    try {
      await createSplit(request);
      onClose();
      resetForm();
    } catch (err: any) {
      setSubmitError(err?.message || 'Failed to create split bill');
    }
  };

  const renderCharge = (
    label: string,
    state: { mode: AmountMode; value: string },
    setState: (v: { mode: AmountMode; value: string }) => void,
  ) => (
    <Box sx={{ flex: 1, minWidth: 200 }}>
      <Typography variant="subtitle2" gutterBottom>
        {label}
      </Typography>
      <Box sx={{ display: 'flex', gap: 1 }}>
        <TextField
          size="small"
          type="number"
          value={state.value}
          onChange={(e) => setState({ ...state, value: e.target.value })}
          fullWidth
        />
        <ToggleButtonGroup
          value={state.mode}
          exclusive
          size="small"
          onChange={(_, v) => v && setState({ ...state, mode: v })}
        >
          <ToggleButton value="amount">Rp</ToggleButton>
          <ToggleButton value="percent">%</ToggleButton>
        </ToggleButtonGroup>
      </Box>
    </Box>
  );

  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Create Split Bill</DialogTitle>
      <DialogContent>
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 1 }}>
          <TextField label="Title" value={title} onChange={(e) => setTitle(e.target.value)} fullWidth required />

          <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
            <TextField
              select
              label="Category"
              value={categoryID}
              onChange={(e) => setCategoryID(Number(e.target.value))}
              sx={{ flex: 1, minWidth: 180 }}
              required
            >
              <MenuItem value="">
                <em>Select</em>
              </MenuItem>
              {categories.map((c) => (
                <MenuItem key={c.id} value={c.id}>
                  {c.name}
                </MenuItem>
              ))}
            </TextField>
            <TextField
              label="Date"
              type="date"
              value={date}
              onChange={(e) => setDate(e.target.value)}
              sx={{ flex: 1, minWidth: 180 }}
              required
              InputLabelProps={{ shrink: true }}
            />
          </Box>

          {/* Split method */}
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Split method
            </Typography>
            <ToggleButtonGroup
              value={method}
              exclusive
              size="small"
              onChange={(_, v) => v && setMethod(v)}
            >
              <ToggleButton value="EQUAL">Equal</ToggleButton>
              <ToggleButton value="EXACT">Exact (Rp)</ToggleButton>
              <ToggleButton value="PERCENTAGE">Percentage</ToggleButton>
              <ToggleButton value="ITEM">Per-item</ToggleButton>
            </ToggleButtonGroup>
          </Box>

          {method !== 'ITEM' && (
            <TextField
              label="Subtotal (Rp)"
              type="number"
              value={subtotalInput}
              onChange={(e) => setSubtotalInput(e.target.value)}
              fullWidth
              helperText="Amount before tax and discount"
            />
          )}

          <Divider />
          <Typography variant="subtitle1">Participants</Typography>

          <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap' }}>
            <Autocomplete
              options={users}
              getOptionLabel={(option) => option.name}
              onChange={(_, value) => handleAddUser(value)}
              renderInput={(params) => <TextField {...params} label="Add registered user" />}
              sx={{ flex: 1, minWidth: 220 }}
              value={null}
              blurOnSelect
            />
            <Box sx={{ display: 'flex', gap: 1, flex: 1, minWidth: 220 }}>
              <TextField
                label="Shadow participant name"
                value={shadowName}
                onChange={(e) => setShadowName(e.target.value)}
                fullWidth
                onKeyDown={(e) => e.key === 'Enter' && (e.preventDefault(), handleAddShadow())}
              />
              <Button onClick={handleAddShadow} variant="outlined" sx={{ minWidth: 'auto' }}>
                <AddIcon />
              </Button>
            </Box>
          </Box>

          <List dense>
            {participants.map((p, index) => (
              <ListItem key={index} secondaryAction={
                p.user_id !== currentUser?.id ? (
                  <IconButton edge="end" onClick={() => handleRemoveParticipant(index)}>
                    <DeleteIcon />
                  </IconButton>
                ) : null
              }>
                <ListItemText primary={p.name} secondary={p.user_id ? 'Registered' : 'Shadow'} sx={{ flex: '0 0 40%' }} />
                {method === 'EXACT' && (
                  <TextField
                    size="small"
                    label="Share (Rp)"
                    type="number"
                    value={p.baseShare}
                    onChange={(e) => updateParticipant(index, { baseShare: e.target.value })}
                    sx={{ width: 140, mr: 5 }}
                  />
                )}
                {method === 'PERCENTAGE' && (
                  <TextField
                    size="small"
                    label="%"
                    type="number"
                    value={p.percentage}
                    onChange={(e) => updateParticipant(index, { percentage: e.target.value })}
                    sx={{ width: 100, mr: 5 }}
                  />
                )}
              </ListItem>
            ))}
          </List>

          {/* Items (ITEM method) */}
          {method === 'ITEM' && (
            <Box>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                <Typography variant="subtitle1">Items</Typography>
                <Button size="small" startIcon={<AddIcon />} onClick={addItem}>
                  Add item
                </Button>
              </Box>
              <Stack spacing={1.5}>
                {items.map((it, index) => (
                  <Paper key={index} variant="outlined" sx={{ p: 1.5 }}>
                    <Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', alignItems: 'center' }}>
                      <TextField
                        size="small"
                        label="Name"
                        value={it.name}
                        onChange={(e) => updateItem(index, { name: e.target.value })}
                        sx={{ flex: 2, minWidth: 140 }}
                      />
                      <TextField
                        size="small"
                        label="Price (Rp)"
                        type="number"
                        value={it.price}
                        onChange={(e) => updateItem(index, { price: e.target.value })}
                        sx={{ flex: 1, minWidth: 110 }}
                      />
                      <TextField
                        size="small"
                        label="Qty"
                        type="number"
                        value={it.quantity}
                        onChange={(e) => updateItem(index, { quantity: e.target.value })}
                        sx={{ width: 70 }}
                      />
                      <TextField
                        select
                        size="small"
                        label="Category"
                        value={it.categoryId}
                        onChange={(e) => updateItem(index, { categoryId: e.target.value === '' ? '' : Number(e.target.value) })}
                        sx={{ flex: 1, minWidth: 130 }}
                      >
                        <MenuItem value="">
                          <em>Bill default</em>
                        </MenuItem>
                        {categories.map((c) => (
                          <MenuItem key={c.id} value={c.id}>
                            {c.name}
                          </MenuItem>
                        ))}
                      </TextField>
                      <IconButton onClick={() => removeItem(index)}>
                        <DeleteIcon />
                      </IconButton>
                    </Box>
                    <Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap', mt: 1 }}>
                      {participants.map((p, pi) => (
                        <Chip
                          key={pi}
                          label={p.name}
                          size="small"
                          color={it.participantIndexes.includes(pi) ? 'primary' : 'default'}
                          variant={it.participantIndexes.includes(pi) ? 'filled' : 'outlined'}
                          onClick={() => toggleItemParticipant(index, pi)}
                        />
                      ))}
                    </Box>
                  </Paper>
                ))}
                {items.length === 0 && (
                  <Typography variant="body2" color="text.secondary">
                    Add items and tap participant chips to assign who shares each item.
                  </Typography>
                )}
              </Stack>
            </Box>
          )}

          <Divider />

          {/* Charges — percentages apply to the subtotal after discount */}
          <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap' }}>
            {renderCharge('Discount', discount, setDiscount)}
            {renderCharge('Tax', tax, setTax)}
            {renderCharge('Service charge', service, setService)}
            {renderCharge('Other charge', other, setOther)}
          </Box>

          {/* Preview */}
          <Paper variant="outlined" sx={{ p: 2, bgcolor: 'action.hover' }}>
            <Typography variant="subtitle2" gutterBottom>
              Preview
            </Typography>
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="body2" color="text.secondary">Subtotal</Typography>
              <Typography variant="body2">{formatIDR(preview.subtotal)}</Typography>
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="body2" color="text.secondary">Discount</Typography>
              <Typography variant="body2">- {formatIDR(discountAmount)}</Typography>
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="body2" color="text.secondary">Tax</Typography>
              <Typography variant="body2">+ {formatIDR(taxAmount)}</Typography>
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="body2" color="text.secondary">Service charge</Typography>
              <Typography variant="body2">+ {formatIDR(serviceAmount)}</Typography>
            </Box>
            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
              <Typography variant="body2" color="text.secondary">Other charge</Typography>
              <Typography variant="body2">+ {formatIDR(otherAmount)}</Typography>
            </Box>
            <Divider sx={{ my: 1 }} />
            <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
              <Typography variant="subtitle2">Total</Typography>
              <Typography variant="subtitle2">{formatIDR(preview.total)}</Typography>
            </Box>
            {participants.map((p, i) => (
              <Box key={i} sx={{ display: 'flex', justifyContent: 'space-between' }}>
                <Typography variant="body2">{p.name}</Typography>
                <Typography variant="body2">{formatIDR(preview.finalShares[i] || 0)}</Typography>
              </Box>
            ))}
            {preview.error && (
              <Alert severity="warning" sx={{ mt: 1 }}>
                {preview.error}
              </Alert>
            )}
          </Paper>

          {submitError && <Alert severity="error">{submitError}</Alert>}
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button onClick={handleSubmit} variant="contained" color="primary" disabled={!!preview.error}>
          Create Split
        </Button>
      </DialogActions>
    </Dialog>
  );
};
