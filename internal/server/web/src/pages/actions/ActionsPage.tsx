import ActionTextType from './ActionTextType';
import CreateAction from './CreateAction';
import ItemAction from './ItemAction';
import UpdateEnabledAction from './UpdateEnabledAction';
import UpdateLabelAction from './UpdateLabelAction';
import { useGetActionsQuery } from '../../api/actionsApi';
import ActionFilterQueryParamNames from '../../constants/api/actionFilterQueryParamNames';
import ActionOrder from '../../constants/api/actionOrder';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { TABLE_PAGE_DEFAULT, TABLE_PAGE_DEFAULT_OPTIONS, TABLE_PAGE_SIZE_DEFAULT } from '../../constants/pagination';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { ActionResponse, ActionsRequestParams, ActionType } from '../../types';
import useActionsFilterQueryParams from '../../use/useActionsFilterQueryParams';
import { convertToLowerCaseUnderscore } from '../../utils/apiHelper';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { sortAlphaIgnoringCase } from '../../utils/sortHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Result, Skeleton, Space, Switch, Table, TablePaginationConfig, Tooltip, Typography } from 'antd';
import { ColumnsType } from 'antd/es/table';
import { FilterValue, SorterResult } from 'antd/es/table/interface';
import parse from 'html-react-parser';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

const { Text, Title } = Typography;

const DEFAULT_POLLING_INTERVAL = 5000;

const ActionsPage = () => {
	const [t] = useTranslation('actions');
	const { locale } = useLocaleProviderContext();
	const [pollingInterval, setPollingInterval] = useState<number>(0);

	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { orderBy, order, page, pageSize } = useActionsFilterQueryParams();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetActionsQuery(
		{
			orderBy,
			order,
			page,
			pageSize
		} as ActionsRequestParams,
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
			sort: SorterResult<ActionResponse> | SorterResult<ActionResponse>[]
		) => {
			queryParams.delete(ActionFilterQueryParamNames.PAGE);

			if (pagination.pageSize !== pageSize) {
				queryParams.set(ActionFilterQueryParamNames.PAGE, '1');
			} else {
				queryParams.append(ActionFilterQueryParamNames.PAGE, `${pagination.current}`);
			}

			queryParams.delete(ActionFilterQueryParamNames.PAGE_SIZE);
			queryParams.append(ActionFilterQueryParamNames.PAGE_SIZE, `${pagination.pageSize}`);

			queryParams.delete(ActionFilterQueryParamNames.ORDER);
			queryParams.delete(ActionFilterQueryParamNames.ORDER_BY);

			if ((sort as SorterResult<ActionResponse>)?.field && (sort as SorterResult<ActionResponse>)?.order) {
				queryParams.append(
					ActionFilterQueryParamNames.ORDER_BY,
					convertToLowerCaseUnderscore((sort as SorterResult<ActionResponse>).field as string)
				);
				queryParams.append(
					ActionFilterQueryParamNames.ORDER,
					(sort as SorterResult<ActionResponse>).order === 'descend'
						? ActionOrder.DESC
						: (ActionOrder.ASC as string)
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

	const columns: ColumnsType<ActionResponse> = useMemo(() => {
		return [
			{
				title: t('col_label'),
				dataIndex: 'id',
				key: 'label',
				sorter: (a, b) => sortAlphaIgnoringCase(a.label, b.label),
				render: (id: string, entity: ActionResponse) => {
					return <UpdateLabelAction id={id} label={entity.label} />;
				}
			},
			{
				title: t('col_enabled'),
				dataIndex: 'id',
				key: 'enabled',
				sorter: (a, b) => sortAlphaIgnoringCase(a.label, b.label),
				render: (id: string, entity: ActionResponse) => {
					return <UpdateEnabledAction id={id} enabled={entity.enabled} />;
				}
			},
			{
				title: t('col_type'),
				dataIndex: 'type',
				key: 'type',
				sorter: (a, b) => sortAlphaIgnoringCase(a.label, b.label),
				render: (type: ActionType) => {
					return <ActionTextType type={type} />;
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
			}
		];
	}, [locale, t]);

	return (
		<>
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<PageHeader
				className="pl-0"
				title={
					<Title level={4} ellipsis>
						{t('title')}
						<Tooltip placement="bottom" title={parse(t('help'))}>
							<Button icon={<QuestionCircleOutlined />} type="link" />
						</Tooltip>
					</Title>
				}
				extra={
					<Space>
						<Space>
							<Text>{t('auto_refresh')}</Text>
							<Switch
								checkedChildren={t('on')}
								unCheckedChildren={t('off')}
								onChange={onAutoRefreshChange}
								value={pollingInterval > 0}
							/>
						</Space>
						<Tooltip title={t('reload_tooltip')} placement="bottom">
							<Button
								icon={<ReloadOutlined />}
								type="link"
								onClick={invokeReload}
								loading={isFetching}
								disabled={isFetching || isLoading}
							/>
						</Tooltip>
					</Space>
				}
			/>
			<CreateAction />
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && data.data.content.length === 0 && <Result status={404} title={t('no_actions')} />}
			{isSuccess && data.data.content.length > 0 && (
				<Table
					onChange={onTableChange}
					pagination={{
						pageSizeOptions: TABLE_PAGE_DEFAULT_OPTIONS,
						placement: ['bottomCenter'],
						showSizeChanger: true,
						pageSize: data?.data.pageSize || TABLE_PAGE_SIZE_DEFAULT,
						total: data?.data.totalElements,
						current: page || TABLE_PAGE_DEFAULT
					}}
					expandable={{
						expandedRowRender: (record) => <ItemAction e={record} />
					}}
					rowKey="id"
					columns={columns}
					loading={isLoading}
					dataSource={data.data.content}
				/>
			)}
		</>
	);
};

export default ActionsPage;
