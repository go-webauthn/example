import axios from "axios";
import {RegisterResult, U2RRegistrationChallenge} from "../models/U2F";
import {U2FRegisterPath} from "../constants/API";
import u2f from "u2f-api";
import {ServiceResponse} from "../models/API";
import {toData} from "./APIService";

function decodeRegisterRequests(challenge: U2RRegistrationChallenge): u2f.RegisterRequest[] {
    let requests: u2f.RegisterRequest[] = [];

    for (let request of challenge.registerRequests) {
        const registerRequest: u2f.RegisterRequest = {
            appId: challenge.appId,
            challenge: request.challenge,
            version: request.version,
        }

        requests.push(registerRequest);
    }

    return requests;
}

interface U2FError {
    name: string;
    message: string;
    metaData: U2FErrorMeta;
}

interface U2FErrorMeta {
    code: number;
    type: string;
}

export async function performRegistrationCeremony(): Promise<RegisterResult> {
    const response = await axios.get<ServiceResponse<U2RRegistrationChallenge>>(U2FRegisterPath)
    if (response.status !== 200) {
        return RegisterResult.Failure;
    }

    const challenge = toData<U2RRegistrationChallenge>(response);

    if (challenge == null) {
        return RegisterResult.Failure;
    }

    const registerRequests = decodeRegisterRequests(challenge);

    let registerResponse: u2f.RegisterResponse;

    try {
        registerResponse = await u2f.register(registerRequests, [], 60);
    } catch (e) {
        const exception = e as U2FError;
        if (exception != null) {
            switch (exception.metaData.code) {
                case 1:
                    return RegisterResult.FailureOther;
                case 2:
                    return RegisterResult.FailureBadRequest;
                case 3:
                    return RegisterResult.FailureConfigurationUnsupported;
                case 4:
                    return RegisterResult.FailureDeviceIneligible;
                case 5:
                    return RegisterResult.FailureTimeout;
            }

            console.error("Unhandled U2F Code: " + exception.metaData.code + ". Exception: " + e);

            return RegisterResult.Failure;
        }

        console.error("Unhandled Exception: " + e);

        return RegisterResult.Failure;
    }

    const postResponse = await axios.post<any>(U2FRegisterPath, registerResponse);

    if (postResponse.status !== 201) {
        return RegisterResult.Failure;
    }

    return RegisterResult.Success;
}