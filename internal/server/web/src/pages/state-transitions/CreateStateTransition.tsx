import { useGetUpdateStateDefinitionsQuery } from '../../api/updateStateDefinitionsApi';
import { useCreateUpdateStateTransitionMutation } from '../../api/updateStateTransitionsApi';
import { CreateUpdateStateTransitionRequest } from '../../types';
import { useNotification } from '../../use/useNotification';
import { renderIcon } from '../../utils/iconHelper';
import { PlusOutlined } from '@ant-design/icons';
import { Button, Collapse, Divider, Form, Select, Tag } from 'antd';
import { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const COLLAPSE_KEY = 'create';

interface FormValues {
	fromStateId: string;
	toStateId: string;
}

const CreateStateTransition: FC = () => {
	const [t] = useTranslation('state_transition_create');
	const { apiError, simpleInfo } = useNotification();
	const [form] = Form.useForm<FormValues>();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateUpdateStateTransitionMutation();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const { data: stateDefinitionsData, isLoading: isLoadingStates } = useGetUpdateStateDefinitionsQuery();
	const stateDefinitions = stateDefinitionsData?.data?.content ?? [];

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: FormValues) => {
			const body: CreateUpdateStateTransitionRequest = {
				fromStateId: values.fromStateId,
				toStateId: values.toStateId
			};
			save({ body });
		});
	}, [form, save]);

	useEffect(() => {
		if (isSuccess) {
			simpleInfo({
				title: t('created_title'),
				message: t('created_message', {
					from: data?.data.fromState.label,
					to: data?.data.toState.label
				})
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

	const stateOptions = stateDefinitions.map((state) => ({
		value: state.id,
		label: (
			<Tag color={state.color} icon={renderIcon(state.icon, { marginRight: 4 })}>
				{state.label}
			</Tag>
		)
	}));

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
									label={t('from_state')}
									tooltip={t('from_state_help')}
									name="fromStateId"
									rules={[{ required: true, message: t('from_state_required') }]}>
									<Select
										style={{ minWidth: 150 }}
										loading={isLoadingStates}
										options={stateOptions}
										placeholder={t('from_state_placeholder')}
									/>
								</Form.Item>
								<Form.Item
									label={t('to_state')}
									tooltip={t('to_state_help')}
									name="toStateId"
									rules={[{ required: true, message: t('to_state_required') }]}>
									<Select
										style={{ minWidth: 150 }}
										loading={isLoadingStates}
										options={stateOptions}
										placeholder={t('to_state_placeholder')}
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

export default CreateStateTransition;
