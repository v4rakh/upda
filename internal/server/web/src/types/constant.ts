export type ConstantsResponse = {
	data: {
		content: ConstantResponse[];
	};
};

export interface ConstantResponse {
	id: string;
	key: string;
	value?: string;
	createdAt: string;
	updatedAt: string;
}

export interface ConstantSingleResponse {
	data: ConstantResponse;
}

export type CreateConstantRequest = {
	key: string;
	value: string;
};

export type ModifyConstantValueRequest = {
	value: string;
};
