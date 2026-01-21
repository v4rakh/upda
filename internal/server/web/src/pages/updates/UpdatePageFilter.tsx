import CreateFilterPreset from '../../components/filter-presets/CreateFilterPreset';
import FilterResetButton from '../../components/filter-presets/FilterResetButton';
import UpdateFilterQueryParamNames from '../../constants/api/updateFilterQueryParamNames';
import UpdateOrder from '../../constants/api/updateOrder';
import UpdateOrderBy from '../../constants/api/updateOrderBy';
import UpdateSearchIn from '../../constants/api/updateSearchIn';
import { UpdateState } from '../../types';
import { FilterPresetType } from '../../types/filterPreset';
import useUpdatesFilterQueryParams from '../../use/useUpdatesFilterQueryParams';
import useUpdateFiltersActive from '../../use/useUpdatesFiltersActive';
import { FilterOutlined, SearchOutlined } from '@ant-design/icons';
import { Badge, Button, Collapse, Divider, Form, Input, Select, Space } from 'antd';
import { compact, forEach, uniq } from 'lodash';
import { FC, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

const COLLAPSE_KEY = 'filter';

const { Search } = Input;

export interface UpdatePageFilterProps {
	loading: boolean;
}

const UpdatePageFilter: FC<UpdatePageFilterProps> = ({ loading }) => {
	const [t] = useTranslation('updates_filters');
	const [form] = Form.useForm();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { searchTerm, searchIn, orderBy, order, state } = useUpdatesFilterQueryParams();
	const { filtersActive } = useUpdateFiltersActive();

	useEffect(() => {
		form.setFieldsValue({
			searchTerm,
			searchIn: searchIn ?? UpdateSearchIn.APPLICATION,
			state,
			orderBy: orderBy ?? UpdateOrderBy.UPDATED_AT,
			order: order ?? UpdateOrder.DESC
		});
	}, [form, order, orderBy, searchIn, searchTerm, state]);

	const onSearchTermChange = useCallback(
		(value: string | undefined) => {
			if (!value) {
				queryParams.delete(UpdateFilterQueryParamNames.SEARCH_TERM);
			} else {
				queryParams.set(UpdateFilterQueryParamNames.SEARCH_TERM, value ?? undefined);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const onSearchInChange = useCallback(
		(value: UpdateSearchIn | undefined) => {
			if (!value) {
				queryParams.delete(UpdateFilterQueryParamNames.SEARCH_IN);
			} else {
				queryParams.set(UpdateFilterQueryParamNames.SEARCH_IN, value ?? undefined);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const onFilterStateChange = useCallback(
		(values: UpdateState[] | undefined) => {
			const all = queryParams.getAll(UpdateFilterQueryParamNames.STATE);
			forEach(all, (v) => {
				queryParams.delete(UpdateFilterQueryParamNames.STATE, v);
			});

			values = uniq(compact(values));
			forEach(values, (v) => {
				queryParams.append(UpdateFilterQueryParamNames.STATE, v);
			});
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const onOrderByChange = useCallback(
		(value: UpdateOrderBy | undefined) => {
			if (!value) {
				queryParams.delete(UpdateFilterQueryParamNames.ORDER_BY);
			} else {
				queryParams.set(UpdateFilterQueryParamNames.ORDER_BY, value ?? undefined);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const onOrderChange = useCallback(
		(value: UpdateOrder | undefined) => {
			if (!value) {
				queryParams.delete(UpdateFilterQueryParamNames.ORDER);
			} else {
				queryParams.set(UpdateFilterQueryParamNames.ORDER, value);
			}
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const filterButtonActive = useMemo(() => {
		return (
			<Badge dot offset={[-10, 5]}>
				<Button type="link" icon={<FilterOutlined />}>
					{t('filters')}
				</Button>
			</Badge>
		);
	}, [t]);

	const filterButtonInactive = useMemo(() => {
		return (
			<Button type="link" icon={<FilterOutlined />}>
				{t('filters')}
			</Button>
		);
	}, [t]);

	return (
		<Collapse
			onChange={(keys) => {
				setCollapseActiveKeys(keys);
			}}
			expandIconPlacement="end"
			bordered={false}
			ghost
			size="small"
			activeKey={collapseActiveKeys}
			items={[
				{
					key: COLLAPSE_KEY,
					showArrow: false,
					label: filtersActive ? filterButtonActive : filterButtonInactive,
					children: (
						<Space orientation="vertical">
							<Form layout="inline" form={form} disabled={loading}>
								<Form.Item label={t('search_term')} name="searchTerm" tooltip={t('search_term_help')}>
									<Search
										variant="filled"
										loading={loading}
										maxLength={255}
										placeholder={t('search_term_placeholder')}
										allowClear
										enterButton={<Button type="link" icon={<SearchOutlined />} />}
										onSearch={onSearchTermChange}
									/>
								</Form.Item>
								<Form.Item label={t('search_in')} name="searchIn" tooltip={t('search_in_help')}>
									<Select
										variant="filled"
										style={{ width: 120 }}
										onChange={onSearchInChange}
										options={[
											{
												value: UpdateSearchIn.APPLICATION,
												label: t(`search_in_${UpdateSearchIn.APPLICATION.toLowerCase()}`)
											},
											{
												value: UpdateSearchIn.PROVIDER,
												label: t(`search_in_${UpdateSearchIn.PROVIDER.toLowerCase()}`)
											},
											{
												value: UpdateSearchIn.HOST,
												label: t(`search_in_${UpdateSearchIn.HOST.toLowerCase()}`)
											}
										]}
									/>
								</Form.Item>
								<Form.Item label={t('state')} name="state" tooltip={t('state_help')}>
									<Select
										variant="filled"
										mode="multiple"
										allowClear
										placeholder={t('state_placeholder')}
										style={{ width: '100%', minWidth: 200 }}
										onChange={onFilterStateChange}
										options={[
											{
												value: UpdateState.PENDING,
												label: t(`state_${UpdateState.PENDING.toLowerCase()}`)
											},
											{
												value: UpdateState.IGNORED,
												label: t(`state_${UpdateState.IGNORED.toLowerCase()}`)
											},
											{
												value: UpdateState.APPROVED,
												label: t(`state_${UpdateState.APPROVED.toLowerCase()}`)
											}
										]}
									/>
								</Form.Item>
								<Form.Item label={t('order_by')} name="orderBy" tooltip={t('order_by_help')}>
									<Select
										variant="filled"
										style={{ width: 120 }}
										onChange={onOrderByChange}
										options={[
											{
												value: UpdateOrderBy.ID,
												label: t(`order_by_${UpdateOrderBy.ID.toLowerCase()}`)
											},
											{
												value: UpdateOrderBy.CREATED_AT,
												label: t(`order_by_${UpdateOrderBy.CREATED_AT.toLowerCase()}`)
											},
											{
												value: UpdateOrderBy.UPDATED_AT,
												label: t(`order_by_${UpdateOrderBy.UPDATED_AT.toLowerCase()}`)
											},
											{
												value: UpdateOrderBy.APPLICATION,
												label: t(`order_by_${UpdateOrderBy.APPLICATION.toLowerCase()}`)
											},
											{
												value: UpdateOrderBy.PROVIDER,
												label: t(`order_by_${UpdateOrderBy.PROVIDER.toLowerCase()}`)
											},
											{
												value: UpdateOrderBy.HOST,
												label: t(`order_by_${UpdateOrderBy.HOST.toLowerCase()}`)
											}
										]}
									/>
								</Form.Item>
								<Form.Item name="order">
									<Select
										variant="filled"
										style={{ width: 120 }}
										onChange={onOrderChange}
										options={[
											{
												value: UpdateOrder.DESC,
												label: t(`order_${UpdateOrder.DESC.toLowerCase()}`)
											},
											{
												value: UpdateOrder.ASC,
												label: t(`order_${UpdateOrder.ASC.toLowerCase()}`)
											}
										]}
									/>
								</Form.Item>
								<FilterResetButton enabled={filtersActive} />
							</Form>
							{filtersActive && (
								<>
									<Divider />
									<CreateFilterPreset type={FilterPresetType.UPDATE} />
								</>
							)}
						</Space>
					)
				}
			]}
		/>
	);
};

export default UpdatePageFilter;
