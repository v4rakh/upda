import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { CreateConstantRequest, ModifyConstantValueRequest, ConstantSingleResponse, ConstantsResponse } from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (results?: ConstantsResponse | ConstantSingleResponse | void, error?: FetchBaseQueryError) => {
	if (error) {
		return [];
	}
	return [ApiTags.Constants] as any;
};

export const constantsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getConstants: build.query<ConstantsResponse, void>({
				query: () => {
					return { url: 'constants' };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Constants, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Constants, id }))
						];
					}
					return [];
				}
			}),
			createConstant: build.mutation<ConstantSingleResponse, CreateConstantRequest>({
				query: (body) => ({ url: 'constants', method: 'POST', body }),
				invalidatesTags
			}),
			modifyValueConstant: build.mutation<
				ConstantSingleResponse,
				{ id: string; body: ModifyConstantValueRequest }
			>({
				query: ({ id, body }) => ({ url: `constants/${id}/value`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Constants, id: arg.id }];
				}
			}),
			deleteConstant: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `constants/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const {
	useGetConstantsQuery,
	useDeleteConstantMutation,
	useCreateConstantMutation,
	useModifyValueConstantMutation,
	useLazyGetConstantsQuery
} = constantsApi;
