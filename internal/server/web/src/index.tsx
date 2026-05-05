import './i18n';
import App from './App';
import LocaleContextProvider from './providers/LocaleContextProvider';
import ThemeProvider from './providers/ThemeProvider';
import store from './store';
import { App as AntDesignApp } from 'antd';
import { StrictMode } from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router';
import './style/app-theme.less';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

root.render(
	<Provider store={store}>
		<StrictMode>
			<Router>
				<ThemeProvider>
					<AntDesignApp>
						<LocaleContextProvider>
							<App />
						</LocaleContextProvider>
					</AntDesignApp>
				</ThemeProvider>
			</Router>
		</StrictMode>
	</Provider>
);
