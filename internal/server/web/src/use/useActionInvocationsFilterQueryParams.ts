import ActionInvocationFilterQueryParamNames from '../constants/api/actionInvocationFilterQueryParamNames';
import { TABLE_PAGE_DEFAULT, TABLE_PAGE_SIZE_DEFAULT } from '../constants/pagination';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router';

const useActionInvocationsFilterQueryParams = () => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | number | string) => {
		return val ? val : undefined;
	}, []);

	return {
		orderBy: replaceNullValue(queryParams.get(ActionInvocationFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(ActionInvocationFilterQueryParamNames.ORDER)),
		page: queryParams.get(ActionInvocationFilterQueryParamNames.PAGE)
			? parseInt(queryParams.get(ActionInvocationFilterQueryParamNames.PAGE) as string)
			: TABLE_PAGE_DEFAULT,
		pageSize: queryParams.get(ActionInvocationFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(ActionInvocationFilterQueryParamNames.PAGE_SIZE) as string)
			: TABLE_PAGE_SIZE_DEFAULT
	};
};

export default useActionInvocationsFilterQueryParams;
