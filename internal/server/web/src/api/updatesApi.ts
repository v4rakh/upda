import { injectEndpoints } from './index';
import UpdateFilterQueryParamNames from '../constants/api/updateFilterQueryParamNames';
import ApiTags from '../constants/apiTags';
import { ModifyUpdateStateRequest, UpdateSingleResponse, UpdatesRequestParams, UpdatesResponse } from '../types';
import { forEach } from 'lodash';

const TAG_LIST_ID = 'LIST';

export const updatesApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getUpdates: build.query<UpdatesResponse, UpdatesRequestParams>({
				query: ({ ...args }) => {
					const { page, pageSize, order, orderBy, state, searchIn, searchTerm } = args;

					const params = new URLSearchParams();
					if (state) {
						forEach(state, (s) => {
							params.append(UpdateFilterQueryParamNames.STATE, s);
						});
					}
					if (searchIn) {
						params.append(UpdateFilterQueryParamNames.SEARCH_IN, `${searchIn}`);
					}
					if (searchTerm) {
						params.append(UpdateFilterQueryParamNames.SEARCH_TERM, `${searchTerm}`);
					}
					if (page) {
						params.append(UpdateFilterQueryParamNames.PAGE, `${page}`);
					}
					if (pageSize) {
						params.append(UpdateFilterQueryParamNames.PAGE_SIZE, `${pageSize}`);
					}
					if (order) {
						params.append(UpdateFilterQueryParamNames.ORDER, `${order}`);
					}
					if (orderBy) {
						params.append(UpdateFilterQueryParamNames.ORDER_BY, `${orderBy}`);
					}
					return { url: `updates?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Updates, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Updates, id }))
						];
					}
					return [];
				}
			}),
			getUpdateById: build.query<UpdateSingleResponse, { id: string }>({
				query: ({ id }) => ({ url: `updates/${id}` }),
				providesTags: (result, error) => {
					if (!error && result?.data) {
						return [{ type: ApiTags.Updates, id: result.data.id }];
					}
					return [];
				}
			}),
			modifyUpdateState: build.mutation<UpdateSingleResponse, { id: string; body: ModifyUpdateStateRequest }>({
				query: ({ id, body }) => ({ url: `updates/${id}/state`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Updates, id: arg.id }];
				}
			}),
			deleteUpdate: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `updates/${id}`, method: 'DELETE' }),
				invalidatesTags: (error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Updates }];
				}
			})
		};
	}
});

export const { useGetUpdatesQuery, useGetUpdateByIdQuery, useModifyUpdateStateMutation, useDeleteUpdateMutation } =
	updatesApi;
