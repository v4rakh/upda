export enum FilterPresetType {
	UPDATE = 'update'
}

export type FilterPresetsResponse = {
	data: {
		content: FilterPresetResponse[];
	};
};

export interface FilterPresetResponse {
	id: string;
	type: FilterPresetType;
	label: string;
	color?: string;
	parameters: string;
}

export interface FilterPresetSingleResponse {
	data: FilterPresetResponse;
}

export type CreateFilterPresetRequest = {
	type: FilterPresetType;
	label: string;
	color?: string;
	parameters: string;
};
