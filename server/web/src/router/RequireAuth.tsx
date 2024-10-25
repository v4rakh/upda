import AppPaths from '../constants/appPaths';
import { useAuthenticatedSelector } from '../selectors/authSelectors';
import { getPageFullPath } from '../utils/urlHelper';
import { FC, ReactNode } from 'react';
import { Navigate } from 'react-router-dom';

export const RequireAuth: FC<{ children: ReactNode | ReactNode[] }> = ({ children }): JSX.Element => {
	const isAuthenticated = useAuthenticatedSelector();

	if (isAuthenticated) {
		return <>{children}</>;
	} else {
		return <Navigate to={getPageFullPath(AppPaths.LOGIN)} replace={true} />;
	}
};
