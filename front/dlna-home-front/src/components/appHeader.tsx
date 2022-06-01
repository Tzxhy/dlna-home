import Brightness4Icon from '@mui/icons-material/Brightness4';
import Brightness7Icon from '@mui/icons-material/Brightness7';
import AppBar from '@mui/material/AppBar';
import IconButton from '@mui/material/IconButton';
import {
    useTheme,
} from '@mui/material/styles';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import {
    useContext,
} from 'react';

import {
    ColorModeContext,
} from '../context/theme';

export default function MyAppBar() {
    const theme = useTheme();
    const colorMode = useContext(ColorModeContext);
    return <AppBar position="static">
        <Toolbar variant="dense" sx={{
            display: 'flex',
            justifyContent: 'space-between',
        }}>
            <Typography variant="h6" color="inherit">媒体推送</Typography>
            <IconButton sx={{
                ml: 1,
            }} onClick={colorMode.toggleColorMode} color="inherit">
                {theme.palette.mode === 'dark' ? <Brightness7Icon /> : <Brightness4Icon />}
            </IconButton>
        </Toolbar>
    </AppBar>;
}
