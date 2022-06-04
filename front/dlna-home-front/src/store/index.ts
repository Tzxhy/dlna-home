
import React from 'react';

import {
    AudioItem,
    PlayList, PlayMode,
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
        show: false,
        currentListPid: '',
        currentItem: {} as AudioItem,
        mode: PlayMode.PLAY_MODE_RANDOM,
    } as {
        status: 'stop' | 'play';
        show: boolean;
        currentItem: AudioItem;
        mode: PlayMode;
    },
};

const AppContext = React.createContext<[typeof defaultStore, React.Dispatch<Action>]>([defaultStore] as any);

export default AppContext;


type Action = {
    type: string;
    data?: any;
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
                ...action.data,
            },
        };
    case 'update-player':
        return {
            ...state,
            player: {
                ...state.player,
                ...action.data,
            },
        };
    case 'player-play':
        return {
            ...state,
            player: {
                ...state.player,
                status: 'play',
            },
        };
    case 'player-pause':
        return {
            ...state,
            player: {
                ...state.player,
                status: 'stop',
            },
        };
    default:
        return state;
    }

};
