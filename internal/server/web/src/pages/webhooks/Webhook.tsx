import {
	useDeleteWebhookMutation,
	useModifyIgnoreHostReplacementWebhookMutation,
	useModifyIgnoreHostWebhookMutation,
	useModifyLabelWebhookMutation
} from '../../api/webhooksApi';
import DateTimeStyle from '../../constants/dateTimeStyle';
import getConfiguration from '../../getConfiguration';
import { useLocaleProviderContext } from '../../providers/LocaleContextProvider';
import { WebhookResponse, WebhookType } from '../../types';
import { useNotification } from '../../use/useNotification';
import { formatDateTimeWithTimeZone } from '../../utils/datetimeHelper';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { CheckOutlined, CloseOutlined, DeleteOutlined, FieldTimeOutlined } from '@ant-design/icons';
import {
	Button,
	Card,
	Collapse,
	CollapseProps,
	Descriptions,
	DescriptionsProps,
	Popconfirm,
	Switch,
	Tag,
	Tooltip,
	Typography
} from 'antd';
import parse from 'html-react-parser';
import linkifyHtml from 'linkify-html';
import { FC, ReactNode, useCallback, useEffect, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface WebhookProps {
	entity: WebhookResponse;
}

const Webhook: FC<WebhookProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('webhook');
	const { apiError } = useNotification();
	const { locale } = useLocaleProviderContext();

	const [deleteWebhook, { isLoading: isDeleteLoading, isError: isErrorDelete, error: deleteError }] =
		useDeleteWebhookMutation();

	const [modifyIgnoreHost, { isLoading: isIgnoreHostLoading, isError: isErrorIgnoreHost, error: ignoreHostError }] =
		useModifyIgnoreHostWebhookMutation();

	const [
		modifyIgnoreHostReplacement,
		{
			isLoading: isIgnoreHostReplacementLoading,
			isError: isErrorIgnoreHostReplacement,
			isSuccess: isSuccessIgnoreHostReplacement,
			error: ignoreHostReplacementError
		}
	] = useModifyIgnoreHostReplacementWebhookMutation();

	const [
		modifyLabel,
		{ isLoading: isLabelLoading, isError: isErrorLabel, isSuccess: isSuccessLabel, error: labelError }
	] = useModifyLabelWebhookMutation();

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
		if (isErrorLabel) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_label'),
					unAuthorized: t('error_unauthorized_update_label'),
					forbidden: t('error_forbidden_update_label'),
					default: t('error_default_update_label')
				},
				error: labelError
			});
		}
	}, [apiError, isErrorLabel, labelError, t]);

	useEffect(() => {
		if (isErrorIgnoreHost) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_ignore_host'),
					unAuthorized: t('error_unauthorized_update_ignore_host'),
					forbidden: t('error_forbidden_update_ignore_host'),
					default: t('error_default_update_ignore_host')
				},
				error: ignoreHostError
			});
		}
	}, [isErrorIgnoreHost, ignoreHostError, t, apiError]);

	useEffect(() => {
		if (isErrorIgnoreHostReplacement) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_ignore_host_replacement'),
					unAuthorized: t('error_unauthorized_update_ignore_host_replacement'),
					forbidden: t('error_forbidden_update_ignore_host_replacement'),
					default: t('error_default_update_ignore_host_replacement')
				},
				error: ignoreHostReplacementError
			});
		}
	}, [apiError, ignoreHostReplacementError, isErrorIgnoreHostReplacement, t]);

	const onIgnoreHostChange = useCallback(
		(checked: boolean) => {
			modifyIgnoreHost({ id: entity.id, body: { ignoreHost: checked } });
		},
		[entity.id, modifyIgnoreHost]
	);

	const submitIgnoreHostReplacementChange = useCallback(
		(value?: string) => {
			if (value && value !== entity.ignoreHostReplacement && value !== '') {
				modifyIgnoreHostReplacement({ id: entity.id, body: { ignoreHostReplacement: value } });
			}
		},
		[entity.id, entity.ignoreHostReplacement, modifyIgnoreHostReplacement]
	);

	const submitLabelChange = useCallback(
		(value?: string) => {
			if (value && value !== entity.label && value !== '') {
				modifyLabel({ id: entity.id, body: { label: value } });
			}
		},
		[entity.id, entity.label, modifyLabel]
	);

	const commandPreview = useMemo(() => {
		const url = `upda webhook send --url ${window.location.protocol}//${window.location.host} --webhook-id ${entity.id} --webhook-token "$TOKEN" --application "Test Application" --application-version "3.0.23"`;
		const preview = parse(linkifyHtml(url));
		const items: CollapseProps['items'] = [
			{
				label: <Text>{t('preview_cli_show')}</Text>,
				children: (
					<Text
						copyable={{
							text: url
						}}
						style={{ fontFamily: 'monospace' }}>
						{preview}
					</Text>
				)
			}
		];
		return <Collapse size="small" items={items} bordered={false} expandIconPlacement="end" />;
	}, [entity.id, t]);

	const urlPreview = useMemo(() => {
		const url = `${getConfiguration().VITE_API_URL}webhooks/${entity.id}`;
		const preview = parse(linkifyHtml(url));
		const items: CollapseProps['items'] = [
			{
				label: <Text>{t('url_show')}</Text>,
				children: (
					<Text copyable={{ text: url }} style={{ fontFamily: 'monospace' }}>
						{preview}
					</Text>
				)
			}
		];
		return <Collapse size="small" items={items} bordered={false} expandIconPlacement="end" />;
	}, [entity.id, t]);

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
				<Button key="del" icon={<DeleteOutlined />} type="text" danger />
			</Tooltip>
		</Popconfirm>
	);
	buttons.push(delAction);

	const webhookDescriptions = useMemo(() => {
		const descriptions: DescriptionsProps['items'] = [
			{
				key: 'label',
				label: (
					<Text type="secondary" ellipsis>
						{t('label')}
					</Text>
				),
				children: (
					<InlineInputValueEditor
						initialValue={entity.label}
						allowBlank={false}
						isLoading={isLabelLoading}
						isSuccess={isSuccessLabel}
						isError={isErrorLabel}
						onSubmit={submitLabelChange}
					/>
				)
			},
			{
				key: 'ignore_host',
				label: (
					<Text type="secondary" ellipsis>
						{t('ignore_host')}
					</Text>
				),
				children: (
					<Switch
						onChange={onIgnoreHostChange}
						loading={isIgnoreHostLoading}
						checkedChildren={<CheckOutlined />}
						unCheckedChildren={<CloseOutlined />}
						checked={entity.ignoreHost}
					/>
				)
			},
			{
				key: 'ignore_host_replacement',
				label: (
					<Text type="secondary" ellipsis>
						{t('ignore_host_replacement')}
					</Text>
				),
				children: (
					<InlineInputValueEditor
						initialValue={entity.ignoreHostReplacement}
						allowBlank={false}
						isLoading={isIgnoreHostReplacementLoading}
						isSuccess={isSuccessIgnoreHostReplacement}
						isError={isErrorIgnoreHostReplacement}
						onSubmit={submitIgnoreHostReplacementChange}
					/>
				)
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
			},
			{
				key: 'url',
				label: (
					<Text type="secondary" ellipsis>
						{t('url')}
					</Text>
				),
				children: urlPreview
			},
			{
				key: 'preview_cli',
				label: (
					<Text type="secondary" ellipsis>
						{t('preview_cli')}
					</Text>
				),
				children: commandPreview
			}
		];

		return (
			<Descriptions
				items={descriptions}
				colon={false}
				layout="vertical"
				size="small"
				column={{ xs: 1, sm: 2, md: 2, lg: 3, xl: 3, xxl: 4 }}
			/>
		);
	}, [
		commandPreview,
		entity.createdAt,
		entity.ignoreHost,
		entity.ignoreHostReplacement,
		entity.label,
		entity.updatedAt,
		isErrorIgnoreHostReplacement,
		isErrorLabel,
		isIgnoreHostLoading,
		isIgnoreHostReplacementLoading,
		isLabelLoading,
		isSuccessIgnoreHostReplacement,
		isSuccessLabel,
		locale,
		onIgnoreHostChange,
		submitIgnoreHostReplacementChange,
		submitLabelChange,
		t,
		urlPreview
	]);

	return (
		<>
			<Card
				size="small"
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
				{webhookDescriptions}
			</Card>
		</>
	);
};

export default Webhook;
