
import './App.css';

import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PlaylistPlayIcon from '@mui/icons-material/PlaylistPlay';
import {
    TabPanelProps,
} from '@mui/lab';

import AppBar from '@mui/material/AppBar';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import Box from '@mui/material/Box';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import Menu from '@mui/material/Menu';
import List from '@mui/material/List';
import MenuItem from '@mui/material/MenuItem';
import Paper from '@mui/material/Paper';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import {
	useContext,
    useState,
	MouseEvent,
} from 'react';

import CastConnectedRoundedIcon from '@mui/icons-material/CastConnectedRounded';
import AppContext from './store';
import Dialog from '@mui/material/Dialog/Dialog';
import DialogTitle from '@mui/material/DialogTitle/DialogTitle';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';

function TabPanel(props: Omit<TabPanelProps, 'value'> & {value: number; index: number;}) {
    const {
        children,
        value,
        index,
    } = props;

    return (
        <Box
            sx={{
                flexGrow: 1,
				pb: '56px',
            }}
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
        >
            {value === index && (
                <Box sx={{
                    p: 3,
                }}>
                    <Typography>{children}</Typography>
                </Box>
            )}
        </Box>
    );
}

export default function App() {

	const ctx = useContext(AppContext)[0]
	console.log('ctx: ', ctx);

    const [tab, setTab] = useState(0);

    const onTabChange = (_: any, t: number) => {
        setTab(t);
    };

	const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
	const open = Boolean(anchorEl);
	const handleClickTvs = (event: MouseEvent<HTMLButtonElement>) => {
		setAnchorEl(event.currentTarget);
	};
	const handleClose = (item: {name: string; url: string;}) => {
		setAnchorEl(null);
	};
    return <Box sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
    }}>
        <AppBar position="static">
            <Toolbar variant="dense">
                <Typography variant="h6" color="inherit" component="div">
                媒体推送
                </Typography>
            </Toolbar>
        </AppBar>

        <TabPanel value={tab} index={0}>
            Item One
        </TabPanel>
        <TabPanel value={tab} index={1}>
            Item Two
        </TabPanel>
        <Paper sx={{
            position: 'fixed',
            bottom: 0,
            left: 0,
            right: 0,
        }} elevation={4}>
            <BottomNavigation
                showLabels
                value={tab}
                onChange={onTabChange}
            >
                <BottomNavigationAction label="播放列表" icon={<PlaylistPlayIcon />} />
                <BottomNavigationAction label="直播" icon={<PlayArrowIcon />} />
            </BottomNavigation>
        </Paper>


		<IconButton sx={{
			position: 'fixed',
			right: t => t.spacing(2),
			bottom: '80px',
			borderRadius: '50%',
			border: 1,
			borderColor: 'grey.200',
			color: 'primary.main',
			boxShadow: 3,
			p: 1,
		}}
			onClick={handleClickTvs}
		>
			<CastConnectedRoundedIcon fontSize='large'/>
		</IconButton>
		
		<Dialog onClose={handleClose} open={open}>
      		<DialogTitle fontSize={24}>选择播放设备：</DialogTitle>
			  <List>
				{
					ctx.devices.map((i) => <ListItem key={i.url}  disablePadding onClick={() => handleClose(i)}>
					<ListItemButton>
					  <ListItemText primary={i.name} />
					</ListItemButton>
				  </ListItem>)
				}
			</List>
		</Dialog>
		
    </Box>;
}
