import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';
import eslint from 'vite-plugin-eslint2';
import stylelint from 'vite-plugin-stylelint';
import svgrPlugin from 'vite-plugin-svgr';
import viteTsconfigPaths from 'vite-tsconfig-paths';

const stylelintOptions = {
	dev: true
};

const eslintOptions = {
	dev: true
};

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
	const env = loadEnv(mode, process.cwd(), '');
	return {
		plugins: [react(), viteTsconfigPaths(), svgrPlugin(), stylelint(stylelintOptions), eslint(eslintOptions)],
		base: './',
		server: {
			open: false,
			host: env.VITE_DEV_HOST,
			port: parseInt(env.VITE_DEV_PORT),
			proxy: {
				'/api': {
					target: env.VITE_DEV_PROXY_TARGET,
					changeOrigin: true
				}
			}
		},
		build: {
			outDir: 'build'
		},
		// @ts-ignore
		test: {
			globals: true,
			environment: 'jsdom',
			coverage: {
				provider: 'v8',
				reporter: ['text', 'html'],
				exclude: ['node_modules/']
			}
		}
	};
});
