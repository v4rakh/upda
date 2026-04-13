import { UpdateStateDefinition, UpdateStateValue } from '../types';

/**
 * Returns color for a state from state definitions
 * @param state the state name
 * @param definitions array of state definitions
 */
export const getUpdateStateColorFromDefinitions = (
	state: UpdateStateValue,
	definitions: UpdateStateDefinition[] | undefined
): string => {
	const stateDef = definitions?.find((s) => s.name === state);
	return stateDef?.color ?? 'gray';
};

/**
 * Returns label for a state from state definitions
 * @param state the state name
 * @param definitions array of state definitions
 */
export const getUpdateStateLabelFromDefinitions = (
	state: UpdateStateValue,
	definitions: UpdateStateDefinition[] | undefined
): string => {
	const stateDef = definitions?.find((s) => s.name === state);
	return stateDef?.label ?? state;
};

/**
 * Returns icon for a state from state definitions
 * @param state the state name
 * @param definitions array of state definitions
 */
export const getUpdateStateIconFromDefinitions = (
	state: UpdateStateValue,
	definitions: UpdateStateDefinition[] | undefined
): string | undefined => {
	const stateDef = definitions?.find((s) => s.name === state);
	return stateDef?.icon;
};
