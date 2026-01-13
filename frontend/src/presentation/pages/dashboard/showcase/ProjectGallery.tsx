import React from 'react';
import {
    Grid,
    Card,
    CardMedia,
    CardContent,
    Typography,
    Box,
    LinearProgress,
    AvatarGroup,
    Avatar,
    Chip,
    Skeleton,
    alpha
} from '@mui/material';

interface Project {
    id: number;
    title: string;
    category: string;
    progress: number;
    status: 'In Progress' | 'Completed' | 'On Hold';
    image: string;
    members: string[];
}

const projects: Project[] = [
    {
        id: 1,
        title: 'Modern UI Revamp',
        category: 'Design',
        progress: 75,
        status: 'In Progress',
        image: 'https://images.unsplash.com/photo-1558655146-d09347e92766?w=500&auto=format&fit=crop&q=60',
        members: ['JD', 'AS', 'MK']
    },
    {
        id: 2,
        title: 'Mobile App API',
        category: 'Development',
        progress: 100,
        status: 'Completed',
        image: 'https://images.unsplash.com/photo-1551033406-611cf9a28f67?w=500&auto=format&fit=crop&q=60',
        members: ['BW', 'CL']
    },
    {
        id: 3,
        title: 'Q1 Marketing Campaign',
        category: 'Marketing',
        progress: 30,
        status: 'On Hold',
        image: 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=500&auto=format&fit=crop&q=60',
        members: ['EM', 'SP', 'RT', 'KL']
    },
];

const getStatusColor = (status: string) => {
    switch (status) {
        case 'Completed': return 'success';
        case 'In Progress': return 'primary';
        case 'On Hold': return 'warning';
        default: return 'default';
    }
};

export const ProjectGallery: React.FC<{ isLoading: boolean }> = ({ isLoading }) => {
    return (
        <Box>
            <Typography variant="h6" sx={{ fontWeight: 700, mb: 3 }}>
                Active Projects
            </Typography>
            <Grid container spacing={3}>
                {(isLoading ? [...Array(3)] : projects).map((project, index) => (
                    <Grid key={project?.id || index} size={{ xs: 12, sm: 6, md: 4 }}>
                        <Card sx={{ borderRadius: 4, overflow: 'hidden', boxShadow: '0 4px 20px rgba(0,0,0,0.05)' }}>
                            {isLoading ? (
                                <Skeleton variant="rectangular" height={140} />
                            ) : (
                                <CardMedia
                                    component="img"
                                    height="140"
                                    image={project.image}
                                    alt={project.title}
                                />
                            )}
                            <CardContent sx={{ p: 3 }}>
                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                                    {isLoading ? (
                                        <Skeleton variant="text" width={60} />
                                    ) : (
                                        <Typography variant="overline" color="text.secondary" sx={{ fontWeight: 700 }}>
                                            {project.category}
                                        </Typography>
                                    )}
                                    {isLoading ? (
                                        <Skeleton variant="rounded" width={80} height={20} />
                                    ) : (
                                        <Chip
                                            label={project.status}
                                            size="small"
                                            color={getStatusColor(project.status)}
                                            sx={{ fontWeight: 600, fontSize: '0.65rem' }}
                                        />
                                    )}
                                </Box>

                                {isLoading ? (
                                    <Skeleton variant="text" width="90%" height={32} />
                                ) : (
                                    <Typography variant="h6" sx={{ fontWeight: 700, mb: 2 }}>
                                        {project.title}
                                    </Typography>
                                )}

                                <Box sx={{ mb: 3 }}>
                                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 0.5 }}>
                                        <Typography variant="caption" color="text.secondary">Progress</Typography>
                                        <Typography variant="caption" sx={{ fontWeight: 700 }}>{isLoading ? '...' : project.progress}%</Typography>
                                    </Box>
                                    {isLoading ? (
                                        <Skeleton variant="rectangular" height={6} sx={{ borderRadius: 1 }} />
                                    ) : (
                                        <LinearProgress
                                            variant="determinate"
                                            value={project.progress}
                                            sx={{ height: 6, borderRadius: 3, bgcolor: (theme) => alpha(theme.palette.divider, 0.1) }}
                                        />
                                    )}
                                </Box>

                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                    <Typography variant="caption" color="text.secondary">Team Members</Typography>
                                    {isLoading ? (
                                        <Skeleton variant="circular" width={24} height={24} />
                                    ) : (
                                        <AvatarGroup max={3} sx={{ '& .MuiAvatar-root': { width: 24, height: 24, fontSize: '0.75rem' } }}>
                                            {project.members.map((m: string, i: number) => (
                                                <Avatar key={i}>{m}</Avatar>
                                            ))}
                                        </AvatarGroup>
                                    )}
                                </Box>
                            </CardContent>
                        </Card>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
};
