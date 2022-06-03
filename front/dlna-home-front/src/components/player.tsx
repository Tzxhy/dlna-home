import ArrowBackIosNewRoundedIcon from '@mui/icons-material/ArrowBackIosNewRounded';
import FavoriteIcon from '@mui/icons-material/Favorite';
import FavoriteBorderIcon from '@mui/icons-material/FavoriteBorder';
import FormatListBulletedOutlinedIcon from '@mui/icons-material/FormatListBulletedOutlined';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import PauseCircleOutlineOutlinedIcon from '@mui/icons-material/PauseCircleOutlineOutlined';
import PlayCircleFilledWhiteOutlinedIcon from '@mui/icons-material/PlayCircleFilledWhiteOutlined';
import RepeatOneOutlinedIcon from '@mui/icons-material/RepeatOneOutlined';
import RepeatOutlinedIcon from '@mui/icons-material/RepeatOutlined';
import ShuffleOutlinedIcon from '@mui/icons-material/ShuffleOutlined';
import SkipNextOutlinedIcon from '@mui/icons-material/SkipNextOutlined';
import SkipPreviousOutlinedIcon from '@mui/icons-material/SkipPreviousOutlined';
import IconButton from '@mui/material/IconButton';
import Slider from '@mui/material/Slider';
import Typography from '@mui/material/Typography';
import Box from '@mui/system/Box';
import Container from '@mui/system/Container';
import {
    memo,
    useContext,
    useState,
} from 'react';

import bg from '../assets/img/music.jpg';
import AppContext from '../store/index.ts';

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
        }}><ArrowBackIosNewRoundedIcon /></IconButton>
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
const RepeatIcon = memo(function(props: {type: number}) {
    const type = props.type;
    if (type === 1) {
        return <RepeatOneOutlinedIcon />;
    }
    if (type === 2) {
        return <RepeatOutlinedIcon />;
    }
    if (type === 3) {
        return <ShuffleOutlinedIcon />;
    }
});

export default function Player() {
    const ctx = useContext(AppContext);
    const store = ctx[0];
    const dispatch = ctx[1];
    const show = store.player.show;

    const status = store.player.status;

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
        }}>
            <Box>
                <Slider
                    size="small"
                    defaultValue={70}
                    aria-label="Small"
                    valueLabelDisplay="auto"
                />

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
                    <IconButton onClick={() => {

                    }}><RepeatIcon type={3} /></IconButton>
                    <IconButton><SkipPreviousOutlinedIcon /></IconButton>
                    <IconButton size='large'>{
                        status === 'stop' ?
                            <PlayCircleFilledWhiteOutlinedIcon fontSize='large' />
                            : <PauseCircleOutlineOutlinedIcon fontSize='large' />
                    }</IconButton>
                    <IconButton><SkipNextOutlinedIcon /></IconButton>
                    <IconButton><FormatListBulletedOutlinedIcon /></IconButton>

                </Box>
                <Box sx={{
                    height: 100,
                    width: 30,
                }}>
                    <Slider
                        aria-label="Temperature"
                        orientation="vertical"
                        defaultValue={30}
                        valueLabelDisplay="auto"
                    />
                </Box>
            </Box>


        </Box>

    </Container>;
}
