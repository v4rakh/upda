import { useDeleteUpdateStateTransitionMutation } from '../../api/updateStateTransitionsApi';
import { useNotification } from '../../use/useNotification';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteStateTransitionProps {
	id: string;
}

const DeleteStateTransition: FC<DeleteStateTransitionProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('state_transition_delete');
	const { apiError } = useNotification();

	const [deleteStateTransition, { isLoading, isError, error }] = useDeleteUpdateStateTransitionMutation();

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
			onConfirm={() => deleteStateTransition({ id: id })}
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

export default DeleteStateTransition;
