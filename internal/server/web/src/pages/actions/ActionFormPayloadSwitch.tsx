import ActionFormPayloadShoutrrr from './ActionFormPayloadShoutrrr';
import { ActionType } from '../../types';
import { FC, ReactNode } from 'react';

export interface ActionFormPayloadSwitchProps {
	isLoading: boolean;
	type: ActionType;
}

const ActionFormPayloadSwitch: FC<ActionFormPayloadSwitchProps> = ({ isLoading, type }): ReactNode => {
	return <>{ActionType.SHOUTRRR == type && <ActionFormPayloadShoutrrr isLoading={isLoading} />}</>;
};

export default ActionFormPayloadSwitch;
