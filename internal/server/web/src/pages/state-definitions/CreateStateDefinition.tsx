import IconSelector from './IconSelector';
import { useCreateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { CreateUpdateStateDefinitionRequest } from '../../types';
import { useNotification } from '../../use/useNotification';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Collapse, ColorPicker, Divider, Form, Input, Switch } from 'antd';
import { Color } from 'antd/es/color-picker';
import { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const COLLAPSE_KEY = 'create';

interface FormValues {
	name: string;
	label: string;
	color: Color | string;
	icon: string;
	description?: string;
	isInitial: boolean;
	skipOnNewVersion: boolean;
}

const CreateStateDefinition: FC = () => {
	const [t] = useTranslation('state_definition_create');
	const { apiError, simpleInfo } = useNotification();
	const [form] = Form.useForm<FormValues>();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateUpdateStateDefinitionMutation();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: FormValues) => {
			const colorValue = typeof values.color === 'string' ? values.color : values.color.toHexString();
			const body: CreateUpdateStateDefinitionRequest = {
				name: values.name,
				label: values.label,
				color: colorValue,
				icon: values.icon,
				description: values.description,
				isInitial: values.isInitial ?? false,
				skipOnNewVersion: values.skipOnNewVersion ?? false
			};
			save({ body });
		});
	}, [form, save]);

	useEffect(() => {
		if (isSuccess) {
			simpleInfo({
				title: t('created_title'),
				message: t('created_message', { name: data?.data.name })
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
	}, [error, isError, t, apiError]);

	// @ts-ignore
	// @ts-ignore
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
							<Form
								layout="inline"
								form={form}
								initialValues={{ isInitial: false, skipOnNewVersion: false, icon: 'TagOutlined' }}>
								<Form.Item
									label={t('name')}
									tooltip={t('name_help')}
									name="name"
									rules={[
										{ required: true, message: t('name_required') },
										{ min: 1, max: 255, message: t('name_size') }
									]}>
									<Input variant="filled" allowClear placeholder={t('name_placeholder')} />
								</Form.Item>
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
								<Form.Item
									label={t('color')}
									tooltip={t('color_help')}
									name="color"
									rules={[{ required: true, message: t('color_required') }]}>
									<ColorPicker
										format="hex"
										styles={{
											popupOverlayInner: { padding: 12 }
										}}
									/>
								</Form.Item>
								<Form.Item
									label={t('icon')}
									tooltip={t('icon_help')}
									name="icon"
									rules={[{ required: true, message: t('icon_required') }]}>
									<IconSelector />
								</Form.Item>
								<Form.Item label={t('description')} tooltip={t('description_help')} name="description">
									<Input variant="filled" allowClear placeholder={t('description_placeholder')} />
								</Form.Item>
								<Form.Item
									label={t('is_initial')}
									tooltip={t('is_initial_help')}
									name="isInitial"
									valuePropName="checked">
									<Switch />
								</Form.Item>
								<Form.Item
									label={t('skip_on_new_version')}
									tooltip={t('skip_on_new_version_help')}
									name="skipOnNewVersion"
									valuePropName="checked">
									<Switch />
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

export default CreateStateDefinition;
