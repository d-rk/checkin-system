import {DeleteIcon, EditIcon, TimeIcon} from '@chakra-ui/icons';
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
import React, {FC} from 'react';
import {User} from '../../api/checkInSystemApi';

type Props = {
  users: User[];
  onShowCheckIns: (userId: number) => Promise<void>;
  onEdit: (userId: number) => Promise<void>;
  onDelete: (userId: number) => Promise<void>;
};

const UserList: FC<Props> = ({users, onShowCheckIns, onEdit, onDelete}) => {
  return (
    <TableContainer>
      <Table>
        <TableCaption>Users</TableCaption>
        <Thead>
          <Tr>
            <Th>Name</Th>
            <Th>Member ID</Th>
            <Th>RFID UID</Th>
            <Th textAlign="right">Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {users.map(user => (
            <Tr key={user.id}>
              <Td>{user.name}</Td>
              <Td>{user.member_id}</Td>
              <Td>{user.rfid_uid}</Td>
              <Td textAlign="right">
                <IconButton
                  aria-label="View CheckIns"
                  title="View CheckIns"
                  icon={<TimeIcon />}
                  onClick={() => onShowCheckIns(user.id)}
                />
                <IconButton
                  aria-label="Edit User"
                  title="Edit User"
                  icon={<EditIcon />}
                  onClick={() => onEdit(user.id)}
                />
                <IconButton
                  colorScheme="red"
                  aria-label="Delete User"
                  title="Delete User"
                  icon={<DeleteIcon />}
                  onClick={() => onDelete(user.id)}
                />
              </Td>
            </Tr>
          ))}
          {users.length === 0 && (
            <Tr>
              <Td>-</Td>
              <Td></Td>
              <Td></Td>
            </Tr>
          )}
        </Tbody>
      </Table>
    </TableContainer>
  );
};

export default UserList;
