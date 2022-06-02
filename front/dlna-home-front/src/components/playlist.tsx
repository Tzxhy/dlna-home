
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline';
import PlaylistAddIcon from '@mui/icons-material/PlaylistAdd';
import {
    Button,
    Divider,
    Menu,
    MenuItem,
    TextField,
    Typography,
} from '@mui/material';
import IconButton from '@mui/material/IconButton';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import Box from '@mui/system/Box';
import Container from '@mui/system/Container';
import {
    useContext,
    useEffect,
    useRef,
    useState,
} from 'react';

import {
    actionRemote,
    createPlayList,
    deletePlayListApi,
    deleteSingleResourceApi,
    getPlayList,
    renamePlayList,
    startPlayResource,
    updatePlayListApi,
} from '../api';
import {
    PlayList,
} from '../api/index.ts';
import {
    showDialog,
} from '../plugin/dialog';
import AppContext from '../store';
import TabHeader from './tabHeader';

export default function Playlist() {

    const ctx = useContext(AppContext);
    const store = ctx[0];
    const dispatch = ctx[1];
    const refreshPlayList = async () => {
        const data = await getPlayList();
        dispatch({
            type: 'set-play-list',
            data: data.list,
        });
    };
    useEffect(() => {
        refreshPlayList();
    }, []);

    const [currentViewPlayList, setCurrentViewPlayList] = useState('');

    const playlist = (<><TabHeader
        title='播放列表'
        note='针对播放列表的操作。可以新增，修改，删除等'
        ext={<IconButton sx={{
            p: 0,
        }} onClick={() => {
            const formRef = {
                name: '',
                list: '',
            };
            showDialog({
                title: '新增播放列表',
                body: <>
                    <TextField
                        autoFocus
                        margin="dense"
                        label="名称"
                        type="text"
                        fullWidth
                        variant="standard"
                        onChange={e => formRef.name = e.target.value}
                    />
                    <TextField
                        margin="dense"
                        label="播放列表"
                        type="text"
                        multiline
                        maxRows={4}
                        fullWidth
                        variant="standard"
                        onChange={e => formRef.list = e.target.value}
                    />
                </>,
                onOk: async (close) => {
                    const {
                        name,
                        list,
                    } = formRef;
                    if (!name || !list) {
                        setFormErrorTips('完善表格');
                        return;
                    }
                    let listObj: {name: string; url: string}[];
                    try {
                        listObj = JSON.parse(list);
                    } catch(e) {
                        return;
                    }
                    await createPlayList(name, listObj);
                    await refreshPlayList();
                    close();
                },
            });
        }}><PlaylistAddIcon /></IconButton>}
    />
    <Typography variant='body1' sx={{
        mt: 2,
        color: 'text.secondary',
    }}>当前共有：{store.playList.length}个列表</Typography>
    <List>
        {
            store.playList.map((i, idx) => {
                return <Box key={i.pid} sx={{
                    color: 'text.primary',
                }}>
                    <ListItem disablePadding onClick={() => {
                        setCurrentViewPlayList(i.pid);
                    }}>
                        <ListItemButton>
                            <ListItemText primary={(idx + 1) + '. ' + i.name} />
                        </ListItemButton>

                    </ListItem>
                    <Divider light />
                </Box>;
            })
        }

    </List></>);

    const getTargetPlayItem = (pid: string) => {
        return store.playList.find(i => i.pid === pid)!;
    };

    const requestPlayStart = async () => {
        if (!store.currentDevice.url) {
            return;
        }
        await actionRemote(currentViewPlayList, 'start', store.currentDevice.url);
    };

    const [showMoreActionAnchor, setShowMoreActionAnchor] = useState<HTMLElement | null>(null);
    const [showItemMoreAnchor, setShowItemMoreAnchor] = useState<HTMLElement | null>(null);
    const showItemMoreAnchorRef = useRef<PlayList['list']['list'][number]>();


    const playListDetail = (currentViewPlayList: string) => (<>
        <TabHeader
            showBack
            onClickBack={() => {
                setCurrentViewPlayList('');
            }}
            title={getTargetPlayItem(currentViewPlayList).name}
            ext={<>
                <IconButton onClick={(e) => {
                    console.log('e: ', e);
                    setShowMoreActionAnchor(e.target as HTMLElement);
                }} sx={{
                    p: 0,
                }}><MoreHorizIcon /></IconButton>
                <Menu
                    anchorEl={showMoreActionAnchor}
                    open={!!showMoreActionAnchor}
                    onClose={() => {
                        setShowMoreActionAnchor(null);
                    }}
                >
                    <MenuItem onClick={() => {
                        setShowMoreActionAnchor(null);
                        let newPlayListName = getTargetPlayItem(currentViewPlayList).name;

                        showDialog({
                            title: '重命名',
                            body: <TextField
                                autoFocus
                                margin="dense"
                                label="新的名称"
                                type="text"
                                defaultValue={newPlayListName}
                                onChange={(e) => newPlayListName = e.target.value}
                                fullWidth
                                variant="standard"
                            />,
                            onOk: async (close) => {
                                if (!newPlayListName) return;
                                await renamePlayList(currentViewPlayList, newPlayListName);
                                refreshPlayList();
                                close();
                            },
                        });
                    }}>重命名</MenuItem>
                    <MenuItem onClick={() => {
                        setShowMoreActionAnchor(null);
                        showDialog({
                            title: '确认删除' + `"${getTargetPlayItem(currentViewPlayList).name}"?`,
                            onOk: async (close) => {
                                await deletePlayListApi(currentViewPlayList);
                                close();
                                setCurrentViewPlayList('');
                                setTimeout(() => {
                                    refreshPlayList();
                                });
                            },
                        });
                    }}>删除</MenuItem>
                    <MenuItem onClick={() => {
                        setShowMoreActionAnchor(null);

                        let newPlayListName = '';
                        let newPlayListList = '';

                        showDialog({
                            title: '替换更新',
                            body: <>
                                <TextField
                                    autoFocus
                                    margin="dense"
                                    label="名称"
                                    type="text"
                                    defaultValue={newPlayListName}
                                    onChange={(e) => newPlayListName = e.target.value}
                                    fullWidth
                                    variant="standard"
                                />
                                <TextField
                                    margin="dense"
                                    label="列表"
                                    type="text"
                                    defaultValue={newPlayListList}
                                    multiline
                                    maxRows={4}
                                    onChange={(e) => newPlayListList = e.target.value}
                                    fullWidth
                                    variant="standard"
                                />
                            </>,
                            onOk: async (close) => {
                                if (!newPlayListName || !newPlayListList) return;
                                let listObj: {name: string; url: string}[];
                                try {
                                    listObj = JSON.parse(newPlayListList);
                                } catch(e) {
                                    return;
                                }
                                await updatePlayListApi(currentViewPlayList, newPlayListName, listObj);
                                refreshPlayList();
                                close();
                            },
                        });
                    }}>替换更新</MenuItem>
                    <MenuItem >增量添加</MenuItem>
                </Menu>
            </>}
        />
        <Button
            onClick={requestPlayStart}
            variant="text" startIcon={<PlayCircleOutlineIcon />}>播放当前列表</Button>
        <Typography variant='body1' sx={{
            mt: 2,
            color: 'text.secondary',
        }}>当前列表共有：{getTargetPlayItem(currentViewPlayList).list.length}个媒体</Typography>

        <List>
            {
                getTargetPlayItem(currentViewPlayList).list.map((i, idx) => {
                    return <Box key={i.aid} sx={{
                        color: 'text.primary',
                    }}>
                        <ListItem disablePadding>
                            <ListItemButton onClick={(e) => {
                                showItemMoreAnchorRef.current = i;
                                setShowItemMoreAnchor(e.target);
                            }}>
                                <ListItemText primary={(idx + 1) + '. ' + i.name} />
                            </ListItemButton>

                        </ListItem>
                        <Divider light />
                    </Box>;
                })
            }

        </List>
        <Menu
            anchorEl={showItemMoreAnchor}
            open={!!showItemMoreAnchor}
            onClose={() => {
                setShowItemMoreAnchor(null);
            }}
        >
            <MenuItem onClick={() => {
                setShowItemMoreAnchor(null);

                showDialog({
                    title: `确认删除${showItemMoreAnchorRef.current.name}`,
                    onOk: async (close) => {
                        await deleteSingleResourceApi(showItemMoreAnchorRef.current.aid);
                        refreshPlayList();
                        close();
                    },
                });
            }}>删除</MenuItem>
            <MenuItem onClick={() => {
                setShowItemMoreAnchor(null);

                startPlayResource(store.currentDevice.url, currentViewPlayList, showItemMoreAnchorRef.current.aid);
            }}>从这里播放</MenuItem>

        </Menu>
    </>);
    return <Container>
        {
            currentViewPlayList ? playListDetail(currentViewPlayList) : playlist
        }
    </Container>;
}
