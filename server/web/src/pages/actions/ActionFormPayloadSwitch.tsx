import ActionFormPayloadShoutrrr from './ActionFormPayloadShoutrrr';
import { ActionType } from '../../types';
import { FC } from 'react';

export interface ActionFormPayloadSwitchProps {
	isLoading: boolean;
	type: ActionType;
}

const ActionFormPayloadSwitch: FC<ActionFormPayloadSwitchProps> = ({ isLoading, type }): JSX.Element => {
	return <>{ActionType.SHOUTRRR == type && <ActionFormPayloadShoutrrr isLoading={isLoading} />}</>;
};

export default ActionFormPayloadSwitch;
