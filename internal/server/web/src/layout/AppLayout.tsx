import AppFooter from './AppFooter';
import HealthHandler from './HealthHandler';
import Menu from './Menu';
import { Layout } from 'antd';
import { FC } from 'react';
import { Outlet } from 'react-router';

const AppLayout: FC = () => {
	return (
		<Layout style={{ minHeight: '100vh' }}>
			<Menu />
			<Layout>
				<HealthHandler>
					<Layout.Content style={{ margin: '24px 16px 0', overflow: 'initial' }}>
						<Outlet />
					</Layout.Content>
					<AppFooter />
				</HealthHandler>
			</Layout>
		</Layout>
	);
};
export default AppLayout;
