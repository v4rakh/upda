import { injectEndpoints } from './index';
import EventFilterQueryParamNames from '../constants/api/eventFilterQueryParamNames';
import ApiTags from '../constants/apiTags';
import { EventSingleResponse, EventsRequestParams, EventsResponse } from '../types/event';
import ApiVersion from '../constants/apiVersion';

const TAG_LIST_ID = 'LIST';

export const eventsApi = injectEndpoints({
	endpoints: (build) => {
		return {
			getEvents: build.query<EventsResponse, EventsRequestParams>({
				query: ({ size, skip, order, orderBy, updateId }) => {
					const params = new URLSearchParams();
					if (size) {
						params.append(EventFilterQueryParamNames.SIZE, `${size}`);
					}
					if (skip) {
						params.append(EventFilterQueryParamNames.SKIP, `${skip}`);
					}
					if (order) {
						params.append(EventFilterQueryParamNames.ORDER, `${order}`);
					}
					if (orderBy) {
						params.append(EventFilterQueryParamNames.ORDER_BY, `${orderBy}`);
					}
					if (updateId) {
						params.append(EventFilterQueryParamNames.UPDATE_ID, `${updateId}`);
					}
					return { url: `${ApiVersion.V1}/events?${params.toString()}` };
				},
				providesTags: (result, error) => {
					if (!error && result?.data.content) {
						return [
							{ type: ApiTags.Events, id: TAG_LIST_ID },
							...result.data.content.map(({ id }) => ({ type: ApiTags.Events, id }))
						];
					}
					return [];
				}
			}),
			getEventById: build.query<EventSingleResponse, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/events/${id}` }),
				providesTags: (result, error) => {
					if (!error && result?.data) {
						return [{ type: ApiTags.Events, id: result.data.id }];
					}
					return [];
				}
			}),
			deleteEvent: build.mutation<void, { id: string }>({
				query: ({ id }) => ({ url: `${ApiVersion.V1}/events/${id}`, method: 'DELETE' }),
				invalidatesTags: (result, error, arg) => {
					if (error) {
						return [];
					}
					return [{ type: ApiTags.Events, id: arg.id }];
				}
			})
		};
	}
});

export const { useLazyGetEventsQuery, useGetEventByIdQuery, useDeleteEventMutation } = eventsApi;
