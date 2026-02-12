import { useLargestScreenSize } from './useLargestScreenSize';
import { useCallback } from 'react';

export interface UseResponsiveGridSizeProps {
	xxl?: number;
	xl?: number;
	lg?: number;
	md?: number;
	sm?: number;
	xs?: number;
}

export interface UseResponsiveGridSize {
	gridSize: number;
}

export const useResponsiveGridSize = ({
	xxl = 6,
	xl = 4,
	lg = 3,
	md = 3,
	sm = 2,
	xs = 1
}: UseResponsiveGridSizeProps): UseResponsiveGridSize => {
	const { largest } = useLargestScreenSize();

	const determineGridSize = useCallback(() => {
		switch (largest) {
			case 'xxl':
				return xxl;
			case 'xl':
				return xl;
			case 'lg':
				return lg;
			case 'md':
				return md;
			case 'sm':
				return sm;
			case 'xs':
				return xs;
			default:
				return xxl;
		}
	}, [largest]);

	return {
		gridSize: determineGridSize()
	};
};
