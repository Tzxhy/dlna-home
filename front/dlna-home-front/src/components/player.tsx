import {
    VolumeDownRounded, VolumeUpRounded,
} from '@mui/icons-material';
import FavoriteIcon from '@mui/icons-material/Favorite';
// import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import FormatListBulletedOutlinedIcon from '@mui/icons-material/FormatListBulletedOutlined';
import KeyboardArrowDownRoundedIcon from '@mui/icons-material/KeyboardArrowDownRounded';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import PauseCircleOutlineOutlinedIcon from '@mui/icons-material/PauseCircleOutlineOutlined';
import PlayCircleFilledWhiteOutlinedIcon from '@mui/icons-material/PlayCircleFilledWhiteOutlined';
import RepeatOneOutlinedIcon from '@mui/icons-material/RepeatOneOutlined';
import RepeatOutlinedIcon from '@mui/icons-material/RepeatOutlined';
import ShuffleOutlinedIcon from '@mui/icons-material/ShuffleOutlined';
import SkipNextOutlinedIcon from '@mui/icons-material/SkipNextOutlined';
import SkipPreviousOutlinedIcon from '@mui/icons-material/SkipPreviousOutlined';
import {
    Stack, useTheme,
} from '@mui/material';
import IconButton from '@mui/material/IconButton';
import Slider from '@mui/material/Slider';
import Typography from '@mui/material/Typography';
import {
    styled,
} from '@mui/system';
import Box from '@mui/system/Box';
import Container from '@mui/system/Container';
import {
    memo,
    useContext,
    useEffect,
    useRef,
} from 'react';

import {
    changePlayModeApi,
    getPositionApi,
    getStatusApi,
    nextSongApi,
    pauseSong,
    PlayMode,
    playSong,
    prevSongApi,
    setVolumeApi,
} from '../api';
import bg from '../assets/img/music.jpg';
import AppContext from '../store';

const Header = memo(function() {
    const ctx = useContext(AppContext);
    const dispatch = ctx[1];
    return <Box sx={{
        my: 2,
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        borderBottom: 1,
        borderColor: 'divider',
    }}>
        <IconButton onClick={() => {
            dispatch({
                type: 'hide-player',
            });
        }}><KeyboardArrowDownRoundedIcon /></IconButton>
        <Box sx={{
            flex: 1,
            textAlign: 'center',
        }}><Typography variant='h3' sx={{
                color: 'text.primary',
                fontSize: 20,
                fontWeight: 600,
            }}>音乐播放器</Typography></Box>
        <IconButton onClick={() => {
        }} sx={{
            p: 0,
        }}><MoreHorizIcon /></IconButton>
    </Box>;
});

// const (
// 	PLAY_MODE_SEQ         = iota // 顺序播放
// 	PLAY_MODE_REPEAT_ONE         // 单曲循环
// 	PLAY_MODE_LIST_REPEAT        // 列表循环
// 	PLAY_MODE_RANDOM             // 乱序播放
// )
const RepeatIcon = memo(function(props: {type: PlayMode}) {
    const type = props.type;
    if (type === PlayMode.PLAY_MODE_REPEAT_ONE) {
        return <RepeatOneOutlinedIcon />;
    }
    if (type === PlayMode.PLAY_MODE_LIST_REPEAT) {
        return <RepeatOutlinedIcon />;
    }
    if (type === PlayMode.PLAY_MODE_RANDOM) {
        return <ShuffleOutlinedIcon />;
    }
    return <ShuffleOutlinedIcon />;
});

const PlayModeList = [
    PlayMode.PLAY_MODE_REPEAT_ONE,
    PlayMode.PLAY_MODE_LIST_REPEAT,
    PlayMode.PLAY_MODE_RANDOM,
];

function formatDuration(value: number) {
    const minute = Math.floor(value / 60);
    const secondLeft = value - minute * 60;
    return `${minute}:${secondLeft < 9 ? `0${secondLeft}` : secondLeft}`;
}

const TinyText = styled(Typography)(
    ({
        theme,
    }) => ({
        color: theme.palette.text.primary,
        fontSize: '0.75rem',
        fontWeight: 500,
        letterSpacing: 0.2,
    })
);

export default function Player() {
    const ctx = useContext(AppContext);
    const store = ctx[0];
    const dispatch = ctx[1];
    const show = store.player.show;
    const position = store.player.position;

    const status = store.player.status;
    const device = store.currentDevice?.url;
    const deviceRef = useRef<string>('');
    deviceRef.current = device;
    const currentPlayMode = store.player.mode;

    async function getPlayerStatus() {
        const data = (await getStatusApi()).data;
        Object.keys(data).forEach(key => {
            if (key === deviceRef.current) {

                dispatch({
                    type: 'update-player',
                    data: {
                        currentItem: {
                            ...data[key as keyof typeof data].current_item,
                        },
                    },
                });
            }
        });

    }
    async function getPosition() {
        const position = await getPositionApi(deviceRef.current);
        dispatch({
            type: 'update-position',
            data: position.position,
        });
    }
    useEffect(() => {
        getPlayerStatus();

        const i = setInterval(() => {
            if (document.visibilityState === 'hidden') return;
            getPlayerStatus();
            getPosition();
        }, 5000);
        return () => {
            clearInterval(i);
        };
    }, []);

    const theme = useTheme();

    const lightIconColor = theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.4)' : 'rgba(0,0,0,0.4)';

    return <Container sx={{
        display: 'flex',
        flexDirection: 'column',
        position: 'fixed',
        left: 0,
        top: show ? 0 : '100vh',
        height: '100%',
        right: 0,
        zIndex: 9,
        bgcolor: 'background.paper',
        transition: 'top .3s',
    }}>
        <Header />
        <img src={bg} width="100%" />
        <Box sx={{
            flex: 1,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
        }}>

            <Typography variant='h3' sx={{
                color: 'text.primary',
                fontSize: 18,
                fontWeight: 500,
            }}>{store.player.currentItem?.name}</Typography>
            <IconButton><FavoriteIcon /></IconButton>

        </Box>

        <Box sx={{
            pb: 2,
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
        }}>
            <Slider
                aria-label="time-indicator"
                size="small"
                value={position.rel_time}
                min={0}
                step={1}
                max={position.track_duration}
                // onChange={(_, value) => dispatch({
                //     type: 'update-position',
                //     data: {
                //         rel_time: value,
                //         track_duration: position.track_duration,
                //     },
                // })}
                sx={{
                    color: theme.palette.mode === 'dark' ? '#fff' : 'rgba(0,0,0,0.87)',
                    height: 4,
                    '& .MuiSlider-thumb': {
                        width: 8,
                        height: 8,
                        transition: '0.3s cubic-bezier(.47,1.64,.41,.8)',
                        '&:before': {
                            boxShadow: '0 2px 12px 0 rgba(0,0,0,0.4)',
                        },
                        '&:hover, &.Mui-focusVisible': {
                            boxShadow: `0px 0px 0px 8px ${
                                theme.palette.mode === 'dark'
                                    ? 'rgb(255 255 255 / 16%)'
                                    : 'rgb(0 0 0 / 16%)'
                            }`,
                        },
                        '&.Mui-active': {
                            width: 20,
                            height: 20,
                        },
                    },
                    '& .MuiSlider-rail': {
                        opacity: 0.28,
                    },
                }}
            />

            <Box
                sx={{
                    width: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    mt: -2,
                }}
            >
                <TinyText>{formatDuration(position.rel_time)}</TinyText>
                <TinyText>-{formatDuration(position.track_duration - position.rel_time)}</TinyText>
            </Box>
        </Box>
        <Box sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
        }}>
            <Box sx={{
                display: 'flex',
                flex: 1,
                alignItems: 'center',
                justifyContent: 'space-around',
            }}>
                <IconButton onClick={async () => {
                    const cIdx = PlayModeList.indexOf(currentPlayMode);
                    const nIdx = (cIdx + 1) % PlayModeList.length;
                    const nMode = PlayModeList[nIdx];
                    await changePlayModeApi(device, nMode);
                    dispatch({
                        type: 'update-player',
                        data: {
                            mode: nMode,
                        },
                    });
                }}><RepeatIcon type={currentPlayMode} /></IconButton>
                <IconButton onClick={async () => {
                    await prevSongApi(device);
                    await getPlayerStatus();
                }}><SkipPreviousOutlinedIcon /></IconButton>
                <IconButton size='large' onClick={() => {
                    if (status === 'stop') {
                        playSong(device);
                        dispatch({
                            type: 'player-play',
                        });
                    } else {
                        pauseSong(device);
                        dispatch({
                            type: 'player-pause',
                        });
                    }
                }}>{
                        status === 'stop' ?
                            <PlayCircleFilledWhiteOutlinedIcon fontSize='large' />
                            : <PauseCircleOutlineOutlinedIcon fontSize='large' />
                    }</IconButton>
                <IconButton onClick={async () => {
                    await nextSongApi(device);
                    await getPlayerStatus();
                }}><SkipNextOutlinedIcon /></IconButton>
                <IconButton><FormatListBulletedOutlinedIcon /></IconButton>
            </Box>

        </Box>
        <Box sx={{
            width: '100%',
            pb: 4,
        }}>
            <Stack spacing={2} direction="row" sx={{
                mb: 1,
                px: 1,
            }} alignItems="center">
                <VolumeDownRounded htmlColor={lightIconColor} />
                <Slider
                    aria-label="Volume"
                    defaultValue={30}
                    onChangeCommitted={(e, v) => {
                        setVolumeApi(device, v as any as number);
                    }}
                    sx={{
                        color: theme.palette.mode === 'dark' ? '#fff' : 'rgba(0,0,0,0.87)',
                        '& .MuiSlider-track': {
                            border: 'none',
                        },
                        '& .MuiSlider-thumb': {
                            width: 24,
                            height: 24,
                            backgroundColor: '#fff',
                            '&:before': {
                                boxShadow: '0 4px 8px rgba(0,0,0,0.4)',
                            },
                            '&:hover, &.Mui-focusVisible, &.Mui-active': {
                                boxShadow: 'none',
                            },
                        },
                    }}
                />
                <VolumeUpRounded htmlColor={lightIconColor} />
            </Stack>
        </Box>
    </Container>;
}
