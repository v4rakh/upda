import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';
import { Button, Form, Input } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

export interface ActionFormShoutrrrProps {
	isLoading: boolean;
}

const ActionFormPayloadShoutrrr: FC<ActionFormShoutrrrProps> = ({ isLoading }): ReactNode => {
	const [t] = useTranslation('action_form_shoutrrrr');

	return (
		<>
			<Form.List
				name="urls"
				rules={[
					{
						validator: async (_, urls) => {
							if (!urls || urls.length < 1) {
								return Promise.reject(new Error(t('urls_validate_minimum')));
							}
						}
					}
				]}>
				{(fields, { add, remove }, { errors }) => (
					<>
						{fields.map((field, index) => (
							<Form.Item
								required={true}
								key={field.key}
								label={`${t('url')} ${index + 1}`}
								tooltip={t('url_help')}>
								<Form.Item
									{...field}
									key={field.key}
									validateTrigger={['onChange', 'onBlur']}
									rules={[
										{
											required: true,
											whitespace: true,
											message: t('urls_validate_not_blank')
										}
									]}
									noStyle>
									<Input
										disabled={isLoading}
										placeholder={t('urls_placeholder')}
										variant="filled"
										style={{ width: '90%' }}
										allowClear
									/>
								</Form.Item>
								{!isLoading && fields.length > 1 ? (
									<MinusCircleOutlined
										style={{ marginLeft: '5%' }}
										onClick={() => remove(field.name)}
									/>
								) : null}
							</Form.Item>
						))}
						<Form.Item>
							<Button disabled={isLoading} type="link" onClick={() => add()} icon={<PlusOutlined />}>
								{t('urls_new')}
							</Button>
							<Form.ErrorList errors={errors} />
						</Form.Item>
					</>
				)}
			</Form.List>
			<Form.Item
				name="body"
				label={t('body_label')}
				required={true}
				tooltip={t('body_help')}
				rules={[
					{ required: true, message: t('body_required') },
					{ min: 1, message: t('body_size') }
				]}>
				<Input
					placeholder={t('body_placeholder')}
					disabled={isLoading}
					variant="filled"
					allowClear
					style={{ width: '95%' }}
				/>
			</Form.Item>
		</>
	);
};

export default ActionFormPayloadShoutrrr;
