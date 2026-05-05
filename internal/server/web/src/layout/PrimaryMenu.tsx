import classes from './style/PrimaryMenu.module.less';

import { useAuth } from '../auth/AuthContext';
import AppPaths from '../constants/appPaths';
import { getPageFullPath } from '../utils/urlHelper';
import {
	BarsOutlined,
	BuildOutlined,
	ClockCircleOutlined,
	FolderOutlined,
	LinkOutlined,
	LoadingOutlined,
	LockOutlined,
	LogoutOutlined,
	PartitionOutlined,
	SwapOutlined,
	UnorderedListOutlined,
	UserOutlined
} from '@ant-design/icons';
import { Menu, Typography } from 'antd';
import { TFunction } from 'i18next';
import { forEach, noop } from 'lodash';
import { FC, ReactNode, useCallback, useMemo, useState } from 'react';
import { useNavigate } from 'react-router';

export interface PrimaryMenuProps {
	t: TFunction;
}

const { Text } = Typography;

const PrimaryMenu: FC<PrimaryMenuProps> = ({ t }): ReactNode => {
	const navigate = useNavigate();
	const { logout, profile, isAuthenticated } = useAuth();
	const [isLoggingOut, setIsLoggingOut] = useState(false);

	const preLogout = useCallback(() => {
		setIsLoggingOut(true);
	}, []);

	const postLogout = useCallback(() => {
		setIsLoggingOut(false);
	}, []);

	const onLogoutClicked = useCallback(() => {
		logout({ preLogout: preLogout, postLogout: postLogout });
	}, [logout, postLogout, preLogout]);

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
				label: t('constants'),
				key: AppPaths.CONSTANTS,
				icon: <FolderOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.CONSTANTS))
			},
			isAuthenticated && {
				label: t('state_definitions'),
				key: AppPaths.STATE_DEFINITIONS,
				icon: <PartitionOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.STATE_DEFINITIONS))
			},
			isAuthenticated && {
				label: t('state_transitions'),
				key: AppPaths.STATE_TRANSITIONS,
				icon: <SwapOutlined />,
				onClick: () => navigate(getPageFullPath(AppPaths.STATE_TRANSITIONS))
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
				icon: isLoggingOut ? <LoadingOutlined /> : <LogoutOutlined />,
				label: (
					<Text strong ellipsis className={classes.username}>
						{profile?.preferredUsername}
					</Text>
				),
				onClick: isLoggingOut ? noop : onLogoutClicked
			});
		}

		return items;
	}, [isAuthenticated, t, navigate, isLoggingOut, profile?.preferredUsername, onLogoutClicked]);
	return (
		<Menu
			theme="dark"
			selectable={false}
			mode="horizontal"
			items={primaryNavs}
			style={{ flex: 1, minWidth: 0, borderBottom: 'none' }}
		/>
	);
};

export default PrimaryMenu;
