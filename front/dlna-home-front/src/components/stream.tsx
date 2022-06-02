import PlayCircleFilledWhiteIcon from '@mui/icons-material/PlayCircleFilledWhite';
import {
    Box,
    IconButton,
    TextField,
} from '@mui/material';
import {
    Container,
} from '@mui/system';
import {
    useContext,
    useState,
} from 'react';

import {
    startStream,
} from '../api';
import {
    showSnackbar,
} from '../plugin/snackbar.tsx';
import AppContext from '../store';
import TabHeader from './tabHeader';


export default function Stream() {
    const ctx = useContext(AppContext);
    const store = ctx[0];

    const [url, setUrl] = useState('');

    const clickPlay = async () => {
        if (!store.currentDevice.url) {
            showSnackbar({
                type: 'error',
                title: '请先选择设备',
                timeout: 3000,
            });
            return;
        }
        if (!url) {
            showSnackbar({
                type: 'error',
                title: '请填入url',
                timeout: 3000,
            });
            return;
        }
        const ok = await startStream(store.currentDevice.url, url);
        showSnackbar({
            type: ok ? 'success' : 'error',
            title: ok ? '请求成功' : '请求失败',
            timeout: 3000,
        });
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
    </Container>;
}
