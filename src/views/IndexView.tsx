import React, {useEffect, useState} from "react";

import u2fApi from "u2f-api";
import {
    Button,
    Container,
    Grid,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableRow,
    TextField,
    Typography
} from "@mui/material";
import { logout, login } from "../services/ClientService";
import {
    isWebauthnPlatformAuthenticatorAvailable, isWebauthnSecure,
    isWebauthnSupported,
} from "../services/WebauthnService";
import { getInfo } from "../services/APIService";
import SecurityKey from "../components/SecurityKey";

export interface Props {}

const IndexView = function (props: Props) {
    const [username, setUsername] = useState("");
    const [isSecure, setIsSecure] = useState(false);
    const [u2fSupported, setU2FSupported] = useState(false);
    const [webauthnSupported, setWebauthnSupported] = useState(false);
    const [platformAuthenticator, setPlatformAuthenticator] = useState(false);
    const [loginUsername, setLoginUsername] = useState("");

    useEffect(() => {
        (async () => {
            const info = await getInfo();

            if (info != null) {
                setBothUsername(info.username);
            }
        })()
    }, [setUsername, setLoginUsername]);

    useEffect(() => {
        setIsSecure(isWebauthnSecure());
    }, [setIsSecure]);

    useEffect(() => {
        setWebauthnSupported(isWebauthnSupported());

        (async () => {
            const wpa = await isWebauthnPlatformAuthenticatorAvailable();
            setPlatformAuthenticator(wpa);
        })()

    }, [setPlatformAuthenticator, setWebauthnSupported]);


    useEffect(() => {
        u2fApi.ensureSupport().then(
            () => setU2FSupported(true),
            () => setU2FSupported(false),
        );
    }, [setU2FSupported]);

    const setBothUsername = (username: string) => {
        setUsername(username);
        setLoginUsername(username);
    };

    const handleUsernameChangeEvent = (e: React.ChangeEvent<HTMLInputElement>) => {
        setUsername(e.target.value);
    };

    const handleLoginClick = async () => {
        const success = await login(username);
        if (success) {
            setLoginUsername(username);
        }
    };

    const handleLogoutClick = async () => {
        const success = await logout();
        if (success) {
            setLoginUsername("");
        }
    };

    return (
        <Grid
            id="privacy-root"
            container
            spacing={0}
            alignItems="center"
            justifyContent="center"
        >
            <Container maxWidth="md">
                <Grid container>
                    <Grid item xs={12}>
                        <Table sx={{minWidth: 650}}>
                            <TableHead>
                                <TableRow>
                                    <TableCell>Attribute</TableCell>
                                    <TableCell>Value</TableCell>
                                </TableRow>
                            </TableHead>
                            <TableBody>
                                <TableRow>
                                    <TableCell>Secure</TableCell>
                                    <TableCell>{isSecure ? "Yes" : "No" }</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>U2F</TableCell>
                                    <TableCell>{u2fSupported ? "Supported" : "Not Supported" }</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>Webauthn</TableCell>
                                    <TableCell>{webauthnSupported ? "Supported" : "Not Supported" }</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>Platform Authenticator</TableCell>
                                    <TableCell>{platformAuthenticator ? "Supported" : "Not Supported" }</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>Identity</TableCell>
                                    <TableCell>{loginUsername === "" ? "anonymous" : loginUsername }</TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>
                    </Grid>
                </Grid>
                <Grid container>
                    <Grid item xs={12}>
                        <Typography variant={'h4'}>Login Form</Typography>
                    </Grid>
                    <Grid item xs={12}>
                        <TextField id="username" label="Username" variant="filled" value={username} onChange={handleUsernameChangeEvent} disabled={loginUsername !== ""}/>
                    </Grid>
                    <Grid item xs={12}>
                        { loginUsername === "" ?
                            <Button onClick={async () => {await handleLoginClick();}}>Login</Button>
                            : <Button onClick={async () => {await handleLogoutClick();}}>Logout</Button>
                        }
                    </Grid>
                </Grid>
                <SecurityKey U2FSupported={u2fSupported} WebauthnSupported={webauthnSupported} LoggedIn={loginUsername !== ""} setBothUsername={setBothUsername} />
            </Container>
        </Grid>
    );
};

export default IndexView;

/*
const useStyles = makeStyles((theme: Theme) => ({
    mainContainer: {
        border: "1px solid #d6d6d6",
        borderRadius: "10px",
        padding: theme.spacing(4),
        marginTop: theme.spacing(2),
        marginBottom: theme.spacing(2),
    },
    section: {
        padding: theme.spacing(2),
    },
}));

 */
