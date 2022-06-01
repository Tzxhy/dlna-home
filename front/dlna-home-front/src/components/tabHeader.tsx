
import ArrowBackIosNewIcon from '@mui/icons-material/ArrowBackIosNew';
import {
    Box,
    IconButton,
    Typography,
} from '@mui/material';
import React from 'react';

type TabHeaderProps = {
    title: string;
    note?: string;
    ext?: React.ReactElement;
    showBack?: boolean;
    onClickBack?: (e: Event) => void;
}
export default function TabHeader(props: TabHeaderProps) {

    return <Box>
        <Box sx={{
            display: 'flex',
            justifyContent: 'space-between',
            pb: 1,
            borderBottom: 1,
            borderColor: 'divider',
        }}>
            <Box sx={{
                display: 'flex',
                alignItems: 'center',
            }}>
                {
                    // @ts-ignore
                    props.showBack ? <IconButton onClick={props.onClickBack}><ArrowBackIosNewIcon /></IconButton> : null
                }
                <Typography variant="h3" sx={{
                    fontSize: 20,
                    color: 'text.primary',
                }}>{props.title}</Typography>
            </Box>
            {
                props.ext ? props.ext : null
            }
        </Box>
        {
            props.note ? <Typography variant="body1" sx={{
                mt: 1,
                fontSize: 16,
                pb: 1,
                color: 'text.secondary',
            }}>{props.note}</Typography> : null
        }
    </Box>;
}
