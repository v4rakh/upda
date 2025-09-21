import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { InfoResponse } from '../types';
import ApiVersion from '../constants/apiVersion';

export const infoApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getInfo: build.query<InfoResponse, void>({
				query: () => ({ url: `${ApiVersion.V1}/info` }),
				providesTags: [ApiTags.Info]
			})
		};
	}
});

export const { useGetInfoQuery } = infoApi;
