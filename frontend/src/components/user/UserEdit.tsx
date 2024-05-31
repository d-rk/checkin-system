import {useDisclosure, useToast} from '@chakra-ui/react';
import React, {
  forwardRef,
  ForwardRefRenderFunction,
  useImperativeHandle,
  useState,
} from 'react';
import {
  addUser,
  getUser,
  updateUser,
  updateUserPassword,
  User,
  UserFields,
  useUserGroups,
} from '../../api/checkInSystemApi';
import {errorToast} from '../../utils/toast';
import {ModalDialog} from '../ModalDialog';
import {UserEditForm} from './UserEditForm';

export interface UserEditRef {
  show: (id?: number) => Promise<void>;
}

type Props = {
  onUserEdited: (user: User) => void;
};

const UserEditComponent: ForwardRefRenderFunction<UserEditRef, Props> = (
  {onUserEdited},
  ref
) => {
  const toast = useToast();

  const {data: groups, error: fetchGroupsError} = useUserGroups();

  const {isOpen, onOpen, onClose} = useDisclosure();

  const [user, setUser] = useState<User>();

  useImperativeHandle(ref, () => ({
    show: async id => {
      if (id) {
        try {
          const user = await getUser(id);
          setUser(user.data);
        } catch (error) {
          toast(errorToast('unable to get user', error));
          return;
        }
      } else {
        setUser(undefined);
      }

      onOpen();
    },
  }));

  if (fetchGroupsError) {
    toast(errorToast('unable to fetch user groups', fetchGroupsError));
  }

  async function handleSubmit(userFields: UserFields) {
    try {
      let updateResponse;
      if (user) {
        updateResponse = await updateUser(user.id, userFields);
        if (userFields.password) {
          await updateUserPassword(user.id, userFields.password);
        }
      } else {
        updateResponse = await addUser(userFields);
      }
      onUserEdited(updateResponse.data);
      onClose();
    } catch (error) {
      toast(errorToast('unable to save user', error));
    }
  }

  return (
    <ModalDialog title="Add User" isOpen={isOpen} onClose={onClose}>
      <UserEditForm onSubmit={handleSubmit} user={user} groups={groups ?? []} />
    </ModalDialog>
  );
};

export const UserEdit = forwardRef(UserEditComponent);
