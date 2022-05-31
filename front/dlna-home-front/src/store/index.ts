import React from 'react';

export const defaultStore = {
	devices: [
		{
			name: '小爱同学1',
			url: 'http://xiaoai1',
		},
		{
			name: '小爱同学2',
			url: 'http://xiaoai2',
		},
	]
}

export const reducer = () => {

}

const AppContext = React.createContext(defaultStore);

export default AppContext;

