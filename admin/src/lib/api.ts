import axios from 'axios';

const api = axios.create({
  baseURL: import.meta.env.SERVER_URL,
  withCredentials: true,
  headers: {
    'Accept': 'application/json',
  }
});

export default api;