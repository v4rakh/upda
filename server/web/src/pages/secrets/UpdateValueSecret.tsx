import { useModifyValueSecretMutation } from '../../api/secretsApi';
import { apiNotification } from '../common/apiNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateValueSecretProps {
	id: string;
	entityValue?: string;
}

const UpdateValueSecret: FC<UpdateValueSecretProps> = ({ id, entityValue }): ReactNode => {
	const [t] = useTranslation('secret_update_value');

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyValueSecretMutation();

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
			placeholder={t('placeholder')}
			allowBlank={false}
			resetOnSuccess={true}
			resetOnError={true}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			onSubmit={submitValueChange}
		/>
	);
};

export default UpdateValueSecret;
