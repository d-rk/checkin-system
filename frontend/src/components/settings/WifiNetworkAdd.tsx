import {useDisclosure, useToast} from '@chakra-ui/react';
import React, {
  forwardRef,
  ForwardRefRenderFunction,
  useImperativeHandle,
} from 'react';
import {WifiNetwork} from '../../api/checkInSystemApi';
import {errorToast} from '../../utils/toast';
import {ModalDialog} from '../ModalDialog';
import {WifiNetworkAddForm} from './WifiNetworkAddForm';

export interface WifiNetworkAddRef {
  show: () => Promise<void>;
}

type Props = {
  onAddNetwork: (network: WifiNetwork) => Promise<void>;
};

const WifiNetworkAddComponent: ForwardRefRenderFunction<
  WifiNetworkAddRef,
  Props
> = ({onAddNetwork}, ref) => {
  const toast = useToast();

  const {isOpen, onOpen, onClose} = useDisclosure();

  useImperativeHandle(ref, () => ({
    show: async () => {
      onOpen();
    },
  }));

  async function handleSubmit(network: WifiNetwork) {
    try {
      await onAddNetwork(network);
      onClose();
    } catch (error) {
      toast(errorToast('unable to add network', error));
    }
  }

  return (
    <ModalDialog title="Add Wifi Network" isOpen={isOpen} onClose={onClose}>
      <WifiNetworkAddForm onSubmit={handleSubmit} />
    </ModalDialog>
  );
};

export const WifiNetworkAdd = forwardRef(WifiNetworkAddComponent);
