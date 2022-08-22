import axios, {AxiosResponse} from 'axios';
import useSWR, {SWRResponse} from 'swr';
import {Websocket, WebsocketBuilder} from 'websocket-ts';

export type User = {
  id: number;
  name: string;
  rfid_uid: string;
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
