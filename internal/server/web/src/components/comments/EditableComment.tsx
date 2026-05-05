import { useModifyCommentContentMutation } from '../../api/commentsApi';
import { useTheme } from '../../providers/ThemeProvider';
import { CommentResponse, CommentSingleResponse } from '../../types';
import { useNotification } from '../../use/useNotification';
import { EditOutlined } from '@ant-design/icons';
import MDEditor from '@uiw/react-md-editor';
import { Button, Form, Space, Typography } from 'antd';
import React, { FC, ReactNode, useCallback, useEffect, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import ReactMarkdown from 'react-markdown';
import rehypeSanitize from 'rehype-sanitize';
import remarkGfm from 'remark-gfm';

export interface EditableCommentProps {
	comment: CommentResponse;
	onEditSuccess?: (response: CommentSingleResponse) => void;
}

const { Text } = Typography;

const EditableComment: FC<EditableCommentProps> = ({ comment, onEditSuccess }): ReactNode => {
	const [t] = useTranslation('comment_edit');
	const { apiError } = useNotification();
	const { isDarkTheme } = useTheme();
	const [form] = Form.useForm();
	const [isEditing, setIsEditing] = useState<boolean>(false);
	const [modify, { data, isLoading, isError, isSuccess, error }] = useModifyCommentContentMutation();

	const editorTheme = useMemo(() => {
		return isDarkTheme ? 'dark' : 'light';
	}, [isDarkTheme]);

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					badRequest: t('error_bad_request'),
					notFound: t('error_unable'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					default: t('error_default')
				},
				error: error
			});
		}

		if (isSuccess && onEditSuccess) {
			onEditSuccess(data);
		}
	}, [isError, error, t, apiError, isSuccess, onEditSuccess, data]);

	const handleEdit = useCallback(() => {
		form.setFieldsValue({ content: comment.content });
		setIsEditing(true);
	}, [comment.content, form]);

	const handleCancel = useCallback(() => {
		setIsEditing(false);
	}, []);

	const onSubmit = useCallback(() => {
		form.validateFields().then((values) => {
			modify({ id: comment.id, body: values });
			handleCancel();
		});
	}, [comment.id, form, handleCancel, modify]);

	return (
		<>
			{isEditing && (
				<Form form={form} layout="horizontal" name="edit_form">
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
						<Space orientation="horizontal">
							<Button type="primary" onClick={onSubmit} loading={isLoading} disabled={isLoading}>
								{t('submit')}
							</Button>
							<Button onClick={handleCancel} disabled={isLoading}>
								{t('cancel')}
							</Button>
						</Space>
					</Form.Item>
				</Form>
			)}
			{!isEditing && (
				<>
					<Text>
						<ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeSanitize]}>
							{comment.content}
						</ReactMarkdown>
					</Text>
					<Button type="link" onClick={handleEdit} icon={<EditOutlined />} />
				</>
			)}
		</>
	);
};

export default EditableComment;
