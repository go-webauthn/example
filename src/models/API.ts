export interface ErrorResponse {
    status: "KO";
    message: string;
}

export interface Response<T> {
    status: "OK";
    data: T;
}

export interface OptionalDataResponse<T> {
    status: "OK";
    data?: T;
}

export type OptionalDataServiceResponse<T> = OptionalDataResponse<T> | ErrorResponse;
export type ServiceResponse<T> = Response<T> | ErrorResponse;
export type SignInResponse = { redirect: string } | undefined;

export interface LoginBody {
    username: string;
}