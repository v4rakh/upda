import { RequireAuth } from './RequireAuth';
import AuthType from '../auth/AuthType';
import AppPathParamNames from '../constants/appPathParamNames';
import AppPaths from '../constants/appPaths';
import getConfiguration from '../getConfiguration';
import AppLayout from '../layout/AppLayout';
import ActionInvocationsPage from '../pages/action-invocations/ActionInvocationsPage';
import ActionsPage from '../pages/actions/ActionsPage';
import ConstantsPage from '../pages/constants/ConstantsPage';
import ErrorPage404 from '../pages/error-pages/ErrorPage404';
import EventsPage from '../pages/events/EventsPage';
import Home from '../pages/home/Home';
import SessionLogin from '../pages/login/SessionLogin';
import SecretsPage from '../pages/secrets/SecretsPage';
import UpdateSinglePage from '../pages/updates/UpdateSinglePage';
import UpdatesPage from '../pages/updates/UpdatesPage';
import WebhooksPage from '../pages/webhooks/WebhooksPage';
import { getAppBasePath } from '../utils/urlHelper';
import React from 'react';
import { Route, Routes } from 'react-router';

const AppRouter = () => {
	return (
		<Routes>
			<Route path={getAppBasePath()} element={<AppLayout />}>
				<Route index element={<Home />} />
				{getConfiguration().VITE_AUTH_TYPE === AuthType.SESSION && (
					<Route path={AppPaths.LOGIN} element={<SessionLogin />} />
				)}
				<Route
					path={AppPaths.UPDATES}
					element={
						<RequireAuth>
							<UpdatesPage />
						</RequireAuth>
					}
				/>
				<Route
					path={`${AppPaths.UPDATES}/:${AppPathParamNames.UPDATE_ID}`}
					element={
						<RequireAuth>
							<UpdateSinglePage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.WEBHOOKS}
					element={
						<RequireAuth>
							<WebhooksPage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.ACTIONS}
					element={
						<RequireAuth>
							<ActionsPage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.ACTION_INVOCATIONS}
					element={
						<RequireAuth>
							<ActionInvocationsPage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.SECRETS}
					element={
						<RequireAuth>
							<SecretsPage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.CONSTANTS}
					element={
						<RequireAuth>
							<ConstantsPage />
						</RequireAuth>
					}
				/>
				<Route
					path={AppPaths.EVENTS}
					element={
						<RequireAuth>
							<EventsPage />
						</RequireAuth>
					}
				/>
				<Route path="*" element={<ErrorPage404 />} />
			</Route>
		</Routes>
	);
};
export default AppRouter;
