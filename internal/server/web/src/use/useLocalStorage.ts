import { useCallback, useState } from 'react';

/**
 * A hook that persists state to localStorage with automatic JSON serialization.
 * Falls back to in-memory state if localStorage is unavailable.
 */
export function useLocalStorage<T>(key: string, initialValue: T): [T, (value: T | ((prev: T) => T)) => void] {
	const [storedValue, setStoredValue] = useState<T>(() => {
		try {
			const item = window.localStorage.getItem(key);
			return item !== null ? (JSON.parse(item) as T) : initialValue;
		} catch {
			return initialValue;
		}
	});

	const setValue = useCallback(
		(value: T | ((prev: T) => T)) => {
			setStoredValue((prev) => {
				const valueToStore = value instanceof Function ? value(prev) : value;
				try {
					window.localStorage.setItem(key, JSON.stringify(valueToStore));
				} catch {
					// localStorage unavailable, continue with in-memory state
				}
				return valueToStore;
			});
		},
		[key]
	);

	return [storedValue, setValue];
}
