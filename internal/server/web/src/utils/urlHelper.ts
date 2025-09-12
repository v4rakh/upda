import getConfiguration from '../getConfiguration';
import { forEach } from 'lodash';

export const getAppBasePath = () => {
	return getConfiguration().VITE_BASE_PATH;
};

export const getPageFullPath = (path: string) => {
	return `${getAppBasePath()}${path}`;
};

export const getPageFullPathWithQueryParameters = (path: string, queryParameters: Record<string, string>) => {
	const params = new URLSearchParams();
	forEach(queryParameters, (value, key) => {
		if (key && value) {
			params.append(key, value);
		}
	});
	return `${getPageFullPath(path)}?${params.toString()}`;
};
