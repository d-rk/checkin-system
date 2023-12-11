export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '';

export const WEBSOCKET_BASE_URL =
  API_BASE_URL !== ''
    ? API_BASE_URL.replace(/^http/, 'ws')
    : location.origin.replace(/^http/, 'ws');
