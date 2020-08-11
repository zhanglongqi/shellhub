import http from '@/store/helpers/http';

const putUser = async (data) => http().put('/user', {
  username: data.username,
  email: data.email,
  oldPassword: data.oldPassword,
  newPassword: data.newPassword,
});

export { putUser as default };
