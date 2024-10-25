import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { CreateSecretRequest, ModifySecretValueRequest, SecretSingleResponse, SecretsResponse } from '../types';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (results?: SecretsResponse | SecretSingleResponse | void, error?: FetchBaseQueryError) => {
	if (error) {
		return [];
	}
	return [ApiTags.Secrets] as any;
};

export const secretsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getSecrets: build.query<SecretsResponse, void>({
				query: () => {
					return { url: 'secrets' };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Secrets, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Secrets, id }))
						];
					}
					return [];
				}
			}),
			createSecret: build.mutation<SecretSingleResponse, CreateSecretRequest>({
				query: (body) => ({ url: 'secrets', method: 'POST', body }),
				invalidatesTags
			}),
			modifyValueSecret: build.mutation<SecretSingleResponse, { id: string; body: ModifySecretValueRequest }>({
				query: ({ id, body }) => ({ url: `secrets/${id}/value`, method: 'PATCH', body }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Secrets, id: arg.id }];
				}
			}),
			deleteSecret: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `secrets/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const { useGetSecretsQuery, useDeleteSecretMutation, useCreateSecretMutation, useModifyValueSecretMutation } =
	secretsApi;
