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
  User,
  UserFields,
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

  async function handleSubmit(userFields: UserFields) {
    try {
      let updateResponse;
      if (user) {
        updateResponse = await updateUser(user.id, userFields);
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
      <UserEditForm onSubmit={handleSubmit} user={user} />
    </ModalDialog>
  );
};

export const UserEdit = forwardRef(UserEditComponent);
