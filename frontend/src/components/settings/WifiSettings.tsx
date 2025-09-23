import {
  Box,
  Flex,
  Heading,
  IconButton,
  SimpleGrid,
  Spacer,
} from '@chakra-ui/react';
import React, {FC, useRef} from 'react';
import {WifiNetwork, WifiStatus} from '../../api/checkInSystemApi';
import WifiNetworkList from './WifiNetworkList';
import {AddIcon} from '@chakra-ui/icons';
import {WifiNetworkAdd, WifiNetworkAddRef} from './WifiNetworkAdd';
import {WifiMode} from './WifiMode';

type Props = {
  networks: WifiNetwork[];
  onAdd: (network: WifiNetwork) => Promise<void>;
  onRemove: (ssid: string) => Promise<void>;
  onToggleWifiMode: () => Promise<void>;
  wifiStatus?: WifiStatus;
};

export const WifiSettings: FC<Props> = ({
  networks,
  onAdd,
  onRemove,
  onToggleWifiMode,
  wifiStatus,
}) => {
  const wifiNetworkAddRef = useRef<WifiNetworkAddRef>(null);

  const onAddWifiNetwork = async () => {
    wifiNetworkAddRef.current?.show();
  };

  return (
    <>
      <SimpleGrid spacing={5} columns={1}>
        <Heading>WiFi Settings</Heading>
        <Box>
          <Flex>
            <Spacer />
            <Box p="4">
              <IconButton
                colorScheme="blue"
                aria-label="Add Wifi Network"
                title="Add Wifi Network"
                icon={<AddIcon />}
                onClick={() => onAddWifiNetwork()}
              />
            </Box>
          </Flex>
          <WifiNetworkList
            networks={networks || []}
            connectedSsid={wifiStatus?.ssid}
            onRemove={onRemove}
          />
          <WifiNetworkAdd ref={wifiNetworkAddRef} onAddNetwork={onAdd} />
        </Box>
        <WifiMode
          isHotspot={wifiStatus?.mode === 'hotspot'}
          onToggle={onToggleWifiMode}
        />
      </SimpleGrid>
    </>
  );
};
