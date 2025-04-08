import { PaginatedRequestParams, PaginatedResponse } from './common';
import UpdateOrder from '../constants/api/updateOrder';
import UpdateOrderBy from '../constants/api/updateOrderBy';

export enum UpdateState {
	PENDING = 'pending',
	APPROVED = 'approved',
	IGNORED = 'ignored'
}

export type UpdatesResponse = {
	data: {
		content: UpdateResponse[];
		orderBy: UpdateOrderBy;
		order: UpdateOrder;
	} & PaginatedResponse;
};

export type UpdatesRequestParams = {
	searchTerm?: string;
	searchIn?: string;
	state?: UpdateState[];
	orderBy?: UpdateOrderBy;
	order?: UpdateOrder;
} & PaginatedRequestParams;

export interface UpdateResponse {
	id: string;
	application: string;
	provider: string;
	host: string;
	version: string;
	state: UpdateState;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateMetadataResponse {
	metadata: any;
}

export interface UpdateSingleResponse {
	data: UpdateResponse & UpdateMetadataResponse;
}

export type ModifyUpdateStateRequest = {
	state: UpdateState;
};
