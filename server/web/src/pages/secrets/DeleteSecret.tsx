import { useDeleteSecretMutation } from '../../api/secretsApi';
import { apiNotification } from '../common/apiNotification';
import { DeleteOutlined } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteSecretProps {
	id: string;
}

const DeleteSecret: FC<DeleteSecretProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('secret_delete');

	const [deleteSecret, { isLoading, isError, error }] = useDeleteSecretMutation();

	useEffect(() => {
		if (isError) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_delete'),
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					default: t('error_default_delete')
				},
				error: error
			});
		}
	}, [isError, error, t]);

	return (
		<Popconfirm
			title={t('delete_title')}
			onConfirm={() => deleteSecret({ id: id })}
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

export default DeleteSecret;
