import { PaginatedRequestParams, PaginatedResponse } from './common';
import ActionInvocationOrder from '../constants/api/actionInvocationOrder';
import ActionInvocationOrderBy from '../constants/api/actionInvocationOrderBy';

export enum ActionInvocationState {
	CREATED = 'created',
	RUNNING = 'running',
	RETRYING = 'retrying',
	ERROR = 'error',
	SUCCESS = 'success'
}

export type ActionInvocationsResponse = {
	data: {
		content: ActionInvocationResponse[];
		orderBy: ActionInvocationOrderBy;
		order: ActionInvocationOrder;
	} & PaginatedResponse;
};

export type ActionInvocationsRequestParams = {
	orderBy?: ActionInvocationOrderBy;
	order?: ActionInvocationOrder;
} & PaginatedRequestParams;

export interface ActionInvocationResponse {
	id: string;
	retryCount: number;
	state: ActionInvocationState;
	message?: string;
	actionId: string;
	eventId: string;
	createdAt: string;
	updatedAt: string;
}

export interface ActionInvocationSingleResponse {
	data: ActionInvocationResponse;
}
