import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { HealthResponse } from '../types';
import ApiVersion from '../constants/apiVersion';

export const healthApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getHealth: build.query<HealthResponse, void>({
				query: () => ({ url: `${ApiVersion.V1}/health` }),
				providesTags: [ApiTags.Health]
			})
		};
	}
});

export const { useGetHealthQuery } = healthApi;
