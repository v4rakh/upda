import PrimaryMenu from './PrimaryMenu';
import { useTheme } from '../providers/ThemeProvider';
import { MoonOutlined, SunOutlined } from '@ant-design/icons';
import { Button, Layout, Tooltip, theme } from 'antd';
import { useTranslation } from 'react-i18next';

const { Header } = Layout;

const Menu = () => {
	const [t] = useTranslation('common');
	const { isDarkTheme, toggleTheme } = useTheme();
	const { token } = theme.useToken();

	return (
		<Header style={{ display: 'flex', alignItems: 'center' }}>
			<PrimaryMenu t={t} />
			<Tooltip title={isDarkTheme ? t('toggle_light_theme') : t('toggle_dark_theme')}>
				<Button
					type="text"
					style={{ color: token.colorTextLightSolid, marginLeft: 'auto' }}
					icon={isDarkTheme ? <SunOutlined /> : <MoonOutlined />}
					onClick={toggleTheme}
				/>
			</Tooltip>
		</Header>
	);
};

export default Menu;
