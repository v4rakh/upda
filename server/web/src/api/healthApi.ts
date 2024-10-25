import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { HealthResponse } from '../types';

export const healthApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getHealth: build.query<HealthResponse, void>({
				query: () => ({ url: 'health' }),
				providesTags: [ApiTags.Health]
			})
		};
	}
});

export const { useGetHealthQuery } = healthApi;
