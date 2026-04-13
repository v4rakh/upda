import { EventName } from '../../types/event';
import { Typography } from 'antd';
import parse from 'html-react-parser';
import { FC, ReactNode, useMemo } from 'react';
import { useTranslation } from 'react-i18next';

const { Text } = Typography;

export interface EventTextProps {
	name: EventName;
	payload: Record<string, string>;
}

const EventText: FC<EventTextProps> = ({ name, payload }): ReactNode => {
	const [t] = useTranslation('event_text');

	// Provide fallbacks for label fields (backwards compatibility with old events)
	const enrichedPayload = useMemo(
		() => ({
			...payload,
			stateLabel: payload.stateLabel || payload.state,
			statePriorLabel: payload.statePriorLabel || payload.statePrior
		}),
		[payload]
	);

	return <Text>{parse(t(`${name.toLowerCase()}`, enrichedPayload))}</Text>;
};

export default EventText;
