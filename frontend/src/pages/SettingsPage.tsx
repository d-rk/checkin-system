import {Box, Center, SimpleGrid, useToast,} from '@chakra-ui/react';
import React, {FC} from 'react';
import {
    addWifiNetwork,
    Clock,
    removeWifiNetwork,
    setHardwareClock,
    toggleWifiMode,
    useClock,
    useWifiMode,
    useWifiNetworks,
    WifiNetwork,
} from '../api/checkInSystemApi';
import {errorToast, successToast} from '../utils/toast';
import {LoadingPage} from './LoadingPage';
import {ClockSettingsForm} from "../components/settings/ClockSettingsForm";
import {parseISO, startOfMinute} from "date-fns";
import {WifiSettings} from "../components/settings/WifiSettings";

export const SettingsPage: FC = () => {
    const toast = useToast();
    const now = startOfMinute(new Date());
    const [isNetworkDown, setIsNetworkDown] = React.useState<boolean>(false);
    const {data: clock, isLoading: isClockLoading, error: clockError, mutate: mutateClock} = useClock(now);
    const {data: networks, isLoading: networksLoading, error: networksError, mutate: mutateNetworks} = useWifiNetworks();
    const {data: wifiMode, isLoading: modeLoading, error: modeError, mutate: mutateMode} = useWifiMode();

    const handleClockSet = async (newClock: Clock) : Promise<Clock> => {
        try {
            await setHardwareClock(parseISO(newClock.timestamp));
            const updatedClock = await mutateClock();
            toast(successToast('hardware clock updated'));
            return updatedClock as Clock;
        } catch (error) {
            toast(errorToast('unable set hardware clock', error));
            throw error
        }
    }

    const handleWifiAdd = async (network: WifiNetwork) => {
        try {
            await addWifiNetwork(network);
            await mutateNetworks();
        } catch (error) {
            toast(errorToast('unable to add network', error));
        }
    }

    const handleWifiRemove = async (ssid: string) => {
        console.log(`add network: ${JSON.stringify(ssid)}`);
        if (await confirm(`Really remove wifi network ${ssid}?`)) {
            try {
                await removeWifiNetwork(ssid);
                await mutateNetworks();
            } catch (error) {
                toast(errorToast('unable to remove network', error));
            }
        }
    }

    const handleToggleWifiMode = async () => {
        if (await confirm(`You will probably lose connection and will have to connect with the correct ip afterwards. continue?`)) {
            try {
                setIsNetworkDown(true);
                await toggleWifiMode();
                await mutateMode();
            } catch (error) {
                toast(errorToast('unable to toggle wifi mode', error));
            } finally {
                setIsNetworkDown(false);
            }
        }
    }

    if (clockError) {
        toast(errorToast('unable read clock', clockError));
    }
    if (networksError) {
        toast(errorToast('unable read wifi networks', networksError));
    }
    if (modeError) {
        toast(errorToast('unable read wifi mode', modeError));
    }

    if (isClockLoading || networksLoading || modeLoading || isNetworkDown) {
        return <LoadingPage message={isNetworkDown ? "no connection" : undefined}/>;
    }

    return (
        <Center>
            <Box>
                <SimpleGrid spacing={20} columns={2}>
                    <ClockSettingsForm currentClock={clock!} onSubmit={handleClockSet}/>
                    <WifiSettings networks={networks || []} onAdd={handleWifiAdd} onRemove={handleWifiRemove}
                                  isHotspot={wifiMode!} onToggleWifiMode={handleToggleWifiMode}/>
                </SimpleGrid>
            </Box>
        </Center>
    );
};
