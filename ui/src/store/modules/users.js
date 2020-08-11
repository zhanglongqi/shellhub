import putUser from '@/store/api/users';

export default {
  namespaced: true,

  state: {
  },

  getters: {
  },

  mutations: {
  },

  actions: {
    async put(context, data) {
      try {
        await putUser(data);
      } catch (err) {
        console.log('Deu ruim!');
        throw err;
      }
    },
  },
};
