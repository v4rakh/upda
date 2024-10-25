import moment from 'moment-timezone';

export const formatDateTimeWithTimeZone = (date: number | string, timezone = moment.tz.guess()) => {
	const d = new Intl.DateTimeFormat('en-US', { dateStyle: 'medium', timeStyle: 'long', timeZone: timezone });
	return d.format(new Date(date));
};
