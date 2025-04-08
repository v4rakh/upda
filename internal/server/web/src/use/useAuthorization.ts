import AppPaths from '../constants/appPaths';
import { useAuthSelector } from '../selectors/authSelectors';
import { updateAuth } from '../slices/authSlice';
import { useAppDispatch } from '../store';
import { getPageFullPath } from '../utils/urlHelper';
import { useCallback } from 'react';
import { useNavigate } from 'react-router';

export interface AuthHook {
	logout: () => void;
	getUserName: () => string | null;
}

export const useAuthorization = (): AuthHook => {
	const navigate = useNavigate();
	const dispatch = useAppDispatch();
	const auth = useAuthSelector();

	return {
		getUserName: useCallback((): string | null => {
			return auth.username;
		}, [auth.username]),
		logout: useCallback((): void => {
			dispatch(updateAuth({ username: null, password: null }));
			navigate(getPageFullPath(AppPaths.LOGIN));
		}, [dispatch, navigate])
	};
};
