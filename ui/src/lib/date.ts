import { formatDistanceToNow  as fnsFormat } from 'date-fns'

export const format = (date: string) => {
  const d = new Date(Date.parse(date))
  return fnsFormat(d) + ' ago'
}
