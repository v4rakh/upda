import { useLargestScreenSize } from './useLargestScreenSize';
import { useEffect, useState } from 'react';

export interface UseResponsiveGridSizeProps {
	initialSize?: number;
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
	initialSize = 3,
	xxl = 6,
	xl = 4,
	lg = 3,
	md = 3,
	sm = 2,
	xs = 1
}: UseResponsiveGridSizeProps): UseResponsiveGridSize => {
	const { largest } = useLargestScreenSize();
	const [gridSize, setGridSize] = useState<number>(initialSize);

	useEffect(() => {
		switch (largest) {
			case 'xxl':
				setGridSize(xxl);
				break;
			case 'xl':
				setGridSize(xl);
				break;
			case 'lg':
				setGridSize(lg);
				break;
			case 'md':
				setGridSize(md);
				break;
			case 'sm':
				setGridSize(sm);
				break;
			case 'xs':
				setGridSize(xs);
				break;
		}
	}, [largest]);

	return {
		gridSize
	};
};
