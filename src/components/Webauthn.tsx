import React from "react";
import {performAssertionCeremony, performAttestationCeremony} from "../services/WebauthnService";
import {AssertionResult, AttestationResult} from "../models/Webauthn";
import {Button, Grid, Typography} from "@mui/material";
import {getInfo} from "../services/APIService";

interface Props {
    setDebugMessage: React.Dispatch<React.SetStateAction<string>>;
    Discoverable: boolean;
    setBothUsername(username: string): void;
}

const Webauthn = function(props: Props) {
    const handleDiscoverableLoginSuccess = async () => {
        const info = await getInfo();

        if (info != null) {
            props.setBothUsername(info.username);
        }
    };

    const handleAttestationClick = async (discoverable: boolean = false) => {
        props.setDebugMessage("Attempting Webauthn Attestation");


        const result = await performAttestationCeremony(discoverable);

        switch (result) {
            case AttestationResult.Success:
                props.setDebugMessage("Successful attestation.");
                break;
            case AttestationResult.FailureSupport:
                props.setDebugMessage("Your browser does not appear to support the configuration.");
                break;
            case AttestationResult.FailureSyntax:
                props.setDebugMessage("The attestation challenge was rejected as malformed or incompatible by your browser.");
                break;
            case AttestationResult.FailureWebauthnNotSupported:
                props.setDebugMessage("Your browser does not support the WebAuthN protocol.");
                break;
            case AttestationResult.FailureUserConsent:
                props.setDebugMessage("You cancelled the attestation request.");
                break;
            case AttestationResult.FailureUserVerificationOrResidentKey:
                props.setDebugMessage("Your device does not support user verification or resident keys but this was required.");
                break;
            case AttestationResult.FailureExcluded:
                props.setDebugMessage("You have registered this device already.");
                break;
            case AttestationResult.FailureUnknown:
                props.setDebugMessage("An unknown error occurred.");
                break;
        }
    };

    const handleAssertionClick = async () => {
        props.setDebugMessage("Attempting Webauthn Assertion");


        const result = await performAssertionCeremony(props.Discoverable);

        switch (result) {
            case AssertionResult.Success:
                props.setDebugMessage("Successful assertion.");

                if (props.Discoverable) {
                    await handleDiscoverableLoginSuccess();
                }
                break;
            case AssertionResult.FailureUserConsent:
                props.setDebugMessage("You cancelled the request.");
                break;
            case AssertionResult.FailureU2FFacetID:
                props.setDebugMessage("The server responded with an invalid Facet ID for the URL.");
                break;
            case AssertionResult.FailureSyntax:
                props.setDebugMessage("The assertion challenge was rejected as malformed or incompatible by your browser.");
                break;
            case AssertionResult.FailureWebauthnNotSupported:
                props.setDebugMessage("Your browser does not support the WebAuthN protocol.");
                break;
            case AssertionResult.FailureUnknownSecurity:
                props.setDebugMessage("An unknown security error occurred.");
                break;
            case AssertionResult.FailureUnknown:
                props.setDebugMessage("An unknown error occurred.");
                break;
            default:
                props.setDebugMessage("An unexpected error occurred.");
                break;
        }
    };

    return (
        <Grid container>
            <Grid item xs={12}>
                <Typography variant={'h5'}>{ props.Discoverable ? "Webauthn (Discoverable)" :  "Webauthn"}</Typography>
            </Grid>
            <Grid item xs={12}>
                { !props.Discoverable ? <Button onClick={async () => {await handleAttestationClick();}}>Attestation</Button> : null }
                { !props.Discoverable ? <Button onClick={async () => {await handleAttestationClick(true);}}>Attestation (Discoverable)</Button> : null }
                <Button onClick={async () => {await handleAssertionClick();}}>Assertion</Button>
            </Grid>
        </Grid>
    );
}

export default Webauthn;