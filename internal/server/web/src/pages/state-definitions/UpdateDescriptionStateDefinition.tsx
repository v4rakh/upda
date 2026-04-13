import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateDescriptionStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateDescriptionStateDefinition: FC<UpdateDescriptionStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_description');
	const { apiError } = useNotification();

	const [update, { isSuccess, isError, error, isLoading }] = useUpdateUpdateStateDefinitionMutation();

	const onSubmit = useCallback(
		(val?: string) => {
			update({
				id: entity.id,
				body: {
					name: entity.name,
					label: entity.label,
					color: entity.color,
					icon: entity.icon,
					description: val || undefined,
					isInitial: entity.isInitial,
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

	return (
		<InlineInputValueEditor
			initialValue={entity.description ?? ''}
			allowBlank={true}
			isLoading={isLoading}
			isError={isError}
			isSuccess={isSuccess}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateDescriptionStateDefinition;
