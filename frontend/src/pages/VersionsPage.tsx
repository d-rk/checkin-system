import {Box, Center, useToast} from '@chakra-ui/react';
import React, {FC} from 'react';
import {useVersion} from '../api/checkInSystemApi';
import {errorToast} from '../utils/toast';
import {LoadingPage} from './LoadingPage';
import {VersionsTable} from '../components/settings/VersionsTable';

export const VersionsPage: FC = () => {
  const toast = useToast();
  const {data: version, isLoading, error} = useVersion();

  if (error) {
    toast(errorToast('unable read version', error));
  }

  if (isLoading) {
    return <LoadingPage />;
  }

  return (
    <Center>
      <Box>
        <VersionsTable backend={version} />
      </Box>
    </Center>
  );
};
