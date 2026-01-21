import { useCreateSecretMutation } from '../../api/secretsApi';
import { CreateSecretRequest } from '../../types';
import { useNotification } from '../../use/useNotification';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Collapse, Divider, Form, Input } from 'antd';
import parse from 'html-react-parser';
import { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const COLLAPSE_KEY = 'create';

const CreateSecret: FC = () => {
	const [t] = useTranslation('secret_create');
	const { apiError, simpleInfo } = useNotification();
	const [form] = Form.useForm();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateSecretMutation();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: CreateSecretRequest) => {
			save(values);
		});
	}, [form, save]);

	useEffect(() => {
		if (isSuccess) {
			simpleInfo({
				title: t('created_title'),
				message: parse(t('created_message', { value: data?.data.value })),
				duration: 20
			});
			reset();
			form.resetFields();
			setCollapseActiveKeys([]);
		}
	}, [isSuccess, data, t, reset, form, simpleInfo]);

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					conflict: t('error_conflict'),
					notFound: t('error_unable'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					badRequest: t('error_bad_request'),
					default: t('error_default')
				},
				error
			});
		}
	}, [error, isSuccess, isError, t, apiError]);

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
					label: (
						<Button type="link" icon={<PlusOutlined />}>
							{t('create')}
						</Button>
					),
					children: (
						<>
							<Form layout="inline" form={form}>
								<Form.Item
									label={t('key')}
									tooltip={t('key_help')}
									name="key"
									rules={[
										{ required: true, message: t('key_required') },
										{ min: 1, max: 255, message: t('key_size') }
									]}>
									<Input variant="filled" allowClear placeholder={t('key_placeholder')} />
								</Form.Item>
								<Form.Item
									label={t('value')}
									tooltip={t('value_help')}
									name="value"
									rules={[
										{ required: true, message: t('value_required') },
										{ min: 1, message: t('value_size') }
									]}>
									<Input.Password variant="filled" allowClear placeholder={t('value_placeholder')} />
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

export default CreateSecret;
