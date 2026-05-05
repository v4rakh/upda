import getConfiguration from '../getConfiguration';

export const footerEnabled = (): boolean => {
	return !!(getConfiguration().VITE_ENABLE_FOOTER && getConfiguration().VITE_ENABLE_FOOTER > 0);
};
