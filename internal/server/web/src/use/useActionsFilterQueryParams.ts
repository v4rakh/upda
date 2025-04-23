import ActionFilterQueryParamNames from '../constants/api/actionFilterQueryParamNames';
import { TABLE_PAGE_DEFAULT, TABLE_PAGE_SIZE_DEFAULT } from '../constants/pagination';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router';

const useActionsFilterQueryParams = () => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | number | string) => {
		return val ? val : undefined;
	}, []);

	return {
		orderBy: replaceNullValue(queryParams.get(ActionFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(ActionFilterQueryParamNames.ORDER)),
		page: queryParams.get(ActionFilterQueryParamNames.PAGE)
			? parseInt(queryParams.get(ActionFilterQueryParamNames.PAGE) as string)
			: TABLE_PAGE_DEFAULT,
		pageSize: queryParams.get(ActionFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(ActionFilterQueryParamNames.PAGE_SIZE) as string)
			: TABLE_PAGE_SIZE_DEFAULT
	};
};

export default useActionsFilterQueryParams;
