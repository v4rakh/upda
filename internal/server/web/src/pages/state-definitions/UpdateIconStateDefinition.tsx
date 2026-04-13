import IconSelector from './IconSelector';
import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { DEFAULT_ICON } from '../../utils/iconHelper';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateIconStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateIconStateDefinition: FC<UpdateIconStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_icon');
	const { apiError } = useNotification();

	const [update, { isError, error, isLoading }] = useUpdateUpdateStateDefinitionMutation();

	const onChange = useCallback(
		(iconName: string | undefined) => {
			if (iconName !== entity.icon) {
				update({
					id: entity.id,
					body: {
						name: entity.name,
						label: entity.label,
						color: entity.color,
						icon: iconName || DEFAULT_ICON,
						description: entity.description,
						isInitial: entity.isInitial,
						skipOnNewVersion: entity.skipOnNewVersion,
						sortOrder: entity.sortOrder
					}
				});
			}
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

	return <IconSelector value={entity.icon} onChange={onChange} disabled={isLoading} />;
};

export default UpdateIconStateDefinition;
