import { api } from '../api';
import authReducer from '../slices/authSlice';
import { isDevelopment } from '../utils/envHelper';
import { combineReducers, configureStore } from '@reduxjs/toolkit';
import { useDispatch } from 'react-redux';

const rootReducer = combineReducers({
	auth: authReducer,
	[api.reducerPath]: api.reducer
});

export type RootState = ReturnType<typeof rootReducer>;

const store = configureStore({
	reducer: rootReducer,
	devTools: isDevelopment(),
	middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(api.middleware)
});

export type AppDispatch = typeof store.dispatch;

export const useAppDispatch = (): AppDispatch => useDispatch<AppDispatch>();

export default store;
