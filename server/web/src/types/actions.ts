import { PaginatedRequestParams, PaginatedResponse } from './common';
import { EventName } from './event';
import ActionOrder from '../constants/api/actionOrder';
import ActionOrderBy from '../constants/api/actionOrderBy';

export enum ActionType {
	SHOUTRRR = 'shoutrrr'
}

export type ActionsResponse = {
	data: {
		content: ActionResponse[];
		orderBy: ActionOrderBy;
		order: ActionOrder;
	} & PaginatedResponse;
};

export type ActionsRequestParams = {
	orderBy?: ActionOrderBy;
	order?: ActionOrder;
} & PaginatedRequestParams;

export interface ActionResponse {
	id: string;
	label: string;
	type: ActionType;
	payload: ActionPayloadShoutrrr;
	matchEvent?: EventName;
	matchHost?: string;
	matchApplication?: string;
	matchProvider?: string;
	enabled: boolean;
	createdAt: string;
	updatedAt: string;
}

export interface ActionTestResponse {
	success: boolean;
	message: string;
}

export interface ActionTestSingleResponse {
	data: ActionTestResponse;
}

export interface ActionPayloadShoutrrr {
	urls: string[];
	body: string;
}

export interface ActionSingleResponse {
	data: ActionResponse;
}

export type ModifyActionLabelRequest = {
	label: string;
};

export type ModifyActionMatchEventRequest = {
	matchEvent?: EventName;
};

export type ModifyActionMatchApplicationRequest = {
	matchApplication?: string;
};

export type ModifyActionMatchHostRequest = {
	matchHost?: string;
};

export type ModifyActionMatchProviderRequest = {
	matchProvider?: string;
};

export type ModifyActionPayloadRequest = {
	type: ActionType;
	payload: ActionPayloadShoutrrr;
};

export type ModifyActionEnabledRequest = {
	enabled: boolean;
};

export type CreateActionRequest = {
	label: string;
	type: ActionType;
	matchEvent?: string;
	matchApplication?: string;
	matchHost?: string;
	matchProvider?: string;
	payload: ActionPayloadShoutrrr;
	enabled: boolean;
};

export type TestActionRequest = {
	application: string;
	provider: string;
	host: string;
	version: string;
	state: string;
};
