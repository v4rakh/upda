import getConfiguration from './getConfiguration';
import AppRouter from './router/AppRouter';
import { isDevelopment } from './utils/envHelper';
import { useEffect } from 'react';

const App = () => {
	useEffect(() => {
		if (isDevelopment()) {
			document.title = getConfiguration().VITE_APP_TITLE;
		}
	}, []);

	return <AppRouter />;
};
export default App;
