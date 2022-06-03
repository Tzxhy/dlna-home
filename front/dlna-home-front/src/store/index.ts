
import React from 'react';

import {
    PlayList,
} from '../api';

export const defaultStore = {
    currentDevice: {
        name: '',
        url: '',
    },
    devices: [] as {name: string; url: string;}[],
    playList: [] as PlayList['list'],
    player: {
        status: 'stop',
        show: true,
        currentListPid: '',
        currentItem: {
            aid: 'BAqptfQX',
            create_date: 1654165880476,
            name: '八扇屏.mp3',
            pid: 'mRWYIQMR',
            url: 'http://192.168.0.108:8080/api/v1/share/download?sid=wrNEAogo&fid=UaXcxeWWZf',
        },
    },
};

const AppContext = React.createContext<[typeof defaultStore, React.Dispatch<Action>]>([defaultStore] as any);

export default AppContext;


type Action = {
    type: string;
    data: any;
}
export const reducer = (state: typeof defaultStore, action: Action) => {
    switch (action.type) {
    case 'set-device-list':
        return {
            ...state,
            devices: action.data,
        };
    case 'set-device':
        return {
            ...state,
            currentDevice: action.data,
        };
    case 'set-play-list':
        return {
            ...state,
            playList: action.data,
        };
    case 'hide-player':
        return {
            ...state,
            player: {
                ...state.player,
                show: false,
                ...action.data,
            },
        };
    case 'show-player':
        return {
            ...state,
            player: {
                ...state.player,
                show: true,
            },
        };
    default:
        return state;
    }

};
