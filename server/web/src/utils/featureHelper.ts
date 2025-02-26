import getConfiguration from '../getConfiguration';

export const footerEnabled = (): boolean => {
	return !!(getConfiguration().VITE_ENABLE_FOOTER && getConfiguration().VITE_ENABLE_FOOTER > 0);
};

export const darkThemeEnabled = (): boolean => {
	return !!(getConfiguration().VITE_ENABLE_DARK_THEME && getConfiguration().VITE_ENABLE_DARK_THEME > 0);
};
