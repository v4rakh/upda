import UpdateFilterLink from './UpdateFilterLink';
import UpdateMetadata from './UpdateMetadata';
import UpdateStateTag from './UpdateStateTag';
import { useGetUpdateByIdQuery, useModifyUpdateStateMutation } from '../../api/updatesApi';
import { useGetUpdateStateDefinitionsQuery } from '../../api/updateStateDefinitionsApi';
import { useGetUpdateStateTransitionsQuery } from '../../api/updateStateTransitionsApi';
import Comments from '../../components/comments/Comments';
import UpdateSearchIn from '../../constants/api/updateSearchIn';
import UpdateTabQueryParamNames from '../../constants/api/updateTabQueryParamNames';
import ApiErrorCodes from '../../constants/apiErrorCodes';
import AppPathParamNames from '../../constants/appPathParamNames';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import UpdateTabNames from '../../constants/updateTabNames';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import useUpdateTabQueryParams from '../../use/useUpdateTabQueryParams';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { renderIcon } from '../../utils/iconHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import AppBreadcrumb from '../common/AppBreadcrumb';
import EventsTree from '../events/EventsTree';
import {
	ClockCircleOutlined,
	CloseOutlined,
	CommentOutlined,
	InfoCircleOutlined,
	MoreOutlined,
	SwapOutlined
} from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Descriptions, FloatButton, Result, Skeleton, Tabs, Tooltip, Typography } from 'antd';
import React, { FC, ReactNode, useCallback, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams, useSearchParams } from 'react-router';

const { Text } = Typography;
const { Group } = FloatButton;
const { Item } = Descriptions;

const UpdateSinglePage: FC = (): ReactNode => {
	const [t] = useTranslation('updates_single');
	const [tUpdate] = useTranslation('update');
	const { locale } = useLocaleProviderContext();
	const { apiError } = useNotification();

	const { [AppPathParamNames.UPDATE_ID]: updateId } = useParams();

	const [queryParams, setSearchQueryParams] = useSearchParams();
	const { tab } = useUpdateTabQueryParams();

	const onTabChange = useCallback(
		(value: string | undefined) => {
			queryParams.set(UpdateTabQueryParamNames.TAB, value ?? UpdateTabNames.DETAILS);
			setSearchQueryParams(queryParams);
		},
		[queryParams, setSearchQueryParams]
	);

	const { isLoading, isError, isSuccess, data, error } = useGetUpdateByIdQuery(
		{
			id: updateId || ''
		},
		{
			skip: updateId === undefined
		}
	);

	const { data: statesData } = useGetUpdateStateDefinitionsQuery();
	const { data: transitionsData } = useGetUpdateStateTransitionsQuery();

	const [modifyUpdateState, { isLoading: isModifyLoading, isError: isErrorModify, error: modifyError }] =
		useModifyUpdateStateMutation();

	useEffect(() => {
		if (isErrorModify) {
			apiError({
				i18n: {
					notFound: t('error_unable_modify_state'),
					unAuthorized: t('error_unauthorized_modify_state'),
					forbidden: t('error_forbidden_modify_state'),
					default: t('error_default_modify_state')
				},
				error: modifyError
			});
		}
	}, [apiError, isErrorModify, modifyError, t]);

	// Get allowed transitions for the current state
	const allowedTransitions = useMemo(() => {
		if (!transitionsData?.data?.content || !statesData?.data?.content || !data?.data?.state) {
			return [];
		}

		const currentStateDef = statesData.data.content.find((s) => s.name === data.data.state);
		if (!currentStateDef) {
			return [];
		}

		// Find transitions from current state
		const transitions = transitionsData.data.content.filter((t) => t.fromState.id === currentStateDef.id);

		// If no transitions defined from this state, allow all other states (fallback)
		if (transitions.length === 0) {
			return statesData.data.content.filter((s) => s.name !== data.data.state);
		}

		return transitions.map((t) => t.toState);
	}, [transitionsData, statesData, data?.data?.state]);

	const tabs = useMemo(
		() => [
			{
				key: UpdateTabNames.DETAILS,
				label: t('details'),
				icon: <InfoCircleOutlined />,
				children: data ? (
					<>
						<Descriptions
							layout="vertical"
							size="small"
							column={{ xs: 1, sm: 2, md: 4, lg: 4, xl: 4, xxl: 4 }}
							colon={false}>
							<Item label={t('state')}>
								<UpdateStateTag state={data.data.state} />
							</Item>
							<Item label={t('application')}>
								<UpdateFilterLink label={data.data.application} searchIn={UpdateSearchIn.APPLICATION} />
							</Item>
							<Item label={t('version')}>
								<Text code>{data.data.version}</Text>
							</Item>
							<Item label={t('host')}>
								<UpdateFilterLink label={data.data.host} searchIn={UpdateSearchIn.HOST} />
							</Item>
							<Item label={t('provider')}>
								<UpdateFilterLink label={data.data.provider} searchIn={UpdateSearchIn.PROVIDER} />
							</Item>
							<Item label={t('created')}>
								{formatDateTimeWithTimeZone(
									data.data.createdAt,
									DateTimeStyle.LONG,
									DateTimeStyle.MEDIUM,
									locale
								)}
							</Item>
							<Item label={t('updated')}>
								{formatDateTimeWithTimeZone(
									data.data.updatedAt,
									DateTimeStyle.LONG,
									DateTimeStyle.MEDIUM,
									locale
								)}
							</Item>
						</Descriptions>
						<Descriptions
							layout="vertical"
							size="small"
							column={{ xs: 1, sm: 2, md: 4, lg: 4, xl: 4, xxl: 4 }}
							colon={false}>
							<Item label={t('metadata')}>
								<UpdateMetadata metadata={data.data.metadata} />
							</Item>
						</Descriptions>
					</>
				) : (
					<></>
				)
			},
			{
				key: UpdateTabNames.COMMENTS,
				label: t('comments'),
				icon: <CommentOutlined />,
				children: updateId && <Comments updateId={updateId} />
			},
			{
				key: UpdateTabNames.EVENTS,
				label: t('events'),
				icon: <ClockCircleOutlined />,
				children: updateId && <EventsTree updateId={updateId} />
			}
		],
		[t, data, locale, updateId]
	);

	const floatingActions = useMemo(() => {
		const transitionButtons = allowedTransitions.map((targetState: UpdateStateDefinition) => (
			<Tooltip
				placement="left"
				title={tUpdate('transition_to', { state: targetState.label, description: targetState.description })}
				key={targetState.id}>
				<FloatButton
					icon={
						renderIcon(targetState.icon, { color: targetState.color }) || (
							<SwapOutlined style={{ color: targetState.color }} />
						)
					}
					onClick={() => {
						if (updateId) {
							modifyUpdateState({ id: updateId, body: { state: targetState.name } });
						}
					}}
				/>
			</Tooltip>
		));

		return (
			<Group trigger="click" type="primary" icon={<MoreOutlined />} closeIcon={<CloseOutlined />}>
				{transitionButtons}
			</Group>
		);
	}, [allowedTransitions, tUpdate, updateId, modifyUpdateState]);

	return (
		<>
			<AppBreadcrumb
				items={[
					{
						label: t('title_parent'),
						path: getPageFullPath(AppPaths.UPDATES)
					},
					{
						label: data?.data.application || '',
						active: true,
						path: ''
					}
				]}
			/>
			<PageHeader
				className="pl-0"
				title={
					<Typography.Title level={4} ellipsis>
						{t('title', { name: data?.data.application })}
					</Typography.Title>
				}
			/>
			{isLoading && <Skeleton loading={isLoading} active={isLoading} />}
			{isError && (error as any)?.data?.status !== ApiErrorCodes.NOT_FOUND && (
				<Result status={500} title={t('default_loading_error_message')} />
			)}
			{isError && (error as any)?.data?.status === ApiErrorCodes.NOT_FOUND && (
				<Result status={404} title={t('no_update')} />
			)}
			{isSuccess && (
				<>
					<Tabs activeKey={tab} items={tabs} onChange={onTabChange} destroyOnHidden />
					{!isModifyLoading && floatingActions}
				</>
			)}
		</>
	);
};

export default UpdateSinglePage;
