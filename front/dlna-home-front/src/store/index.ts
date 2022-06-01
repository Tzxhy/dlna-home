import React from 'react';

import {
    PlayList,
} from '../api';

export const defaultStore = {
    currentDevice: {
        name: '',
        url: '',
    },
    devices: [
        {
            name: '小爱同学1',
            url: 'http://xiaoai1',
        },
        {
            name: '小爱同学2',
            url: 'http://xiaoai2',
        },
    ],
    playList: [
        {
            pid: '',
            name: '所有',
            create_date: 0,
            list: [],
        }
    ] as PlayList['list'],
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

    default:
        return state;
    }

};
