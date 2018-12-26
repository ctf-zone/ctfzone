import api from '../utils/api';

export default {
  state: {
    stats: {}
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload };
    }
  },
  effects: {
    async getStats() {
      const response = await api.get('/user/stats');
      const data = response.data;
      this.set({ stats: data });
      return data;
    },
    async createSolution({ id, flag }) {
      await api.post(`/user/solutions/${id}`, { flag });
    }
  }
};
