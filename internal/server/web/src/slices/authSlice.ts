import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { AuthProfile } from '../types';

export type Auth = {
	isAuthenticated: boolean;
	profile?: AuthProfile;
};
const initialState: Auth = {
	isAuthenticated: false
};

const authSlice = createSlice({
	name: 'auth',
	initialState,
	reducers: {
		updateAuth: (_, { payload: { isAuthenticated, profile } }: PayloadAction<Auth>): Auth => {
			return { isAuthenticated, profile };
		}
	}
});

export const { updateAuth } = authSlice.actions;

export default authSlice.reducer;
