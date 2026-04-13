import CreateStateTransition from './CreateStateTransition';
import DeleteStateTransition from './DeleteStateTransition';
import { useGetUpdateStateTransitionsQuery } from '../../api/updateStateTransitionsApi';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { UpdateStateTransition } from '../../types';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { renderIcon } from '../../utils/iconHelper';
import { sortAlphaIgnoringCase } from '../../utils/sortHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Col, Result, Row, Skeleton, Table, Tag, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import parse from 'html-react-parser';
import { FC, useCallback, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

const StateTransitionsPage: FC = () => {
	const [t] = useTranslation('state_transitions');
	const { locale } = useLocaleProviderContext();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetUpdateStateTransitionsQuery();

	const sortedData = useMemo(() => {
		if (!data?.data?.content) return [];
		return [...data.data.content].sort((a, b) => {
			const fromStateOrder = a.fromState.sortOrder - b.fromState.sortOrder;
			if (fromStateOrder !== 0) return fromStateOrder;
			return a.toState.sortOrder - b.toState.sortOrder;
		});
	}, [data]);

	const invokeReload = useCallback(() => {
		refetch();
	}, [refetch]);

	const columns: ColumnsType<UpdateStateTransition> = useMemo(() => {
		return [
			{
				title: t('col_from_state'),
				dataIndex: 'fromState',
				key: 'fromState',
				sorter: (a, b) => a.fromState.sortOrder - b.fromState.sortOrder,
				render: (_, entity: UpdateStateTransition) => {
					return (
						<Tag
							color={entity.fromState.color}
							icon={renderIcon(entity.fromState.icon, { marginRight: 4 })}>
							{entity.fromState.label}
						</Tag>
					);
				}
			},
			{
				title: t('col_to_state'),
				dataIndex: 'toState',
				key: 'toState',
				sorter: (a, b) => a.toState.sortOrder - b.toState.sortOrder,
				render: (_, entity: UpdateStateTransition) => {
					return (
						<Tag color={entity.toState.color} icon={renderIcon(entity.toState.icon, { marginRight: 4 })}>
							{entity.toState.label}
						</Tag>
					);
				}
			},
			{
				title: t('col_created_at'),
				dataIndex: 'createdAt',
				key: 'createdAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortAlphaIgnoringCase(a.createdAt, b.createdAt),
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('col_updated_at'),
				dataIndex: 'updatedAt',
				key: 'updatedAt',
				ellipsis: true,
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortAlphaIgnoringCase(a.updatedAt, b.updatedAt),
				render: (value) => formatDateTimeWithTimeZone(value, DateTimeStyle.LONG, DateTimeStyle.MEDIUM, locale)
			},
			{
				title: t('actions'),
				dataIndex: 'id',
				key: 'actions',
				ellipsis: false,
				render: (id: string) => {
					return <DeleteStateTransition id={id} />;
				}
			}
		];
	}, [locale, t]);

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
				extra={
					<Tooltip title={t('reload_tooltip')} placement="bottom">
						<Button
							icon={<ReloadOutlined />}
							type="link"
							onClick={invokeReload}
							loading={isFetching}
							disabled={isFetching || isLoading}
						/>
					</Tooltip>
				}
			/>
			<CreateStateTransition />
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && sortedData.length === 0 && <Result status={404} title={t('no_state_transitions')} />}
			{isSuccess && sortedData.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} lg={24}>
						<Table
							rowKey="id"
							columns={columns}
							loading={isLoading}
							dataSource={sortedData}
							pagination={false}
						/>
					</Col>
				</Row>
			)}
		</>
	);
};

export default StateTransitionsPage;
