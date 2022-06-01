import {
    Button,Dialog, DialogActions, DialogContent, DialogTitle,
} from '@mui/material';
import ReactDOM from 'react-dom/client';


type ElementDialogConf = {
    title?: string;
    body?: React.ReactElement;
    onOk?: () => void;
    onCancel?: () => void;
}
function MyDialog(props: ElementDialogConf) {
    return <Dialog open={true} onClose={props.onCancel}>
        {
            props.title ? <DialogTitle>{props.title}</DialogTitle> : null
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


type DialogConf = {
    title?: string;
    body?: React.ReactElement;
    onOk?: (callback: ()=>void) => void;
    onCancel?: () => void;
}
export function showDialog(params: DialogConf) {
    const div = document.createElement('div');
    const root = ReactDOM.createRoot(div);
    const onClose = () => {
        const cb = () => {

            root.unmount();
            div.remove();
        };
        params.onCancel?.(cb);
        root.unmount();
        div.remove();

    };

    const onOk = () => {
        const cb = () => {

            root.unmount();
            div.remove();
        };
        params.onOk?.(cb);
    };
    root.render(
        <MyDialog {...params} onCancel={onClose} onOk={onOk} />
    );
    document.body.appendChild(div);
}
