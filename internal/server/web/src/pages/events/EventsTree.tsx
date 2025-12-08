import Event from './Event';
import { useLazyGetEventsQuery } from '../../api/eventsApi';
import { EventName, EventResponse, EventsRequestParams } from '../../types/event';
import useEventsFilterQueryParams from '../../use/useEventsFilterQueryParams';
import { DownOutlined, ReloadOutlined } from '@ant-design/icons';
import { Button, Col, Flex, Result, Row, Skeleton, Timeline, Tooltip } from 'antd';
import { concat, filter, unionBy, values } from 'lodash';
import React, { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

export interface EventsTreeProps {
	updateId?: string;
}

const EventsTree: FC<EventsTreeProps> = ({ updateId }) => {
	const [t] = useTranslation('events');

	const { orderBy, order, size, skip } = useEventsFilterQueryParams();
	const [trigger, result] = useLazyGetEventsQuery();

	const [events, setEvents] = useState<EventResponse[] | undefined>(undefined);
	const [hasMore, setHasMore] = useState<boolean>(false);

	const fetchData = useCallback(
		async (s: number, o: number) => {
			const res = await trigger({ size: s, skip: o, order, orderBy, updateId } as EventsRequestParams, false);
			if (res.isSuccess && res.data.data.content) {
				const merged = unionBy(events, res.data.data.content, 'id');
				setEvents(values(merged));
				setHasMore(res.data.data.hasNext);
			}
		},
		[trigger, order, orderBy, updateId, events]
	);

	useEffect(() => {
		if (!events) {
			fetchData(size, skip);
		}
	}, [events, fetchData, size, skip]);

	const removeEvent = useCallback((id: string) => {
		setEvents((prevEvents) => filter(prevEvents, (e) => e.id !== id));
	}, []);

	const invokeReload = useCallback(() => {
		setHasMore(false);
		setEvents(undefined);
	}, []);

	const onLoadMore = useCallback(async () => {
		if (result?.data?.data.hasNext) {
			await fetchData(result.data.data.size, result.data.data.skip + result.data.data.size);
		}
	}, [fetchData, result]);

	const renderTimelineItemColor = useCallback((name: EventName): string => {
		switch (name) {
			case EventName.UPDATE_DELETED:
				return 'red';
			case EventName.UPDATE_CREATED:
				return 'green';
			default:
				return 'blue';
		}
	}, []);

	return (
		<>
			<Flex justify="end" align="center">
				<Tooltip title={t('reload_tooltip')} placement="bottom">
					<Button
						icon={<ReloadOutlined />}
						type="link"
						onClick={invokeReload}
						loading={result.isFetching}
						disabled={result.isLoading || result.isFetching}
					/>
				</Tooltip>
			</Flex>
			{(result.isLoading || result.isFetching) && (
				<Skeleton
					loading={result.isLoading || result.isFetching}
					active={result.isLoading || result.isFetching}
				/>
			)}
			{result.isError && <Result status="error" title={t('error_default_loading')} />}
			{result.isSuccess && events?.length === 0 && <Result status={404} title={t('no_events')} />}
			{result.isSuccess && events && events.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} sm={16}>
						<Timeline
							mode="alternate"
							reverse={false}
							items={concat(
								[
									...events.map((event) => {
										return {
											content: (
												<Event
													key={event.id}
													entity={event}
													onDeleteSuccess={() => removeEvent(event.id)}
												/>
											),
											color: renderTimelineItemColor(event.name)
										};
									})
								],
								hasMore
									? [
											{
												content: (
													<Button
														size="small"
														type="link"
														onClick={onLoadMore}
														icon={<DownOutlined />}
														loading={result.isLoading || result.isFetching}
														disabled={result.isLoading || result.isFetching}>
														{t('load_more')}
													</Button>
												),
												color: 'blue'
											}
										]
									: []
							)}
						/>
					</Col>
				</Row>
			)}
		</>
	);
};

export default EventsTree;
