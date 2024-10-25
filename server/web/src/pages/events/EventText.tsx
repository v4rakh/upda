import { EventName } from '../../types/event';
import { Typography } from 'antd';
import parse from 'html-react-parser';
import { FC } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface EventTextProps {
	name: EventName;
	payload: Record<string, string>;
}

const EventText: FC<EventTextProps> = ({ name, payload }): JSX.Element => {
	const [t] = useTranslation('event_text');

	return <Text>{parse(t(`${name.toLowerCase()}`, payload))}</Text>;
};

export default EventText;
