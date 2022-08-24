import axios, {AxiosResponse} from 'axios';
import {format} from 'date-fns';
import useSWR, {SWRResponse} from 'swr';
import {Websocket, WebsocketBuilder} from 'websocket-ts';

export type UserFields = {
  name: string;
  rfid_uid: string;
};

export type User = UserFields & {
  id: number;
};

export type CheckIn = {
  id: number;
  timestamp: string;
  user_id: number;
};

export type CheckInWithUser = CheckIn & {
  user: User;
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

export const useUserCheckInList = (
  userId: number
): SWRResponse<CheckIn[], Error> => {
  return useSWR<CheckIn[], Error>(`/api/v1/users/${userId}/checkins`, fetcher);
};

export const createWebsocket = (
  listener: (payload: any) => void
): Websocket => {
  return new WebsocketBuilder('ws://localhost:8080/websocket')
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
