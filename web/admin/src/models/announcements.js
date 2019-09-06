import api from '../utils/api'
import parseLinkHeader from 'parse-link-header'
import { stringify } from 'qs'

export default {
  state: {
    items: [],
    item: {},
    links: {},
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async list({ link = null, filters = {} }) {
      const query = stringify(filters)
      const response = await api.get(link || query ? `announcements?${query}` : 'announcements')
      const links = parseLinkHeader(response.headers.link) || {}
      const items = response.data

      this.set({ items, links })
      return items
    },
    async get(announcementId) {
      const response = await api.get(`announcements/${announcementId}`)
      const item = response.data
      this.set({ item })
      return item
    },
    async delete(announcementId) {
      await api.delete(`announcements/${announcementId}`)
    },
    async create(data) {
      const response = await api.post(`announcements`, data)
      const item = response.data
      this.set({ item })
      return item
    },
    async update(data) {
      const { id, ...rest } = data
      const response = await api.put(`announcements/${id}`, rest)
      const item = response.data
      this.set({ item })
      return item
    },
  },
}
