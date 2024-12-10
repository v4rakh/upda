import { useModifyMatchProviderActionMutation } from '../../api/actionsApi';
import { apiNotification } from '../common/apiNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateMatchProviderActionProps {
	id: string;
	matchProvider?: string;
}

const UpdateMatchProviderAction: FC<UpdateMatchProviderActionProps> = ({ id, matchProvider }): ReactNode => {
	const [t] = useTranslation('action_update_match_provider');

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyMatchProviderActionMutation();

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
			if (!value || value === '') {
				modify({ id: id, body: { matchProvider: undefined } });
				return;
			}
			modify({ id: id, body: { matchProvider: value } });
		},
		[id, modify]
	);

	return (
		<InlineInputValueEditor
			initialValue={matchProvider}
			placeholder={t('all')}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateMatchProviderAction;
