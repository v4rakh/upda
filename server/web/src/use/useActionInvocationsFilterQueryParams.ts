import ActionInvocationFilterQueryParamNames from '../constants/api/actionInvocationFilterQueryParamNames';
import { PAGE_DEFAULT, PAGE_SIZE_DEFAULT } from '../constants/pagination';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';

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
			: PAGE_DEFAULT,
		pageSize: queryParams.get(ActionInvocationFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(ActionInvocationFilterQueryParamNames.PAGE_SIZE) as string)
			: PAGE_SIZE_DEFAULT
	};
};

export default useActionInvocationsFilterQueryParams;
