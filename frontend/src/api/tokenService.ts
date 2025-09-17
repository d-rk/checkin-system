import {BearerToken} from "./checkInSystemApi";
import axios from "axios";

const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';
const TOKEN_EXPIRY_KEY = 'tokenExpiry';

export const storeTokens = (tokenData: BearerToken) => {
  localStorage.setItem(ACCESS_TOKEN_KEY, tokenData.token);
  localStorage.setItem(REFRESH_TOKEN_KEY, tokenData.refreshToken);

  if (tokenData.expiresIn) {
    const expiryTime = Date.now() + tokenData.expiresIn * 1000;
    localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString());
  }
};

export const getStoredAccessToken = (): string | null => {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
};

export const getStoredRefreshToken = (): string | null => {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
};

export const isTokenExpired = (): boolean => {
  const expiryTime = localStorage.getItem(TOKEN_EXPIRY_KEY);
  if (!expiryTime) return true;

  return Date.now() > parseInt(expiryTime);
};

export const clearTokens = () => {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(TOKEN_EXPIRY_KEY);
};

export const setAuthHeader = (token: string) => {
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  } else {
    delete axios.defaults.headers.common['Authorization'];
  }
};

// Initialize auth header from stored token
export const initializeAuth = () => {
  const token = getStoredAccessToken();
  if (token) {
    setAuthHeader(token);
  }
};

