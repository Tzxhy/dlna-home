
import './App.css';

import ConnectedTvIcon from '@mui/icons-material/ConnectedTv';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import PlaylistPlayIcon from '@mui/icons-material/PlaylistPlay';
import RefreshRoundedIcon from '@mui/icons-material/RefreshRounded';
import TvIcon from '@mui/icons-material/Tv';
import {
    TabPanelProps,
} from '@mui/lab';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import Box from '@mui/material/Box';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import Paper from '@mui/material/Paper';
import {
    MouseEvent,
    useContext,
    useEffect,
    useState,
} from 'react';

import {
    getDeviceList,
} from './api';
import MyAppBar from './components/appHeader';
import Player from './components/player';
import PlayerBar from './components/playerBar';
import Playlist from './components/playlist';
import Stream from './components/stream';
import {
    showDialog,
} from './plugin/dialog';
import AppContext from './store';

function TabPanel(props: Omit<TabPanelProps, 'value'> & {value: number; index: number;}) {
    const {
        children,
        value,
        index,
        sx,
    } = props;

    return (
        <Box
            sx={{
                flexGrow: 1,
                pb: '56px',
                ...sx,
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
                    {children}
                </Box>
            )}
        </Box>
    );
}

function DeviceDialog(props: {close: () => void;}) {
    const ctx = useContext(AppContext)[0];
    const dispatch = useContext(AppContext)[1];
    const currentDevice = useContext(AppContext)[0].currentDevice;
    return <List>
        {
            ctx.devices.map((i) => <ListItem key={i.url} disablePadding onClick={() => {
                dispatch({
                    type: 'set-device',
                    data: i,
                });
                setTimeout(() => {
                    props.close();
                }, 500);
            }}>
                <ListItemButton selected={i.url === currentDevice.url}>
                    <ListItemText primary={i.name} />
                </ListItemButton>
            </ListItem>)
        }
    </List>;
}

export default function App() {

    const ctx = useContext(AppContext)[0];
    const selectedDevice = !!ctx.currentDevice.url;
    const dispatch = useContext(AppContext)[1];

    const [tab, setTab] = useState(0);

    const onTabChange = (_: any, t: number) => {
        setTab(t);
    };

    const handleClickTvs = (_event: MouseEvent<HTMLButtonElement>) => {

        const close = showDialog({
            title: <>
                ??????????????? <IconButton onClick={refreshDeviceList}><RefreshRoundedIcon /></IconButton>
            </>,
            body: () => <DeviceDialog close={close} />,
            showBtns: false,
        });
    };

    async function refreshDeviceList() {
        const data = await getDeviceList();
        if (data?.data) {
            const keys = Object.keys(data.data);
            dispatch({
                type: 'set-device-list',
                data: keys.map(key => ({
                    name: key,
                    url: data.data[key],
                })),
            });

        }
    }

    useEffect(() => {
        refreshDeviceList();
    }, []);

    return <Box sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        bgcolor: 'background.paper',
    }}>
        <MyAppBar />

        <TabPanel value={tab} index={0} sx={{
            overflow: 'auto',
        }}>
            <Playlist />
        </TabPanel>
        <TabPanel value={tab} index={1} sx={{
            overflow: 'auto',
        }}>
            <Stream />
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
                <BottomNavigationAction label="????????????" icon={<PlaylistPlayIcon />} />
                <BottomNavigationAction label="??????" icon={<PlayArrowIcon />} />
            </BottomNavigation>
        </Paper>


        <IconButton sx={{
            position: 'fixed',
            right: (t: any) => t.spacing(2),
            bottom: '136px',
            borderRadius: '50%',
            border: 1,
            borderColor: 'divider',
            color: 'primary.main',
            boxShadow: 3,
            p: 1,
        }}
        onClick={handleClickTvs}
        >
            {
                selectedDevice ? <ConnectedTvIcon fontSize='large' />
                    : <TvIcon fontSize='large'/>
            }

        </IconButton>

        <PlayerBar />

        <Player />

    </Box>;
}
