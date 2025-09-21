import { AuthProviderSession } from './AuthProviderSession';
import AuthType from './AuthType';
import { FC, ReactNode } from 'react';

export const AuthProvider: FC<{ authType: AuthType; children: ReactNode }> = ({ authType, children }) => {
	if (AuthType.SESSION === authType) {
		return <AuthProviderSession>{children}</AuthProviderSession>;
	} else {
		throw new Error('A valid AuthProviderType must be defined');
	}
};
