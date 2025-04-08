import { injectEndpoints } from './index';
import ActionInvocationFilterQueryParamNames from '../constants/api/actionInvocationFilterQueryParamNames';
import ApiTags from '../constants/apiTags';
import { ActionInvocationsRequestParams, ActionInvocationSingleResponse, ActionInvocationsResponse } from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (
	results?: ActionInvocationsResponse | ActionInvocationSingleResponse | void,
	error?: FetchBaseQueryError
) => {
	if (error) {
		return [];
	}
	return [ApiTags.ActionInvocations] as any;
};

export const actionInvocationsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getActionInvocations: build.query<ActionInvocationsResponse, ActionInvocationsRequestParams>({
				query: ({ page, pageSize, order, orderBy }) => {
					const params = new URLSearchParams();
					if (page) {
						params.append(ActionInvocationFilterQueryParamNames.PAGE, `${page}`);
					}
					if (pageSize) {
						params.append(ActionInvocationFilterQueryParamNames.PAGE_SIZE, `${pageSize}`);
					}
					if (order) {
						params.append(ActionInvocationFilterQueryParamNames.ORDER, `${order}`);
					}
					if (orderBy) {
						params.append(ActionInvocationFilterQueryParamNames.ORDER_BY, `${orderBy}`);
					}
					return { url: `action-invocations?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.ActionInvocations, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.ActionInvocations, id }))
						];
					}
					return [];
				}
			}),
			getActionInvocationById: build.query<ActionInvocationSingleResponse, { id: string }>({
				query: ({ id }) => ({ url: `action-invocations/${id}` }),
				providesTags: (result, error) => {
					if (!error && result?.data) {
						return [{ type: ApiTags.ActionInvocations, id: result.data.id }];
					}
					return [];
				}
			}),
			deleteActionInvocation: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `action-invocations/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const { useGetActionInvocationsQuery, useGetActionInvocationByIdQuery, useDeleteActionInvocationMutation } =
	actionInvocationsApi;
