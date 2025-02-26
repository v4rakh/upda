import { useModifyEnabledActionMutation } from '../../api/actionsApi';
import { useNotification } from '../../use/useNotification';
import { CheckOutlined, CloseOutlined } from '@ant-design/icons';
import { Switch } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateEnabledActionProps {
	id: string;
	enabled: boolean;
}

const UpdateEnabledAction: FC<UpdateEnabledActionProps> = ({ id, enabled }): ReactNode => {
	const [t] = useTranslation('action_update_enabled');
	const { apiError } = useNotification();

	const [modify, { isLoading, isError, error }] = useModifyEnabledActionMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_value'),
					unAuthorized: t('error_unauthorized_update_value'),
					forbidden: t('error_forbidden_update_value'),
					default: t('error_default_update_value')
				},
				error: isError
			});
		}
	}, [isError, error, t, apiError]);

	const onEnabledChange = useCallback(
		(checked: boolean) => {
			modify({ id: id, body: { enabled: checked } });
		},
		[id, modify]
	);

	return (
		<Switch
			onChange={onEnabledChange}
			loading={isLoading}
			checkedChildren={<CheckOutlined />}
			unCheckedChildren={<CloseOutlined />}
			checked={enabled}
		/>
	);
};

export default UpdateEnabledAction;
