import {UseToastOptions} from '@chakra-ui/react';
import axios from 'axios';

type ApiErrorResponse = {
  error: string;
};

export const successToast = (message: string): UseToastOptions => {
  return {
    title: 'Operation successful',
    description: message,
    status: 'success',
    duration: 3000,
    isClosable: true,
  };
};

export const errorToast = (
  message: string,
  error?: unknown
): UseToastOptions => {
  if (axios.isAxiosError(error)) {
    // add message from api call
    let details = (error.response?.data as ApiErrorResponse).error;

    if (!details) {
      details = 'unexpected error';
    }
    message = `${message}: ${details}`;
  }

  console.log(`error occured: ${JSON.stringify(error)}`);

  return {
    title: 'Error occurred',
    description: message,
    status: 'error',
    duration: 9000,
    isClosable: true,
  };
};
