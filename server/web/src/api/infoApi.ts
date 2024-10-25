import { injectEndpoints } from './index';
import ApiTags from '../constants/apiTags';
import { InfoResponse } from '../types';

export const infoApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getInfo: build.query<InfoResponse, void>({
				query: () => ({ url: 'info' }),
				providesTags: [ApiTags.Info]
			})
		};
	}
});

export const { useGetInfoQuery } = infoApi;
