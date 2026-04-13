import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import ApiVersion from '../constants/apiVersion';
import {
	CreateUpdateStateDefinitionRequest,
	ModifyUpdateStateDefinitionRequest,
	ReorderUpdateStateDefinitionsRequest,
	UpdateStateDefinitionSingleResponse,
	UpdateStateDefinitionsResponse
} from '../types';

const TAG_LIST_ID = 'LIST';

export const updateStateDefinitionsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getUpdateStateDefinitions: build.query<UpdateStateDefinitionsResponse, void>({
				query: () => ({ url: `${ApiVersion.V1}/update-state-definitions` }),
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.UpdateStateDefinitions, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.UpdateStateDefinitions, id }))
						];
					}
					return [];
				}
			}),
			getUpdateStateDefinitionById: build.query<UpdateStateDefinitionSingleResponse, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/update-state-definitions/${id}` }),
				providesTags: (result, error) => {
					if (!error && result?.data) {
						return [{ type: ApiTags.UpdateStateDefinitions, id: result.data.id }];
					}
					return [];
				}
			}),
			createUpdateStateDefinition: build.mutation<
				UpdateStateDefinitionSingleResponse,
				{ body: CreateUpdateStateDefinitionRequest }
			>({
				query: ({ body }) => ({
					url: `${ApiVersion.V1}/update-state-definitions`,
					method: 'POST',
					body
				}),
				invalidatesTags: (result, error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.UpdateStateDefinitions, id: TAG_LIST_ID }];
				}
			}),
			updateUpdateStateDefinition: build.mutation<
				UpdateStateDefinitionSingleResponse,
				{ id: string; body: ModifyUpdateStateDefinitionRequest }
			>({
				query: ({ id, body }) => ({
					url: `${ApiVersion.V1}/update-state-definitions/${id}`,
					method: 'PUT',
					body
				}),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [
						{ type: ApiTags.UpdateStateDefinitions, id: arg.id },
						{ type: ApiTags.UpdateStateDefinitions, id: TAG_LIST_ID }
					];
				}
			}),
			deleteUpdateStateDefinition: build.mutation<void, { id: string }>({
				query: ({ id }) => ({
					url: `${ApiVersion.V1}/update-state-definitions/${id}`,
					method: 'DELETE'
				}),
				invalidatesTags: (result, error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.UpdateStateDefinitions, id: TAG_LIST_ID }];
				}
			}),
			reorderUpdateStateDefinitions: build.mutation<void, ReorderUpdateStateDefinitionsRequest>({
				query: (body) => ({
					url: `${ApiVersion.V1}/update-state-definitions/reorder`,
					method: 'PATCH',
					body
				}),
				invalidatesTags: (result, error) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.UpdateStateDefinitions, id: TAG_LIST_ID }];
				}
			})
		};
	}
});

export const {
	useGetUpdateStateDefinitionsQuery,
	useGetUpdateStateDefinitionByIdQuery,
	useCreateUpdateStateDefinitionMutation,
	useUpdateUpdateStateDefinitionMutation,
	useDeleteUpdateStateDefinitionMutation,
	useReorderUpdateStateDefinitionsMutation
} = updateStateDefinitionsApi;
