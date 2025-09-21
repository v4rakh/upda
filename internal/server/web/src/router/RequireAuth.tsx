import { useAuth } from '../auth/AuthContext';
import AppPaths from '../constants/appPaths';
import { getPageFullPath } from '../utils/urlHelper';
import { FC, ReactNode } from 'react';
import { Navigate } from 'react-router';

export const RequireAuth: FC<{ children: ReactNode | ReactNode[] }> = ({ children }): ReactNode => {
	const { isAuthenticated } = useAuth();

	if (isAuthenticated) {
		return <>{children}</>;
	} else {
		return <Navigate to={getPageFullPath(AppPaths.LOGIN)} replace={true} />;
	}
};
