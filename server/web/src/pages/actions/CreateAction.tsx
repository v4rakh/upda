import ActionFormPayloadSwitch from './ActionFormPayloadSwitch';
import ActionFormType from './ActionFormType';
import ActionSelectEvent from './ActionSelectEvent';
import { useCreateActionMutation } from '../../api/actionsApi';
import { ActionPayloadShoutrrr, ActionType } from '../../types';
import { apiNotification } from '../common/apiNotification';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Collapse, Divider, Form, Input, Switch } from 'antd';
import { ReactNode, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const COLLAPSE_KEY = 'edit_type_payload';

const CreateAction = (): ReactNode => {
	const [t] = useTranslation('action_create');
	const [form] = Form.useForm();
	const type = Form.useWatch('type', form);
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const [save, { isSuccess, isError, reset, error, isLoading }] = useCreateActionMutation();

	useEffect(() => {
		if (isSuccess) {
			reset();
			form.resetFields();
			setCollapseActiveKeys([]);
		}
	}, [form, isSuccess, reset]);

	useEffect(() => {
		if (isError) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					badRequest: t('error_bad_request'),
					default: t('error_default')
				},
				error: error
			});
		}
	}, [isError, error, t]);

	const onSubmit = useCallback(() => {
		form.validateFields().then(
			(
				values: ActionPayloadShoutrrr & {
					label: string;
					type: ActionType;
					enabled: boolean;
					matchEvent?: string;
					matchApplication?: string;
					matchProvider?: string;
					matchHost?: string;
				}
			) => {
				save({
					label: values.label,
					type: values.type,
					enabled: values.enabled,
					matchEvent: values.matchEvent,
					matchApplication: values.matchApplication,
					matchProvider: values.matchProvider,
					matchHost: values.matchHost,
					payload: { urls: values.urls, body: values.body }
				});
			}
		);
	}, [form, save]);

	return (
		<Collapse
			onChange={(keys) => {
				setCollapseActiveKeys(keys);
			}}
			expandIcon={() => <></>}
			expandIconPosition="end"
			bordered={false}
			ghost
			size="small"
			activeKey={collapseActiveKeys}
			items={[
				{
					key: COLLAPSE_KEY,
					label: (
						<Button type="link" icon={<PlusOutlined />}>
							{t('create')}
						</Button>
					),
					children: (
						<>
							<Form
								form={form}
								layout="vertical"
								initialValues={{ type: ActionType.SHOUTRRR, enabled: true }}>
								<Form.Item
									label={t('label')}
									tooltip={t('label_help')}
									name="label"
									rules={[
										{ required: true, message: t('label_required') },
										{ min: 1, max: 255, message: t('label_size') }
									]}>
									<Input variant="filled" allowClear placeholder={t('label_placeholder')} />
								</Form.Item>
								<Form.Item label={t('match_event')} tooltip={t('match_event_help')} name="matchEvent">
									<ActionSelectEvent loading={isLoading} />
								</Form.Item>
								<Form.Item
									label={t('match_application')}
									tooltip={t('match_application_help')}
									name="matchApplication">
									<Input placeholder={t('all')} variant="filled" />
								</Form.Item>
								<Form.Item
									label={t('match_provider')}
									tooltip={t('match_provider_help')}
									name="matchProvider">
									<Input placeholder={t('all')} variant="filled" />
								</Form.Item>
								<Form.Item label={t('match_host')} tooltip={t('match_host_help')} name="matchHost">
									<Input placeholder={t('all')} variant="filled" />
								</Form.Item>
								<ActionFormType isLoading={isLoading} initialValue={ActionType.SHOUTRRR} />
								<ActionFormPayloadSwitch isLoading={isLoading} type={type} />
								<Form.Item
									label={t('enabled_label')}
									tooltip={t('enabled_help')}
									required={true}
									name="enabled"
									valuePropName="checked">
									<Switch
										loading={isLoading}
										checkedChildren={t('yes')}
										unCheckedChildren={t('no')}
									/>
								</Form.Item>
								<Form.Item>
									<Button
										type="primary"
										onClick={() => onSubmit()}
										loading={isLoading}
										disabled={isLoading}>
										{t('submit')}
									</Button>
								</Form.Item>
							</Form>
							<Divider />
						</>
					)
				}
			]}
		/>
	);
};

export default CreateAction;
