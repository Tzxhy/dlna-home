import './index.css';

import React, { useReducer, useState } from 'react';
import ReactDOM from 'react-dom/client';

import App from './App';
import Store, {defaultStore} from './store';
function Main() {
	const [o, updateStore] = useState(defaultStore);
	const [state, dispatch] = useReducer(() => {}, o)
	return <React.StrictMode>
	<Store.Provider value={[state, dispatch]}>
		<App />
	</Store.Provider>
</React.StrictMode>
}
ReactDOM.createRoot(document.getElementById('root')!).render(
    <Main />
);
