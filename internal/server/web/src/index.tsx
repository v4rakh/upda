import './i18n';
import App from './App';
import LocaleContextProvider from './providers/LocaleContextProvider';
import store from './store';
import { darkThemeEnabled } from './utils/featureHelper';
import { App as AntDesignApp, ConfigProvider, theme } from 'antd';
import { StrictMode } from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router';
import './style/app-theme.less';
import '@ant-design/v5-patch-for-react-19';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
const algorithm = darkThemeEnabled() ? theme.darkAlgorithm : theme.defaultAlgorithm;

root.render(
	<Provider store={store}>
		<StrictMode>
			<Router>
				<ConfigProvider
					theme={{
						algorithm: algorithm
					}}>
					<AntDesignApp>
						<LocaleContextProvider>
							<App />
						</LocaleContextProvider>
					</AntDesignApp>
				</ConfigProvider>
			</Router>
		</StrictMode>
	</Provider>
);
