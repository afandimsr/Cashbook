import React from 'react';
import {
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  Button,
  Chip,
} from '@mui/material';
import { CheckCircle as CheckCircleIcon } from '@mui/icons-material';
import type { SplitParticipant } from '../../../../../domain/entities/SharedExpense';
import { useSplitStore } from '../../../../../state/splitStore';
import { formatIDR } from '../../../../utils/formatCurrency';

interface SettlementListProps {
  billID: string;
  participants: SplitParticipant[];
  isPayer: boolean;
}

export const SettlementList: React.FC<SettlementListProps> = ({ billID, participants, isPayer }) => {
  const { settleParticipant } = useSplitStore();

  const handleSettle = async (participantID: string) => {
    try {
      await settleParticipant(billID, participantID);
    } catch (err) {
      console.error(err);
    }
  };

  return (
    <List dense>
      {participants.map((p) => (
        <ListItem key={p.id}>
          <ListItemText
            primary={p.user_name || p.shadow_name || `User #${p.user_id}`}
            secondary={`Share: ${formatIDR(p.share_amount)}`}
          />
          <ListItemSecondaryAction>
            {p.is_paid ? (
              <Chip
                icon={<CheckCircleIcon />}
                label="Paid"
                color="success"
                size="small"
                variant="outlined"
              />
            ) : isPayer ? (
              <Button
                variant="contained"
                size="small"
                onClick={() => handleSettle(p.id)}
                color="primary"
              >
                Mark Paid
              </Button>
            ) : (
              <Chip label="Unpaid" size="small" variant="outlined" />
            )}
          </ListItemSecondaryAction>
        </ListItem>
      ))}
    </List>
  );
};
