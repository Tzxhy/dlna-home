import './App.css';

import {
    AutoCenter,
    Button,
    Card,
    FloatingBubble,
    Form,
    Input,
    List,
    Modal,
    Popup,
    PullToRefresh,
    Selector,
    Space,
    TextArea,
    Toast,
} from 'antd-mobile';
import {
    SelectorOption,
} from 'antd-mobile/es/components/selector';
import {
    MovieOutline,
    PlayOutline, SoundMuteOutline, SoundOutline, StopOutline,
} from 'antd-mobile-icons';
import {
    useEffect,
    useRef,
    useState,
} from 'react';

import {
    actionRemote,
    changePlayModeApi,
    createPlayList,
    deletePlayListApi,
    getDeviceList,
    getPlayList,
    getStatusApi,
    getVolumeApi,
    nextSongApi,
    PlayList,
    prevSongApi,
    setVolumeApi,
    updatePlayListApi,
} from './api';


function App() {
    const [devices, setDevices] = useState<SelectorOption<string>[]>([]);
    const [playlist, setPlaylist] = useState<SelectorOption<string>[]>([]);
    const [playListSelected, setPlayListSelected] = useState<string[]>([]);
    const [deviceSelected, setDeviceSelected] = useState<string[]>([]);
    const playListRef = useRef<PlayList['list']>();
    const [playListDetail, setPlayListDetail] = useState<PlayList['list'][number]['list']>([]);

    const [volume, setVolume] = useState(-1);

    const refresh = async () => {
        setDevices([]);
        setPlaylist([]);
        setPlayListSelected([]);
        setDeviceSelected([]);
        setPlayListDetail([]);
        const device = (await getDeviceList()).data;
        const newList = Object.keys(device).map(deviceName => ({
            label: deviceName,
            value: device[deviceName] as string,
        }));
        setDevices(
            newList
        );

        const play = (await getPlayList()).list;
        playListRef.current = play;
        const newPlayList = play.map(item => ({
            label: item.name,
            value: item.pid,
        }));
        setPlaylist(newPlayList);
    };
    useEffect(() => {
        refresh();
        getStatusApi();
    }, []);

    const onPlayListChange = (arr: string[]) => {
        if (!playListRef.current || !arr.length) return;
        const v = arr[0];
        setPlayListSelected(arr ?? []);
        const item = playListRef.current.find(i => i.pid === v)!;
        setPlayListDetail(arr.length ? item.list : []);
    };

    const addPlayList = () => {
        let playListName = '';
        let playDetail = '';
        Modal.confirm({
            content: <Form
                layout='vertical'
            >
                <Form.Item
                    name='name'
                    label='名称'
                    rules={[{
                        required: true,
                        message: '不能为空',
                    }]}
                >
                    <Input
                        placeholder='请输入名称'
                        onChange={val => {
                            playListName = val;
                        }}
                    />
                </Form.Item>
                <Form.Item
                    name='list'
                    label='列表'
                    rules={[{
                        required: true,
                        message: '不能为空',
                    }]}
                >
                    <TextArea
                        placeholder='请输入歌曲列表'
                        onChange={val => {
                            playDetail = val;
                        }}
                    />
                </Form.Item>
            </Form>,
            onConfirm: async () => {
                let list = [];
                try {
                    list = JSON.parse(playDetail);
                } catch(e) {
                    Toast.show({
                        icon: 'fail',
                        content: '解析失败',
                    });
                    return;
                }
                await createPlayList(playListName, list);
                refresh();
            },
        });
    };
    const updatePlayList = () => {
        if (!playListSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择播放列表',
            });
            return;
        }
        const c = playListSelected[0];
        const originalItem = playListRef.current!.find(i => i.pid === c)!;
        const originalName = originalItem.name ?? '';
        let playListName = originalName;
        let playDetail = JSON.stringify(originalItem.list, null, 4);
        Modal.confirm({
            content: <Form
                layout='vertical'
            >
                <Form.Item
                    name='name'
                    label='名称'
                    rules={[{
                        required: true,
                        message: '不能为空',
                    }]}
                >
                    <Input
                        placeholder='请输入名称'
                        defaultValue={playListName}
                        onChange={val => {
                            playListName = val;
                        }}
                    />
                </Form.Item>
                <Form.Item
                    name='list'
                    label='列表'
                    rules={[{
                        required: true,
                        message: '不能为空',
                    }]}
                >
                    <TextArea
                        placeholder='请输入歌曲列表'
                        defaultValue={playDetail}
                        onChange={val => {
                            playDetail = val;
                        }}
                    />
                </Form.Item>
            </Form>,
            onConfirm: async () => {
                let list = [];
                try {
                    list = JSON.parse(playDetail);
                } catch(e) {
                    Toast.show({
                        icon: 'fail',
                        content: '解析失败',
                    });
                    return;
                }
                await updatePlayListApi(playListSelected[0], playListName, list);
                refresh();
            },
        });
    };

    const deletePlayList = () => {
        if (!playListSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择播放列表',
            });
            return;
        }

        Modal.confirm({
            content: `确认删除“${playlist.find(i => i.value === playListSelected[0])!.label}”?`,
            onConfirm: async () => {
                deletePlayListApi(playListSelected[0]);
                refresh();
            },
        });
    };

    const play = () => {
        if (!deviceSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择设备',
            });
            return;
        }
        if (!playListSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择播放列表',
            });
            return;
        }

        actionRemote(playListSelected[0], 'start', deviceSelected[0], playModeSelected[0]);
    };

    const stop = () => {
        if (!deviceSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择设备',
            });
            return;
        }
        actionRemote('', 'stop', deviceSelected[0]);
    };

    const _setVolumeApi = async (level: number) => {
        const d = await setVolumeApi(deviceSelected[0], level);
        if (d) {
            setVolume(level);
        }
    };

    const subVolume = async () => {
        let newV = volume - 5;
        if (newV < 5) newV = 5;
        _setVolumeApi(newV);
    };

    const onDeviceChanged = async(arr: string[]) => {
        if (!arr.length) return;
        setDeviceSelected(arr);
        setDeviceChooserVisible(false);
        if (arr.length) {
            const d = await getVolumeApi(arr[0]);
            if (d?.ok) {
                setVolume(d.level);
            } else {
                setVolume(-1);
            }
        } else {
            setVolume(-1);
        }
    };

    const addVolume = async () => {
        let newV = volume + 5;
        if (newV > 100) newV = 100;
        _setVolumeApi(newV);
    };

    const nextSong = () => {
        nextSongApi(deviceSelected[0]);
    };

    const prevSong = () => {
        prevSongApi(deviceSelected[0]);
    };

    const [playModeSelected, setPlayModeSelected] = useState<number[]>([0]);

    const onPlayModeChanged = (arr: number[]) => {
        if (!arr.length) return;
        setPlayModeSelected(arr);
        changePlayModeApi(deviceSelected[0], arr[0]);
    };

    const playModeOptions = [{
        label: '顺序播放',
        value: 0,
    }, {
        label: '单曲循环',
        value: 1,
    }, {
        label: '列表循环',
        value: 2,
    }, {
        label: '乱序播放',
        value: 3,
    }];

    const [deviceChooserVisible, setDeviceChooserVisible] = useState(false);

    return (
        <>
            <PullToRefresh
                onRefresh={refresh}
            >
                <FloatingBubble
                    axis='xy'
                    magnetic='x'
                    style={{
                        '--initial-position-bottom': '24px',
                        '--initial-position-right': '24px',
                        '--edge-distance': '24px',
                    }}
                >
                    <MovieOutline fontSize={32} onClick={() => setDeviceChooserVisible(true)} />
                </FloatingBubble>
                <Popup
                    visible={deviceChooserVisible}
                    onMaskClick={() => {
                        setDeviceChooserVisible(false);
                    }}
                    bodyStyle={{
                        borderTopLeftRadius: '8px',
                        borderTopRightRadius: '8px',
                        minHeight: '20vh',
                    }}
                >
                    <Card title='设备选择'>
                        <Selector
                            value={deviceSelected}
                            options={devices}
                            onChange={(arr) => onDeviceChanged(arr)}
                        />
                    </Card>
                </Popup>

                <Card title='播放列表选择' extra={<Button size='small' onClick={addPlayList}>创建播放列表</Button>}>
                    <Space block direction='vertical'>

                        <Selector
                            value={playListSelected}
                            options={playlist}
                            onChange={onPlayListChange}
                        />
                        {
                            playListSelected[0] ? <Space direction='vertical' block>
                                <Button block onClick={updatePlayList}>更新列表内容</Button>
                                <Button block onClick={deletePlayList}>删除播放列表</Button>
                            </Space> : null
                        }
                    </Space>

                </Card>
                {
                    deviceSelected.length > 0 ? <>
                        <Card title='执行操作'>
                            <Space wrap block direction='horizontal'>
                                <Selector
                                    value={playModeSelected}
                                    options={playModeOptions}
                                    onChange={(arr) => onPlayModeChanged(arr)}
                                />
                                {
                                    playListSelected.length > 0 ?
                                        <Button onClick={play}><PlayOutline />播放</Button> : null
                                }
                                <Button onClick={stop}><StopOutline />停止</Button>
                                <Button onClick={prevSong}><SoundMuteOutline />上一曲</Button>
                                <Button onClick={nextSong}><SoundOutline />下一曲</Button>
                                <Space wrap block align='center'>
                                    <AutoCenter>当前音量：{volume}</AutoCenter>
                                    <Button onClick={addVolume}><SoundMuteOutline />音量增</Button>
                                    <Button onClick={subVolume}><SoundOutline />音量减</Button>
                                </Space>


                            </Space>

                        </Card>
                        <Card title='当前播放列表歌单'>
                            <List header={'共' + playListDetail.length + '首'}>
                                {
                                    playListDetail.map(i => <List.Item key={i.aid}>{i.name}</List.Item>)
                                }
                            </List>
                        </Card>
                    </> : null
                }
            </PullToRefresh>

        </>
    );
}

export default App;
