import useActionAutoSuggestion from '../../use/useActionAutoSuggestion';
import { MinusCircleOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons';
import { Alert, Button, Flex, Form, FormInstance, Mentions } from 'antd';
import parse from 'html-react-parser';
import { FC, KeyboardEvent, ReactNode, useCallback } from 'react';
import { useTranslation } from 'react-i18next';

export interface ActionFormShoutrrrProps {
	isLoading: boolean;
	form: FormInstance;
}

const KEY_ENTER = 'Enter';
const PREFIX_TRIGGER = ['<'];

const ActionFormPayloadShoutrrr: FC<ActionFormShoutrrrProps> = ({ isLoading, form }): ReactNode => {
	const [t] = useTranslation('action_form_shoutrrrr');
	const { mentionOptions, reloadMentionOptions } = useActionAutoSuggestion();

	// Prevents line breaks
	const handleKeyDown = useCallback((e: KeyboardEvent<HTMLTextAreaElement>) => {
		if (KEY_ENTER === e.key) {
			e.preventDefault();
		}
	}, []);

	return (
		<>
			<Alert showIcon type="info" style={{ marginBottom: '2vh' }} title={parse(t('global_help_banner'))} />
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
									<Flex>
										<Mentions
											defaultValue={form.getFieldValue('urls')[index]}
											onChange={(val) => {
												const currentUrls = form.getFieldValue('urls');
												currentUrls[index] = val;
												form.setFieldValue('urls', currentUrls);
											}}
											split={''}
											prefix={PREFIX_TRIGGER}
											autoSize={{ minRows: 1, maxRows: 5 }}
											allowClear
											options={mentionOptions}
											onKeyDown={handleKeyDown}
											placeholder={t('urls_placeholder')}
											disabled={isLoading}
											variant="filled"
										/>
										<Button type="link" icon={<ReloadOutlined />} onClick={reloadMentionOptions} />
										{!isLoading && fields.length > 1 ? (
											<MinusCircleOutlined
												style={{ marginLeft: '5%' }}
												onClick={() => remove(field.name)}
											/>
										) : null}
									</Flex>
								</Form.Item>
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
				<Flex>
					<Mentions
						defaultValue={form.getFieldValue('body')}
						onChange={(val) => form.setFieldValue('body', val)}
						split={''}
						prefix={PREFIX_TRIGGER}
						autoSize={{ minRows: 1, maxRows: 5 }}
						allowClear
						options={mentionOptions}
						// onKeyDown={handleKeyDown}
						placeholder={t('body_placeholder')}
						disabled={isLoading}
						variant="filled"
					/>
					<Button type="link" icon={<ReloadOutlined />} onClick={reloadMentionOptions} />
				</Flex>
			</Form.Item>
		</>
	);
};

export default ActionFormPayloadShoutrrr;
