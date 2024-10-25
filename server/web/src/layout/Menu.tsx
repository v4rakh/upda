import PrimaryMenu from './PrimaryMenu';
import { Layout } from 'antd';
import { useTranslation } from 'react-i18next';

const { Header } = Layout;

const Menu = () => {
	const [t] = useTranslation('common');

	return (
		<Header style={{ display: 'flex' }}>
			<PrimaryMenu t={t} />
		</Header>
	);
};

export default Menu;
