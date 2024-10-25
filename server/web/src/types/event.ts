import { WindowedResponse } from './common';
import EventOrder from '../constants/api/eventOrder';
import EventOrderBy from '../constants/api/eventOrderBy';

export enum EventName {
	UPDATE_CREATED = 'update_created',
	UPDATE_UPDATED = 'update_updated',
	UPDATE_UPDATED_STATE = 'update_updated_state',
	UPDATE_UPDATED_VERSION = 'update_updated_version',
	UPDATE_DELETED = 'update_deleted'
}

export type EventsResponse = {
	data: {
		content: EventResponse[];
		orderBy: EventOrderBy;
		order: EventOrder;
	} & WindowedResponse;
};

export type EventsRequestParams = {
	orderBy?: EventOrderBy;
	order?: EventOrder;
} & WindowedResponse;

export interface EventResponse {
	id: string;
	name: EventName;
	createdAt: string;
	updatedAt: string;
	payload: Record<string, string>;
}

export interface EventSingleResponse {
	data: EventResponse;
}
