import {
    Button,Dialog, DialogActions, DialogContent, DialogTitle,
} from '@mui/material';
import {
    useState,
} from 'react';


type ElementDialogConf = {
    show: boolean;
    title?: string | React.ReactElement;
    body?: string | React.ReactElement | (() => React.ReactElement);
    onOk?: () => void;
    onCancel?: () => void;
	showBtns?: boolean;
}
function MyDialog(props: ElementDialogConf) {
    const showBtns = props.showBtns ?? true;

    return <Dialog open={props.show} onClose={props.onCancel}>
        {
            props.title ? <DialogTitle>{props.title}</DialogTitle> : null
        }
        <DialogContent>
            {
                props.body ? (typeof props.body === 'function' ? props.body() : props.body) : null
            }
        </DialogContent>
        {
            showBtns ? <DialogActions>
                <Button onClick={props.onCancel}>取消</Button>
                <Button onClick={props.onOk}>确定</Button>
            </DialogActions> : null
        }
    </Dialog>;
}

let close: () => void;
let show: () => void;
let props: DialogConf;
export function MyDialogAdapter() {
    const [_show, setShow] = useState(false);
    close = () => setShow(false);
    show = () => setShow(true);
    // @ts-ignore
    return <MyDialog show={_show} {...props} />;
}

type DialogConf = {
    title?: string | React.ReactElement;
    body?: string | React.ReactElement | (() => React.ReactElement);
    onOk?: (callback: ()=>void) => void;
    onCancel?: () => void;
	showBtns?: boolean;
}
export function showDialog(params: DialogConf) {
    props = {
        ...params,
        onCancel: () => {
            params.onCancel?.();
            close();
        },
        onOk: () => {
            const cb = () => {
                close();
            };
            params.onOk?.(cb);
        },
    };
    show();
    return close;
}
