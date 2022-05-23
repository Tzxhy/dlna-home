import './App.css';

import {
    Button,
    Card,
    Input,
    List,
    Modal,
    PullToRefresh,
    Selector,
    TextArea,
    Toast,
} from 'antd-mobile';
import {
    SelectorOption,
} from 'antd-mobile/es/components/selector';
import {
    useEffect,
    useRef,
    useState,
} from 'react';

import {
    actionRemote,
    createPlayList,
    deletePlayListApi,
    getDeviceList, getPlayList, PlayList, updatePlayListApi,
} from './api';


function App() {
    const [devices, setDevices] = useState<SelectorOption<string>[]>([]);
    const [playlist, setPlaylist] = useState<SelectorOption<string>[]>([]);
    const [playListSelected, setPlayListSelected] = useState<string[]>([]);
    const [deviceSelected, setDeviceSelected] = useState<string[]>([]);
    const playListRef = useRef<PlayList['list']>();
    const [playListDetail, setPlayListDetail] = useState<PlayList['list'][number]['list']>([]);

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
        const v = arr[0];
        if (!playListRef.current) return;
        setPlayListSelected(arr);
        const item = playListRef.current.find(i => i.pid === v)!;
        setPlayListDetail(item.list);
    };

    const addPlayList = () => {
        let playListName = '';
        let playDetail = '';
        Modal.confirm({
            content: <div>
                <Input
                    placeholder='请输入名称'
                    onChange={val => {
                        playListName = val;
                    }}
                />
                <TextArea
                    placeholder='请输入歌曲列表'
                    onChange={val => {
                        playDetail = val;
                    }}
                />
            </div>,
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
        let playListName = '';
        let playDetail = '';
        Modal.confirm({
            content: <div>
                <Input
                    placeholder='请输入名称'
                    onChange={val => {
                        playListName = val;
                    }}
                />
                <TextArea
                    placeholder='请输入歌曲列表'
                    onChange={val => {
                        playDetail = val;
                    }}
                />
            </div>,
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
            content: `确认删除${playlist.find(i => i.value === playListSelected[0])!.label}?`,
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
        if (!playListSelected.length) {
            Toast.show({
                icon: 'fail',
                content: '先选择播放列表',
            });
            return;
        }
        actionRemote(playListSelected[0], 'stop', deviceSelected[0]);
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
                        onChange={(arr) => setDeviceSelected(arr)}
                    />
                </Card>
                <Card title='播放列表选择'>
                    <Selector
                        value={playListSelected}
                        options={playlist}
                        onChange={onPlayListChange}
                    />
                    <Button style={{
                        marginTop: 16,
                    }} color='primary' block onClick={addPlayList}>创建播放列表</Button>
                    <Button style={{
                        marginTop: 16,
                    }} color='primary' block onClick={updatePlayList}>更新列表内容</Button>
                    <Button style={{
                        marginTop: 16,
                    }} color='danger' block onClick={deletePlayList}>删除播放列表</Button>
                </Card>
                <Card title='执行操作'>
                    <Button color='success' onClick={play}>播放</Button>
                    <Button color='danger' onClick={stop}>停止</Button>
                </Card>
                <Card title='当前播放列表歌单'>
                    <List header={'共' + playListDetail.length + '首'}>
                        {
                            playListDetail.map(i => <List.Item key={i.aid}>{i.name}</List.Item>)
                        }
                    </List>
                </Card>
            </PullToRefresh>

        </>
    );
}

export default App;
