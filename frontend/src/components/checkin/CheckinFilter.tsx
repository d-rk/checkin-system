import {DownloadIcon} from '@chakra-ui/icons';
import {Box, Flex, IconButton, Spacer} from '@chakra-ui/react';
import React, {FC, useCallback} from 'react';
import {Datepicker} from '../datepicker/Datepicker';
import {CalendarDate} from '@uselessdev/datepicker';

type Props = {
  date?: Date;
  onDateChange: (date: Date) => void;
  onDownload: () => void;
};

export const CheckInFilter: FC<Props> = ({date, onDateChange, onDownload}) => {
  const handleDateChange = useCallback(
    (date: CalendarDate) => {
      onDateChange(date as Date);
    },
    [onDateChange]
  );

  return (
    <Flex>
      <Datepicker date={date} onDateChange={handleDateChange} />
      <Spacer />
      <Box p="4">
        <IconButton
          colorScheme="blue"
          aria-label="Download .csv"
          title="Download .csv"
          icon={<DownloadIcon />}
          onClick={onDownload}
        />
      </Box>
    </Flex>
  );
};
