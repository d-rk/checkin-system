import axios, {AxiosResponse} from 'axios';
import {format} from 'date-fns';
import FileDownload from 'js-file-download';
import useSWR, {SWRResponse} from 'swr';
import {Websocket, WebsocketBuilder} from 'websocket-ts';
import {WEBSOCKET_BASE_URL} from './config';

export type UserFields = {
  name: string;
  rfid_uid: string;
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

export const isCheckInMessage = (message: any): message is CheckInMessage => {
  return (message as CheckInMessage).rfid_uid !== undefined;
};

const fetcher = async (url: string) => {
  const result = await axios.get(url);
  return result.data;
};

export const useUserList = (): SWRResponse<User[], Error> => {
  return useSWR<User[], Error>('/api/v1/users', fetcher);
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
  return axios.put(`/api/v1/users/${userId}`, user);
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
  return new WebsocketBuilder(`${WEBSOCKET_BASE_URL}/api/v1/websocket`)
    .onOpen(() => {
      console.log('opened');
    })
    .onClose(() => {
      console.log('closed');
    })
    .onError(() => {
      console.log('error');
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
