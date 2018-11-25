import { stringify } from 'qs'

import api from '~/utils/api'

export default {
  state: {
    items: [],
    hints: {},
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload }
    },
    setHints(state, { challengeId, hints }) {
      return {
        ...state,
        hints: {
          ...state.hints,
          [challengeId]: hints,
        },
      }
    },
  },
  effects: {
    async list({ link = null, filters = {} }) {
      const query = stringify(filters)
      const response = await api.get(link ||
        query
          ? `annonncements?${query}`
          : 'announcements'
      )
      const items = response.data

      this.set({ items })

      return items
    },
    async getHints(challengeId) {
      const response = await api.get(`announcements?challengeId=${challengeId}`)
      const hints = response.data
      this.setHints({ challengeId, hints })
      return hints
    },
  },
}
