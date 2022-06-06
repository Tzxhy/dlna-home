import axios from 'axios';
export type DeviceList = {
    data: Record<string, string>
}
export async function getDeviceList(): Promise<DeviceList> {
    return axios.get<any>('/api/v1/device-list').then(d => d.data);
}
export async function startStream(renderer_url: string, url: string): Promise<boolean> {
    return axios.post<any>('/api/v1/start-one', {
        renderer_url,
        url,
    }).then(d => !!d.data?.ok);
}
export type AudioItem = {
    pid: string;
    aid: string;
    url: string;
    name: string;
    create_date: number;
}
export type PlayList = {
    list: {
        pid: string;
        name: string;
        create_date: number;
        list: AudioItem[]
    }[]
}
export async function getPlayList(): Promise<PlayList> {
    return axios.get<any>('/api/v1/playlist').then(d => d.data);
}

export async function updatePlayListApi(pid: string, name: string, list: {name: string; url: string}[]): Promise<void> {
    await axios.post('/api/v1/update-playlist', {
        pid: pid,
        name,
        list,
    });
}

export async function renamePlayList(pid: string, new_name: string): Promise<void> {
    await axios.post('/api/v1/rename-playlist', {
        pid: pid,
        new_name,
    });
}


export async function createPlayList(name: string, list: {name: string; url: string}[]): Promise<void> {
    const create = await axios.post<any>('/api/v1/create-playlist', {
        name,
    }).then(d => d.data);
    await updatePlayListApi(create.pid, name, list);
}


export async function actionRemote(
    pid: string,
    action_name: 'start' | 'stop', renderer_url: string, play_mode = PlayMode.PLAY_MODE_RANDOM): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        pid,
        action_name,
        renderer_url,
        play_mode,
    }).then(d => d.data);
}

export async function deletePlayListApi(pid: string): Promise<void> {
    await axios.post<any>('/api/v1/delete-playlist', {
        pid,
    }).then(d => d.data);
}
export async function deleteSingleResourceApi(aid: string): Promise<void> {
    await axios.post<any>('/api/v1/delete-single-resource', {
        aid,
    }).then(d => d.data);
}

export async function startPlayResource(renderer_url: string, pid: string, aid: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        action_name: 'jump',
        pid,
        aid,
        renderer_url,
    }).then(d => d.data);
}

type StatusResp = Record<string, {
    data: Record<string, {
        status: 'play' | 'stop';
        renderer_url: string;
        current_item: AudioItem;
    }>;
}>
export async function getStatusApi(): Promise<StatusResp> {
    return axios.get<any>('/api/v1/get-status?' + Date.now()).then(d => d.data);
}

export async function setVolumeApi(renderer_url: string, level: number): Promise<boolean> {
    return axios.post<any>('/api/v1/volume', {
        renderer_url,
        level,
    }).then(d => {
        return d.data?.ok ?? false as boolean;
    });
}

export async function nextSongApi(renderer_url: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        renderer_url,
        action_name: 'next',
    });
}
export async function prevSongApi(renderer_url: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        renderer_url,
        action_name: 'previous',
    });
}

export async function playSong(renderer_url: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        renderer_url,
        action_name: 'play',
    });
}

export async function pauseSong(renderer_url: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        renderer_url,
        action_name: 'pause',
    });
}

export enum PlayMode {
    PLAY_MODE_SEQ,
	PLAY_MODE_REPEAT_ONE, // 单曲循环
	PLAY_MODE_LIST_REPEAT, // 列表循环
	PLAY_MODE_RANDOM,
}
export async function changePlayModeApi(renderer_url: string, mode: PlayMode): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        renderer_url,
        action_name: 'changePlayMode',
        play_mode: mode,
    });
}
