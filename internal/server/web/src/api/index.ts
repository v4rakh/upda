import ApiTags from '../constants/apiTags';
import getConfiguration from '../getConfiguration';
import { updateAuth } from '../slices/authSlice';
import { BaseQueryApi, createApi, FetchArgs, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import AuthType from '../auth/AuthType';

const baseQuery = fetchBaseQuery({
	baseUrl: getConfiguration().VITE_API_URL,
	credentials: getConfiguration().VITE_AUTH_TYPE === AuthType.SESSION ? 'include' : undefined
});

const baseQueryWithReAuth = async (args: string | FetchArgs, api: BaseQueryApi, extraOptions: any) => {
	let result = await baseQuery(args, api, extraOptions);

	if ((getConfiguration().VITE_AUTH_TYPE === AuthType.SESSION && result?.meta?.response?.status) === 401) {
		api.dispatch(updateAuth({ isAuthenticated: false }));
		result = await baseQuery(args, api, extraOptions);
	}

	return result;
};

export const api = createApi({
	reducerPath: 'api',
	baseQuery: baseQueryWithReAuth,
	refetchOnMountOrArgChange: true,
	tagTypes: Object.values(ApiTags),
	endpoints: () => ({})
});

export const { injectEndpoints } = api;
