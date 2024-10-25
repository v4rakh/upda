import { api } from '../api';
import { STORE } from '../constants/localStorageKeys';
import authReducer from '../slices/authSlice';
import { isDevelopment } from '../utils/envHelper';
import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { pick } from 'lodash';
import { useDispatch } from 'react-redux';

const saveToLocalStorage = (state: RootState): void => {
	try {
		const serializedState = JSON.stringify(pick(state, 'auth'));
		localStorage.setItem(STORE, serializedState);
	} catch (e) {
		throw new Error('Could not persist to local storage');
	}
};

const loadFromLocalStorage = (): RootState | null => {
	try {
		const serializedState = localStorage.getItem(STORE);
		if (serializedState === null) {
			return null;
		}

		return JSON.parse(serializedState);
	} catch (e) {
		throw new Error('Could not load from local storage');
	}
};

const rootReducer = combineReducers({
	auth: authReducer,
	[api.reducerPath]: api.reducer
});

export type RootState = ReturnType<typeof rootReducer>;

const persistedState = loadFromLocalStorage();

const store = configureStore({
	reducer: rootReducer,
	devTools: isDevelopment(),
	preloadedState: persistedState ?? undefined,
	middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware)
});

store.subscribe(() => saveToLocalStorage(store.getState()));

export type AppDispatch = typeof store.dispatch;

export const useAppDispatch = (): AppDispatch => useDispatch<AppDispatch>();

export default store;
