import { injectEndpoints } from './index';
import { LoginRequest } from '../types';
import ApiVersion from '../constants/apiVersion';

export const loginApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getProbeLogin: build.mutation<void, Partial<LoginRequest>>({
				query: (body) => ({
					url: `${ApiVersion.V1}/login`, // requires an endpoint which return 204 on successful login via basic auth
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
