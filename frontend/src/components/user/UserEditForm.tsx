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
import React, {FC} from 'react';
import {useForm} from 'react-hook-form';
import {
  createWebsocket,
  isCheckInMessage,
  User,
  UserFields,
} from '../../api/checkInSystemApi';

type Props = {
  user?: User;
  onSubmit: (inputs: UserFields) => void;
};

export const UserEditForm: FC<Props> = ({user, onSubmit}) => {
  const {
    register,
    handleSubmit,
    setValue,
    formState: {errors, isSubmitting},
  } = useForm<UserFields>({
    defaultValues: {
      name: user?.name,
      rfid_uid: user?.rfid_uid,
      member_id: user?.member_id,
    },
  });

  const [rfidViaWebsocket, setRfidViaWebsocket] = React.useState(true);

  React.useMemo(
    () =>
      createWebsocket((payload: any) => {
        if (rfidViaWebsocket && isCheckInMessage(payload)) {
          setValue('rfid_uid', payload.rfid_uid);
        }
      }),
    []
  );

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
              {...register('name', {
                required: 'field is required',
              })}
              placeholder="enter name"
            />
            <FormHelperText>Name to identify user</FormHelperText>
            <FormErrorMessage>
              {errors.name && errors.name.message}
            </FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={errors?.member_id !== undefined}>
            <FormLabel>Member ID</FormLabel>
            <Input
              {...register('member_id', {})}
              placeholder="enter member id"
            />
            <FormHelperText>Member id of user</FormHelperText>
            <FormErrorMessage>
              {errors.member_id && errors.member_id.message}
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
