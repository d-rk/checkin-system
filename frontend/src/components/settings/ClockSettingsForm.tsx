import {
    Button,
    Checkbox,
    Divider,
    FormControl,
    FormErrorMessage,
    FormLabel,
    Heading,
    Input,
    SimpleGrid,
} from '@chakra-ui/react';
import React, {FC} from 'react';
import {useForm} from 'react-hook-form';
import {Clock} from '../../api/checkInSystemApi';

import {formatISO} from "date-fns";
import {fromLocaleClock, toLocaleClock} from "../../utils/time";

type Props = {
  currentClock: Clock;
  onSubmit: (newClock: Clock) => Promise<Clock>;
};

export const ClockSettingsForm: FC<Props> = ({currentClock, onSubmit}) => {

  const {
    register,
    handleSubmit,
    reset,
    formState: {errors, dirtyFields, isSubmitting},
  } = useForm<Clock>({
    defaultValues: toLocaleClock(currentClock),
  });

  const [manualClock, setManualClock] = React.useState<boolean>(false);

    const toggleManualClock = () => {
        setManualClock(prevChecked => !prevChecked);
    };

  const onSubmitInternal = async (newClock: Clock) => {
    if (dirtyFields.timestamp) {
      const updatedClock = await onSubmit(fromLocaleClock(newClock));
      reset(toLocaleClock(updatedClock));
    } else {
      const updatedClock = await onSubmit({...fromLocaleClock(newClock), timestamp: formatISO(new Date())});
      reset(toLocaleClock(updatedClock));
    }
  };

  return (
    <>
      <form onSubmit={handleSubmit(onSubmitInternal)}>
        <Heading>Clock Settings</Heading>
        <Divider />
        <SimpleGrid spacing={2} columns={1}>
          <FormControl isInvalid={errors?.refTimestamp !== undefined}>
            <FormLabel>Browser Time</FormLabel>
            <Input
              {...register('refTimestamp', {
                required: 'field is required',
              })}
              placeholder="enter refTimestamp"
              minW={235}
              disabled
              readOnly
            />
            <FormErrorMessage>
              {errors.refTimestamp && errors.refTimestamp.message}
            </FormErrorMessage>
          </FormControl>
          <FormControl isInvalid={errors?.timestamp !== undefined}>
            <FormLabel>Hardware Clock</FormLabel>
            <Input
                {...register('timestamp', {
                  required: 'field is required',
                })}
                placeholder="enter timestamp"
                disabled={!manualClock}
                readOnly={!manualClock}
            />
            <FormErrorMessage>
              {errors.timestamp && errors.timestamp.message}
            </FormErrorMessage>
          </FormControl>
            <Checkbox onChange={toggleManualClock} isChecked={manualClock}>
                Manually enter clock time
            </Checkbox>
        </SimpleGrid>

        <Button
          mt={4}
          colorScheme="blue"
          isLoading={isSubmitting}
          type="submit"
        >
          Update Hardware Clock
        </Button>
      </form>
    </>
  );
};
