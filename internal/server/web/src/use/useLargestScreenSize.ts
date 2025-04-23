import { Breakpoint, ScreenMap } from 'antd/es/_util/responsiveObserver';
import useBreakpoint from 'antd/es/grid/hooks/useBreakpoint';
import { useCallback } from 'react';

export interface UseLargestScreenSize {
	largest: Breakpoint;
}

export const useLargestScreenSize = (): UseLargestScreenSize => {
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

	return {
		largest: findLargestScreen(screens)
	};
};
