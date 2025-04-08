/// <reference types="vite/client" />

interface ImportMetaEnv {
	readonly VITE_API_URL: string;
	readonly VITE_TITLE: string;
	readonly VITE_ENABLE_DARK_THEME: number;
	readonly VITE_ENABLE_FOOTER: number;
}

interface ImportMeta {
	readonly env: ImportMetaEnv;
}
