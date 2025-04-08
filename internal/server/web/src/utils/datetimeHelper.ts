import moment from 'moment-timezone';
import DateTimeStyle from '../constants/dateTimeStyle';

export const formatDateTimeWithTimeZone = (
	date: number | string,
	timeStyle: DateTimeStyle,
	dateStyle: DateTimeStyle,
	locale: string,
	timezone = moment.tz.guess()
) => {
	const dtf = new Intl.DateTimeFormat(locale, {
		timeStyle: timeStyle,
		dateStyle: dateStyle,
		timeZone: timezone
	});
	return dtf.format(new Date(date));
};
