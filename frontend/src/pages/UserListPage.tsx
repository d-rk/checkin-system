import {AddIcon} from '@chakra-ui/icons';
import {
  Box,
  Center,
  Flex,
  IconButton,
  Spacer,
  useToast,
} from '@chakra-ui/react';
import React, {FC, useRef} from 'react';
import {useNavigate} from 'react-router-dom';
import {
  deleteUser,
  User,
  useUserGroups,
  useUserList,
} from '../api/checkInSystemApi';
import useModals from '../components/useModals';
import {UserEdit, UserEditRef} from '../components/user/UserEdit';
import UserList from '../components/user/UserList';
import {errorToast} from '../utils/toast';
import { LoadingPage } from "./LoadingPage";

export const UserListPage: FC = () => {
  const navigate = useNavigate();
  const toast = useToast();
  const {data: users, isLoading, error, mutate} = useUserList();
  const {mutate: mutateUserGroups} = useUserGroups();
  const userEditRef = useRef<UserEditRef>(null);

  const {confirm} = useModals();

  const onShowCheckIns = async (userId: number) => {
    navigate(`/users/${userId}/checkins`);
  };

  const onCreateUser = async () => {
    userEditRef.current?.show();
  };

  const onEditUser = async (userId: number) => {
    userEditRef.current?.show(userId);
  };

  const onUserEdited = (editUser: User) => {
    const otherUsers = users?.filter(user => user.id !== editUser.id) || [];
    mutate([...otherUsers, editUser]);
    mutateUserGroups();
  };

  const onDeleteUser = async (userId: number) => {
    if (await confirm('Are you sure?')) {
      try {
        await deleteUser(userId);
        mutate(users?.filter(u => u.id !== userId));
        mutateUserGroups();
      } catch (error) {
        toast(errorToast('unable to delete user', error));
      }
    }
  };

  if (error) {
    toast(errorToast('unable to list users', error));
  }

  if (isLoading) {
    return <LoadingPage />
  }

  return (
    <Center>
      <Box>
        <Flex>
          <Spacer />
          <Box p="4">
            <IconButton
              colorScheme="blue"
              aria-label="Add User"
              title="Add User"
              icon={<AddIcon />}
              onClick={() => onCreateUser()}
            />
          </Box>
        </Flex>
        <UserList
          users={users ?? []}
          onShowCheckIns={onShowCheckIns}
          onEdit={onEditUser}
          onDelete={onDeleteUser}
        />
        <UserEdit ref={userEditRef} onUserEdited={onUserEdited} />
      </Box>
    </Center>
  );
};
