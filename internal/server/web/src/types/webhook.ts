import { PaginatedRequestParams, PaginatedResponse } from './common';
import WebhookOrder from '../constants/api/webhookOrder';
import WebhookOrderBy from '../constants/api/webhookOrderBy';

export enum WebhookType {
	GENERIC = 'generic',
	DIUN = 'diun'
}

export type WebhooksResponse = {
	data: {
		content: WebhookResponse[];
		orderBy: WebhookOrderBy;
		order: WebhookOrder;
	} & PaginatedResponse;
};

export type WebhooksRequestParams = {
	orderBy?: WebhookOrderBy;
	order?: WebhookOrder;
} & PaginatedRequestParams;

export interface WebhookResponse {
	id: string;
	label: string;
	type: WebhookType;
	token: string;
	ignoreHost: boolean;
	ignoreHostReplacement: string;
	createdAt: string;
	updatedAt: string;
}

export interface WebhookSingleResponse {
	data: WebhookResponse;
}

export type CreateWebhookRequest = {
	label: string;
	type: WebhookType;
	ignoreHost: boolean;
	ignoreHostReplacement: string;
};

export type ModifyWebhookLabelRequest = {
	label: string;
};

export type ModifyWebhookIgnoreHostRequest = {
	ignoreHost: boolean;
};

export type ModifyWebhookIgnoreHostReplacementRequest = {
	ignoreHostReplacement: string;
};
