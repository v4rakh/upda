import AuthType from './auth/AuthType';

declare global {
	interface Window {
		runtime_config: Configuration;
	}
}

interface Configuration {
	VITE_API_URL: string;
	VITE_BASE_PATH: string;
	VITE_TITLE: string;
	VITE_ENABLE_FOOTER: number;
	VITE_AUTH_TYPE: AuthType;
}

/**
 * Derive configuration values depending on environment:
 * - load from vite env in case of development
 * - load from frozen window object otherwise
 */
const getConfiguration = (): Configuration => {
	if (window && window.runtime_config && Object.keys(window.runtime_config).length > 0) {
		return window.runtime_config;
	} else if (import.meta.env.DEV) {
		return import.meta.env as unknown as Configuration;
	}

	throw new Error('Cannot bootstrap configuration from window or environment');
};

export default getConfiguration;
