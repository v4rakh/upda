import { useGetInfoQuery } from '../api/infoApi';
import { footerEnabled } from '../utils/featureHelper';
import { Layout, Space, Typography } from 'antd';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

const AppFooter = () => {
	const [t] = useTranslation('version');
	const { data } = useGetInfoQuery();

	if (footerEnabled()) {
		return (
			<Layout.Footer style={{ textAlign: 'center' }}>
				<Space>
					<Text>&copy; {new Date().getFullYear()}</Text>
					<Text>
						{t('version', {
							version: data?.data.version ? data?.data.version.slice(0, 7) : t('unknown')
						})}
					</Text>
				</Space>
			</Layout.Footer>
		);
	} else {
		return <></>;
	}
};
export default AppFooter;
