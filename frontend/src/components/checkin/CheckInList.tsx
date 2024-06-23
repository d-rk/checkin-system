import {DeleteIcon} from '@chakra-ui/icons';
import {
  IconButton,
  Table,
  TableCaption,
  TableContainer,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
} from '@chakra-ui/react';
import {formatInTimeZone} from 'date-fns-tz';
import React, {FC} from 'react';
import {CheckInWithUser} from '../../api/checkInSystemApi';

type Props = {
  checkIns: CheckInWithUser[];
  onDelete: (checkinId: number) => Promise<void>;
};

const CheckInList: FC<Props> = ({checkIns, onDelete}) => {
  return (
    <TableContainer>
      <Table>
        <TableCaption>CheckIns</TableCaption>
        <Thead>
          <Tr>
            <Th>Time</Th>
            <Th>User</Th>
            <Th>Group</Th>
            <Th textAlign="right">Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {checkIns.map(checkIn => (
            <Tr key={checkIn.id}>
              <Td>
                {formatInTimeZone(
                  new Date(checkIn.timestamp),
                  'Europe/Berlin',
                  'HH:mm:ss'
                )}
              </Td>
              <Td>{checkIn.user.name}</Td>
              <Td>{checkIn.user.group}</Td>
              <Td textAlign="right">
                <IconButton
                  colorScheme="red"
                  aria-label="Delete CheckIn"
                  title="Delete CheckIn"
                  icon={<DeleteIcon />}
                  onClick={() => onDelete(checkIn.id)}
                />
              </Td>
            </Tr>
          ))}
          {checkIns.length === 0 && (
            <Tr>
              <Td>-</Td>
              <Td></Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default CheckInList;
