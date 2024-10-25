import { UpdateState } from '../types';

/**
 * Returns color for a state
 * @param state the state
 */
export const getUpdateStateColor = (state: UpdateState): string => {
	let color = 'white';
	switch (state) {
		case UpdateState.PENDING:
			color = 'deepskyblue';
			break;
		case UpdateState.APPROVED:
			color = 'limegreen';
			break;
		case UpdateState.IGNORED:
			color = 'goldenrod';
			break;
	}

	return color;
};
