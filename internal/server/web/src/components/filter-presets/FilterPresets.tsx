import FilterResetButton from './FilterResetButton';
import { useDeleteFilterPresetMutation, useGetFilterPresetsByTypeQuery } from '../../api/filterPresetsApi';
import { FilterPresetResponse, FilterPresetType } from '../../types/filterPreset';
import { useNotification } from '../../use/useNotification';
import { Flex, Skeleton, Tag } from 'antd';
import { map } from 'lodash';
import { ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

type FilterPresetsProps = {
	type: FilterPresetType;
	showFilterReset?: boolean;
	filtersActive?: boolean;
};

const FilterPresets = ({ type, showFilterReset, filtersActive }: FilterPresetsProps): ReactNode => {
	const [t] = useTranslation('filter_presets');
	const [, setSearchParams] = useSearchParams();
	const { apiError } = useNotification();

	const {
		data: getData,
		isLoading: isGetLoading,
		error: getError,
		isError: isGetError,
		isSuccess: isGetSuccess,
		isFetching: isGetFetching
	} = useGetFilterPresetsByTypeQuery(type);

	const [callDelete, { isError: isDeleteError, error: deleteError }] = useDeleteFilterPresetMutation();

	useEffect(() => {
		if (isGetError) {
			apiError({
				i18n: {
					unAuthorized: t('error_unauthorized_get'),
					forbidden: t('error_forbidden_get'),
					badRequest: t('error_bad_request_get'),
					default: t('error_default_get')
				},
				error: getError
			});
		}
	}, [t, apiError, isGetError, getError]);

	useEffect(() => {
		if (isDeleteError) {
			apiError({
				i18n: {
					unAuthorized: t('error_unauthorized_delete'),
					forbidden: t('error_forbidden_delete'),
					notFound: t('error_not_found_delete'),
					default: t('error_default_delete')
				},
				error: isDeleteError
			});
		}
	}, [t, apiError, isDeleteError, deleteError]);

	const onClick = useCallback(
		(preset: FilterPresetResponse) => {
			setSearchParams(new URLSearchParams(preset.parameters));
		},
		[setSearchParams]
	);

	const onClose = useCallback(
		(preset: FilterPresetResponse) => {
			callDelete({ id: preset.id });
		},
		[callDelete]
	);

	return (
		<>
			{isGetLoading || (isGetFetching && <Skeleton />)}
			{isGetSuccess && getData.data.content.length > 0 && (
				<Flex justify="start" align="center">
					{map(getData.data.content, (preset) => {
						return (
							<Tag
								bordered={false}
								key={preset.id}
								color={preset.color}
								onClick={() => onClick(preset)}
								closable
								onClose={() => onClose(preset)}>
								{preset.label}
							</Tag>
						);
					})}
					{showFilterReset && <FilterResetButton enabled={filtersActive} />}
				</Flex>
			)}
		</>
	);
};

export default FilterPresets;
