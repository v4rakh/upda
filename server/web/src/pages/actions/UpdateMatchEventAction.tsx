import ActionSelectEvent from './ActionSelectEvent';
import { useModifyMatchEventActionMutation } from '../../api/actionsApi';
import { EventName } from '../../types/event';
import { useNotification } from '../../use/useNotification';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateMatchEventActionProps {
	id: string;
	matchEvent?: EventName;
}

const UpdateMatchEventAction: FC<UpdateMatchEventActionProps> = ({ id, matchEvent }): ReactNode => {
	const [t] = useTranslation('action_update_match_event');
	const { apiError } = useNotification();

	const [modify, { isLoading, isError, error }] = useModifyMatchEventActionMutation();

	const onChange = useCallback(
		(selected?: string) => {
			if (!selected || selected === '') {
				modify({ id: id, body: { matchEvent: undefined } });
				return;
			}
			modify({ id: id, body: { matchEvent: selected as EventName } });
		},
		[id, modify]
	);

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

	return <ActionSelectEvent loading={isLoading} onChange={onChange} name={matchEvent} />;
};

export default UpdateMatchEventAction;
