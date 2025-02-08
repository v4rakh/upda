import UpdateStateTag from './UpdateStateTag';
import { useDeleteUpdateMutation, useModifyUpdateStateMutation } from '../../api/updatesApi';
import AppPaths from '../../constants/appPaths';
import DateTimeStyle from '../../constants/dateTimeStyle';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { UpdateResponse, UpdateState } from '../../types';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { getUpdateStateColor } from '../../utils/updateHelper';
import { getPageFullPath } from '../../utils/urlHelper';
import { apiNotification } from '../common/apiNotification';
import {
	CheckCircleOutlined,
	DeleteOutlined,
	FieldTimeOutlined,
	InfoCircleOutlined,
	InteractionOutlined,
	StopOutlined
} from '@ant-design/icons';
import { Badge, Button, Card, Col, Popconfirm, Row, Space, Tooltip, Typography } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

const { Text } = Typography;

export interface UpdateProps {
	entity: UpdateResponse;
}

const Update: FC<UpdateProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('update');
	const { locale } = useLocaleProviderContext();
	const navigate = useNavigate();

	const redirectToDetails = useCallback(() => {
		navigate(getPageFullPath(`${AppPaths.UPDATES}/${entity.id}`));
	}, [entity.id, navigate]);

	const [deleteUpdate, { isLoading: isDeleteLoading, isError: isErrorDelete, error: deleteError }] =
		useDeleteUpdateMutation();
	const [modifyUpdateState, { isLoading: isModifyLoading, isError: isErrorModify, error: modifyError }] =
		useModifyUpdateStateMutation();

	useEffect(() => {
		if (isErrorDelete) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					default: t('error_default_delete')
				},
				error: deleteError
			});
		}
	}, [isErrorDelete, deleteError, t]);

	useEffect(() => {
		if (isErrorModify) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_modify_state'),
					unAuthorized: t('error_unauthorized_modify_state'),
					forbidden: t('error_forbidden_modify_state'),
					default: t('error_default_modify_state')
				},
				error: modifyError
			});
		}
	}, [isErrorModify, modifyError, t]);

	const buttons = [];
	const ackAction = (
		<Tooltip placement="bottom" title={t('help_approve')}>
			<Button
				key="ack"
				disabled={isModifyLoading || isDeleteLoading}
				loading={isModifyLoading}
				icon={<CheckCircleOutlined style={{ color: 'green' }} />}
				type="text"
				onClick={() => modifyUpdateState({ id: entity.id, body: { state: UpdateState.APPROVED } })}
			/>
		</Tooltip>
	);
	const ignoreAction = (
		<Tooltip placement="bottom" title={t('help_ignore')}>
			<Button
				key="ignore"
				disabled={isModifyLoading || isDeleteLoading}
				loading={isModifyLoading}
				icon={<StopOutlined style={{ color: 'orange' }} />}
				type="text"
				onClick={() => modifyUpdateState({ id: entity.id, body: { state: UpdateState.IGNORED } })}
			/>
		</Tooltip>
	);
	const pendingAction = (
		<Tooltip placement="bottom" title={t('help_pending')}>
			<Button
				key="pending"
				disabled={isModifyLoading || isDeleteLoading}
				loading={isModifyLoading}
				icon={<InteractionOutlined style={{ color: 'blue' }} />}
				type="text"
				onClick={() => modifyUpdateState({ id: entity.id, body: { state: UpdateState.PENDING } })}
			/>
		</Tooltip>
	);

	const color = getUpdateStateColor(entity.state);
	switch (entity.state) {
		case UpdateState.PENDING:
			buttons.push(ackAction);
			buttons.push(ignoreAction);
			break;
		case UpdateState.APPROVED:
			buttons.push(pendingAction);
			buttons.push(ignoreAction);
			break;
		case UpdateState.IGNORED:
			buttons.push(ackAction);
			buttons.push(pendingAction);
			break;
	}

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

	return (
		<>
			<Card
				style={{ borderColor: color }}
				styles={{ header: { backgroundColor: color, borderColor: color } }}
				loading={isDeleteLoading || isModifyLoading}
				key={entity.id}
				title={
					<Space>
						{UpdateState.PENDING === entity.state && (
							<Tooltip title={t('handle_pending')}>
								<Badge dot />
							</Tooltip>
						)}
						<Button type="text" ghost onClick={redirectToDetails}>
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
										DateTimeStyle.MEDIUM,
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
										DateTimeStyle.MEDIUM,
										locale
									),
									updated: formatDateTimeWithTimeZone(
										entity.updatedAt,
										DateTimeStyle.LONG,
										DateTimeStyle.MEDIUM,
										locale
									)
								})}>
								<FieldTimeOutlined />
							</Tooltip>
						)}
					</>
				}
				actions={buttons}>
				<Space direction="vertical">
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('version')}
								</Text>
								<Text>{entity.version}</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('provider')}
								</Text>
								<Text>{entity.provider}</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('host')}
								</Text>
								<Text>{entity.host}</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('state')}
								</Text>
								<UpdateStateTag state={entity.state} />
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('created')}
								</Text>
								<Text>
									{formatDateTimeWithTimeZone(
										entity.createdAt,
										DateTimeStyle.LONG,
										DateTimeStyle.MEDIUM,
										locale
									)}
								</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('updated')}
								</Text>
								<Text>
									{formatDateTimeWithTimeZone(
										entity.updatedAt,
										DateTimeStyle.LONG,
										DateTimeStyle.MEDIUM,
										locale
									)}
								</Text>
							</Space>
						</Col>
					</Row>
				</Space>
			</Card>
		</>
	);
};

export default Update;
