import UpdateFilterLink from './UpdateFilterLink';
import UpdateStateTag from './UpdateStateTag';
import { useDeleteUpdateMutation, useModifyUpdateStateMutation } from '../../api/updatesApi';
import { useGetUpdateStateDefinitionsQuery } from '../../api/updateStateDefinitionsApi';
import { useGetUpdateStateTransitionsQuery } from '../../api/updateStateTransitionsApi';
import UpdateSearchIn from '../../constants/api/updateSearchIn';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { UpdateResponse, UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { renderIcon } from '../../utils/iconHelper';
import { getUpdateStateColorFromDefinitions } from '../../utils/updateHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import { SwapOutlined, DeleteOutlined, FieldTimeOutlined, InfoCircleOutlined } from '@ant-design/icons';
import { Badge, Button, Card, Descriptions, DescriptionsProps, Popconfirm, Space, Tooltip, Typography } from 'antd';
import React, { FC, ReactNode, useCallback, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

const { Text } = Typography;

export interface UpdateProps {
	entity: UpdateResponse;
}

const Update: FC<UpdateProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('update');
	const { apiError } = useNotification();
	const { locale } = useLocaleProviderContext();
	const navigate = useNavigate();

	const { data: statesData } = useGetUpdateStateDefinitionsQuery();
	const { data: transitionsData } = useGetUpdateStateTransitionsQuery();

	const redirectToDetails = useCallback(() => {
		navigate(getPageFullPath(`${AppPaths.UPDATES}/${entity.id}`));
	}, [entity.id, navigate]);

	const [deleteUpdate, { isLoading: isDeleteLoading, isError: isErrorDelete, error: deleteError }] =
		useDeleteUpdateMutation();
	const [modifyUpdateState, { isLoading: isModifyLoading, isError: isErrorModify, error: modifyError }] =
		useModifyUpdateStateMutation();

	useEffect(() => {
		if (isErrorDelete) {
			apiError({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					default: t('error_default_delete')
				},
				error: deleteError
			});
		}
	}, [isErrorDelete, deleteError, t, apiError]);

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
		if (!transitionsData?.data?.content || !statesData?.data?.content) {
			return [];
		}

		const currentStateDef = statesData.data.content.find((s) => s.name === entity.state);
		if (!currentStateDef) {
			return [];
		}

		// Find transitions from current state
		const transitions = transitionsData.data.content.filter((t) => t.fromState.id === currentStateDef.id);

		// If no transitions defined from this state, allow all other states (fallback)
		if (transitions.length === 0) {
			return statesData.data.content.filter((s) => s.name !== entity.state);
		}

		return transitions.map((t) => t.toState);
	}, [transitionsData, statesData, entity.state]);

	const stateDefinitions = statesData?.data?.content;
	const color = getUpdateStateColorFromDefinitions(entity.state, stateDefinitions);
	const currentStateDef = stateDefinitions?.find((s) => s.name === entity.state);
	const isInitialState = currentStateDef?.isInitial ?? false;

	// Build transition buttons dynamically
	const transitionButtons = useMemo(() => {
		return allowedTransitions.map((targetState: UpdateStateDefinition) => (
			<Tooltip
				placement="bottom"
				title={t('transition_to', { state: targetState.label, description: targetState.description })}
				key={targetState.id}>
				<Button
					disabled={isModifyLoading || isDeleteLoading}
					loading={isModifyLoading}
					icon={
						renderIcon(targetState.icon, { color: targetState.color }) || (
							<SwapOutlined style={{ color: targetState.color }} />
						)
					}
					type="text"
					onClick={() => modifyUpdateState({ id: entity.id, body: { state: targetState.name } })}
				/>
			</Tooltip>
		));
	}, [allowedTransitions, isModifyLoading, isDeleteLoading, modifyUpdateState, entity.id, t]);

	const buttons = [...transitionButtons];

	const delAction = (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => deleteUpdate({ id: entity.id })}
			okText={t('delete')}
			placement="bottom"
			cancelText={t('cancel')}
			okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
			<Tooltip title={t('help_delete')} placement="bottom">
				<Button
					key="del"
					icon={<DeleteOutlined style={{ color: 'red' }} />}
					type="text"
					loading={isDeleteLoading}
					danger
				/>
			</Tooltip>
		</Popconfirm>
	);
	const detailsAction = (
		<Tooltip placement="bottom" title={t('help_details')}>
			<Button
				key="details"
				icon={<InfoCircleOutlined style={{ color: 'blue' }} />}
				type="text"
				onClick={redirectToDetails}
			/>
		</Tooltip>
	);
	buttons.push(delAction);
	buttons.push(detailsAction);

	const updateDescriptions = useMemo(() => {
		const descriptions: DescriptionsProps['items'] = [
			{
				key: 'version',
				label: (
					<Text type="secondary" ellipsis>
						{t('version')}
					</Text>
				),
				children: <Text>{entity.version}</Text>
			},
			{
				key: 'provider',
				label: (
					<Text type="secondary" ellipsis>
						{t('provider')}
					</Text>
				),
				children: <UpdateFilterLink label={entity.provider} searchIn={UpdateSearchIn.PROVIDER} />
			},
			{
				key: 'host',
				label: (
					<Text type="secondary" ellipsis>
						{t('host')}
					</Text>
				),
				children: <UpdateFilterLink label={entity.host} searchIn={UpdateSearchIn.HOST} />
			},
			{
				key: 'state',
				label: (
					<Text type="secondary" ellipsis>
						{t('state')}
					</Text>
				),
				children: <UpdateStateTag state={entity.state} />
			},
			{
				key: 'created',
				label: (
					<Text type="secondary" ellipsis>
						{t('created')}
					</Text>
				),
				children: (
					<Text>
						{formatDateTimeWithTimeZone(entity.createdAt, DateTimeStyle.SHORT, DateTimeStyle.SHORT, locale)}
					</Text>
				)
			},
			{
				key: 'updated',
				label: (
					<Text type="secondary" ellipsis>
						{t('updated')}
					</Text>
				),
				children: (
					<Text>
						{formatDateTimeWithTimeZone(entity.updatedAt, DateTimeStyle.SHORT, DateTimeStyle.SHORT, locale)}
					</Text>
				)
			}
		];

		return (
			<Descriptions
				items={descriptions}
				colon={false}
				layout="vertical"
				size="small"
				column={{ xs: 2, sm: 3, md: 3, lg: 4, xl: 4, xxl: 4 }}
			/>
		);
	}, [entity.createdAt, entity.host, entity.provider, entity.state, entity.updatedAt, entity.version, locale, t]);

	return (
		<>
			<Card
				size="small"
				style={{ borderColor: color }}
				styles={{ header: { backgroundColor: color, borderColor: color } }}
				loading={isDeleteLoading || isModifyLoading}
				key={entity.id}
				title={
					<Space>
						{isInitialState && (
							<Tooltip title={t('recently_ingested')}>
								<Badge dot />
							</Tooltip>
						)}
						<Button type="link" onClick={redirectToDetails}>
							<Text ellipsis={{ tooltip: entity.application }}>{entity.application}</Text>
						</Button>
					</Space>
				}
				extra={
					<>
						{entity.createdAt == entity.updatedAt && (
							<Tooltip
								placement="bottom"
								title={t('created_at', {
									created: formatDateTimeWithTimeZone(
										entity.createdAt,
										DateTimeStyle.LONG,
										DateTimeStyle.LONG,
										locale
									)
								})}>
								<FieldTimeOutlined />
							</Tooltip>
						)}
						{entity.createdAt !== entity.updatedAt && (
							<Tooltip
								placement="bottom"
								title={t('created_at_diff', {
									created: formatDateTimeWithTimeZone(
										entity.createdAt,
										DateTimeStyle.LONG,
										DateTimeStyle.LONG,
										locale
									),
									updated: formatDateTimeWithTimeZone(
										entity.updatedAt,
										DateTimeStyle.LONG,
										DateTimeStyle.LONG,
										locale
									)
								})}>
								<FieldTimeOutlined />
							</Tooltip>
						)}
					</>
				}
				actions={buttons}>
				{updateDescriptions}
			</Card>
		</>
	);
};

export default Update;
