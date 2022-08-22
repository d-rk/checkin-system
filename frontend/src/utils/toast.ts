import {UseToastOptions} from '@chakra-ui/react';

export const errorToast = (
  message: string,
  error: unknown
): UseToastOptions => {
  console.log(`error occured: ${JSON.stringify(error)}`);

  return {
    title: 'Error occurred',
    description: message,
    status: 'error',
    duration: 9000,
    isClosable: true,
  };
};
