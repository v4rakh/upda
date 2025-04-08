import { Auth } from '../slices/authSlice';
import { RootState } from '../store';
import { useSelector } from 'react-redux';

export const useAuthSelector = (): Auth => useSelector((state: RootState) => state.auth);
export const useAuthenticatedSelector = (): boolean =>
	useSelector((state: RootState) => state.auth.username !== null && state.auth.password !== null);
