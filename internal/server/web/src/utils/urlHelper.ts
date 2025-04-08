import getConfiguration from '../getConfiguration';

export const getAppBasePath = () => {
	return getConfiguration().VITE_BASE_PATH;
};

export const getPageFullPath = (path: string) => {
	return `${getAppBasePath()}${path}`;
};
