import * as AntdIcons from '@ant-design/icons';
import { ReactNode } from 'react';

// Get all icon components from @ant-design/icons
// Filter out non-icon exports (like createFromIconfontCN, etc.)
const iconComponents = Object.entries(AntdIcons).filter(([name]) => name.endsWith('Outlined')) as [
	string,
	React.ComponentType<{ style?: React.CSSProperties }>
][];

// Create a map of icon name to component
export const iconMap: Record<string, React.ComponentType<{ style?: React.CSSProperties }>> = Object.fromEntries(
	iconComponents
);

// Get all available icon names
export const iconNames: string[] = iconComponents.map(([name]) => name).sort();

// Render an icon by name
export const renderIcon = (iconName?: string, style?: React.CSSProperties): ReactNode => {
	if (!iconName) return null;
	const IconComponent = iconMap[iconName];
	if (!IconComponent) return null;
	return <IconComponent style={style} />;
};

// Default icon name
export const DEFAULT_ICON = 'TagOutlined';
