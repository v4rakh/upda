import ApiTags from '../constants/apiTags';
import getConfiguration from '../getConfiguration';
import { updateAuth } from '../slices/authSlice';
import { RootState } from '../store';
import { BaseQueryApi, createApi, FetchArgs, fetchBaseQuery } from '@reduxjs/toolkit/query/react';

const baseQuery = fetchBaseQuery({
	baseUrl: getConfiguration().VITE_API_URL,
	prepareHeaders: (headers, { getState }) => {
		const state = getState() as RootState;
		const username = state.auth.username;
		const password = state.auth.password;
		const authHeader = window.btoa(`${username}:${password}`);

		if (username && password && authHeader) {
			headers.set('Authorization', `Basic ${authHeader}`);
		}
		return headers;
	}
});

const baseQueryWithReAuth = async (args: string | FetchArgs, api: BaseQueryApi, extraOptions: any) => {
	let result = await baseQuery(args, api, extraOptions);

	if (result?.meta?.response?.status === 401) {
		api.dispatch(updateAuth({ username: null, password: null }));
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
