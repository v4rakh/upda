declare global {
	interface Window {
		runtime_config: Configuration;
	}
}

interface Configuration {
	VITE_API_URL: string;
	VITE_BASE_PATH: string;
	VITE_TITLE: string;
	VITE_ENABLE_DARK_THEME: number;
	VITE_ENABLE_FOOTER: number;
}

/**
 * Derive configuration values depending on environment:
 * - load from process.env in case of development
 * - load from window object otherwise
 *
 * development check must be fully written as isDevelopment uses configuration itself
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
