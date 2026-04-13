import { PaginatedRequestParams, PaginatedResponse } from './common';
import UpdateOrder from '../constants/api/updateOrder';
import UpdateOrderBy from '../constants/api/updateOrderBy';

// Type alias for dynamic state values from the API
export type UpdateStateValue = string;

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
	state?: UpdateStateValue[];
	orderBy?: UpdateOrderBy;
	order?: UpdateOrder;
} & PaginatedRequestParams;

export interface UpdateResponse {
	id: string;
	application: string;
	provider: string;
	host: string;
	version: string;
	state: UpdateStateValue;
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
	state: UpdateStateValue;
};
