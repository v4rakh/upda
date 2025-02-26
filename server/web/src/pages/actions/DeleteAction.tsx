import { useDeleteActionMutation } from '../../api/actionsApi';
import { useNotification } from '../../use/useNotification';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteActionProps {
	id: string;
}

const DeleteAction: FC<DeleteActionProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('action_delete');
	const { apiError } = useNotification();

	const [deleteAction, { isLoading, isError, error }] = useDeleteActionMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					default: t('error_default_delete')
				},
				error: error
			});
		}
	}, [isError, error, t, apiError]);

	return (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => deleteAction({ id: id })}
			okText={t('delete')}
			placement="right"
			cancelText={t('cancel')}
			okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
			<Tooltip title={t('help_delete')} placement="right">
				<Button loading={isLoading} key="del" icon={<DeleteOutlined />} type="text" danger>
					{t('delete')}
				</Button>
			</Tooltip>
		</Popconfirm>
	);
};

export default DeleteAction;
