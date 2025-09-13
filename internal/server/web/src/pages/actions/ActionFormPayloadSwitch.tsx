import ActionFormPayloadShoutrrr from './ActionFormPayloadShoutrrr';
import { ActionType } from '../../types';
import type { FormInstance } from 'antd';
import { FC, ReactNode } from 'react';

export interface ActionFormPayloadSwitchProps {
	isLoading: boolean;
	type: ActionType;
	form: FormInstance;
}

const ActionFormPayloadSwitch: FC<ActionFormPayloadSwitchProps> = ({ isLoading, type, form }): ReactNode => {
	return <>{ActionType.SHOUTRRR == type && <ActionFormPayloadShoutrrr isLoading={isLoading} form={form} />}</>;
};

export default ActionFormPayloadSwitch;
