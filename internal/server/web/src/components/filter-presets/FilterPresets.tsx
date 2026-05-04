import FilterResetButton from './FilterResetButton';
import { useDeleteFilterPresetMutation, useGetFilterPresetsByTypeQuery } from '../../api/filterPresetsApi';
import { FilterPresetResponse, FilterPresetType } from '../../types/filterPreset';
import { useNotification } from '../../use/useNotification';
import { CloseOutlined } from '@ant-design/icons';
import { Flex, Popconfirm, Skeleton, Tag } from 'antd';
import { map } from 'lodash';
import { ReactNode, useCallback, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

interface FilterPresetsProps {
	type: FilterPresetType;
	showFilterReset?: boolean;
	filtersActive?: boolean;
}

const FilterPresets = ({ type, showFilterReset, filtersActive }: FilterPresetsProps): ReactNode => {
	const [t] = useTranslation('filter_presets');
	const [searchParams, setSearchParams] = useSearchParams();
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
		async (preset: FilterPresetResponse) => {
			// Check if the current search params match the preset being deleted
			const presetParams = new URLSearchParams(preset.parameters);

			// Compare parameters by checking if all preset params match current params
			let isCurrentlyActive = true;
			if (presetParams.toString() !== searchParams.toString()) {
				// Try comparing sorted entries
				const presetEntries = Array.from(presetParams.entries()).sort((a, b) => a[0].localeCompare(b[0]));
				const currentEntries = Array.from(searchParams.entries()).sort((a, b) => a[0].localeCompare(b[0]));

				if (presetEntries.length !== currentEntries.length) {
					isCurrentlyActive = false;
				} else {
					isCurrentlyActive = presetEntries.every((entry, index) => {
						const currentEntry = currentEntries[index];
						return entry[0] === currentEntry[0] && entry[1] === currentEntry[1];
					});
				}
			}

			// Delete the preset
			await callDelete({ id: preset.id });

			// If the deleted preset was currently active, clear the filters
			if (isCurrentlyActive) {
				setSearchParams({});
			}
		},
		[callDelete, searchParams, setSearchParams]
	);

	return (
		<>
			{isGetLoading || (isGetFetching && <Skeleton />)}
			{isGetSuccess && getData.data.content.length > 0 && (
				<Flex justify="start" align="center" gap="small">
					{map(getData.data.content, (preset) => {
						return (
							<Tag
								key={preset.id}
								closable
								variant="filled"
								color={preset.color}
								onClick={() => onClick(preset)}
								onClose={(e) => {
									e.preventDefault();
								}}
								closeIcon={
									<Popconfirm
										title={t('delete_confirm_title')}
										description={t('delete_confirm_description')}
										onConfirm={() => onClose(preset)}
										okText={t('delete_confirm_ok')}
										cancelText={t('delete_confirm_cancel')}
										okButtonProps={{ danger: true }}>
										<CloseOutlined />
									</Popconfirm>
								}
								style={{ cursor: 'pointer' }}>
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
