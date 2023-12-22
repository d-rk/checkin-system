import {
  Box,
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
import {Controller, useForm} from 'react-hook-form';
import {
  createWebsocket,
  isCheckInMessage,
  User,
  UserFields,
} from '../../api/checkInSystemApi';
import {CreatableSelect} from 'chakra-react-select';

type Props = {
  user?: User;
  groups: string[];
  onSubmit: (inputs: UserFields) => void;
};

type SelectOption = {
  label: string;
  value: string | null;
};

export const UserEditForm: FC<Props> = ({user, groups, onSubmit}) => {
  const {
    register,
    handleSubmit,
    control,
    setValue,
    formState: {errors, isSubmitting},
  } = useForm<UserFields>({
    defaultValues: {
      name: user?.name,
      rfidUid: user?.rfidUid,
      role: user?.role,
      memberId: user?.memberId,
    },
  });

  const [isAdmin, setIsAdmin] = React.useState(user?.role === 'ADMIN');

  const groupsWithNull: SelectOption[] = [
    {label: '-', value: null},
    ...groups.map(g => {
      return {label: g, value: g};
    }),
  ];

  const [groupOptions, setGroupOptions] = React.useState(groupsWithNull);

  React.useMemo(
    () =>
      createWebsocket((payload: any) => {
        if (isCheckInMessage(payload)) {
          setValue('rfidUid', payload.rfid_uid);
        }
      }),
    [setValue]
  );

  React.useEffect(() => {
    setValue('role', isAdmin ? 'ADMIN' : 'USER');
  }, [setValue, isAdmin]);

  React.useEffect(() => {
    // workaround: if we use the defaultValue on the Controller,
    // clearing the select no longer works
    setValue('group', user?.group);
  }, [setValue, user?.group]);

  const toggleIsAdmin = () => {
    setIsAdmin(prevChecked => !prevChecked);
  };

  const handleCreateGroup = (value: string) => {
    if (!groupOptions.some(g => g.value === value)) {
      setGroupOptions(prevState => {
        return [...prevState, {label: value, value: value}];
      });
    }
    setValue('group', value);
  };

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <SimpleGrid spacing={10} columns={2}>
          <FormControl isInvalid={errors?.name !== undefined}>
            <FormLabel>Name</FormLabel>
            <Input
              {...register('name', {
                required: 'field is required',
              })}
              placeholder="enter name"
            />
            <FormErrorMessage>
              {errors.name && errors.name.message}
            </FormErrorMessage>
          </FormControl>

          <Box />

          <FormControl isInvalid={errors?.memberId !== undefined}>
            <FormLabel>Member ID</FormLabel>
            <Input
              {...register('memberId', {})}
              placeholder="enter member id"
            />
            <FormErrorMessage>
              {errors.memberId && errors.memberId.message}
            </FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={errors?.group !== undefined}>
            <FormLabel>Group</FormLabel>

            <Controller
              control={control}
              name="group"
              render={({field: {onChange, value, ref}}) => (
                <CreatableSelect
                  ref={ref}
                  value={groupOptions.filter(g => g.value === value)}
                  onChange={val => onChange(val?.value)}
                  onCreateOption={handleCreateGroup}
                  options={groupOptions}
                  placeholder="choose group"
                  isClearable
                />
              )}
            />
            <FormErrorMessage>
              {errors.group && errors.group.message}
            </FormErrorMessage>
          </FormControl>

          <FormControl isInvalid={errors?.rfidUid !== undefined}>
            <FormLabel>RFID UID</FormLabel>
            <Input
              {...register('rfidUid', {})}
              placeholder="waiting for rfid token..."
            />
            <FormHelperText>
              Place rfid token near reader or enter id manually
            </FormHelperText>
            <FormErrorMessage>
              {errors.rfidUid && errors.rfidUid.message}
            </FormErrorMessage>
          </FormControl>

          <Box />

          <Box>
            <FormLabel>Admin Access</FormLabel>
            <Checkbox onChange={toggleIsAdmin} isChecked={isAdmin}>
              Administrator
            </Checkbox>
          </Box>

          {isAdmin && (
            <FormControl isInvalid={errors?.password !== undefined}>
              <FormLabel>Password</FormLabel>
              <Input
                {...register(
                  'password',
                  !user ? {required: 'field is required'} : {}
                )}
                placeholder="enter password"
                type="password"
              />
              <FormErrorMessage>
                {errors.password && errors.password.message}
              </FormErrorMessage>
            </FormControl>
          )}
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
