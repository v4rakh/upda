import { useDeleteActionInvocationMutation } from '../../api/actionInvocationsApi';
import { apiNotification } from '../common/apiNotification';
import { DeleteOutlined, DeleteTwoTone } from '@ant-design/icons';
import { Button, Popconfirm, Tooltip } from 'antd';
import { FC, ReactNode, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface DeleteActionInvocationProps {
	id: string;
}

const DeleteActionInvocation: FC<DeleteActionInvocationProps> = ({ id }): ReactNode => {
	const [t] = useTranslation('action_invocation_delete');

	const [deleteActionInvocation, { isLoading, isError, error }] = useDeleteActionInvocationMutation();

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
			onConfirm={() => deleteActionInvocation({ id: id })}
			okText={t('delete')}
			placement="right"
			cancelText={t('cancel')}
			okButtonProps={{ icon: <DeleteOutlined />, type: 'primary', danger: true }}>
			<Tooltip title={t('help_delete')} placement="right">
				<Button loading={isLoading} key="del" icon={<DeleteTwoTone twoToneColor={'red'} />} type="text" danger>
					{t('delete')}
				</Button>
			</Tooltip>
		</Popconfirm>
	);
};

export default DeleteActionInvocation;
