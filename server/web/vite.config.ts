/// <reference types="vitest" />
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';
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
export default defineConfig({
	plugins: [react(), viteTsconfigPaths(), svgrPlugin(), stylelint(stylelintOptions), eslint(eslintOptions)],
	server: {
		open: false,
		host: '127.0.0.1',
		port: 3001
	},
	build: {
		outDir: 'build'
	},
	// @ts-ignore
	test: {
		globals: true,
		environment: 'jsdom',
		setupFiles: './src/setupTests.ts',
		coverage: {
			reporter: ['text', 'html'],
			exclude: ['node_modules/', 'src/setupTests.ts']
		}
	}
});
