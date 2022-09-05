import {DownloadIcon} from '@chakra-ui/icons';
import {Box, Flex, IconButton, Spacer} from '@chakra-ui/react';
import React, {FC} from 'react';
import {Datepicker} from '../datepicker/Datepicker';

type Props = {
  date?: Date;
  onDateChange: (date: Date) => void;
  onDownload: () => void;
};

export const CheckInFilter: FC<Props> = ({date, onDateChange, onDownload}) => {
  return (
    <Flex>
      <Datepicker
        date={date}
        onDateChange={date => onDateChange(date as Date)}
      />
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
