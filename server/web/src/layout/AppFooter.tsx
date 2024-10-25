import { Layout, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

const AppFooter = () => {
	const [t] = useTranslation('version');
	return (
		<Layout.Footer style={{ textAlign: 'center' }}>
			<Space>
				<Typography.Text>&copy; {new Date().getFullYear()}</Typography.Text>
				<Typography.Text>
					{t('version')} {t('number')}
				</Typography.Text>
			</Space>
		</Layout.Footer>
	);
};
export default AppFooter;
