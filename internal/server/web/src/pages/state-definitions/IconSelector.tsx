import { iconMap, iconNames, renderIcon, DEFAULT_ICON } from '../../utils/iconHelper';
import { DownOutlined } from '@ant-design/icons';
import { Button, Empty, Input, Popover, Space, Tooltip } from 'antd';
import { FC, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

export interface IconSelectorProps {
	value?: string;
	onChange?: (iconName: string | undefined) => void;
	disabled?: boolean;
}

const ICONS_PER_PAGE = 500;

const IconSelector: FC<IconSelectorProps> = ({ value, onChange, disabled }) => {
	const [t] = useTranslation('state_definition_create');
	const [open, setOpen] = useState(false);
	const [search, setSearch] = useState('');

	const filteredIcons = useMemo(() => {
		if (!search) return iconNames.slice(0, ICONS_PER_PAGE);
		const lowerSearch = search.toLowerCase();
		return iconNames.filter((name) => name.toLowerCase().includes(lowerSearch)).slice(0, ICONS_PER_PAGE);
	}, [search]);

	const handleSelect = useCallback(
		(iconName: string) => {
			onChange?.(iconName);
			setOpen(false);
			setSearch('');
		},
		[onChange]
	);

	const handleClear = useCallback(() => {
		onChange?.(DEFAULT_ICON);
		setOpen(false);
		setSearch('');
	}, [onChange]);

	const content = (
		<div style={{ width: 320 }}>
			<Input
				placeholder={t('icon_search_placeholder')}
				value={search}
				onChange={(e) => setSearch(e.target.value)}
				allowClear
				style={{ marginBottom: 8 }}
				autoFocus
			/>
			<div
				style={{
					maxHeight: 240,
					overflowY: 'auto',
					display: 'grid',
					gridTemplateColumns: 'repeat(8, 1fr)',
					gap: 4
				}}>
				{filteredIcons.length === 0 ? (
					<div style={{ gridColumn: 'span 8' }}>
						<Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t('icon_no_results')} />
					</div>
				) : (
					filteredIcons.map((iconName) => {
						const IconComponent = iconMap[iconName];
						const isSelected = value === iconName;
						return (
							<Tooltip key={iconName} title={iconName} placement="top" mouseEnterDelay={0.5}>
								<Button
									type={isSelected ? 'primary' : 'text'}
									size="small"
									style={{
										width: 32,
										height: 32,
										display: 'flex',
										alignItems: 'center',
										justifyContent: 'center',
										padding: 0
									}}
									onClick={() => handleSelect(iconName)}>
									<IconComponent style={{ fontSize: 16 }} />
								</Button>
							</Tooltip>
						);
					})
				)}
			</div>
			{filteredIcons.length >= ICONS_PER_PAGE && (
				<div style={{ marginTop: 8, fontSize: 12, color: '#999', textAlign: 'center' }}>
					{t('icon_showing_limited', { count: ICONS_PER_PAGE })}
				</div>
			)}
			<div style={{ marginTop: 8, borderTop: '1px solid #f0f0f0', paddingTop: 8 }}>
				<Button size="small" onClick={handleClear}>
					{t('icon_reset_default')}
				</Button>
			</div>
		</div>
	);

	return (
		<Popover
			content={content}
			trigger="click"
			open={open}
			onOpenChange={setOpen}
			placement="bottomLeft"
			arrow={false}>
			<Button
				disabled={disabled}
				style={{
					minWidth: 60,
					backgroundColor: 'rgba(0, 0, 0, 0.04)',
					borderColor: 'transparent'
				}}>
				<Space size={4}>
					{renderIcon(value || DEFAULT_ICON, { fontSize: 14 })}
					<DownOutlined style={{ fontSize: 10 }} />
				</Space>
			</Button>
		</Popover>
	);
};

export default IconSelector;
