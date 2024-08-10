import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import reactControlStatements from 'vite-plugin-react-control-statements';
import { pigment } from '@pigment-css/vite-plugin';
import { TanStackRouterVite } from '@tanstack/router-plugin/vite';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        TanStackRouterVite(),
        react(), 
        reactControlStatements(),
        pigment({
            babelOptions:{
                compact: false
            }
        })
    ],
})
