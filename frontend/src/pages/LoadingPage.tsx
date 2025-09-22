import {Box, Center, Heading, Spinner} from '@chakra-ui/react';
import React, {FC} from 'react';

type Props = {
  message?: string;
};

export const LoadingPage: FC<Props> = ({message}) => {
  return (
    <Center h="calc(100vh)">
      <Box>
        <Spinner size="xl" />
      </Box>
      {message && (
        <Heading pl={10} mt={-2}>
          {message}
        </Heading>
      )}
    </Center>
  );
};
