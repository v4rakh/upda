import { useCreateCommentMutation } from '../../api/commentsApi';
import { CommentSingleResponse, CreateCommentRequest } from '../../types';
import { useNotification } from '../../use/useNotification';
import { darkThemeEnabled } from '../../utils/featureHelper';
import { PlusOutlined } from '@ant-design/icons';
import MDEditor from '@uiw/react-md-editor';
import { Button, Collapse, Divider, Form } from 'antd';
import { ReactNode, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import rehypeSanitize from 'rehype-sanitize';

const COLLAPSE_KEY = 'create';

type CreateCommentProps = {
	updateId: string;
	onCreateSuccess?: (response: CommentSingleResponse) => void;
};

const CreateComment = ({ updateId, onCreateSuccess }: CreateCommentProps): ReactNode => {
	const [t] = useTranslation('comment_create');
	const { apiError } = useNotification();
	const [form] = Form.useForm();
	const [save, { data, isSuccess, isError, reset, error, isLoading }] = useCreateCommentMutation();
	const [collapseActiveKeys, setCollapseActiveKeys] = useState<string[] | string>([]);

	const editorTheme = useMemo(() => {
		return darkThemeEnabled() ? 'dark' : 'light';
	}, []);

	const onSubmit = useCallback(() => {
		form.validateFields().then((values: CreateCommentRequest) => {
			save({ updateId: updateId, body: values });
		});
	}, [form, save, updateId]);

	useEffect(() => {
		if (isSuccess) {
			reset();
			form.resetFields();
			setCollapseActiveKeys([]);

			if (onCreateSuccess) {
				onCreateSuccess(data);
			}
		}
	}, [data, form, isSuccess, onCreateSuccess, reset]);

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
	}, [error, isSuccess, isError, t, apiError]);

	return (
		<Collapse
			onChange={(keys) => {
				setCollapseActiveKeys(keys);
			}}
			expandIconPosition="end"
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
							<Form layout="horizontal" form={form}>
								<Form.Item
									name="content"
									label={t('content')}
									rules={[{ required: true, message: t('content_required') }]}>
									<MDEditor
										preview="live"
										overflow={false}
										data-color-mode={editorTheme}
										previewOptions={{ rehypePlugins: [rehypeSanitize] }}
									/>
								</Form.Item>
								<Form.Item>
									<Button type="primary" onClick={onSubmit} loading={isLoading} disabled={isLoading}>
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

export default CreateComment;
