import {Flex, FormLabel, Select} from '@chakra-ui/react';
import React, {FC, useState} from 'react';
import {Datepicker} from '../datepicker/Datepicker';

type Props = {
  onDateChange: (date: Date) => void;
};

export const CheckInFilter: FC<Props> = ({onDateChange}) => {
  const [date, setDate] = useState<Date>(new Date());

  const handleDateChanged = (date: Date) => {
    setDate(date);
    onDateChange(date);
  };

  return (
    <Flex>
      <Datepicker date={date} onDateChange={handleDateChanged} />
    </Flex>
  );
};
