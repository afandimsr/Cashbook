import React from 'react';
import {
    Card,
    Typography,
    Box,
    Stack,
    Button,
    Switch,
    Slider,
    FormControlLabel,
    Checkbox,
    Radio,
    RadioGroup,
    Divider,
    alpha
} from '@mui/material';

export const UIControls: React.FC = () => {
    return (
        <Card sx={{ p: 3, borderRadius: 4, boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>
                UI Elements Showcase
            </Typography>

            <Typography variant="subtitle2" sx={{ mb: 2, fontWeight: 600 }}>Buttons & Variations</Typography>
            <Stack direction="row" spacing={1} flexWrap="wrap" useFlexGap sx={{ mb: 4 }}>
                <Button variant="contained">Primary</Button>
                <Button variant="contained" color="secondary">Secondary</Button>
                <Button variant="outlined">Outlined</Button>
                <Button variant="text">Text Only</Button>
                <Button variant="contained" disabled>Disabled</Button>
            </Stack>

            <Divider sx={{ my: 3, borderStyle: 'dashed' }} />

            <Typography variant="subtitle2" sx={{ mb: 2, fontWeight: 600 }}>Toggles & Selectors</Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1, mb: 4 }}>
                <FormControlLabel control={<Switch defaultChecked />} label="Active Status" />
                <FormControlLabel control={<Checkbox defaultChecked />} label="Email Notifications" />
                <RadioGroup row defaultValue="a">
                    <FormControlLabel value="a" control={<Radio />} label="Option A" />
                    <FormControlLabel value="b" control={<Radio />} label="Option B" />
                </RadioGroup>
            </Box>

            <Divider sx={{ my: 3, borderStyle: 'dashed' }} />

            <Typography variant="subtitle2" sx={{ mb: 2, fontWeight: 600 }}>Data Range Range</Typography>
            <Box sx={{ px: 1 }}>
                <Slider
                    defaultValue={50}
                    aria-label="Default"
                    valueLabelDisplay="auto"
                    sx={{ color: 'primary.main' }}
                />
            </Box>

            <Box sx={{ mt: 4, p: 2, borderRadius: 2, bgcolor: (theme) => alpha(theme.palette.info.main, 0.05), border: (theme) => `1px solid ${alpha(theme.palette.info.main, 0.1)}` }}>
                <Typography variant="caption" color="info.main" sx={{ display: 'block', fontWeight: 600 }}>
                    TIP: All components are fully themed and accessible.
                </Typography>
            </Box>
        </Card>
    );
};
