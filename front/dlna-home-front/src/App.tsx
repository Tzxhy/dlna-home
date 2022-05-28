import './App.css';

import {
    AutoCenter,
    Button,
    Card,
    Form,
    Input,
    List,
    Modal,
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
    PlayOutline, SoundMuteOutline, SoundOutline, StopOutline,
} from 'antd-mobile-icons';
import {
    useEffect,
    useRef,
    useState,
} from 'react';

import {
    actionRemote,
    createPlayList,
    deletePlayListApi,
    getDeviceList, getPlayList, getVolumeApi, PlayList, setVolumeApi, updatePlayListApi,
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
    }, []);

    const onPlayListChange = (arr: string[]) => {
        if (!playListRef.current) return;
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

        actionRemote(playListSelected[0], 'start', deviceSelected[0]);
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
        setDeviceSelected(arr);
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


    return (
        <>
            <PullToRefresh
                onRefresh={refresh}
            >
                <Card title='设备选择'>
                    <Selector
                        value={deviceSelected}
                        options={devices}
                        onChange={(arr) => onDeviceChanged(arr)}
                    />
                </Card>
                <Card title='播放列表选择'>
                    <Space block direction='vertical'>
                        <Space block direction='vertical'>
                            <Button color='primary' block onClick={addPlayList}>创建播放列表</Button>
                        </Space>
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
                            <Space wrap block direction='vertical'>
                                {
                                    playListSelected.length > 0 ?
                                        <Button block onClick={play}><PlayOutline />播放</Button> : null
                                }
                                <Button block onClick={stop}><StopOutline />停止</Button>
                                <AutoCenter>当前音量：{volume}</AutoCenter>
                                <Button block onClick={addVolume}><SoundMuteOutline />音量增</Button>
                                <Button block onClick={subVolume}><SoundOutline />音量减</Button>
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
