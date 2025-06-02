import { useCreateFilterPresetMutation } from '../../api/filterPresetsApi';
import { CreateFilterPresetRequest, FilterPresetType } from '../../types/filterPreset';
import { useNotification } from '../../use/useNotification';
import { PlusOutlined } from '@ant-design/icons';
import { Button, ColorPicker, Form, Input } from 'antd';
import { Color } from 'antd/lib/color-picker';
import { useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

type CreateFilterPresetProps = {
	type: FilterPresetType;
};

type FilterPresetFormValues = {
	label: string;
	color?: Color;
};

const CreateFilterPreset = ({ type }: CreateFilterPresetProps) => {
	const [t] = useTranslation('filter_presets_create');
	const { apiError } = useNotification();
	const [form] = Form.useForm();
	const [searchParams] = useSearchParams();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateFilterPresetMutation();

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: FilterPresetFormValues) => {
			save({
				type: type,
				label: values.label,
				color: values.color ? values.color.toHexString() : undefined,
				parameters: searchParams.toString()
			} as CreateFilterPresetRequest);
		});
	}, [form, save, type, searchParams]);

	useEffect(() => {
		if (isSuccess) {
			reset();
			form.resetFields();
		}
	}, [isSuccess, data, t, reset, form]);

	useEffect(() => {
		if (isError) {
			apiError({
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
	}, [error, isError, t, apiError]);

	return (
		<Form form={form} layout="inline" onFinish={onSubmit} variant="filled">
			<Form.Item
				label={t('label')}
				name="label"
				rules={[
					{ required: true, message: t('label_required') },
					{ min: 1, max: 255, message: t('label_size') }
				]}>
				<Input minLength={1} maxLength={255} allowClear placeholder={t('label_placeholder')} />
			</Form.Item>
			<Form.Item label={t('color')} name="color">
				<ColorPicker mode="single" format="hex" />
			</Form.Item>
			<Form.Item>
				<Button icon={<PlusOutlined />} type="link" htmlType="submit" loading={isLoading} disabled={isLoading}>
					{t('create')}
				</Button>
			</Form.Item>
		</Form>
	);
};

export default CreateFilterPreset;
