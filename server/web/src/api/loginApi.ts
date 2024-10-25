import { injectEndpoints } from './index';
import { LoginRequest } from '../types';

export const loginApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getProbeLogin: build.mutation<void, Partial<LoginRequest>>({
				query: (body) => ({
					url: 'login', // requires an endpoint which return 204 on successful login via basic auth
					headers: {
						Authorization: `Basic ${window.btoa(body.username + ':' + body.password)}`
					}
				}),
				invalidatesTags: []
			})
		};
	}
});

export const { useGetProbeLoginMutation } = loginApi;
