import axios from 'axios';
export type DeviceList = {
    data: Record<string, string>
}
export async function getDeviceList(): Promise<DeviceList> {
    return axios.get<any>('/api/v1/device-list').then(d => d.data);
}
export type PlayList = {
    list: {
        pid: string;
        name: string;
        create_date: number;
        list: {
            pid: string;
            aid: string;
            url: string;
            name: string;
            create_date: number;
        }[]
    }[]
}
export async function getPlayList(): Promise<PlayList> {
    return axios.get<any>('/api/v1/playlist').then(d => d.data);
}

export async function createPlayList(name: string, list: {name: string; url: string}[]): Promise<void> {
    const create = await axios.post<any>('/api/v1/create-playlist', {
        name,
    }).then(d => d.data);
    await updatePlayListApi(create.pid, name, list);
}
export async function updatePlayListApi(pid: string, name: string, list: {name: string; url: string}[]): Promise<void> {
    await axios.post('/api/v1/update-playlist', {
        pid: pid,
        name,
        list,
    });
}
export async function deletePlayListApi(pid: string): Promise<void> {
    await axios.post<any>('/api/v1/delete-playlist', {
        pid,
    }).then(d => d.data);
}

export async function actionRemote(pid: string, action_name: 'start' | 'stop', renderer_url: string): Promise<void> {
    await axios.post<any>('/api/v1/action', {
        pid,
        action_name,
        renderer_url,
    }).then(d => d.data);
}

export async function setVolumeApi(renderer_url: string, level: number): Promise<boolean> {
    return axios.post<any>('/api/v1/volume', {
        renderer_url,
        level,
    }).then(d => {
        return d.data?.ok ?? false as boolean;
    });
}

export async function getVolumeApi(renderer_url: string): Promise<{
    ok: boolean;
	level: number;
}> {
    return axios.get<any>('/api/v1/volume', {
        params: {
            renderer_url,
        },
    }).then(d => d.data);
}
