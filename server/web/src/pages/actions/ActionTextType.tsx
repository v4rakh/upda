import { ActionType } from '../../types';
import { Typography } from 'antd';
import { FC, ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface ActionTextTypeProps {
	type: ActionType;
}

const ActionTextType: FC<ActionTextTypeProps> = ({ type }): ReactNode => {
	const [t] = useTranslation('action_text_type');

	return <Text>{t(type)}</Text>;
};

export default ActionTextType;
