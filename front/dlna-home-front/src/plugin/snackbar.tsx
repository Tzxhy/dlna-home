import {
    Alert,
    Snackbar,
} from '@mui/material';
import {
    useState,
} from 'react';


type ElementSnackBarConf = {
    show: boolean;
    type: 'success' | 'error';
    title: string;
    timeout: number;
    onClose: () => void;
}
function MySnackBar(props: ElementSnackBarConf) {

    return <Snackbar open={props.show} autoHideDuration={props.timeout} onClose={props.onClose}>
        <Alert onClose={props.onClose} severity={props.type} sx={{
            width: '100%',
        }}>
            {props.title}
        </Alert>
    </Snackbar>;
}

let close: () => void;
let show: () => void;
let props: ElementSnackBarConf;

export function MySnackBarAdapter() {
    const [_show, setShow] = useState(false);
    close = () => setShow(false);
    show = () => setShow(true);
    return <MySnackBar {...props} show={_show} />;
}

type SnackBarConf = {
    type: 'error' | 'success';
    title: string;
    timeout?: number;
}
export function showSnackbar(params: SnackBarConf) {
    // @ts-ignore
    props = {
        ...params,
        onClose: () => {
            close();
        },
    };
    show();
    return close;
}
