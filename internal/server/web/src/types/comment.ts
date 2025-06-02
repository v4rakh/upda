import { PaginatedRequestParams, PaginatedResponse } from './common';

export type CommentsResponse = {
	data: {
		content: CommentResponse[];
	} & PaginatedResponse;
};

export interface CommentResponse {
	id: string;
	author: string;
	content: string;
	updateId: string;
	createdAt: string;
	updatedAt: string;
}

export interface CommentSingleResponse {
	data: CommentResponse;
}

export type CreateCommentRequest = {
	content: string;
};

export type CommentsRequestParams = {
	updateId: string;
} & PaginatedRequestParams;

export type ModifyCommentContentRequest = {
	content: string;
};
