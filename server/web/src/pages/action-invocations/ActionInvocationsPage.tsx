import ItemActionInvocation from './ItemActionInvocation';
import { useGetActionInvocationsQuery } from '../../api/actionInvocationsApi';
import ActionInvocationFilterQueryParamNames from '../../constants/api/actionInvocationFilterQueryParamNames';
import ActionInvocationOrder from '../../constants/api/actionInvocationOrder';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { PAGE_DEFAULT, PAGE_DEFAULT_OPTIONS, PAGE_SIZE_DEFAULT } from '../../constants/pagination';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { ActionInvocationResponse, ActionInvocationsRequestParams, ActionInvocationState } from '../../types';
import useActionInvocationsFilterQueryParams from '../../use/useActionInvocationsFilterQueryParams';
import { convertToLowerCaseUnderscore } from '../../utils/apiHelper';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { sortAlphaIgnoringCase, sortNumber } from '../../utils/sortHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import {
	CheckCircleOutlined,
	ClockCircleOutlined,
	ExclamationCircleOutlined,
	LoadingOutlined,
	QuestionCircleOutlined,
	ReloadOutlined
} from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Result, Skeleton, Space, Switch, Table, TablePaginationConfig, Tag, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import { FilterValue, SorterResult } from 'antd/es/table/interface';
import parse from 'html-react-parser';
import { ReactNode, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

const DEFAULT_POLLING_INTERVAL = 5000;

const ActionInvocationsPage = () => {
	const [t] = useTranslation('action_invocations');
	const { locale } = useLocaleProviderContext();
	const [pollingInterval, setPollingInterval] = useState<number>(DEFAULT_POLLING_INTERVAL);

	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { orderBy, order, page, pageSize } = useActionInvocationsFilterQueryParams();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetActionInvocationsQuery(
		{
			orderBy,
			order,
			page,
			pageSize
		} as ActionInvocationsRequestParams,
		{
			pollingInterval
		}
	);

	const invokeReload = useCallback(() => {
		refetch();
	}, [refetch]);

	const onAutoRefreshChange = (checked: boolean) => {
		if (checked) {
			setPollingInterval(DEFAULT_POLLING_INTERVAL);
		} else {
			setPollingInterval(0);
		}
	};

	const onTableChange = useCallback(
		(
			pagination: TablePaginationConfig,
			filters: Record<string, FilterValue | null>,
			sort: SorterResult<ActionInvocationResponse> | SorterResult<ActionInvocationResponse>[]
		) => {
			queryParams.delete(ActionInvocationFilterQueryParamNames.PAGE);

			if (pagination.pageSize !== pageSize) {
				queryParams.set(ActionInvocationFilterQueryParamNames.PAGE, '1');
			} else {
				queryParams.append(ActionInvocationFilterQueryParamNames.PAGE, `${pagination.current}`);
			}

			queryParams.delete(ActionInvocationFilterQueryParamNames.PAGE_SIZE);
			queryParams.append(ActionInvocationFilterQueryParamNames.PAGE_SIZE, `${pagination.pageSize}`);

			queryParams.delete(ActionInvocationFilterQueryParamNames.ORDER);
			queryParams.delete(ActionInvocationFilterQueryParamNames.ORDER_BY);

			if (
				(sort as SorterResult<ActionInvocationResponse>)?.field &&
				(sort as SorterResult<ActionInvocationResponse>)?.order
			) {
				queryParams.append(
					ActionInvocationFilterQueryParamNames.ORDER_BY,
					convertToLowerCaseUnderscore((sort as SorterResult<ActionInvocationResponse>).field as string)
				);
				queryParams.append(
					ActionInvocationFilterQueryParamNames.ORDER,
					(sort as SorterResult<ActionInvocationResponse>).order === 'descend'
						? ActionInvocationOrder.DESC
						: (ActionInvocationOrder.ASC as string)
				);
			}
			setSearchQueryParams(queryParams);
		},
		[pageSize, queryParams, setSearchQueryParams]
	);

	useEffect(() => {
		if (isError) {
			setPollingInterval(0);
		}
	}, [isError, setPollingInterval, t]);

	const renderState = useCallback(
		(state: ActionInvocationState): ReactNode => {
			switch (state) {
				case ActionInvocationState.CREATED:
					return (
						<Tooltip title={t(`state_${state.toLowerCase()}_description`)}>
							<Tag icon={<ClockCircleOutlined />} color={'gray'}>
								{t(`state_${state.toLowerCase()}`)}
							</Tag>
						</Tooltip>
					);
				case ActionInvocationState.ERROR:
					return (
						<Tooltip title={t(`state_${state.toLowerCase()}_description`)}>
							<Tag icon={<ExclamationCircleOutlined />} color={'red'}>
								{t(`state_${state.toLowerCase()}`)}
							</Tag>
						</Tooltip>
					);
				case ActionInvocationState.RETRYING:
				case ActionInvocationState.RUNNING:
					return (
						<Tooltip title={t(`state_${state.toLowerCase()}_description`)}>
							<Tag icon={<LoadingOutlined />} color={'blue'}>
								{t(`state_${state.toLowerCase()}`)}
							</Tag>
						</Tooltip>
					);
				case ActionInvocationState.SUCCESS:
					return (
						<Tooltip title={t(`state_${state.toLowerCase()}_description`)}>
							<Tag icon={<CheckCircleOutlined />} color={'green'}>
								{t(`state_${state.toLowerCase()}`)}
							</Tag>
						</Tooltip>
					);
			}
		},
		[t]
	);

	const columns: ColumnsType<ActionInvocationResponse> = useMemo(() => {
		return [
			{
				title: t('col_state'),
				dataIndex: 'state',
				key: 'state',
				responsive: ['xs', 'sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortAlphaIgnoringCase(a.state, b.state),
				render: (value) => {
					return renderState(value);
				}
			},
			{
				title: t('col_retry_count'),
				dataIndex: 'retryCount',
				key: 'retryCount',
				responsive: ['sm', 'md', 'lg', 'xl', 'xxl'],
				sorter: (a, b) => sortNumber(a.retryCount, b.retryCount)
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
			}
		];
	}, [locale, renderState, t]);

	return (
		<>
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<PageHeader
				className={'pl-0'}
				title={
					<Typography.Title level={4} ellipsis>
						{t('title')}
						<Tooltip placement="bottom" title={parse(t('help'))}>
							<Button icon={<QuestionCircleOutlined />} type="link" />
						</Tooltip>
					</Typography.Title>
				}
				extra={
					<Space>
						<Space>
							<Typography.Text>{t('auto_refresh')}</Typography.Text>
							<Switch
								checkedChildren={t('on')}
								unCheckedChildren={t('off')}
								onChange={onAutoRefreshChange}
								value={pollingInterval > 0}
							/>
						</Space>
						<Tooltip title={t('reload_tooltip')} placement={'bottom'}>
							<Button
								icon={<ReloadOutlined />}
								type={'link'}
								onClick={invokeReload}
								loading={isFetching}
								disabled={isFetching || isLoading}
							/>
						</Tooltip>
					</Space>
				}
			/>
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && data.data.content.length === 0 && <Result status={404} title={t('no_action_invocations')} />}
			{isSuccess && data.data.content.length > 0 && (
				<Table
					onChange={onTableChange}
					pagination={{
						pageSizeOptions: PAGE_DEFAULT_OPTIONS,
						position: ['bottomCenter'],
						showSizeChanger: true,
						pageSize: data?.data.pageSize || PAGE_SIZE_DEFAULT,
						total: data?.data.totalElements,
						current: page || PAGE_DEFAULT
					}}
					expandable={{
						expandedRowRender: (record) => <ItemActionInvocation e={record} />,
						expandRowByClick: true
					}}
					rowKey={'id'}
					columns={columns}
					loading={isLoading}
					dataSource={data.data.content}
				/>
			)}
		</>
	);
};

export default ActionInvocationsPage;
