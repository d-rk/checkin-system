import axios, {AxiosResponse} from 'axios';
import {format} from 'date-fns';
import FileDownload from 'js-file-download';
import useSWR, {SWRResponse} from 'swr';
import {ExponentialBackoff, Websocket, WebsocketBuilder} from 'websocket-ts';
import {WEBSOCKET_BASE_URL} from './config';

export type UserFields = {
  name: string;
  memberId?: string;
  rfidUid?: string;
  role: string;
  group?: string;
  password?: string;
};

export type User = UserFields & {
  id: number;
};

export type CheckIn = {
  id: number;
  date: string;
  timestamp: string;
  user_id: number;
};

export type CheckInWithUser = CheckIn & {
  user: User;
};

export type CheckInMessage = {
  rfid_uid: string;
  check_in?: CheckIn;
};

export type CheckInDate = {
  date: string;
};

export type Clock = {
  refTimestamp: string;
  timestamp: string;
};

export type BearerToken = {
  token: string;
};

type Credentials = {
  username: string;
  password: string;
};

let lastCredentials: Credentials;

export const isCheckInMessage = (message: any): message is CheckInMessage => {
  return (message as CheckInMessage).rfid_uid !== undefined;
};

const fetcher = async (url: string) => {
  const result = await axios.get(url);
  return result.data;
};

export const apiLogin = async (
  credentials: Credentials
): Promise<BearerToken> => {
  const response = await axios.post('/api/login', credentials);

  if (response.status === 200) {
    lastCredentials = credentials;
  }

  return response.data;
};

export const refreshAccessToken = async () => {
  // not really a refresh of token atm
  return apiLogin(lastCredentials);
};

export const useUserList = (): SWRResponse<User[], Error> => {
  return useSWR<User[], Error>('/api/v1/users', fetcher);
};

export const useUserGroups = (): SWRResponse<string[], Error> => {
  return useSWR<string[], Error>('/api/v1/user-groups', fetcher);
};

export const useClock = (ref: Date): SWRResponse<Clock, Error> => {
  return useSWR<Clock, Error>(
    `/api/v1/clock?ref=${ref.toISOString()}`,
    fetcher
  );
};

export const getAuthenticatedUser = async (): Promise<User> => {
  const response = await axios.get('/api/v1/users/me');
  return response.data;
};

export const getUser = (userId: number): Promise<AxiosResponse<User>> => {
  return axios.get(`/api/v1/users/${userId}`);
};

export const deleteUser = (userId: number) => {
  return axios.delete(`/api/v1/users/${userId}`);
};

export const updateUser = (
  userId: number,
  user: UserFields
): Promise<AxiosResponse<User>> => {
  return axios.put(`/api/v1/users/${userId}`, {id: userId, ...user});
};

export const updateUserPassword = (
  userId: number,
  password: string
): Promise<AxiosResponse<User>> => {
  return axios.put(`/api/v1/users/${userId}/password`, {password: password});
};

export const addUser = (user: UserFields): Promise<AxiosResponse<User>> => {
  return axios.post('/api/v1/users', user);
};

export const useCheckInList = (
  date: Date
): SWRResponse<CheckInWithUser[], Error> => {
  return useSWR<CheckInWithUser[], Error>(
    `/api/v1/checkins/per-day?day=${format(date, 'yyyy-MM-dd')}`,
    fetcher
  );
};

export const deleteCheckIn = (checkInId: number) => {
  return axios.delete(`/api/v1/checkins/${checkInId}`);
};

export const useCheckInDates = (): SWRResponse<CheckInDate[], Error> => {
  return useSWR<CheckInDate[], Error>('/api/v1/checkins/dates', fetcher);
};

export const downloadCheckInList = async (date: Date) => {
  const response = await axios.get(
    `/api/v1/checkins/per-day?day=${format(date, 'yyyy-MM-dd')}`,
    {
      headers: {Accept: 'application/csv'},
    }
  );
  FileDownload(response.data, response.headers['x-filename'] ?? 'export.csv');
};

export const downloadAllCheckIns = async () => {
  const response = await axios.get('/api/v1/checkins/all', {
    headers: {Accept: 'application/csv'},
  });
  FileDownload(
    response.data,
    response.headers['x-filename'] ?? 'export_all.csv'
  );
};

export const useUserCheckInList = (
  userId: number
): SWRResponse<CheckIn[], Error> => {
  return useSWR<CheckIn[], Error>(`/api/v1/users/${userId}/checkins`, fetcher);
};

export const downloadUserCheckInList = async (userId: number) => {
  const response = await axios.get(`/api/v1/users/${userId}/checkins`, {
    headers: {Accept: 'application/csv'},
  });
  FileDownload(response.data, response.headers['x-filename'] ?? 'export.csv');
};

export const createWebsocket = (
  listener: (payload: any) => void
): Websocket => {
  console.log(`using websocket: ${WEBSOCKET_BASE_URL}/websocket`);
  return new WebsocketBuilder(`${WEBSOCKET_BASE_URL}/websocket`)
    .withBackoff(new ExponentialBackoff(100, 7))
    .onOpen(() => {
      console.log('opened');
    })
    .onClose((_, ev: Event) => {
      console.log('closed' + JSON.stringify(ev));
    })
    .onError((_, ev: Event) => {
      console.log('error:' + JSON.stringify(ev));
    })
    .onMessage((_, ev) => {
      const payload = JSON.parse(ev.data);
      listener(payload);
    })
    .onRetry(() => {
      console.log('retry');
    })
    .build();
};
