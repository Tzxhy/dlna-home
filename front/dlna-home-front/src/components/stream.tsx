import PlayCircleFilledWhiteIcon from '@mui/icons-material/PlayCircleFilledWhite';
import {
    Alert,
    AlertColor,
    Box,
    IconButton,
    Snackbar,
    TextField,
} from '@mui/material';
import {
    Container,
} from '@mui/system';
import {
    useContext,
    useRef,
    useState,
} from 'react';

import {
    startStream,
} from '../api';
import AppContext from '../store';
import TabHeader from './tabHeader';


export default function Stream() {
    const ctx = useContext(AppContext);
    const store = ctx[0];

    const [url, setUrl] = useState('');
    const [showBar, setShowBar] = useState(false);
    const barRef = useRef<{
        type: AlertColor;
        tips: string;
    }>({
        type: 'error',
        tips: '',
    });

    const clickPlay = async () => {
        if (!store.currentDevice.url) {
            barRef.current = {
                type: 'error',
                tips: '请先选择设备',
            };
            setShowBar(true);
            return;
        }
        if (!url) {
            barRef.current = {
                type: 'error',
                tips: '请填入url',
            };
            setShowBar(true);
            return;
        }
        const ok = await startStream(store.currentDevice.url, url);
        setShowBar(true);
        if (ok) {
            barRef.current = {
                type: 'success',
                tips: '请求成功',
            };
        } else {
            barRef.current = {
                type: 'error',
                tips: '请求失败',
            };
        }
    };

    return <Container>
        <TabHeader title='直播' note='直播只能播放一个文件，且不能操作媒体设备。注意，url中通常不能包含 "&"符号，大概率会失效（看具体设备，已知小爱音响不识别&符）' />
        <Box sx={{
            mt: 4,
            display: 'flex',
            alignItems: 'center',
            flexDirection: {
                mobile: 'column',
                tablet: 'row',
            },
        }}>
            <TextField
                sx={{
                    flex: 4,
                }}
                required
                onChange={(e) => setUrl(e.target.value)}
                label="资源地址"
                defaultValue=""
            />
            <IconButton sx={{
                flex: 1,
            }} onClick={clickPlay}><PlayCircleFilledWhiteIcon />播放</IconButton>
        </Box>
        <Snackbar
            open={showBar}
            autoHideDuration={6000}
            onClose={() => setShowBar(false)}
        >
            <Alert onClose={() => setShowBar(false)} severity={barRef.current.type} sx={{
                width: '100%',
            }}>
                {
                    barRef.current.tips
                }
            </Alert>
        </Snackbar>
    </Container>;
}
