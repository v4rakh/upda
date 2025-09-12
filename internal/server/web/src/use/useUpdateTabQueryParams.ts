import { useCallback } from 'react';
import { useSearchParams } from 'react-router';
import UpdateTabQueryParamNames from '../constants/api/updateTabQueryParamNames';
import UpdateTabNames from '../constants/updateTabNames';

const useUpdateTabQueryParams = () => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | string) => {
		return val ? val : undefined;
	}, []);

	return {
		tab: replaceNullValue(queryParams.get(UpdateTabQueryParamNames.TAB)) ?? UpdateTabNames.DETAILS
	};
};

export default useUpdateTabQueryParams;
