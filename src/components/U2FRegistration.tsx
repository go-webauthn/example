import React from "react";
import {Button, Grid, Typography} from "@mui/material";
import {performRegistrationCeremony} from "../services/U2FService";
import {RegisterResult} from "../models/U2F";

interface Props {
    setDebugMessage: React.Dispatch<React.SetStateAction<string>>;
}

const U2FRegistration = function (props: Props) {
    const handleU2FRegistrationClick = async () => {
        props.setDebugMessage("Attempting U2F registration");

        const result = await performRegistrationCeremony();

        switch (result) {
            case RegisterResult.Success:
                props.setDebugMessage("Successfully registered U2F device");
                break;
            case RegisterResult.FailureOther:
                props.setDebugMessage("Failed to register U2F device (did you cancel it?)");
                break;
            case RegisterResult.FailureTimeout:
                props.setDebugMessage("Failed to register U2F device due to a timeout");
                break;
            case RegisterResult.FailureConfigurationUnsupported:
                props.setDebugMessage("Failed to register U2F device as the configuration wasn't supported");
                break;
            case RegisterResult.FailureDeviceIneligible:
                props.setDebugMessage("Failed to register U2F device as it is ineligible");
                break;
            case RegisterResult.FailureBadRequest:
                props.setDebugMessage("Failed to register U2F device due to a malformed request");
                break;
            case RegisterResult.Failure:
                props.setDebugMessage("Failed to register U2F device due to an unknown error");
                break;
        }
    }

    return (
        <Grid container>
            <Grid item xs={12}>
                <Typography variant={'h5'}>U2F</Typography>
            </Grid>
            <Grid item xs={12}>
                <Button onClick={async () => {await handleU2FRegistrationClick();}}>Register</Button>
            </Grid>
        </Grid>
    );
}

export default U2FRegistration;