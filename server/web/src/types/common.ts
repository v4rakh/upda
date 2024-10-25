export interface PaginatedResponse {
	page: number;
	pageSize: number;
	totalElements: number;
	totalPages: number;
}

export interface PaginatedRequestParams {
	page?: number;
	pageSize?: number;
}

export interface WindowedResponse {
	size: number;
	skip: number;
	hasNext: boolean;
}
