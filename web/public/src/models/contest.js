import api from '../utils/api';

export default {
  state: {
    status: {}
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload };
    }
  },
  effects: {
    async getStatus() {
      const response = await api.get('/game');
      this.set({ status: response.data });
      return response.data;
    }
  }
};
