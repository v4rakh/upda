import { compact, isEmpty, uniq } from 'lodash';

export const replaceNullValue = (val: null | number | string): undefined | number | string => {
	return val ? val : undefined;
};

export const replaceEmptyValue = (val: string[]) => {
	return isEmpty(val) ? undefined : uniq(compact(val));
};
