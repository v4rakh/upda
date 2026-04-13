export interface UpdateStateDefinition {
	id: string;
	name: string;
	label: string;
	color: string;
	icon: string;
	description?: string;
	isInitial: boolean;
	skipOnNewVersion: boolean;
	sortOrder: number;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateStateDefinitionsResponse {
	data: {
		content: UpdateStateDefinition[];
	};
}

export interface UpdateStateDefinitionSingleResponse {
	data: UpdateStateDefinition;
}

export interface CreateUpdateStateDefinitionRequest {
	name: string;
	label: string;
	color: string;
	icon: string;
	description?: string;
	isInitial: boolean;
	skipOnNewVersion: boolean;
}

export interface ReorderUpdateStateDefinitionItem {
	id: string;
	sortOrder: number;
}

export interface ReorderUpdateStateDefinitionsRequest {
	items: ReorderUpdateStateDefinitionItem[];
}

export interface ModifyUpdateStateDefinitionRequest {
	name: string;
	label: string;
	color: string;
	icon: string;
	description?: string;
	isInitial: boolean;
	skipOnNewVersion: boolean;
	sortOrder: number;
}

export interface UpdateStateTransition {
	id: string;
	fromState: UpdateStateDefinition;
	toState: UpdateStateDefinition;
	createdAt: string;
	updatedAt: string;
}

export interface UpdateStateTransitionsResponse {
	data: {
		content: UpdateStateTransition[];
	};
}

export interface UpdateStateTransitionSingleResponse {
	data: UpdateStateTransition;
}

export interface CreateUpdateStateTransitionRequest {
	fromStateId: string;
	toStateId: string;
}
