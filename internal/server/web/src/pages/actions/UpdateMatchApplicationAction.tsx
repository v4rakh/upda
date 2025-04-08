import { useModifyMatchApplicationActionMutation } from '../../api/actionsApi';
import { useNotification } from '../../use/useNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateMatchApplicationActionProps {
	id: string;
	matchApplication?: string;
}

const UpdateMatchApplicationAction: FC<UpdateMatchApplicationActionProps> = ({ id, matchApplication }): ReactNode => {
	const [t] = useTranslation('action_update_match_application');
	const { apiError } = useNotification();

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyMatchApplicationActionMutation();

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

	const onSubmit = useCallback(
		(value?: string) => {
			if (!value || value === '') {
				modify({ id: id, body: { matchApplication: undefined } });
				return;
			}
			modify({ id: id, body: { matchApplication: value } });
		},
		[id, modify]
	);

	return (
		<InlineInputValueEditor
			initialValue={matchApplication}
			placeholder={t('all')}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateMatchApplicationAction;
