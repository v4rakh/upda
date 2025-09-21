import { AuthProfile } from '../types';
import { createContext, useContext } from 'react';

export interface AuthContextLoginArgs {
	credentials: any;
	preLogin?: () => void;
	postLogin?: () => void;
}

export interface AuthContextLogoutArgs {
	preLogout?: () => void;
	postLogout?: () => void;
}

export interface AuthContextRehydrateArgs {
	preRehydrate?: () => void;
	postRehydrate?: () => void;
}

/**
 * Enforces presence of common functionality related to authentication.
 * Implementations must match their back-end counterpart.
 *
 * @param login logs the session in
 * @param logout logs the session out
 * @param profile returns the session's profile
 * @param rehydrate a function which rehydrates authentication state, usually by calling profile (automatically called in AuthProviderRehydrate)
 * @param isAuthenticated returns true if browser session is authenticated, false otherwise
 */
export interface AuthContextType {
	login: (args: AuthContextLoginArgs) => void;
	logout: (args: AuthContextLogoutArgs) => void;
	rehydrate: (args: AuthContextRehydrateArgs) => void;
	profile: AuthProfile | undefined;
	isAuthenticated: boolean;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error('useAuth must be used within AuthProvider');
	}
	return context;
};
