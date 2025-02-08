import './i18n';
import App from './App';
import LocaleContextProvider from './providers/LocaleContextProvider';
import store from './store';
import { StrictMode } from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router';
import './style/app-theme.less';
import '@ant-design/v5-patch-for-react-19';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

root.render(
	<Provider store={store}>
		<StrictMode>
			<Router>
				<LocaleContextProvider>
					<App />
				</LocaleContextProvider>
			</Router>
		</StrictMode>
	</Provider>
);
