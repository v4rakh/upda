import ActionFormPayloadSwitch from './ActionFormPayloadSwitch';
import ActionFormType from './ActionFormType';
import { useModifyTypeAndPayloadActionMutation } from '../../api/actionsApi';
import { ActionPayloadShoutrrr, ActionType } from '../../types';
import { EventName } from '../../types/event';
import { useNotification } from '../../use/useNotification';
import { SettingOutlined } from '@ant-design/icons';
import { Button, Collapse, Form, Space, Typography } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdatePayloadActionProps {
	id: string;
	type: ActionType;
	payload: ActionPayloadShoutrrr | undefined;
	matchEvent?: EventName;
}

const { Text } = Typography;

const COLLAPSE_KEY = 'update_type_and_payload';

const UpdateTypeAndPayloadAction: FC<UpdatePayloadActionProps> = ({ id, type, payload, matchEvent }): ReactNode => {
	const [t] = useTranslation('action_update_payload');
	const { apiError } = useNotification();
	const [form] = Form.useForm();
	const typeValue = Form.useWatch('type', form);

	const [modify, { isLoading, isError, error }] = useModifyTypeAndPayloadActionMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_value'),
					unAuthorized: t('error_unauthorized_update_value'),
					forbidden: t('error_forbidden_update_value'),
					default: t('error_default_update_value')
				},
				error: error
			});
		}
	}, [isError, error, t, apiError]);

	const onSubmit = useCallback(
		(values: ActionPayloadShoutrrr & { type: ActionType }) => {
			modify({ id: id, body: { type: values.type, payload: { urls: values.urls, body: values.body } } });
		},
		[id, modify]
	);

	return (
		<Collapse
			bordered={false}
			ghost={false}
			size="small"
			items={[
				{
					key: COLLAPSE_KEY,
					label: (
						<Space>
							<SettingOutlined />
							<Text>{t('update_type_and_payload')}</Text>
						</Space>
					),
					children: (
						<Form
							form={form}
							style={{ maxWidth: '100%' }}
							onFinish={onSubmit}
							initialValues={{ type, ...payload }}
							layout="vertical">
							<ActionFormType isLoading={isLoading} initialValue={typeValue} />
							<ActionFormPayloadSwitch isLoading={isLoading} type={typeValue} form={form} matchEvent={matchEvent} />
							<Form.Item>
								<Button type="primary" htmlType="submit" disabled={isLoading} loading={isLoading}>
									{t('submit')}
								</Button>
							</Form.Item>
						</Form>
					)
				}
			]}
		/>
	);
};

export default UpdateTypeAndPayloadAction;
