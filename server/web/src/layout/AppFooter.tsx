import { Layout, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

const AppFooter = () => {
	const [t] = useTranslation('version');
	return (
		<Layout.Footer style={{ textAlign: 'center' }}>
			<Space>
				<Text>&copy; {new Date().getFullYear()}</Text>
				<Text>
					{t('version')} {t('number')}
				</Text>
			</Space>
		</Layout.Footer>
	);
};
export default AppFooter;
