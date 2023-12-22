import {
  Box,
  Input,
  Popover,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  useDisclosure,
  useOutsideClick,
} from '@chakra-ui/react';
import {
  Calendar,
  CalendarControls,
  CalendarDate,
  CalendarDays,
  CalendarMonth,
  CalendarMonthName,
  CalendarMonths,
  CalendarNextButton,
  CalendarPrevButton,
  CalendarWeek,
} from '@uselessdev/datepicker';
import {format, isValid} from 'date-fns';
import React, {FC} from 'react';
import {CheckInDay} from './CheckInDay';

type Props = {
  date?: CalendarDate;
  onDateChange: (date: CalendarDate) => void;
};

export const Datepicker: FC<Props> = ({date, onDateChange}) => {
  const [value, setValue] = React.useState('');

  const {isOpen, onOpen, onClose} = useDisclosure();

  const initialRef = React.useRef(null);
  const calendarRef = React.useRef(null);

  const handleSelectDate = (date: CalendarDate) => {
    onDateChange(date);
    setValue(() => (isValid(date) ? format(date, 'yyyy-MM-dd') : ''));
    onClose();
  };

  const match = (value: string) => value.match(/(\d{4})-(\d{2})-(\d{2})/);

  const handleInputChange = ({target}: React.ChangeEvent<HTMLInputElement>) => {
    setValue(target.value);

    if (match(target.value)) {
      onClose();
    }
  };

  useOutsideClick({
    ref: calendarRef,
    handler: onClose,
    enabled: isOpen,
  });

  React.useEffect(() => {
    if (date) {
      setValue(format(date, 'yyyy-MM-dd'));
    } else {
      setValue('');
    }
  }, [date]);

  React.useEffect(() => {
    if (match(value)) {
      const date = new Date(value);

      return onDateChange(date);
    }
  }, [value, onDateChange]);

  return (
    <>
      <Popover
        placement="auto-start"
        isOpen={isOpen}
        onClose={onClose}
        initialFocusRef={initialRef}
        isLazy
      >
        <PopoverTrigger>
          <Box onClick={onOpen} ref={initialRef} p="4">
            <Input
              placeholder="yyyy-MM-dd"
              w="min-content"
              value={value}
              onChange={handleInputChange}
            />
          </Box>
        </PopoverTrigger>

        <PopoverContent
          p={0}
          w="min-content"
          border="none"
          outline="none"
          _focus={{boxShadow: 'none'}}
          ref={calendarRef}
        >
          <Calendar
            value={{start: date}}
            onSelectDate={date => handleSelectDate(date as CalendarDate)}
            singleDateSelection
            highlightToday
            disableFutureDates
          >
            <PopoverBody p={0}>
              <CalendarControls>
                <CalendarPrevButton />
                <CalendarNextButton />
              </CalendarControls>

              <CalendarMonths>
                <CalendarMonth>
                  <CalendarMonthName />
                  <CalendarWeek />
                  <CalendarDays>
                    <CheckInDay />
                  </CalendarDays>
                </CalendarMonth>
              </CalendarMonths>
            </PopoverBody>
          </Calendar>
        </PopoverContent>
      </Popover>
    </>
  );
};
