import dayjs from '@2dfire/share/Date'
export default {
	date(value, reg) {
		if (!value) return ""
		return dayjs(value).format(reg)
	}
}
