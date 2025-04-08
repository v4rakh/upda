import { useModifyValueConstantMutation } from '../../api/constantsApi';
import { useNotification } from '../../use/useNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateValueConstantProps {
	id: string;
	entityValue?: string;
}

const UpdateValueConstant: FC<UpdateValueConstantProps> = ({ id, entityValue }): ReactNode => {
	const [t] = useTranslation('constant_update_value');
	const { apiError } = useNotification();

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyValueConstantMutation();

	useEffect(() => {
		if (isError) {
			apiError({
				i18n: {
					notFound: t('error_unable_update_value'),
					unAuthorized: t('error_unauthorized_update_value'),
					forbidden: t('error_forbidden_update_value'),
					default: t('error_default_update_value')
				},
				error: error
			});
		}
	}, [isError, error, t, apiError]);

	const submitValueChange = useCallback(
		(value?: string) => {
			if (value && value !== entityValue && value !== '') {
				modify({ id: id, body: { value: value } });
			}
		},
		[entityValue, id, modify]
	);

	return (
		<InlineInputValueEditor
			initialValue={entityValue}
			allowBlank={false}
			resetOnSuccess={false}
			resetOnError={false}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			onSubmit={submitValueChange}
		/>
	);
};

export default UpdateValueConstant;
