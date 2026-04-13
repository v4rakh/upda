import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import ApiVersion from '../constants/apiVersion';
import {
	CreateUpdateStateTransitionRequest,
	UpdateStateTransitionSingleResponse,
	UpdateStateTransitionsResponse
} from '../types';

const TAG_LIST_ID = 'LIST';

export const updateStateTransitionsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getUpdateStateTransitions: build.query<UpdateStateTransitionsResponse, void>({
				query: () => ({ url: `${ApiVersion.V1}/update-state-transitions` }),
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.UpdateStateTransitions, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.UpdateStateTransitions, id }))
						];
					}
					return [];
				}
			}),
			getUpdateStateTransitionsByFromStateId: build.query<UpdateStateTransitionsResponse, { stateId: string }>({
				query: ({ stateId }) => ({ url: `${ApiVersion.V1}/update-state-transitions/from/${stateId}` }),
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.UpdateStateTransitions, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.UpdateStateTransitions, id }))
						];
					}
					return [];
				}
			}),
			createUpdateStateTransition: build.mutation<
				UpdateStateTransitionSingleResponse,
				{ body: CreateUpdateStateTransitionRequest }
			>({
				query: ({ body }) => ({
					url: `${ApiVersion.V1}/update-state-transitions`,
					method: 'POST',
					body
				}),
				invalidatesTags: (result, error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.UpdateStateTransitions, id: TAG_LIST_ID }];
				}
			}),
			deleteUpdateStateTransition: build.mutation<void, { id: string }>({
				query: ({ id }) => ({
					url: `${ApiVersion.V1}/update-state-transitions/${id}`,
					method: 'DELETE'
				}),
				invalidatesTags: (result, error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.UpdateStateTransitions, id: TAG_LIST_ID }];
				}
			})
		};
	}
});

export const {
	useGetUpdateStateTransitionsQuery,
	useGetUpdateStateTransitionsByFromStateIdQuery,
	useCreateUpdateStateTransitionMutation,
	useDeleteUpdateStateTransitionMutation
} = updateStateTransitionsApi;
