import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { FetchBaseQueryError } from '@reduxjs/toolkit/query';
import {
	CreateFilterPresetRequest,
	FilterPresetSingleResponse,
	FilterPresetsResponse,
	FilterPresetType
} from '../types/filterPreset';
import ApiVersion from '../constants/apiVersion';

const TAG_LIST_ID = 'LIST';

const invalidatesTags = (
	results?: FilterPresetsResponse | FilterPresetSingleResponse | void,
	error?: FetchBaseQueryError
) => {
	if (error) {
		return [];
	}
	return [ApiTags.FilterPresets] as any;
};

export const filterPresetsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getFilterPresetsByType: build.query<FilterPresetsResponse, FilterPresetType>({
				query: (type) => {
					return { url: `${ApiVersion.V1}/filter-presets/${type}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.FilterPresets, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.FilterPresets, id }))
						];
					}
					return [];
				}
			}),
			createFilterPreset: build.mutation<FilterPresetSingleResponse, CreateFilterPresetRequest>({
				query: (body) => ({ url: '${ApiVersion.V1}/filter-presets', method: 'POST', body }),
				invalidatesTags
			}),
			deleteFilterPreset: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/filter-presets/${id}`, method: 'DELETE' }),
				invalidatesTags
			})
		};
	}
});

export const { useGetFilterPresetsByTypeQuery, useCreateFilterPresetMutation, useDeleteFilterPresetMutation } =
	filterPresetsApi;
