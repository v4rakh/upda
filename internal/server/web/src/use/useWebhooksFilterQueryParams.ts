import WebhookFilterQueryParamNames from '../constants/api/webhookFilterQueryParamNames';
import { CARD_PAGE_DEFAULT, CARD_PAGE_SIZE_DEFAULT } from '../constants/pagination';
import { useCallback } from 'react';
import { useSearchParams } from 'react-router';

const useWebhooksFilterQueryParams = () => {
	const [queryParams] = useSearchParams();

	const replaceNullValue = useCallback((val: null | number | string) => {
		return val ? val : undefined;
	}, []);

	return {
		orderBy: replaceNullValue(queryParams.get(WebhookFilterQueryParamNames.ORDER_BY)),
		order: replaceNullValue(queryParams.get(WebhookFilterQueryParamNames.ORDER)),
		page: queryParams.get(WebhookFilterQueryParamNames.PAGE)
			? parseInt(queryParams.get(WebhookFilterQueryParamNames.PAGE) as string)
			: CARD_PAGE_DEFAULT,
		pageSize: queryParams.get(WebhookFilterQueryParamNames.PAGE_SIZE)
			? parseInt(queryParams.get(WebhookFilterQueryParamNames.PAGE_SIZE) as string)
			: CARD_PAGE_SIZE_DEFAULT
	};
};

export default useWebhooksFilterQueryParams;
