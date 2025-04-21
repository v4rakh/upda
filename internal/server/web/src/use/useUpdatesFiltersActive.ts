import useUpdatesFilterQueryParams from './useUpdatesFilterQueryParams';
import { useCallback } from 'react';
import UpdateOrderBy from '../constants/api/updateOrderBy';
import UpdateOrder from '../constants/api/updateOrder';
import UpdateSearchIn from '../constants/api/updateSearchIn';

export interface UseUpdatesFiltersActive {
	filtersActive: boolean;
}

const useUpdateFiltersActive = (): UseUpdatesFiltersActive => {
	const { searchTerm, searchIn, orderBy, order, state } = useUpdatesFilterQueryParams();

	const isActive = useCallback(() => {
		return (
			(searchTerm && searchTerm !== '') ||
			(searchIn && searchIn !== UpdateSearchIn.APPLICATION) ||
			(state && state.length > 0) ||
			(orderBy && orderBy !== UpdateOrderBy.UPDATED_AT) ||
			(order && order !== UpdateOrder.DESC) ||
			false
		);
	}, [searchTerm, searchIn, orderBy, order, state]);

	return { filtersActive: isActive() };
};

export default useUpdateFiltersActive;
