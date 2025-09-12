import EventFilterQueryParamNames from '../constants/api/eventFilterQueryParamNames';
import { WINDOW_SIZE_DEFAULT, WINDOW_SKIP_DEFAULT } from '../constants/pagination';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router';

const useEventsFilterQueryParams = () => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | number | string) => {
		return val ? val : undefined;
	}, []);

	return {
		updateId: replaceNullValue(queryParams.get(EventFilterQueryParamNames.UPDATE_ID)),
		orderBy: replaceNullValue(queryParams.get(EventFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(EventFilterQueryParamNames.ORDER)),
		size: queryParams.get(EventFilterQueryParamNames.SIZE)
			? parseInt(queryParams.get(EventFilterQueryParamNames.SIZE) as string)
			: WINDOW_SIZE_DEFAULT,
		skip: queryParams.get(EventFilterQueryParamNames.SKIP)
			? parseInt(queryParams.get(EventFilterQueryParamNames.SKIP) as string)
			: WINDOW_SKIP_DEFAULT
	};
};

export default useEventsFilterQueryParams;
