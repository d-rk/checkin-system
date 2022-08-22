import React from 'react';
import {useDisclosure, useToast} from '@chakra-ui/react';
import {
  forwardRef,
  ForwardRefRenderFunction,
  useImperativeHandle,
  useState,
} from 'react';
import {getUser, User} from '../../api/checkInSystemApi';
import {errorToast} from '../../utils/toast';
import {ModalDialog} from '../ModalDialog';
import {Inputs, UserEditForm} from './UserEditForm';

export interface UserEditRef {
  show: (id?: number) => Promise<void>;
}

type Props = {};

const UserEditComponent: ForwardRefRenderFunction<UserEditRef, Props> = (
  props,
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

  function handleSubmit(userInput: Inputs) {
    alert(JSON.stringify(userInput));
  }

  return (
    <ModalDialog title="Add User" isOpen={isOpen} onClose={onClose}>
      <UserEditForm onSubmit={handleSubmit} user={user} />
    </ModalDialog>
  );
};

export const UserEdit = forwardRef(UserEditComponent);
