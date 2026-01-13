import React from 'react';
import {
    Grid,
    Card,
    Typography,
    TextField,
    Box,
    FormControl,
    InputLabel,
    Select,
    MenuItem,
    Autocomplete,
    Rating,
    ToggleButtonGroup,
    ToggleButton,
    Slider,
    Switch,
    Checkbox,
    RadioGroup,
    FormControlLabel,
    Radio,
    Stack,
    Skeleton,
    alpha
} from '@mui/material';
import FormatAlignLeftIcon from '@mui/icons-material/FormatAlignLeft';
import FormatAlignCenterIcon from '@mui/icons-material/FormatAlignCenter';
import FormatAlignRightIcon from '@mui/icons-material/FormatAlignRight';

export const InputShowcase: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    const top100Films = [
        { label: 'The Shawshank Redemption', year: 1994 },
        { label: 'The Godfather', year: 1972 },
        { label: 'The Dark Knight', year: 2008 },
    ];

    if (isLoading) {
        return (
            <Grid container spacing={3}>
                {[...Array(6)].map((_, i) => (
                    <Grid key={i} size={{ xs: 12, md: 6 }}>
                        <Card sx={{ p: 3, borderRadius: 4 }}>
                            <Skeleton variant="text" width="40%" height={32} sx={{ mb: 2 }} />
                            <Skeleton variant="rectangular" height={56} sx={{ borderRadius: 1, mb: 2 }} />
                            <Skeleton variant="rectangular" height={56} sx={{ borderRadius: 1 }} />
                        </Card>
                    </Grid>
                ))}
            </Grid>
        );
    }

    return (
        <Grid container spacing={3}>
            {/* Text Fields */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Text Fields</Typography>
                    <Stack spacing={3}>
                        <TextField fullWidth label="Standard TextField" variant="outlined" />
                        <TextField fullWidth label="With Helper Text" helperText="Please enter your name" />
                        <TextField fullWidth label="Password" type="password" />
                        <TextField fullWidth label="Error State" error helperText="Incorrect entry." />
                    </Stack>
                </Card>
            </Grid>

            {/* Select & Autocomplete */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Select & Choice</Typography>
                    <Stack spacing={3}>
                        <FormControl fullWidth>
                            <InputLabel>Age</InputLabel>
                            <Select label="Age" defaultValue={20}>
                                <MenuItem value={10}>Ten</MenuItem>
                                <MenuItem value={20}>Twenty</MenuItem>
                                <MenuItem value={30}>Thirty</MenuItem>
                            </Select>
                        </FormControl>
                        <Autocomplete
                            disablePortal
                            options={top100Films}
                            renderInput={(params) => <TextField {...params} label="Movie Search" />}
                        />
                        <Box>
                            <Typography component="legend" variant="body2" color="text.secondary">Rating</Typography>
                            <Rating name="simple-controlled" defaultValue={4} />
                        </Box>
                    </Stack>
                </Card>
            </Grid>

            {/* Toggles & Selection */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Checkboxes & Radios</Typography>
                    <Grid container spacing={2}>
                        <Grid size={{ xs: 6 }}>
                            <Typography variant="subtitle2" sx={{ mb: 1 }}>Checkboxes</Typography>
                            <Box sx={{ display: 'flex', flexDirection: 'column' }}>
                                <FormControlLabel control={<Checkbox defaultChecked />} label="Option 1" />
                                <FormControlLabel control={<Checkbox />} label="Option 2" />
                            </Box>
                        </Grid>
                        <Grid size={{ xs: 6 }}>
                            <Typography variant="subtitle2" sx={{ mb: 1 }}>Radio Group</Typography>
                            <RadioGroup defaultValue="female">
                                <FormControlLabel value="female" control={<Radio />} label="Female" />
                                <FormControlLabel value="male" control={<Radio />} label="Male" />
                            </RadioGroup>
                        </Grid>
                    </Grid>
                    <Divider sx={{ my: 3 }} />
                    <Typography variant="subtitle2" sx={{ mb: 1 }}>Switches</Typography>
                    <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 2 }}>
                        <FormControlLabel control={<Switch defaultChecked />} label="Label 1" />
                        <FormControlLabel control={<Switch color="secondary" />} label="Label 2" />
                    </Box>
                </Card>
            </Grid>

            {/* Advanced Inputs */}
            <Grid size={{ xs: 12, md: 6 }}>
                <Card sx={{ p: 3, borderRadius: 4, height: '100%', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>Advanced Controls</Typography>
                    <Typography variant="subtitle2" sx={{ mb: 2 }}>Sliders</Typography>
                    <Slider defaultValue={30} aria-label="Default" valueLabelDisplay="auto" sx={{ mb: 4 }} />
                    <Slider defaultValue={20} disabled aria-label="Disabled" sx={{ mb: 4 }} />

                    <Typography variant="subtitle2" sx={{ mb: 2 }}>Toggle Buttons</Typography>
                    <ToggleButtonGroup value="left" exclusive aria-label="text alignment">
                        <ToggleButton value="left" aria-label="left aligned">
                            <FormatAlignLeftIcon />
                        </ToggleButton>
                        <ToggleButton value="center" aria-label="centered">
                            <FormatAlignCenterIcon />
                        </ToggleButton>
                        <ToggleButton value="right" aria-label="right aligned">
                            <FormatAlignRightIcon />
                        </ToggleButton>
                    </ToggleButtonGroup>
                </Card>
            </Grid>
        </Grid>
    );
};

// Internal Divider for clean layout
const Divider = ({ sx }: { sx?: any }) => (
    <Box sx={{ height: '1px', bgcolor: (theme) => alpha(theme.palette.divider, 0.1), ...sx }} />
);
