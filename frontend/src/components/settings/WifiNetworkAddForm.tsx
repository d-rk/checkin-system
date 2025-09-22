import {
  Button,
  FormControl,
  FormErrorMessage,
  FormLabel,
  Input,
  InputGroup,
  InputRightElement,
  SimpleGrid,
} from '@chakra-ui/react';
import React, {FC} from 'react';
import {useForm} from 'react-hook-form';
import {WifiNetwork} from '../../api/checkInSystemApi';

type Props = {
  onSubmit: (network: WifiNetwork) => void;
};

export const WifiNetworkAddForm: FC<Props> = ({onSubmit}) => {
  const {
    register,
    handleSubmit,
    formState: {errors, isSubmitting},
  } = useForm<WifiNetwork>();

  const [showPassword, setShowPassword] = React.useState(false);

  return (
    <>
      <form onSubmit={handleSubmit(onSubmit)}>
        <SimpleGrid spacing={10} columns={{base: 1, md: 2}}>
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
            <InputGroup size="md">
              <Input
                pr="4.5rem"
                {...register('password', {required: 'field is required'})}
                type={showPassword ? 'text' : 'password'}
                placeholder="Enter password"
                autoComplete="new-password"
              />
              <InputRightElement width="4.5rem">
                <Button
                  h="1.75rem"
                  size="sm"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? 'Hide' : 'Show'}
                </Button>
              </InputRightElement>
            </InputGroup>
            <FormErrorMessage>
              {errors.password && errors.password.message}
            </FormErrorMessage>
          </FormControl>
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
