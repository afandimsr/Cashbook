import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  Stack,
  CircularProgress,
  Alert,
  IconButton,
  Collapse,
  Divider,
} from '@mui/material';
import {
  Add as AddIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  CallMade as OwedIcon,
  CallReceived as IOweIcon,
} from '@mui/icons-material';
import { useSplitStore } from '../../../../state/splitStore';
import { SplitBillModal } from './components/SplitBillModal';
import { SettlementList } from './components/SettlementList';
import { useAuthStore } from '../../../../state/authStore';
import { formatIDR } from '../../../utils/formatCurrency';

export const SharedExpensePage: React.FC = () => {
  const { splits, summary, isLoading, error, fetchSplits, fetchSummary } = useSplitStore();
  const [modalOpen, setModalOpen] = useState(false);
  const [expandedId, setExpandedId] = useState<string | null>(null);

  const currentUser = useAuthStore((state) => state.user);

  useEffect(() => {
    fetchSplits();
    fetchSummary();
  }, [fetchSplits, fetchSummary]);

  const toggleExpand = (id: string) => {
    setExpandedId(expandedId === id ? null : id);
  };

  return (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 3 }}>
        <Typography variant="h4" fontWeight="bold">
          Shared Expenses
        </Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => setModalOpen(true)}
        >
          Create Split
        </Button>
      </Box>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid size={{ xs: 12, md: 6 }}>
          <Card sx={{ bgcolor: 'success.light', color: 'success.contrastText' }}>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <OwedIcon fontSize="large" />
                <Box>
                  <Typography variant="subtitle2">Total Owed To Me</Typography>
                  <Typography variant="h5" fontWeight="bold">
                    {formatIDR(summary?.total_owed_to_me)}
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>
        <Grid size={{ xs: 12, md: 6 }}>
          <Card sx={{ bgcolor: 'error.light', color: 'error.contrastText' }}>
            <CardContent>
              <Stack direction="row" spacing={2} alignItems="center">
                <IOweIcon fontSize="large" />
                <Box>
                  <Typography variant="subtitle2">Total I Owe</Typography>
                  <Typography variant="h5" fontWeight="bold">
                    {formatIDR(summary?.total_i_owe)}
                  </Typography>
                </Box>
              </Stack>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {error && <Alert severity="error" sx={{ mb: 3 }}>{error}</Alert>}

      {isLoading && splits.length === 0 ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 5 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Stack spacing={2}>
          {splits.map((split) => {
            const isPayer = split.payer_id === currentUser?.id;
            const isExpanded = expandedId === split.id;

            return (
              <Card key={split.id} variant="outlined">
                <CardContent sx={{ pb: 1 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <Box>
                      <Typography variant="h6">{split.title}</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {new Date(split.date).toLocaleDateString()} • {formatIDR(split.total_amount)}
                      </Typography>
                    </Box>
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography
                        variant="subtitle2"
                        color={isPayer ? 'success.main' : 'error.main'}
                        fontWeight="bold"
                      >
                        {isPayer ? 'You Paid' : 'You Owe'}
                      </Typography>
                      <IconButton onClick={() => toggleExpand(split.id)}>
                        {isExpanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
                      </IconButton>
                    </Box>
                  </Box>
                </CardContent>
                <Collapse in={isExpanded}>
                  <Divider />
                  <Box sx={{ p: 2, bgcolor: 'action.hover' }}>
                    <Box sx={{ mb: 1.5 }}>
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography variant="body2" color="text.secondary">Subtotal</Typography>
                        <Typography variant="body2">{formatIDR(split.subtotal)}</Typography>
                      </Box>
                      {split.discount_amount > 0 && (
                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                          <Typography variant="body2" color="text.secondary">Discount</Typography>
                          <Typography variant="body2">- {formatIDR(split.discount_amount)}</Typography>
                        </Box>
                      )}
                      {split.tax_amount > 0 && (
                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                          <Typography variant="body2" color="text.secondary">Tax</Typography>
                          <Typography variant="body2">+ {formatIDR(split.tax_amount)}</Typography>
                        </Box>
                      )}
                      {split.service_charge > 0 && (
                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                          <Typography variant="body2" color="text.secondary">Service charge</Typography>
                          <Typography variant="body2">+ {formatIDR(split.service_charge)}</Typography>
                        </Box>
                      )}
                      {split.other_charge > 0 && (
                        <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                          <Typography variant="body2" color="text.secondary">Other charge</Typography>
                          <Typography variant="body2">+ {formatIDR(split.other_charge)}</Typography>
                        </Box>
                      )}
                      <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                        <Typography variant="body2" fontWeight="bold">Total</Typography>
                        <Typography variant="body2" fontWeight="bold">{formatIDR(split.total_amount)}</Typography>
                      </Box>
                    </Box>
                    <Divider sx={{ mb: 1.5 }} />
                    <Typography variant="subtitle2" gutterBottom>
                      Breakdown
                    </Typography>
                    <SettlementList
                      billID={split.id}
                      participants={split.participants || []}
                      isPayer={isPayer}
                    />
                  </Box>
                </Collapse>
              </Card>
            );
          })}
          {splits.length === 0 && (
            <Typography variant="body1" color="text.secondary" textAlign="center" sx={{ py: 5 }}>
              No shared expenses found. Create one to get started!
            </Typography>
          )}
        </Stack>
      )}

      <SplitBillModal open={modalOpen} onClose={() => setModalOpen(false)} />
    </Box>
  );
};
