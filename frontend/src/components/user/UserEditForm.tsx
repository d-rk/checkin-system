import {
  Button,
  Checkbox,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  Input,
  SimpleGrid,
} from '@chakra-ui/react';
import React, {FC, useEffect} from 'react';
import {useForm} from 'react-hook-form';
import {Websocket} from 'websocket-ts/lib/websocket';
import {createWebsocket, User} from '../../api/checkInSystemApi';

export type Inputs = {
  name: string;
  rfid_uid: string;
};

type Props = {
  user?: User;
  onSubmit: (inputs: Inputs) => void;
};

export const UserEditForm: FC<Props> = ({user, onSubmit}) => {
  const {
    register,
    handleSubmit,
    formState: {errors, isSubmitting},
  } = useForm<Inputs>();

  const [name, setName] = React.useState<string>();
  const [rfidUid, setRfidUid] = React.useState<string>();
  const [rfidViaWebsocket, setRfidViaWebsocket] = React.useState(true);

  React.useState<Websocket>(
    createWebsocket((payload: any) => {
      if (rfidViaWebsocket && payload?.rfid_uid) {
        setRfidUid(payload.rfid_uid);
      }
    })
  );

  useEffect(() => {
    setName(user?.name);
    setRfidUid(user?.rfid_uid);
  }, [user]);

  const toggleRfidViaWebsocket = () => {
    setRfidViaWebsocket(prevChecked => !prevChecked);
  };

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <SimpleGrid spacing={10}>
          <FormControl isInvalid={errors?.name !== undefined}>
            <FormLabel>Name</FormLabel>
            <Input
              {...register('name', {required: 'field is required'})}
              placeholder="enter name"
              value={name}
              onChange={event => setName(event.target.value)}
            />
            <FormHelperText>Name to identify user</FormHelperText>
            <FormErrorMessage>
              {errors.name && errors.name.message}
            </FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={errors?.rfid_uid !== undefined}>
            <FormLabel>RFID UID</FormLabel>
            <Input
              {...register('rfid_uid', {required: 'field is required'})}
              placeholder={
                rfidViaWebsocket
                  ? 'waiting for rfid token...'
                  : 'id of the rfid token'
              }
              disabled={rfidViaWebsocket}
              value={rfidUid}
              onChange={event => setRfidUid(event.target.value)}
            />
            <FormHelperText>
              {rfidViaWebsocket
                ? 'Place rfid token near reader'
                : 'Enter the id of the rfid token'}
            </FormHelperText>
            <FormErrorMessage>
              {errors.rfid_uid && errors.rfid_uid.message}
            </FormErrorMessage>
          </FormControl>

          <Checkbox onChange={toggleRfidViaWebsocket}>
            Enter rfid_uid manually
          </Checkbox>
        </SimpleGrid>

        <Button
          mt={4}
          colorScheme="blue"
          isLoading={isSubmitting}
          type="submit"
        >
          Save
        </Button>
      </form>
    </>
  );
};
