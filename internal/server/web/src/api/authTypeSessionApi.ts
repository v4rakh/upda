import { injectEndpoints } from './index';
import { AuthProfileResponse, AuthTypeSessionLoginRequest } from '../types';
import ApiVersion from '../constants/apiVersion';

export const authTypeSessionApi = injectEndpoints({
	endpoints: (build) => {
		return {
			login: build.mutation<void, AuthTypeSessionLoginRequest>({
				query: (body) => ({ url: `${ApiVersion.V1}/login`, method: 'POST', body })
			}),
			logout: build.query<void, void>({
				query: () => ({ url: `${ApiVersion.V1}/logout` })
			}),
			profile: build.query<AuthProfileResponse, void>({
				query: () => ({ url: `${ApiVersion.V1}/profile` })
			})
		};
	}
});

export const { useLoginMutation, useLazyLogoutQuery, useLazyProfileQuery } = authTypeSessionApi;
