import { useCreateWebhookMutation } from '../../api/webhooksApi';
import { CreateWebhookRequest, WebhookType } from '../../types';
import { apiNotification } from '../common/apiNotification';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Collapse, Divider, Form, Input, Select, Switch } from 'antd';
import parse from 'html-react-parser';
import { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const COLLAPSE_KEY = 'create';

const CreateWebhook: FC = () => {
	const [t] = useTranslation('webhook_create');
	const [form] = Form.useForm();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateWebhookMutation();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: CreateWebhookRequest) => {
			save(values);
		});
	}, [form, save]);

	useEffect(() => {
		if (isSuccess) {
			apiNotification.simpleInfo({
				title: t('created_title'),
				message: parse(t('created_message', { token: data?.data.token })),
				duration: 60
			});
			reset();
			form.resetFields();
			setCollapseActiveKeys([]);
		}
	}, [isSuccess, data, t, reset, form]);

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
				error
			});
		}
	}, [error, isError, t]);

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
								layout="inline"
								form={form}
								initialValues={{ type: WebhookType.GENERIC, ignore_host: false }}>
								<Form.Item
									label={t('label')}
									name="label"
									rules={[
										{ required: true, message: t('label_required') },
										{ min: 1, max: 255, message: t('label_size') }
									]}>
									<Input variant="filled" allowClear placeholder={t('label_placeholder')} />
								</Form.Item>
								<Form.Item label={t('type')} name="type" required={true}>
									<Select
										variant="filled"
										style={{ width: 100 }}
										options={[
											{
												value: WebhookType.GENERIC,
												label: t(`type_${WebhookType.GENERIC.toLowerCase()}`)
											},
											{
												value: WebhookType.DIUN,
												label: t(`type_${WebhookType.DIUN.toLowerCase()}`)
											}
										]}
									/>
								</Form.Item>

								<Form.Item label={t('ignore_host')} name="ignoreHost" valuePropName="checked">
									<Switch loading={isLoading} />
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

export default CreateWebhook;
