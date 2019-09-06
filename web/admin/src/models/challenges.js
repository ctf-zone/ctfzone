import api from '../utils/api'
import parseLinkHeader from 'parse-link-header'
import { stringify } from 'qs'

export default {
  state: {
    challenges: [],
    challenge: {},
    links: {},
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
  },
  effects: {
    async list({ link = null, filters = {} }) {
      const query = stringify(filters, { allowDots: true, arrayFormat: 'repeat' })
      const response = await api.get(link || 'challenges?' + query)
      const links = parseLinkHeader(response.headers.link) || {}
      const challenges = response.data

      this.set({ challenges, links })
      return challenges
    },
    async get(challengeId) {
      const response = await api.get(`challenges/${challengeId}`)
      const challenge = response.data
      this.set({ challenge })
      return challenge
    },
    async delete(challengeId) {
      await api.delete(`challenges/${challengeId}`)
    },
    async create(data) {
      const response = await api.post(`challenges`, data)
      const challenge = response.data
      this.set({ challenge })
      return challenge
    },
    async update(data) {
      const { id, ...rest } = data
      const response = await api.put(`challenges/${id}`, rest)
      const challenge = response.data
      this.set({ challenge })
      return challenge
    },
  },
}
