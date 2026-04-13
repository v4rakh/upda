import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import InlineInputValueEditor from '../common/InlineInputValueEditor';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateNameStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateNameStateDefinition: FC<UpdateNameStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_name');
	const { apiError } = useNotification();

	const [update, { isSuccess, isError, error, isLoading }] = useUpdateUpdateStateDefinitionMutation();

	const onSubmit = useCallback(
		(val?: string) => {
			if (val !== undefined) {
				update({
					id: entity.id,
					body: {
						name: val,
						label: entity.label,
						color: entity.color,
						icon: entity.icon,
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
					conflict: t('error_conflict_update_value'),
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
			initialValue={entity.name}
			allowBlank={false}
			isLoading={isLoading}
			isError={isError}
			isSuccess={isSuccess}
			onSubmit={onSubmit}
		/>
	);
};

export default UpdateNameStateDefinition;
