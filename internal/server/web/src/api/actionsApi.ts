import { injectEndpoints } from './index';
import ActionFilterQueryParamNames from '../constants/api/actionFilterQueryParamNames';
import ApiTags from '../constants/apiTags';
import {
	ActionSingleResponse,
	ActionsRequestParams,
	ActionsResponse,
	ActionTestSingleResponse,
	CreateActionRequest,
	ModifyActionEnabledRequest,
	ModifyActionLabelRequest,
	ModifyActionMatchApplicationRequest,
	ModifyActionMatchEventRequest,
	ModifyActionMatchHostRequest,
	ModifyActionMatchProviderRequest,
	ModifyActionPayloadRequest,
	TestActionRequest
} from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import ApiVersion from '../constants/apiVersion';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (results?: ActionsResponse | ActionSingleResponse | void, error?: FetchBaseQueryError) => {
	if (error) {
		return [];
	}
	return [ApiTags.Actions] as any;
};

export const actionsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getActions: build.query<ActionsResponse, ActionsRequestParams>({
				query: ({ page, pageSize, order, orderBy }) => {
					const params = new URLSearchParams();
					if (page) {
						params.append(ActionFilterQueryParamNames.PAGE, `${page}`);
					}
					if (pageSize) {
						params.append(ActionFilterQueryParamNames.PAGE_SIZE, `${pageSize}`);
					}
					if (order) {
						params.append(ActionFilterQueryParamNames.ORDER, `${order}`);
					}
					if (orderBy) {
						params.append(ActionFilterQueryParamNames.ORDER_BY, `${orderBy}`);
					}
					return { url: `${ApiVersion.V1}/actions?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Actions, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Actions, id }))
						];
					}
					return [];
				}
			}),
			getActionById: build.query<ActionSingleResponse, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/actions/${id}` }),
				providesTags: (result, error) => {
					if (!error && result?.data) {
						return [{ type: ApiTags.Actions, id: result.data.id }];
					}
					return [];
				}
			}),
			createAction: build.mutation<ActionSingleResponse, CreateActionRequest>({
				query: (body) => ({ url: `${ApiVersion.V1}/actions`, method: 'POST', body }),
				invalidatesTags
			}),
			testAction: build.mutation<ActionTestSingleResponse, { id: string; body: TestActionRequest }>({
				query: ({ id, body }) => ({ url: `${ApiVersion.V1}/actions/${id}/test`, method: 'POST', body })
			}),
			modifyLabelAction: build.mutation<ActionSingleResponse, { id: string; body: ModifyActionLabelRequest }>({
				query: ({ id, body }) => ({ url: `${ApiVersion.V1}/actions/${id}/label`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyMatchApplicationAction: build.mutation<
				ActionSingleResponse,
				{ id: string; body: ModifyActionMatchApplicationRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/actions/${id}/match-application`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyMatchHostAction: build.mutation<
				ActionSingleResponse,
				{ id: string; body: ModifyActionMatchHostRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/actions/${id}/match-host`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyMatchEventAction: build.mutation<
				ActionSingleResponse,
				{ id: string; body: ModifyActionMatchEventRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/actions/${id}/match-event`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyMatchProviderAction: build.mutation<
				ActionSingleResponse,
				{ id: string; body: ModifyActionMatchProviderRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/actions/${id}/match-provider`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyTypeAndPayloadAction: build.mutation<
				ActionSingleResponse,
				{ id: string; body: ModifyActionPayloadRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/actions/${id}/payload`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Actions, id: arg.id }];
				}
			}),
			modifyEnabledAction: build.mutation<ActionSingleResponse, { id: string; body: ModifyActionEnabledRequest }>(
				{
					query: ({ id, body }) => ({ url: `${ApiVersion.V1}/actions/${id}/enabled`, method: 'PATCH', body }),
					invalidatesTags: (result, error, arg) => {
						if (error) {
							return [];
						}
						return [{ type: ApiTags.Actions, id: arg.id }];
					}
				}
			),
			deleteAction: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/actions/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const {
	useGetActionsQuery,
	useGetActionByIdQuery,
	useDeleteActionMutation,
	useModifyLabelActionMutation,
	useModifyMatchEventActionMutation,
	useModifyMatchApplicationActionMutation,
	useModifyMatchHostActionMutation,
	useModifyMatchProviderActionMutation,
	useModifyTypeAndPayloadActionMutation,
	useModifyEnabledActionMutation,
	useCreateActionMutation,
	useTestActionMutation
} = actionsApi;
