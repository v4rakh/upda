import classes from './style/PrimaryMenu.module.less';

import AppPaths from '../constants/appPaths';
import { useAuthenticatedSelector } from '../selectors/authSelectors';
import { useAuthorization } from '../use/useAuthorization';
import { getPageFullPath } from '../utils/urlHelper';
import {
	BarsOutlined,
	BuildOutlined,
	ClockCircleOutlined,
	LinkOutlined,
	LockOutlined,
	LogoutOutlined,
	UnorderedListOutlined,
	UserOutlined
} from '@ant-design/icons';
import { Menu, Typography } from 'antd';
import { TFunction } from 'i18next';
import { forEach } from 'lodash';
import { FC, ReactNode, useMemo } from 'react';
import { useNavigate } from 'react-router';

export type PrimaryMenuProps = {
	t: TFunction;
};

const PrimaryMenu: FC<PrimaryMenuProps> = ({ t }): ReactNode => {
	const navigate = useNavigate();
	const isAuthenticated = useAuthenticatedSelector();
	const { logout, getUserName } = useAuthorization();
	const username = getUserName();

	const primaryNavs = useMemo(() => {
		const staticItems = [
			isAuthenticated && {
				label: t('updates'),
				key: AppPaths.UPDATES,
				icon: <UnorderedListOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.UPDATES))
			},
			isAuthenticated && {
				label: t('webhooks'),
				key: AppPaths.WEBHOOKS,
				icon: <LinkOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.WEBHOOKS))
			},
			isAuthenticated && {
				label: t('actions'),
				key: AppPaths.ACTIONS,
				icon: <BuildOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.ACTIONS))
			},
			isAuthenticated && {
				label: t('action_invocations'),
				key: AppPaths.ACTION_INVOCATIONS,
				icon: <BarsOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.ACTION_INVOCATIONS))
			},
			isAuthenticated && {
				label: t('secrets'),
				key: AppPaths.SECRETS,
				icon: <LockOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.SECRETS))
			},
			isAuthenticated && {
				label: t('events'),
				key: AppPaths.EVENTS,
				icon: <ClockCircleOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.EVENTS))
			}
		];

		const items = [];
		forEach(staticItems, (s) => {
			items.push(s);
		});

		if (!isAuthenticated) {
			items.push({
				label: t('login'),
				key: AppPaths.LOGIN,
				icon: <UserOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.LOGIN))
			});
		} else {
			items.push({
				key: 'menu_logout',
				icon: <LogoutOutlined />,
				label: (
					<Typography.Text strong ellipsis className={classes.username}>
						{username}
					</Typography.Text>
				),
				onClick: () => {
					logout();
				}
			});
		}

		return items;
	}, [isAuthenticated, t, navigate, username, logout]);
	return (
		<Menu theme="dark" selectable={false} mode="horizontal" items={primaryNavs} style={{ flex: 1, minWidth: 0 }} />
	);
};

export default PrimaryMenu;
