import { Breakpoint, ScreenMap } from 'antd/es/_util/responsiveObserver';
import useBreakpoint from 'antd/es/grid/hooks/useBreakpoint';
import { useCallback, useEffect, useState } from 'react';

export interface UseLargestScreenSize {
	largest: Breakpoint;
}

export const useLargestScreenSize = (): UseLargestScreenSize => {
	const [largest, setLargest] = useState<Breakpoint>('xs');
	const screens = useBreakpoint();

	const findLargestScreen = useCallback((s: ScreenMap): Breakpoint => {
		if (s['xxl']) {
			return 'xxl';
		} else if (s['xl']) {
			return 'xl';
		} else if (s['lg']) {
			return 'lg';
		} else if (s['md']) {
			return 'md';
		} else if (s['sm']) {
			return 'sm';
		} else {
			return 'xs';
		}
	}, []);

	useEffect(() => {
		setLargest(findLargestScreen(screens));
	}, [screens]);

	return {
		largest
	};
};
