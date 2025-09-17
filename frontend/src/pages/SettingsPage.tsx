import {Box, Center, SimpleGrid, useToast,} from '@chakra-ui/react';
import React, {FC} from 'react';
import {Clock, setHardwareClock, useClock, WlanInfo,} from '../api/checkInSystemApi';
import {errorToast, successToast} from '../utils/toast';
import {LoadingPage} from './LoadingPage';
import {ClockSettingsForm} from "../components/settings/ClockSettingsForm";
import {parseISO, startOfMinute} from "date-fns";
import {WlanSettingsForm} from "../components/settings/WlanSettingsForm";

export const SettingsPage: FC = () => {
    const toast = useToast();
    const now = startOfMinute(new Date());
    const {data: clock, isLoading, error, mutate} = useClock(now);

    const handleClockSet = async (newClock: Clock) : Promise<Clock> => {
        try {
            await setHardwareClock(parseISO(newClock.timestamp));
            const updatedClock = await mutate();
            toast(successToast('hardware clock updated'));
            return updatedClock as Clock;
        } catch (error) {
            toast(errorToast('unable set hardware clock', error));
            throw error
        }
    }

    const handleWlanChange = async (newInfo: WlanInfo) => {
        console.log(`set wlan: ${JSON.stringify(newInfo)}`);
    }

    if (error) {
        toast(errorToast('unable read clock', error));
    }

    if (isLoading) {
        return <LoadingPage/>;
    }

    return (
        <Center>
            <Box>
                <SimpleGrid spacing={20} columns={2}>
                    <ClockSettingsForm currentClock={clock!} onSubmit={handleClockSet}/>
                    <WlanSettingsForm current={{ssid: 'tesfs', hotspotMode: true}} onSubmit={handleWlanChange}/>
                </SimpleGrid>
            </Box>
        </Center>
    );
};
