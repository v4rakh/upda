import EventsTree from './EventsTree';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { QuestionCircleOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Tooltip, Typography } from 'antd';
import parse from 'html-react-parser';
import React, { FC } from 'react';
import { useTranslation } from 'react-i18next';

const EventsPage: FC = () => {
	const [t] = useTranslation('events');

	return (
		<>
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<PageHeader
				className="pl-0"
				title={
					<Typography.Title level={4} ellipsis>
						{t('title')}
						<Tooltip placement="bottom" title={parse(t('help'))}>
							<Button icon={<QuestionCircleOutlined />} type="link" />
						</Tooltip>
					</Typography.Title>
				}
			/>
			<EventsTree />
		</>
	);
};

export default EventsPage;
