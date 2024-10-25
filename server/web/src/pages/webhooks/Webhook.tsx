import {
	useDeleteWebhookMutation,
	useModifyIgnoreHostWebhookMutation,
	useModifyLabelWebhookMutation
} from '../../api/webhooksApi';
import getConfiguration from '../../getConfiguration';
import { WebhookResponse, WebhookType } from '../../types';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import { apiNotification } from '../common/apiNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { CheckOutlined, CloseOutlined, DeleteOutlined, DeleteTwoTone, FieldTimeOutlined } from '@ant-design/icons';
import { Button, Card, Col, Popconfirm, Row, Space, Switch, Tag, Tooltip, Typography } from 'antd';
import { FC, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface WebhookProps {
	entity: WebhookResponse;
}

const Webhook: FC<WebhookProps> = ({ entity }): JSX.Element => {
	const [t] = useTranslation('webhook');

	const [deleteWebhook, { isLoading: isDeleteLoading, isError: isErrorDelete, error: deleteError }] =
		useDeleteWebhookMutation();

	const [modifyIgnoreHost, { isLoading: isIgnoreHostLoading, isError: isErrorIgnoreHost, error: ignoreHostError }] =
		useModifyIgnoreHostWebhookMutation();

	const [
		modifyLabel,
		{ isLoading: isLabelLoading, isError: isErrorLabel, isSuccess: isSuccessLabel, error: labelError }
	] = useModifyLabelWebhookMutation();

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
		if (isErrorLabel) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_update_label'),
					unAuthorized: t('error_unauthorized_update_label'),
					forbidden: t('error_forbidden_update_label'),
					default: t('error_default_update_label')
				},
				error: labelError
			});
		}
	}, [isErrorLabel, labelError, t]);

	useEffect(() => {
		if (isErrorIgnoreHost) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_update_ignore_host'),
					unAuthorized: t('error_unauthorized_update_ignore_host'),
					forbidden: t('error_forbidden_update_ignore_host'),
					default: t('error_default_update_ignore_host')
				},
				error: ignoreHostError
			});
		}
	}, [isErrorIgnoreHost, ignoreHostError, t]);

	const onIgnoreHostChange = useCallback(
		(checked: boolean) => {
			modifyIgnoreHost({ id: entity.id, body: { ignoreHost: checked } });
		},
		[entity.id, modifyIgnoreHost]
	);

	const submitLabelChange = useCallback(
		(value?: string) => {
			if (value && value !== entity.label && value !== '') {
				modifyLabel({ id: entity.id, body: { label: value } });
			}
		},
		[entity.id, entity.label, modifyLabel]
	);

	const buttons = [];
	const delAction = (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => deleteWebhook({ id: entity.id })}
			okText={t('delete')}
			placement="bottom"
			cancelText={t('cancel')}
			okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
			<Tooltip title={t('help_delete')} placement="bottom">
				<Button key="del" icon={<DeleteTwoTone twoToneColor={'red'} />} type={'text'} danger />
			</Tooltip>
		</Popconfirm>
	);
	buttons.push(delAction);

	return (
		<>
			<Card
				loading={isDeleteLoading}
				key={entity.id}
				title={
					<Tag color={WebhookType.DIUN === entity.type ? 'blue-inverse' : 'green-inverse'}>
						{t(`type_${entity.type.toLocaleLowerCase()}`)}
					</Tag>
				}
				extra={
					<>
						{entity.createdAt == entity.updatedAt && (
							<Tooltip
								placement={'bottom'}
								title={t('created_at', {
									created: formatDateTimeWithTimeZone(entity.createdAt)
								})}>
								<FieldTimeOutlined />
							</Tooltip>
						)}
						{entity.createdAt !== entity.updatedAt && (
							<Tooltip
								placement={'bottom'}
								title={t('created_at_diff', {
									created: formatDateTimeWithTimeZone(entity.createdAt),
									updated: formatDateTimeWithTimeZone(entity.updatedAt)
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
									{t('label')}
								</Text>
								<InlineInputValueEditor
									initialValue={entity.label}
									allowBlank={false}
									isLoading={isLabelLoading}
									isSuccess={isSuccessLabel}
									isError={isErrorLabel}
									onSubmit={submitLabelChange}
								/>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('url')}
								</Text>
								<Text code copyable>
									{`${getConfiguration().VITE_API_URL}webhooks/${entity.id}`}
								</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('ignore_host')}
								</Text>
								<Switch
									onChange={onIgnoreHostChange}
									loading={isIgnoreHostLoading}
									checkedChildren={<CheckOutlined />}
									unCheckedChildren={<CloseOutlined />}
									checked={entity.ignoreHost}
								/>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('created')}
								</Text>
								<Text>{formatDateTimeWithTimeZone(entity.createdAt)}</Text>
							</Space>
						</Col>
					</Row>
					<Row>
						<Col>
							<Space>
								<Text type="secondary" ellipsis>
									{t('updated')}
								</Text>
								<Text>{formatDateTimeWithTimeZone(entity.updatedAt)}</Text>
							</Space>
						</Col>
					</Row>
				</Space>
			</Card>
		</>
	);
};

export default Webhook;
