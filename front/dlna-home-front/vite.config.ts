import react from '@vitejs/plugin-react';
import {
    defineConfig,
} from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    build: {
        outDir: '../../assets',
    },
    server: {
        port: 3001,
        proxy: {
            '/api': {
                target: 'http://localhost:8081',
            },
        },
    },
});
