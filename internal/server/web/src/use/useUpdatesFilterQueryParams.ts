import UpdateFilterQueryParamNames from '../constants/api/updateFilterQueryParamNames';
import {
	UPDATES_CARD_PAGE_DEFAULT,
	UPDATES_CARD_PAGE_SIZE_DEFAULT,
	TABLE_PAGE_DEFAULT,
	TABLE_PAGE_SIZE_DEFAULT
} from '../constants/pagination';
import { parseInt } from 'lodash';
import { useSearchParams } from 'react-router';
import { replaceEmptyValue, replaceNullValue } from '../utils/queryParamsHelper';

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

	return {
		searchTerm: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.SEARCH_TERM)),
		searchIn: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.SEARCH_IN)),
		state: replaceEmptyValue(queryParams.getAll(UpdateFilterQueryParamNames.STATE)),
		orderBy: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(UpdateFilterQueryParamNames.ORDER)),
		page: queryParams.get(UpdateFilterQueryParamNames.PAGE)
			? parseInt(queryParams.get(UpdateFilterQueryParamNames.PAGE) as string)
			: UPDATES_CARD_PAGE_DEFAULT,
		pageSize: queryParams.get(UpdateFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(UpdateFilterQueryParamNames.PAGE_SIZE) as string)
			: UPDATES_CARD_PAGE_SIZE_DEFAULT
	};
};

export default useUpdatesFilterQueryParams;
