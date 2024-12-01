import './i18n';
import App from './App';
import LocaleContextProvider from './providers/LocaleContextProvider';
import store from './store';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { HashRouter as Router } from 'react-router';
import './style/app-theme.less';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);

root.render(
	<Provider store={store}>
		{/* the waring below appears when using ant menu due to react strict mode.
		 * Using <Button> results in "findDOMNode is deprecated in StrictMode" warning
		 * Fix not yet available. follow thread: https://github.com/ant-design/ant-design/issues/22493 */}
		{/*<React.StrictMode>*/}
		<Router>
			<LocaleContextProvider>
				<App />
			</LocaleContextProvider>
		</Router>
		{/*</React.StrictMode>*/}
	</Provider>
);
