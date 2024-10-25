import DeleteActionInvocation from './DeleteActionInvocation';
import { useGetActionByIdQuery } from '../../api/actionsApi';
import { useGetEventByIdQuery } from '../../api/eventsApi';
import { ActionInvocationResponse } from '../../types';
import ActionTextType from '../actions/ActionTextType';
import EventText from '../events/EventText';
import { ReloadOutlined } from '@ant-design/icons';
import { Button, Descriptions, Result, Skeleton, Tooltip, Typography } from 'antd';
import { FC, useCallback } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;
const { Item } = Descriptions;

export interface ItemActionInvocationProps {
	e: ActionInvocationResponse;
}

const ItemActionInvocation: FC<ItemActionInvocationProps> = ({ e }): JSX.Element => {
	const [t] = useTranslation('action_invocation_item');

	const {
		isLoading: isLoadingAction,
		isError: isErrorAction,
		isSuccess: isSuccessAction,
		isFetching: isFetchingAction,
		data: dataAction,
		refetch: refetchAction
	} = useGetActionByIdQuery({ id: e.actionId });

	const {
		isLoading: isLoadingEvent,
		isError: isErrorEvent,
		isSuccess: isSuccessEvent,
		isFetching: isFetchingEvent,
		data: dataEvent,
		refetch: refetchEvent
	} = useGetEventByIdQuery({ id: e.eventId });

	const invokeActionReload = useCallback(() => {
		refetchAction();
	}, [refetchAction]);

	const invokeEventReload = useCallback(() => {
		refetchEvent();
	}, [refetchEvent]);

	return (
		<>
			{e.message && (
				<Descriptions layout="vertical" size="small" colon={false}>
					<Item label={t('message')}>
						<Text type="danger">{e.message}</Text>
					</Item>
				</Descriptions>
			)}
			{isLoadingAction && <Skeleton loading={isLoadingAction} active={isLoadingAction} />}
			{isErrorAction && (
				<Result
					status="error"
					title={t('error_default_loading_action')}
					extra={
						<Tooltip title={t('reload_tooltip')} placement={'bottom'}>
							<Button
								icon={<ReloadOutlined />}
								type={'link'}
								onClick={invokeActionReload}
								loading={isFetchingAction}
								disabled={isFetchingAction}>
								{t('reload_text')}
							</Button>
						</Tooltip>
					}
				/>
			)}
			{isSuccessAction && dataAction.data && (
				<Descriptions title={t('action')} layout="vertical" size="small" colon={false}>
					<Item label={t('label')}>{dataAction.data.label}</Item>
					<Item label={t('type')}>
						<ActionTextType type={dataAction.data.type} />
					</Item>
				</Descriptions>
			)}
			{isLoadingEvent && <Skeleton loading={isLoadingEvent} active={isLoadingEvent} />}
			{isErrorEvent && (
				<Result
					status="error"
					title={t('error_default_loading_event')}
					extra={
						<Tooltip title={t('reload_tooltip')} placement={'bottom'}>
							<Button
								icon={<ReloadOutlined />}
								type={'link'}
								onClick={invokeEventReload}
								loading={isFetchingEvent}
								disabled={isFetchingEvent}>
								{t('reload_text')}
							</Button>
						</Tooltip>
					}
				/>
			)}
			{isSuccessEvent && dataEvent.data && (
				<Descriptions title={t('event')} layout="vertical" size="small" colon={false}>
					<Item label={t('name')}>
						<EventText name={dataEvent.data.name} payload={dataEvent.data.payload} />
					</Item>
				</Descriptions>
			)}
			<Descriptions size="small" layout="vertical" colon={false}>
				<Item>
					<DeleteActionInvocation id={e.id} />
				</Item>
			</Descriptions>
		</>
	);
};

export default ItemActionInvocation;
