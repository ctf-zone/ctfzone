import api from '../utils/api';

export default {
  state: {
    // TODO: use array
    items: {}
  },
  reducers: {
    set: (state, payload) => {
      return { ...state, ...payload };
    },
    setItem: (state, { id, item }) => {
      return {
        ...state,
        items: {
          ...state.items,
          [id]: item
        }
      };
    }
  },
  effects: {
    async list() {
      const response = await api.get('/challenges');
      const items = response.data.reduce((result, item) => {
        result[item.challenge.id] = item;
        return result;
      }, {});

      this.set({ items });

      return items;
    },
    async get(id) {
      const response = await api.get(`/challenges/${id}`);
      const item = response.data;

      this.setItem({ id: item.challenge.id, item });

      return item;
    }
  }
};
