import {Box, Center, Spinner} from '@chakra-ui/react';
import React, {FC} from 'react';

export const LoadingPage: FC = () => {
  return (
    <Center h="calc(100vh)">
      <Box>
        <Spinner size="xl" />
      </Box>
    </Center>
  );
};
