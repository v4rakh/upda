import Event from './Event';
import { useLazyGetEventsQuery } from '../../api/eventsApi';
import { EventName, EventResponse, EventsRequestParams } from '../../types/event';
import useEventsFilterQueryParams from '../../use/useEventsFilterQueryParams';
import AppBreadcrumb from '../common/AppBreadcrumb';
import { DownOutlined, QuestionCircleOutlined, ReloadOutlined } from '@ant-design/icons';
import { PageHeader } from '@ant-design/pro-layout';
import { Button, Col, Result, Row, Skeleton, Space, Timeline, Tooltip, Typography } from 'antd';
import parse from 'html-react-parser';
import { filter, unionBy, values } from 'lodash';
import React, { FC, useCallback, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

const EventsPage: FC = () => {
	const [t] = useTranslation('events');

	const { orderBy, order, size, skip } = useEventsFilterQueryParams();
	const [trigger, result] = useLazyGetEventsQuery();

	const [events, setEvents] = useState<EventResponse[] | undefined>(undefined);
	const [hasNext, setHasMore] = useState<boolean>(false);

	const fetchData = useCallback(
		async (s: number, o: number) => {
			const res = await trigger({ size: s, skip: o, order, orderBy } as EventsRequestParams, false);
			if (res.isSuccess && res.data.data.content) {
				const merged = unionBy(events, res.data.data.content, 'id');
				setEvents(values(merged));
				setHasMore(res.data.data.hasNext);
			}
		},
		[trigger, order, orderBy, events]
	);

	useEffect(() => {
		if (!events) {
			fetchData(size, skip);
		}
	}, [events, fetchData, size, skip]);

	const removeEvent = useCallback(
		(id: string) => {
			if (events) {
				const removed = filter(events, (e) => e.id !== id);
				setEvents(removed);
			}
		},
		[events]
	);

	const invokeReload = useCallback(() => {
		setHasMore(false);
		setEvents(undefined);
	}, []);

	const onLoadMore = useCallback(() => {
		if (result?.data?.data.hasNext) {
			fetchData(result.data.data.size, result.data.data.skip + result.data.data.size);
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
			<AppBreadcrumb items={[{ label: t('title'), active: true, path: '' }]} />
			<PageHeader
				className="pl-0"
				title={
					<Typography.Title level={4} ellipsis>
						{t('title')}
						<Tooltip placement="bottom" title={parse(t('help'))}>
							<Button icon={<QuestionCircleOutlined />} type="link" />
						</Tooltip>
					</Typography.Title>
				}
				extra={
					<Space>
						<Tooltip title={t('reload_tooltip')} placement="bottom">
							<Button
								icon={<ReloadOutlined />}
								type="link"
								onClick={invokeReload}
								loading={result.isFetching}
								disabled={result.isLoading || result.isFetching}
							/>
						</Tooltip>
					</Space>
				}
			/>
			{result.isLoading && <Skeleton loading={result.isLoading} active={result.isLoading} />}
			{result.isError && <Result status="error" title={t('error_default_loading')} />}
			{result.isSuccess && events?.length === 0 && <Result status={404} title={t('no_events')} />}
			{result.isSuccess && events && events.length > 0 && (
				<Row justify="center" align="middle">
					<Col xs={24} sm={16}>
						<Timeline
							pending={hasNext}
							pendingDot={
								<Button
									size="small"
									type="link"
									onClick={onLoadMore}
									icon={<DownOutlined />}
									loading={result.isLoading || result.isFetching}
									disabled={result.isLoading || result.isFetching}>
									{t('load_more')}
								</Button>
							}
							mode="alternate"
							reverse={false}
							items={[
								...events.map((event) => {
									return {
										children: (
											<Event
												key={event.id}
												entity={event}
												onDeleteSuccess={() => removeEvent(event.id)}
											/>
										),
										color: renderTimelineItemColor(event.name)
									};
								})
							]}
						/>
					</Col>
				</Row>
			)}
		</>
	);
};

export default EventsPage;
