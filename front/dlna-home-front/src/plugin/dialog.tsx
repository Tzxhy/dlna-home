import {
    Button,Dialog, DialogActions, DialogContent, DialogTitle,
} from '@mui/material';
import ReactDOM from 'react-dom/client';


type DialogConf = {
    title?: string;
    body?: React.ReactElement;
    onOk?: () => void;
    onCancel?: () => void;
}
function MyDialog(props: DialogConf) {
    return <Dialog open={true} onClose={props.onCancel}>
        32131212312
        {
            props.title ? <DialogTitle>修改名称</DialogTitle> : null
        }
        <DialogContent>
            {
                props.body ? props.body : null
            }
        </DialogContent>
        <DialogActions>
            <Button onClick={props.onCancel}>取消</Button>
            <Button onClick={props.onOk}>确定</Button>
        </DialogActions>
    </Dialog>;
}

export function showDialog(params: DialogConf) {
    const div = document.createElement('div');
    div.setAttribute('type', '111');
    const root = ReactDOM.createRoot(div);
    const onClose = () => {
        params.onCancel?.();
        root.unmount();
        div.remove();
    };
    root.render(
        <MyDialog {...params} onCancel={onClose} />
    );
    root.unmount();
    document.body.appendChild(div);
}
