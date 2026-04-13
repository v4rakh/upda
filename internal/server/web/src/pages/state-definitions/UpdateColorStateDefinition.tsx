import { useUpdateUpdateStateDefinitionMutation } from '../../api/updateStateDefinitionsApi';
import { UpdateStateDefinition } from '../../types';
import { useNotification } from '../../use/useNotification';
import { ColorPicker, Tag } from 'antd';
import { Color } from 'antd/es/color-picker';
import { FC, ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';

export interface UpdateColorStateDefinitionProps {
	entity: UpdateStateDefinition;
}

const UpdateColorStateDefinition: FC<UpdateColorStateDefinitionProps> = ({ entity }): ReactNode => {
	const [t] = useTranslation('state_definition_update_color');
	const { apiError } = useNotification();

	const [update, { isError, error, isLoading }] = useUpdateUpdateStateDefinitionMutation();

	const onChange = useCallback(
		(color: Color) => {
			const colorValue = color.toHexString();
			if (colorValue !== entity.color) {
				update({
					id: entity.id,
					body: {
						name: entity.name,
						label: entity.label,
						color: colorValue,
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
		<>
			<ColorPicker
				value={entity.color}
				format="hex"
				onChangeComplete={onChange}
				disabled={isLoading}
				styles={{
					popupOverlayInner: { padding: 12 }
				}}
			/>
			<Tag color={entity.color} style={{ marginLeft: 8 }}>
				{entity.color}
			</Tag>
		</>
	);
};

export default UpdateColorStateDefinition;
