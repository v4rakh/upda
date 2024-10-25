export type SecretsResponse = {
	data: {
		content: SecretResponse[];
	};
};

export interface SecretResponse {
	id: string;
	key: string;
	value?: string;
	createdAt: string;
	updatedAt: string;
}

export interface SecretSingleResponse {
	data: SecretResponse;
}

export type CreateSecretRequest = {
	key: string;
	value: string;
};

export type ModifySecretValueRequest = {
	value: string;
};
