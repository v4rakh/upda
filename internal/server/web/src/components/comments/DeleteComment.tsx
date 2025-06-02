import { useDeleteCommentMutation } from '../../api/commentsApi';
import { useNotification } from '../../use/useNotification';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteCommentProps {
	id: string;
	onDeleteSuccess?: () => void;
}

const DeleteComment: FC<DeleteCommentProps> = ({ id, onDeleteSuccess }): ReactNode => {
	const [t] = useTranslation('comment_delete');
	const { apiError } = useNotification();

	const [DeleteComment, { isSuccess, isLoading, isError, error }] = useDeleteCommentMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable'),
					unAuthorized: t('error_unauthorized'),
					forbidden: t('error_forbidden'),
					default: t('error_default')
				},
				error: error
			});
		}

		if (isSuccess && onDeleteSuccess) {
			onDeleteSuccess();
		}
	}, [isError, error, t, apiError, isSuccess, onDeleteSuccess]);

	return (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => DeleteComment({ id: id })}
			okText={t('delete')}
			placement="bottom"
			cancelText={t('cancel')}
			okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
			<Tooltip title={t('help_delete')} placement="bottom">
				<Button loading={isLoading} key="del" icon={<DeleteOutlined />} type="text" danger />
			</Tooltip>
		</Popconfirm>
	);
};

export default DeleteComment;
