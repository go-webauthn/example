import { getEmbeddedVariable } from "../utils/Configuration";

const ExternalURL = getEmbeddedVariable("externalurl")

export const U2FRegisterPath = ExternalURL + "/api/u2f/register";
export const AttestationPath = ExternalURL + "/api/webauthn/attestation";
export const DiscoverableAttestationPath = AttestationPath + "?discoverable=true";
export const AssertionPath = ExternalURL + "/api/webauthn/assertion";
export const DiscoverableAssertionPath = AssertionPath + "?discoverable=true";
export const LoginPath = ExternalURL + "/api/login";
export const LogoutPath = ExternalURL + "/api/logout";
export const InfoPath = ExternalURL + "/api/info";