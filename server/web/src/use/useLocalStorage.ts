import { useState } from 'react';

const useLocalStorage = <T extends string, R>(key: T, initialValue: R): [R, (value: R) => void, () => void] => {
	const [storedValue, setStoredValue] = useState<R>(() => {
		try {
			const item = window.localStorage.getItem(key);
			return item ? JSON.parse(item) : initialValue;
		} catch (error) {
			// eslint-disable-next-line no-console
			console.log(error);
			return initialValue;
		}
	});

	const setValue = (value: R): void => {
		try {
			const valueToStore = value instanceof Function ? value(storedValue) : value;
			setStoredValue(valueToStore);
			window.localStorage.setItem(key, JSON.stringify(valueToStore));
		} catch (error) {
			// eslint-disable-next-line no-console
			console.log(error);
		}
	};

	const removeValue = (): void => {
		try {
			window.localStorage.removeItem(key);
		} catch (error) {
			// eslint-disable-next-line no-console
			console.log(error);
		}
	};

	return [storedValue, setValue, removeValue];
};

export default useLocalStorage;
