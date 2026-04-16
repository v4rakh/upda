import ActionFormPayloadShoutrrr from './ActionFormPayloadShoutrrr';
import { ActionType } from '../../types';
import { EventName } from '../../types/event';
import type { FormInstance } from 'antd';
import { FC, ReactNode } from 'react';

export interface ActionFormPayloadSwitchProps {
	isLoading: boolean;
	type: ActionType;
	form: FormInstance;
	matchEvent?: EventName | string;
}

const ActionFormPayloadSwitch: FC<ActionFormPayloadSwitchProps> = ({ isLoading, type, form, matchEvent }): ReactNode => {
	return (
		<>
			{ActionType.SHOUTRRR == type && (
				<ActionFormPayloadShoutrrr isLoading={isLoading} form={form} matchEvent={matchEvent} />
			)}
		</>
	);
};

export default ActionFormPayloadSwitch;
