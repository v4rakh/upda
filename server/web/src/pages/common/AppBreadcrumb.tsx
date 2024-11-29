import AppPaths from '../../constants/appPaths';
import { getPageFullPath } from '../../utils/urlHelper';
import { HomeOutlined } from '@ant-design/icons';
import { Breadcrumb, Typography } from 'antd';
import { forEach } from 'lodash';
import { FC } from 'react';
import { Link } from 'react-router';

export interface AppBreadcrumbProps {
	items: { label: string; path: string; active?: boolean }[];
	showHome?: boolean;
}

const AppBreadcrumb: FC<AppBreadcrumbProps> = ({ items, showHome = true }) => {
	const shownItems = [];

	if (showHome) {
		shownItems.push({
			title: (
				<Link to={getPageFullPath(AppPaths.HOME)}>
					<HomeOutlined />
				</Link>
			)
		});
	}

	forEach(items, (s) => {
		const active = <Link to={s.path}>{s.label}</Link>;
		const inactive = <Typography.Text ellipsis>{s.label}</Typography.Text>;

		shownItems.push({ title: !s.active ? active : inactive });
	});

	return <Breadcrumb items={shownItems} />;
};

export default AppBreadcrumb;
