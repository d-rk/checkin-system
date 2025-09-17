import {Button, FormControl, FormErrorMessage, FormLabel, Heading, Input, SimpleGrid, Switch,} from '@chakra-ui/react';
import React, {FC} from 'react';
import {Controller, useForm} from 'react-hook-form';
import {WlanInfo,} from '../../api/checkInSystemApi';

type Props = {
  current: WlanInfo;
  onSubmit: (newInfo: WlanInfo) => void;
};

export const WlanSettingsForm: FC<Props> = ({current, onSubmit}) => {
  const {
    register,
      control,
    handleSubmit,
    formState: {errors, isSubmitting},
  } = useForm<WlanInfo>({
    defaultValues: current,
  });

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>

        <SimpleGrid spacing={5} columns={1}>
          <Heading>Wireless Settings</Heading>
          <FormControl>
            <FormLabel>Hotspot Mode</FormLabel>
            <Controller
                control={control}
                name="hotspotMode"
                render={({ field: { onChange, value, ref } }) => (
                    <Switch isChecked={value} onChange={onChange} ref={ref} />
                )}
            />
          </FormControl>
          <FormControl isInvalid={errors?.ssid !== undefined}>
            <FormLabel>SSID</FormLabel>
            <Input
                {...register('ssid', {
                  required: 'field is required',
                })}
                placeholder="enter ssid"
            />
            <FormErrorMessage>
              {errors.ssid && errors.ssid.message}
            </FormErrorMessage>
          </FormControl>
          <FormControl isInvalid={errors?.password !== undefined}>
            <FormLabel>Password</FormLabel>
            <Input
                {...register('password')}
                placeholder="change password"
            />
          </FormControl>
        </SimpleGrid>

        <Button
          mt={4}
          colorScheme="blue"
          isLoading={isSubmitting}
          type="submit"
        >
          Update WLAN Settings
        </Button>
      </form>
    </>
  );
};
