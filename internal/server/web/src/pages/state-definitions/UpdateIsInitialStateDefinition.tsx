import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { Switch } from 'antd';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateIsInitialStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateIsInitialStateDefinition: FC<UpdateIsInitialStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_is_initial');
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
					isInitial: checked,
					skipOnNewVersion: entity.skipOnNewVersion,
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

	return <Switch checked={entity.isInitial} onChange={onChange} loading={isLoading} />;
};

export default UpdateIsInitialStateDefinition;
