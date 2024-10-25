/**
 * Sorts alphabetically and trims and converts to lower case
 * @param a comparator a
 * @param b comparator b
 */
export const sortAlphaIgnoringCase = (a: string, b: string) => {
	if (!a || !b) {
		return 0;
	}

	if (a.toLowerCase().trim() < b.toLowerCase().trim()) {
		return -1;
	}
	if (a.toLowerCase().trim() > b.toLowerCase().trim()) {
		return 1;
	}

	return 0;
};

/**
 * Sorts number values
 * @param a comparator a
 * @param b comparator b
 */
export const sortNumber = (a: number, b: number) => (a === b ? 0 : a < b ? -1 : 1);

/**
 * Sorts boolean values
 * @param a comparator a
 * @param b comparator b
 */
export const sortBool = (a: boolean, b: boolean) => (a === b ? 0 : a ? -1 : 1);
