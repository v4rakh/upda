import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { Switch } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateSkipOnNewVersionStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateSkipOnNewVersionStateDefinition: FC<UpdateSkipOnNewVersionStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_skip_on_new_version');
	const { apiError } = useNotification();

	const [update, { isError, error, isLoading }] = useUpdateUpdateStateDefinitionMutation();

	const onChange = useCallback(
		(checked: boolean) => {
			update({
				id: entity.id,
				body: {
					name: entity.name,
					label: entity.label,
					color: entity.color,
					icon: entity.icon,
					description: entity.description,
					isInitial: entity.isInitial,
					skipOnNewVersion: checked,
					sortOrder: entity.sortOrder
				}
			});
		},
		[update, entity]
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
				error
			});
		}
	}, [error, isError, t, apiError]);

	return <Switch checked={entity.skipOnNewVersion} onChange={onChange} loading={isLoading} />;
};

export default UpdateSkipOnNewVersionStateDefinition;
