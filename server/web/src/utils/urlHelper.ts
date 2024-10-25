export const getAppBasePath = () => {
	return '/';
};

export const getPageFullPath = (path: string) => {
	return `${getAppBasePath()}${path}`;
};
