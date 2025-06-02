import Update from './Update';
import UpdatePageFilter from './UpdatePageFilter';
import { useGetUpdatesQuery } from '../../api/updatesApi';
import FilterPreset from '../../components/filter-presets/FilterPresets';
import UpdateFilterQueryParamNames from '../../constants/api/updateFilterQueryParamNames';
import { CARD_PAGE_DEFAULT, CARD_PAGE_DEFAULT_OPTIONS, CARD_PAGE_SIZE_DEFAULT } from '../../constants/pagination';
import { UpdatesRequestParams } from '../../types';
import { FilterPresetType } from '../../types/filterPreset';
import { useResponsiveGridSize } from '../../use/useResponsiveGridSize';
import useUpdatesFilterQueryParams from '../../use/useUpdatesFilterQueryParams';
import useUpdateFiltersActive from '../../use/useUpdatesFiltersActive';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Col, List, Result, Row, Skeleton, Space, Switch, Tooltip, Typography } from 'antd';
import parse from 'html-react-parser';
import { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

const { Text, Title } = Typography;

const DEFAULT_POLLING_INTERVAL = 10000;

const UpdatesPage: FC = () => {
	const [t] = useTranslation('updates');

	const { gridSize } = useResponsiveGridSize({ xxl: 4, xl: 3, lg: 2, md: 2, sm: 1, xs: 1 });

	const [pollingInterval, setPollingInterval] = useState<number>(0);
	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { orderBy, order, page, pageSize, state, searchIn, searchTerm } = useUpdatesFilterQueryParams();
	const { filtersActive } = useUpdateFiltersActive();

	const { isLoading, isError, refetch, isFetching, isSuccess, data } = useGetUpdatesQuery(
		{
			orderBy,
			order,
			page,
			pageSize,
			state,
			searchIn,
			searchTerm
		} as UpdatesRequestParams,
		{
			pollingInterval
		}
	);

	const invokeReload = useCallback(() => {
		refetch();
	}, [refetch]);

	useEffect(() => {
		if (isError) {
			setPollingInterval(0);
		}
	}, [isError, setPollingInterval, t]);

	const onPaginationChange = useCallback(
		(pageSelected: number, pageSizeSelected: number) => {
			queryParams.delete(UpdateFilterQueryParamNames.PAGE);
			if (pageSize !== pageSizeSelected) {
				queryParams.set(UpdateFilterQueryParamNames.PAGE, '1');
			} else {
				queryParams.append(UpdateFilterQueryParamNames.PAGE, `${pageSelected}`);
			}
			queryParams.delete(UpdateFilterQueryParamNames.PAGE_SIZE);
			queryParams.append(UpdateFilterQueryParamNames.PAGE_SIZE, `${pageSizeSelected}`);

			setSearchQueryParams(queryParams);
		},
		[pageSize, queryParams, setSearchQueryParams]
	);

	const onAutoRefreshChange = useCallback(
		(checked: boolean) => {
			if (checked) {
				setPollingInterval(DEFAULT_POLLING_INTERVAL);
			} else {
				setPollingInterval(0);
			}
		},
		[setPollingInterval]
	);

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
			<UpdatePageFilter loading={isLoading || isFetching} />
			<FilterPreset type={FilterPresetType.UPDATE} showFilterReset filtersActive={filtersActive} />
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && <Result status="error" title={t('error_default_loading')} />}
			{isSuccess && data.data.content.length === 0 && <Result status={404} title={t('no_updates')} />}
			{isSuccess && data.data.content.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} lg={24}>
						<List
							pagination={{
								position: 'bottom',
								align: 'center',
								pageSize: data?.data.pageSize || CARD_PAGE_SIZE_DEFAULT,
								pageSizeOptions: CARD_PAGE_DEFAULT_OPTIONS,
								total: data?.data.totalElements || 0,
								current: page || CARD_PAGE_DEFAULT,
								onChange: onPaginationChange,
								showSizeChanger: true
							}}
							grid={{
								gutter: [
									{ xs: 8, sm: 16, md: 24, lg: 32 },
									{ xs: 8, sm: 16, md: 24, lg: 32 }
								],
								column: data.data.content.length <= gridSize ? data.data.content.length : gridSize
							}}
							dataSource={data.data.content}
							renderItem={(e) => (
								<List.Item>
									<Update key={e.id} entity={e} />
								</List.Item>
							)}></List>
					</Col>
				</Row>
			)}
		</>
	);
};

export default UpdatesPage;
