import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type Auth = {
	username: string | null;
	password: string | null;
};
const initialState: Auth = {
	username: null,
	password: null
};

const authSlice = createSlice({
	name: 'auth',
	initialState,
	reducers: {
		updateAuth: (state, { payload }: PayloadAction<Auth>): Auth => {
			return { username: payload.username, password: payload.password };
		}
	}
});

export const { updateAuth } = authSlice.actions;

export default authSlice.reducer;
