import { CloseOutlined } from '@ant-design/icons';
import { Button, Tooltip } from 'antd';
import { useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

type FilterResetButtonProps = {
	enabled?: boolean;
};

const FilterResetButton = ({ enabled }: FilterResetButtonProps) => {
	const [t] = useTranslation('filter_reset_button');

	const [, setSearchQueryParams] = useSearchParams();

	const onReset = useCallback(() => {
		setSearchQueryParams(new URLSearchParams());
	}, [setSearchQueryParams]);

	return (
		<Tooltip title={t('tooltip_reset')}>
			<Button icon={<CloseOutlined />} danger type="link" onClick={onReset} disabled={!enabled} />
		</Tooltip>
	);
};

export default FilterResetButton;
