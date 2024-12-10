import { useModifyLabelActionMutation } from '../../api/actionsApi';
import { apiNotification } from '../common/apiNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateLabelActionProps {
	id: string;
	label: string;
}

const UpdateLabelAction: FC<UpdateLabelActionProps> = ({ id, label }): ReactNode => {
	const [t] = useTranslation('action_update_label');

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyLabelActionMutation();

	useEffect(() => {
		if (isError) {
			apiNotification.error({
				i18n: {
					notFound: t('error_unable_update_value'),
					unAuthorized: t('error_unauthorized_update_value'),
					forbidden: t('error_forbidden_update_value'),
					default: t('error_default_update_value')
				},
				error: error
			});
		}
	}, [isError, error, t]);

	const onSubmit = useCallback(
		(value?: string) => {
			if (value && value !== '') {
				modify({ id: id, body: { label: value } });
			}
		},
		[id, modify]
	);

	return (
		<InlineInputValueEditor
			initialValue={label}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			allowBlank={false}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateLabelAction;
