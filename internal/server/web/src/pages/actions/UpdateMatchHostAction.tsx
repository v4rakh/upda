import { useModifyMatchHostActionMutation } from '../../api/actionsApi';
import { useNotification } from '../../use/useNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateMatchHostActionProps {
	id: string;
	matchHost?: string;
}

const UpdateMatchHostAction: FC<UpdateMatchHostActionProps> = ({ id, matchHost }): ReactNode => {
	const [t] = useTranslation('action_update_match_host');
	const { apiError } = useNotification();

	const [modify, { isLoading, isError, isSuccess, error }] = useModifyMatchHostActionMutation();

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
				modify({ id: id, body: { matchHost: undefined } });
				return;
			}
			modify({ id: id, body: { matchHost: value } });
		},
		[id, modify]
	);

	return (
		<InlineInputValueEditor
			initialValue={matchHost}
			placeholder={t('all')}
			isLoading={isLoading}
			isSuccess={isSuccess}
			isError={isError}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateMatchHostAction;
