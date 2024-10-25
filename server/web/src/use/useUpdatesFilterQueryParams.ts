import UpdateFilterQueryParamNames from '../constants/api/updateFilterQueryParamNames';
import { PAGE_DEFAULT, PAGE_SIZE_DEFAULT } from '../constants/pagination';
import { compact, isEmpty, parseInt, uniq } from 'lodash';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';

export interface UseUpdatesFilterQueryParams {
	searchTerm: string | undefined | number;
	searchIn: string | undefined | number;
	state: string[] | undefined;
	orderBy: undefined | number | string;
	order: undefined | number | string;
	page: number;
	pageSize: number;
}

const useUpdatesFilterQueryParams = (): UseUpdatesFilterQueryParams => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | number | string): undefined | number | string => {
		return val ? val : undefined;
	}, []);

	const replaceEmptyValue = useCallback((val: string[]) => {
		return isEmpty(val) ? undefined : uniq(compact(val));
	}, []);

	return {
		searchTerm: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.SEARCH_TERM)),
		searchIn: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.SEARCH_IN)),
		state: replaceEmptyValue(queryParams.getAll(UpdateFilterQueryParamNames.STATE)),
		orderBy: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.ORDER)),
		page: queryParams.get(UpdateFilterQueryParamNames.PAGE)
			? parseInt(queryParams.get(UpdateFilterQueryParamNames.PAGE) as string)
			: PAGE_DEFAULT,
		pageSize: queryParams.get(UpdateFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(UpdateFilterQueryParamNames.PAGE_SIZE) as string)
			: PAGE_SIZE_DEFAULT
	};
};

export default useUpdatesFilterQueryParams;
