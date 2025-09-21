import { AuthContext, AuthContextRehydrateArgs, AuthContextLoginArgs, AuthContextLogoutArgs } from './AuthContext';
import { useLazyLogoutQuery, useLazyProfileQuery, useLoginMutation } from '../api/authTypeSessionApi';
import AppPaths from '../constants/appPaths';
import { updateAuth } from '../slices/authSlice';
import { RootState, useAppDispatch } from '../store';
import { AuthProfileResponse, AuthTypeSessionLoginRequest } from '../types';
import { useNotification } from '../use/useNotification';
import { getPageFullPath } from '../utils/urlHelper';
import { noop } from 'lodash';
import React, { FC, ReactNode, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router';

export const AuthProviderSession: FC<{ children: ReactNode }> = ({ children }) => {
	const [t] = useTranslation('auth_provider_session');
	const dispatch = useAppDispatch();
	const navigate = useNavigate();
	const { apiError, simpleError } = useNotification();

	const isAuthenticated = useSelector((state: RootState) => state.auth.isAuthenticated);
	const profileData = useSelector((state: RootState) => state.auth.profile);

	const [triggerLogout] = useLazyLogoutQuery();
	const [triggerLogin] = useLoginMutation();
	const [triggerProfile] = useLazyProfileQuery();

	const profile = useCallback(async (): Promise<AuthProfileResponse | undefined> => {
		const { isSuccess, isError, data } = await triggerProfile();

		if (isSuccess) {
			return data;
		}

		if (isError) {
			simpleError({ title: t('profile_error_title'), message: t('profile_error_description') });
		}

		return undefined;
	}, [triggerProfile, simpleError, t]);

	const login = useCallback(
		async ({ credentials, preLogin = noop, postLogin = noop }: AuthContextLoginArgs) => {
			if (!isSessionLoginRequest(credentials)) {
				throw new Error('Invalid type for auth type session auth');
			}

			preLogin();

			const { error } = await triggerLogin(credentials);

			if (error) {
				apiError({
					i18n: {
						unAuthorized: t('login_error_unauthorized'),
						default: t('login_error_default_message')
					},
					error: error
				});

				postLogin();
				return;
			}

			const profileResult = await profile();

			dispatch(updateAuth({ isAuthenticated: true, profile: profileResult ? profileResult.data : undefined }));
			navigate(getPageFullPath(AppPaths.HOME));
			postLogin();
		},
		[apiError, dispatch, profile, triggerLogin, navigate, t]
	);

	const logout = useCallback(
		async ({ preLogout = noop, postLogout = noop }: AuthContextLogoutArgs) => {
			preLogout();
			const { isSuccess, isError } = await triggerLogout();
			postLogout();

			if (isSuccess) {
				dispatch(updateAuth({ isAuthenticated: false, profile: undefined }));
				navigate(getPageFullPath(AppPaths.LOGIN));
			}

			if (isError) {
				simpleError({ title: t('logout_error_title'), message: t('logout_error_description') });
			}
		},
		[dispatch, triggerLogout, navigate, simpleError, t]
	);

	const rehydrate = useCallback(
		async ({ preRehydrate = noop, postRehydrate = noop }: AuthContextRehydrateArgs) => {
			preRehydrate();
			const { isSuccess, data } = await triggerProfile();

			if (isSuccess) {
				dispatch(updateAuth({ isAuthenticated: true, profile: data?.data }));
			}

			postRehydrate();
		},
		[triggerProfile, dispatch]
	);

	return (
		<AuthContext.Provider
			value={{
				login: login,
				logout: logout,
				profile: profileData,
				rehydrate: rehydrate,
				isAuthenticated: isAuthenticated
			}}>
			{children}
		</AuthContext.Provider>
	);
};

const isSessionLoginRequest = (credentials: unknown): credentials is AuthTypeSessionLoginRequest => {
	return (
		typeof credentials === 'object' &&
		credentials !== null &&
		typeof (credentials as any).username === 'string' &&
		typeof (credentials as any).password === 'string'
	);
};
