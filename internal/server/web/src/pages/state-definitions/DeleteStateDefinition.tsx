import { useDeleteUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { useNotification } from '../../use/useNotification';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteStateDefinitionProps {
	id: string;
}

const DeleteStateDefinition: FC<DeleteStateDefinitionProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('state_definition_delete');
	const { apiError } = useNotification();

	const [deleteStateDefinition, { isLoading, isError, error }] = useDeleteUpdateStateDefinitionMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					badRequest: t('error_bad_request_delete'),
					default: t('error_default_delete')
				},
				error: error
			});
		}
	}, [isError, error, t, apiError]);

	return (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => deleteStateDefinition({ id: id })}
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

export default DeleteStateDefinition;
