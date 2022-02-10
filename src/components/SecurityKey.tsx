import { Grid } from "@mui/material";
import U2FRegistration from "./U2FRegistration";
import React, { useState } from "react";
import Webauthn from "./Webauthn";

interface Props {
    LoggedIn: boolean;
    U2FSupported: boolean;
    WebauthnSupported: boolean;
    setBothUsername(username: string): void;
}

const SecurityKey = function(props: Props) {
    const [debugMessage, setDebugMessage] = useState("");

    return (
        <Grid container>
            { props.U2FSupported && props.LoggedIn ? <U2FRegistration setDebugMessage={setDebugMessage} /> : null }
            { props.WebauthnSupported && props.LoggedIn ? <Webauthn Discoverable={false} setDebugMessage={setDebugMessage} setBothUsername={props.setBothUsername} /> : null }
            { props.WebauthnSupported && !props.LoggedIn ? <Webauthn Discoverable={true} setDebugMessage={setDebugMessage} setBothUsername={props.setBothUsername} /> : null }
            <Grid item xs={12}>
                Debug Message: { debugMessage }
            </Grid>
        </Grid>
    );
}

export default SecurityKey;