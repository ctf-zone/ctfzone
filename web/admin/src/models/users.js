import api from '~/utils/api'
import parseLinkHeader from 'parse-link-header'
import { stringify } from 'qs'

export default {
  state: {
    users: [],
    user: {},
    links: {},
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async list({ link = null, filters = {} }) {
      const { extra, ...restFilters } = filters

      let extraQuery = Object.
        entries(extra || {}).
        reduce((acc, v) => acc += `${v[0]}:${v[1].join(',')};`, '')

      let query = extraQuery.length !== 0
        ? `extra=${extraQuery}&`
        : ''

      query += stringify(restFilters, { allowDots: true })

      const response = await api.get(link || 'users?' + query)
      const links = parseLinkHeader(response.headers.link) || {}

      const users = response.data
      this.set({ users, links })

      return users
    },
    async get(userId) {
      const response = await api.get(`users/${userId}`)
      const user = response.data
      this.set({ user })
      return user
    },
    async delete(userId) {
      await api.delete(`users/${userId}`)
    },
    async create(data) {
      const response = await api.post(`users`, data)
      const user = response.data
      this.set({ user })
      return user
    },
    async update({ id, ...rest }) {
      const response = await api.put(`users/${id}`, rest)
      const user = response.data
      this.set({ user })
      return user
    },
  },
}
