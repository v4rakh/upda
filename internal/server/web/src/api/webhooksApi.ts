import { injectEndpoints } from './index';
import WebhookFilterQueryParamNames from '../constants/api/webhookFilterQueryParamNames';
import ApiTags from '../constants/apiTags';
import {
	CreateWebhookRequest,
	ModifyWebhookIgnoreHostRequest,
	ModifyWebhookLabelRequest,
	WebhookSingleResponse,
	WebhooksRequestParams,
	WebhooksResponse
} from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import ApiVersion from '../constants/apiVersion';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (results?: WebhooksResponse | WebhookSingleResponse | void, error?: FetchBaseQueryError) => {
	if (error) {
		return [];
	}
	return [ApiTags.Webhooks] as any;
};

export const webhooksApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getWebhooks: build.query<WebhooksResponse, WebhooksRequestParams>({
				query: ({ page, pageSize, order, orderBy }) => {
					const params = new URLSearchParams();
					if (page) {
						params.append(WebhookFilterQueryParamNames.PAGE, `${page}`);
					}
					if (pageSize) {
						params.append(WebhookFilterQueryParamNames.PAGE_SIZE, `${pageSize}`);
					}
					if (order) {
						params.append(WebhookFilterQueryParamNames.ORDER, `${order}`);
					}
					if (orderBy) {
						params.append(WebhookFilterQueryParamNames.ORDER_BY, `${orderBy}`);
					}
					return { url: `${ApiVersion.V1}/webhooks?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Webhooks, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Webhooks, id }))
						];
					}
					return [];
				}
			}),
			createWebhook: build.mutation<WebhookSingleResponse, CreateWebhookRequest>({
				query: (body) => ({ url: `${ApiVersion.V1}/webhooks`, method: 'POST', body }),
				invalidatesTags
			}),
			deleteWebhook: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/webhooks/${id}`, method: 'DELETE' }),
				invalidatesTags
			}),
			modifyLabelWebhook: build.mutation<WebhookSingleResponse, { id: string; body: ModifyWebhookLabelRequest }>({
				query: ({ id, body }) => ({ url: `${ApiVersion.V1}/webhooks/${id}/label`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Webhooks, id: arg.id }];
				}
			}),
			modifyIgnoreHostWebhook: build.mutation<
				WebhookSingleResponse,
				{ id: string; body: ModifyWebhookIgnoreHostRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/webhooks/${id}/ignore-host`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Webhooks, id: arg.id }];
				}
			})
		};
	}
});

export const {
	useGetWebhooksQuery,
	useDeleteWebhookMutation,
	useCreateWebhookMutation,
	useModifyLabelWebhookMutation,
	useModifyIgnoreHostWebhookMutation
} = webhooksApi;
