import {DownloadIcon} from '@chakra-ui/icons';
import {
  Box,
  Center,
  Flex,
  Heading,
  IconButton,
  Spacer,
  useToast,
} from '@chakra-ui/react';
import React, {FC, useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import {
  createWebsocket,
  downloadUserCheckInList,
  getUser,
  isCheckInMessage,
  User,
  useUserCheckInList,
} from '../api/checkInSystemApi';
import UserCheckInList from '../components/checkin/UserCheckInList';
import {errorToast} from '../utils/toast';

export const UserCheckInListPage: FC = () => {
  const {userId} = useParams();

  const [user, setUser] = useState<User>();

  useEffect(() => {
    const getUserFunc = async () => {
      if (userId) {
        const u = await getUser(+userId);
        setUser(u.data);
      }
    };
    getUserFunc();
  }, [userId]);

  const toast = useToast();
  const {
    data: checkIns,
    error,
    mutate,
  } = useUserCheckInList(userId ? +userId : -1);

  React.useMemo(
    () =>
      createWebsocket((payload: any) => {
        if (isCheckInMessage(payload) && payload.check_in) {
          if (userId && payload.check_in.user_id === +userId) {
            mutate();
          } else {
            console.log(
              `received checking for different user ${payload.check_in.user_id} != ${userId}`
            );
          }
        } else {
          console.log(`received rfid without user: ${payload.rfid_uid}`);
        }
      }),
    [mutate, userId]
  );

  const handleDownload = () => {
    downloadUserCheckInList(user!.id);
  };

  if (error) {
    toast(errorToast('unable to list checkIns', error));
  }

  return (
    <Center>
      <Box>
        <Flex>
          <Heading as="h5">CheckIns by &quot;{user?.name}&quot;</Heading>
        </Flex>
        <Flex>
          <Spacer />
          <Box p="4">
            <IconButton
              colorScheme="blue"
              aria-label="Download .csv"
              title="Download .csv"
              icon={<DownloadIcon />}
              onClick={handleDownload}
            />
          </Box>
        </Flex>
        <UserCheckInList checkIns={checkIns ?? []} />
      </Box>
    </Center>
  );
};
