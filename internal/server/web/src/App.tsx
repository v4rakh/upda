import { AuthProvider } from './auth/AuthProvider';
import { AuthProviderRehydrate } from './auth/AuthProviderRehydrate';
import getConfiguration from './getConfiguration';
import AppRouter from './router/AppRouter';
import { useEffect } from 'react';

const App = () => {
	useEffect(() => {
		document.title = getConfiguration().VITE_TITLE;
	}, []);

	return (
		<AuthProvider authType={getConfiguration().VITE_AUTH_TYPE}>
			<AuthProviderRehydrate>
				<AppRouter />
			</AuthProviderRehydrate>
		</AuthProvider>
	);
};
export default App;
