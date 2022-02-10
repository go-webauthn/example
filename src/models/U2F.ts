export enum RegisterResult {
    Success = 1,
    Failure,
    FailureOther,
    FailureBadRequest,
    FailureConfigurationUnsupported,
    FailureDeviceIneligible,
    FailureTimeout
}

export interface U2RRegistrationChallenge {
    appId: string;
    registerRequests: [
        {
            version: string;
            challenge: string;
        },
    ];
    registeredKeys: [
        {
            version: string;
            keyHandle: string;
            appId: string;
        }
    ] | undefined;
}

