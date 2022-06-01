import './index.css';

import {
    createTheme, ThemeProvider,
} from '@mui/material/styles';
import React, {
    useMemo,
    useReducer,
} from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';
import {
    ColorModeContext,
} from './context/theme';
import AppContext, {
    defaultStore,
    reducer,
} from './store';


function Main() {

    const [state, dispatch] = useReducer(reducer, defaultStore);
    const [mode, setMode] = React.useState<'light' | 'dark'>('light');
    const colorMode = React.useMemo(
        () => ({
            toggleColorMode: () => {
                setMode((prevMode) => (prevMode === 'light' ? 'dark' : 'light'));
            },
        }),
        []
    );
    const theme = useMemo(
        () =>
            createTheme({
                palette: {
                    mode,

                },
                breakpoints: {
                    values: {
                        // @ts-ignore
                        mobile: 0,
                        tablet: 640,
                        laptop: 1024,
                        desktop: 1280,
                    },
                },
            }),
        [mode]
    );
    return <React.StrictMode>
        <ColorModeContext.Provider value={colorMode}>
            <ThemeProvider theme={theme}>
                <AppContext.Provider value={[state, dispatch]}>
                    <App />
                </AppContext.Provider>
            </ThemeProvider>
        </ColorModeContext.Provider>
    </React.StrictMode>;
}
ReactDOM.createRoot(document.getElementById('root')!).render(
    <Main />
);
